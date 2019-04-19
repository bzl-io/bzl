package fmt

import (
	"bytes"
	stdfmt "fmt"
	"path"
	"path/filepath"

	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/urfave/cli"

	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/buildtools/differ"
	"github.com/bazelbuild/buildtools/tables"
	"github.com/bazelbuild/buildtools/warn"
)

var FmtCommand = &cli.Command{
	Name:  "fmt",
	Usage: "Format build and skylark files",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mode",
			Usage: `Buildifier mode, one of "check", "diff", "fix", "print_if_changed"`,
			Value: "fix",
		},
		cli.StringFlag{
			Name:  "lint",
			Usage: `lint mode, one of "warn", "fix"`,
			Value: "fix",
		},
		cli.StringSliceFlag{
			Name:  "disable",
			Usage: "Disable rewrites on target",
		},
		cli.StringSliceFlag{
			Name:  "allow_sort",
			Usage: "Allow sorting on target",
		},
		cli.StringSliceFlag{
			Name:  "warning",
			Usage: "Warning flags",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Input type (build|workspace|bzl)",
			Value: "build",
		},
		cli.StringFlag{
			Name:  "add_tables",
			Usage: "Path to tables file that will be merged with {SOMETHING}",
			Value: "build",
		},
		cli.BoolFlag{
			Name:  "no_recursive",
			Usage: "Do not process files recursively",
		},
	},
	Action: func(c *cli.Context) error {
		exitCode, err := execute(c)
		if err != nil {
			return cli.NewExitError(stdfmt.Sprintf("fmt failed: %v", err), exitCode)
		}
		return nil
	},
}

var LintCommand = &cli.Command{
	Name:  "lint",
	Usage: "Lint build and skylark files",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mode",
			Usage: `Buildifier mode, one of "check", "diff", "fix", "print_if_changed"`,
			Value: "check",
		},
		cli.StringFlag{
			Name:  "lint",
			Usage: `lint mode, one of "warn", "fix"`,
			Value: "warn",
		},
		cli.StringSliceFlag{
			Name:  "disable",
			Usage: "Disable rewrites on target",
		},
		cli.StringSliceFlag{
			Name:  "allow_sort",
			Usage: "Allow sorting on target",
		},
		cli.StringSliceFlag{
			Name:  "warning",
			Usage: "Warning flags",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Input type (build|workspace|bzl)",
			Value: "build",
		},
		cli.StringFlag{
			Name:  "add_tables",
			Usage: "Path to tables file that will be merged with {SOMETHING}",
			Value: "build",
		},
		cli.BoolFlag{
			Name:  "no_recursive",
			Usage: "Do not process files recursively",
		},
	},
	Action: func(c *cli.Context) error {
		exitCode, err := execute(c)
		if err != nil {
			return cli.NewExitError(stdfmt.Sprintf("fmt failed: %v", err), exitCode)
		}
		return nil
	},
}

// Code is taken from buildifier.
//
// TODO(pcj): Modify upstream buildifier code such that it can be used as a
// library.
//

func execute(c *cli.Context) (int, error) {

	filePath := c.Args().First()
	tablesPath := c.String("tables_path")
	addTablesPath := c.String("add_tables_path")

	inputType := c.String("type")
	mode := c.String("mode")
	lint := c.String("lint")
	dflag := false
	warningsList := c.StringSlice("warning")
	diffProgram := c.String("differ")
	multiDiff := c.Bool("multidiff")
	recursive := !c.Bool("no_recursive")

	// Pass down debug flags into build package.
	// TODO(pcj): what do these to
	build.DisableRewrites = c.StringSlice("disable")
	build.AllowSort = c.StringSlice("allow_sort")

	if err := ValidateInputType(inputType); err != nil {
		return 2, err
	}

	if err := ValidateModes(mode, lint, dflag); err != nil {
		return 2, err
	}

	warningsList, err := ValidateWarnings(strings.Join(warningsList, ","), &warn.AllWarnings, &warn.DefaultWarnings)
	if err != nil {
		return 2, err
	}

	// If the path flag is set, must only be formatting a single file.
	// It doesn't make sense for multiple files to have the same path.
	if (filePath != "" || mode == "print_if_changed") && len(c.Args()) > 1 {
		return 2, stdfmt.Errorf("can only format one file when using -path flag or -mode=print_if_changed")
	}

	if tablesPath != "" {
		if err := tables.ParseAndUpdateJSONDefinitions(tablesPath, false); err != nil {
			return 2, stdfmt.Errorf("failed to parse %s for -tables: %s", tablesPath, err)
		}
	}

	if addTablesPath != "" {
		if err := tables.ParseAndUpdateJSONDefinitions(addTablesPath, true); err != nil {
			return 2, stdfmt.Errorf("failed to parse %s for -add_tables: %s", addTablesPath, err)
		}
	}

	differ, deprecationWarning := differ.Find()
	if diffProgram != "" {
		differ.Cmd = diffProgram
		differ.MultiDiff = multiDiff
	} else {
		if deprecationWarning && mode == "diff" {
			return 2, stdfmt.Errorf("selecting diff program with the BUILDIFIER_DIFF, BUILDIFIER_MULTIDIFF, and DISPLAY environment variables is deprecated, use flags -diff_command and -multi_diff instead")
		}
	}
	diff = differ

	if filePath == "-" {
		// Read from stdin, write to stdout.
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return 2, err
		}
		if mode == "fix" {
			mode = "pipe"
		}
		processFile(mode, filePath, "", data, inputType, lint, warningsList, false)
	} else {
		files := c.Args()
		if len(files) == 0 {
			wd, err := os.Getwd()
			if err != nil {
				return 2, err
			}
			files = []string{wd}
		}
		if recursive {
			var err error
			files, err = ExpandDirectories(files)
			if err != nil {
				return 3, err
			}
		}
		processFiles(mode, filePath, files, inputType, lint, warningsList)
	}

	if err := diff.Run(); err != nil {
		return 2, err
	}

	for _, file := range toRemove {
		os.Remove(file)
	}

	return exitCode, nil
}

func processFiles(mode string, filePath string, files []string, inputType, lint string, warningsList []string) {
	// Decide how many file reads to run in parallel.
	// At most 100, and at most one per 10 input files.
	nworker := 100
	if n := (len(files) + 9) / 10; nworker > n {
		nworker = n
	}
	runtime.GOMAXPROCS(nworker + 1)

	// Start nworker workers reading stripes of the input
	// argument list and sending the resulting data on
	// separate channels. file[k] is read by worker k%nworker
	// and delivered on ch[k%nworker].
	type result struct {
		file string
		data []byte
		err  error
	}

	ch := make([]chan result, nworker)
	for i := 0; i < nworker; i++ {
		ch[i] = make(chan result, 1)
		go func(i int) {
			for j := i; j < len(files); j += nworker {
				file := files[j]
				data, err := ioutil.ReadFile(file)
				ch[i] <- result{file, data, err}
			}
		}(i)
	}

	// Process files. The processing still runs in a single goroutine
	// in sequence. Only the reading of the files has been parallelized.
	// The goal is to optimize for runs where most files are already
	// formatted correctly, so that reading is the bulk of the I/O.
	for i, file := range files {
		res := <-ch[i%nworker]
		if res.file != file {
			stdfmt.Fprintf(os.Stderr, "buildifier: internal phase error: got %s for %s", res.file, file)
			os.Exit(3)
		}
		if res.err != nil {
			stdfmt.Fprintf(os.Stderr, "buildifier: %v\n", res.err)
			exitCode = 3
			continue
		}
		processFile(mode, filePath, file, res.data, inputType, lint, warningsList, len(files) > 1)
	}
}

// exitCode is the code to use when exiting the program.
// The codes used by buildifier are:
//
// 0: success, everything went well
// 1: syntax errors in input
// 2: usage errors: invoked incorrectly
// 3: unexpected runtime errors: file I/O problems or internal bugs
// 4: check mode failed (reformat is needed)
var exitCode = 0

// toRemove is a list of files to remove before exiting.
var toRemove []string

// diff is the differ to use when mode == "diff".
var diff *differ.Differ

// processFile processes a single file containing data.
// It has been read from filename and should be written back if fixing.
func processFile(mode string, filePath, filename string, data []byte, inputType, lint string, warningsList []string, displayFileNames bool) {
	defer func() {
		if err := recover(); err != nil {
			stdfmt.Fprintf(os.Stderr, "buildifier: %s: internal error: %v\n", filename, err)
			exitCode = 3
		}
	}()

	parser := GetParser(inputType)

	f, err := parser(filename, data)
	if err != nil {
		// Do not use buildifier: prefix on this error.
		// Since it is a parse error, it begins with file:line:
		// and we want that to be the first thing in the error.
		stdfmt.Fprintf(os.Stderr, "%v\n", err)
		if exitCode < 1 {
			exitCode = 1
		}
		return
	}

	pkg := GetPackageName(filename)
	verbose := true
	if Lint(f, pkg, lint, warningsList, verbose) {
		exitCode = 4
	}

	if filePath != "" {
		f.Path = filePath
	}

	beforeRewrite := build.Format(f)
	var info build.RewriteInfo
	build.Rewrite(f, &info)
	showlog := true

	ndata := build.Format(f)

	switch mode {
	case "check":
		// check mode: print names of files that need formatting.
		if !bytes.Equal(data, ndata) {
			// Print:
			//	name # list of what changed
			reformat := ""
			if !bytes.Equal(data, beforeRewrite) {
				reformat = " reformat"
			}
			log := ""
			if len(info.Log) > 0 && showlog {
				sort.Strings(info.Log)
				var uniq []string
				last := ""
				for _, s := range info.Log {
					if s != last {
						last = s
						uniq = append(uniq, s)
					}
				}
				log = " " + strings.Join(uniq, " ")
			}
			stdfmt.Printf("%s #%s %s%s\n", filename, reformat, &info, log)
			exitCode = 4
		}
		return

	case "diff":
		// diff mode: run diff on old and new.
		if bytes.Equal(data, ndata) {
			return
		}
		outfile, err := WriteTemp(ndata)
		if err != nil {
			toRemove = append(toRemove, outfile)
			stdfmt.Fprintf(os.Stderr, "buildifier: %v\n", err)
			exitCode = 3
			return
		}
		infile := filename
		if filename == "" {
			// data was read from standard filename.
			// Write it to a temporary file so diff can read it.
			infile, err = WriteTemp(data)
			if err != nil {
				toRemove = append(toRemove, infile)
				stdfmt.Fprintf(os.Stderr, "buildifier: %v\n", err)
				exitCode = 3
				return
			}
		}
		if displayFileNames {
			stdfmt.Fprintf(os.Stderr, "%v:\n", filename)
		}
		if err := diff.Show(infile, outfile); err != nil {
			stdfmt.Fprintf(os.Stderr, "%v\n", err)
			exitCode = 4
		}

	case "pipe":
		// pipe mode - reading from stdin, writing to stdout.
		// ("pipe" is not from the command line; it is set above in main.)
		os.Stdout.Write(ndata)
		return

	case "fix":
		// fix mode: update files in place as needed.
		if bytes.Equal(data, ndata) {
			return
		}

		err := ioutil.WriteFile(filename, ndata, 0666)
		if err != nil {
			stdfmt.Fprintf(os.Stderr, "buildifier: %s\n", err)
			exitCode = 3
			return
		}

		if verbose {
			stdfmt.Fprintf(os.Stderr, "fixed %s\n", filename)
		}
	case "print_if_changed":
		if bytes.Equal(data, ndata) {
			return
		}

		if _, err := os.Stdout.Write(ndata); err != nil {
			stdfmt.Fprintf(os.Stderr, "buildifier: error writing output: %v\n", err)
			exitCode = 3
			return
		}
	}
}

// utils package

func isStarlarkFile(filename string) bool {
	basename := strings.ToLower(filepath.Base(filename))
	ext := filepath.Ext(basename)
	switch ext {
	case ".bzl", ".sky":
		return true
	case ".proto", "protodevel":
		return false
	}
	base := basename[:len(basename)-len(ext)]
	switch {
	case ext == ".build" || base == "build":
		return true
	case ext == ".workspace" || base == "workspace":
		return true
	}
	return false
}

// ExpandDirectories takes a list of file/directory names and returns a list
// with file names by traversing each directory recursively and searching for
// relevant Starlark files.
func ExpandDirectories(args []string) ([]string, error) {
	files := []string{}
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return []string{}, err
		}
		if !info.IsDir() {
			files = append(files, arg)
			continue
		}
		err = filepath.Walk(arg, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return err
			}
			if isStarlarkFile(path) {
				files = append(files, path)
			}
			return err
		})
		if err != nil {
			return []string{}, err
		}
	}
	return files, nil
}

// GetParser returns a parser for a given file type
func GetParser(inputType string) func(filename string, data []byte) (*build.File, error) {
	switch inputType {
	case "build":
		return build.ParseBuild
	case "bzl":
		return build.ParseBzl
	case "auto":
		return build.Parse
	case "workspace":
		return build.ParseWorkspace
	default:
		return build.ParseDefault
	}
}

// WriteTemp writes data to a temporary file and returns the name of the file.
func WriteTemp(data []byte) (file string, err error) {
	f, err := ioutil.TempFile("", "buildifier-tmp-")
	if err != nil {
		return "", stdfmt.Errorf("creating temporary file: %v", err)
	}
	defer f.Close()
	name := f.Name()
	if _, err := f.Write(data); err != nil {
		return "", stdfmt.Errorf("writing temporary file: %v", err)
	}
	return name, nil
}

// GetPackageName returns the package name of a file by searching for a WORKSPACE file
func GetPackageName(filename string) string {
	dirs := filepath.SplitList(path.Dir(filename))
	parent := ""
	index := len(dirs) - 1
	for i, chunk := range dirs {
		parent = path.Join(parent, chunk)
		metadata := path.Join(parent, "METADATA")
		if _, err := os.Stat(metadata); !os.IsNotExist(err) {
			index = i
		}
	}
	return strings.Join(dirs[index+1:], "/")
}

// Lint calls the linter and returns true if there are any unresolved warnings
func Lint(f *build.File, pkg, lint string, warningsList []string, verbose bool) bool {
	switch lint {
	case "warn":
		warnings := warn.FileWarnings(f, pkg, warningsList, false)
		warn.PrintWarnings(f, warnings, false)
		return len(warnings) > 0
	case "fix":
		warn.FixWarnings(f, pkg, warningsList, verbose)
	}
	return false
}

// ValidateInputType validates the value of --type
func ValidateInputType(inputType string) error {
	switch inputType {
	case "build", "bzl", "workspace", "default", "auto":
		return nil

	default:
		return stdfmt.Errorf("unrecognized input type %s; valid types are build, bzl, workspace, default, auto", inputType)
	}
}

// isRecognizedMode checks whether the given mode is one of the valid modes.
func isRecognizedMode(validModes []string, mode string) bool {
	for _, m := range validModes {
		if mode == m {
			return true
		}
	}
	return false
}

// ValidateModes validates flags --mode, --lint, and -d
func ValidateModes(mode, lint string, dflag bool, additionalModes ...string) error {
	if dflag {
		if mode != "" {
			return stdfmt.Errorf("cannot specify both -d and -mode flags")
		}
		mode = "diff"
	}

	// Check mode.
	validModes := []string{"check", "diff", "fix", "print_if_changed"}
	validModes = append(validModes, additionalModes...)

	if mode == "" {
		mode = "fix"
	} else if !isRecognizedMode(validModes, mode) {
		return stdfmt.Errorf("unrecognized mode %s; valid modes are %s", mode, strings.Join(validModes, ", "))
	}

	// Check lint mode.
	switch lint {
	case "":
		lint = "off"

	case "off", "warn":
		// ok

	case "fix":
		if mode != "fix" {
			return stdfmt.Errorf("--lint=fix is only compatible with --mode=fix")
		}

	default:
		return stdfmt.Errorf("unrecognized lint mode %s; valid modes are warn and fix", lint)
	}

	return nil
}

// ValidateWarnings validates the value of the --warnings flag
func ValidateWarnings(warnings string, allWarnings, defaultWarnings *[]string) ([]string, error) {

	// Check lint warnings
	var warningsList []string
	switch warnings {
	case "", "default":
		warningsList = *defaultWarnings
	case "all":
		warningsList = *allWarnings
	default:
		// Either all or no warning categories should start with "+" or "-".
		// If all of them start with "+" or "-", the semantics is
		// "default set of warnings + something - something".
		plus := map[string]bool{}
		minus := map[string]bool{}
		for _, warning := range strings.Split(warnings, ",") {
			if strings.HasPrefix(warning, "+") {
				plus[warning[1:]] = true
			} else if strings.HasPrefix(warning, "-") {
				minus[warning[1:]] = true
			} else {
				warningsList = append(warningsList, warning)
			}
		}
		if len(warningsList) > 0 && (len(plus) > 0 || len(minus) > 0) {
			return []string{}, stdfmt.Errorf("warning categories with modifiers (\"+\" or \"-\") can't me mixed with raw warning categories")
		}
		if len(warningsList) == 0 {
			for _, warning := range *defaultWarnings {
				if !minus[warning] {
					warningsList = append(warningsList, warning)
				}
			}
			for warning := range plus {
				warningsList = append(warningsList, warning)
			}
		}
	}
	return warningsList, nil
}
