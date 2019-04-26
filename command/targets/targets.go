package targets

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/bzl-io/bzl/bazelutil"
	build "github.com/bzl-io/bzl/proto/build"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

var Command = &cli.Command{
	Name:    "target",
	Aliases: []string{"targets"},
	Usage:   "Pretty print query output",
	Action:  execute,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:   "nocolor",
			Usage:  "Don't olorize output (experimental)",
			EnvVar: "BZL_TARGET_NO_COLOR",
		},
		cli.StringSliceFlag{
			Name:   "sort",
			Usage:  `Sort by field kind|label`,
			EnvVar: "BZL_TARGET_SORT",
			Value: &cli.StringSlice{
				"kind",
			},
		},
		cli.StringSliceFlag{
			Name:   "include",
			Usage:  `String that an entry must have in label|kind`,
			EnvVar: "BZL_TARGET_INCLUDE",
		},
		cli.StringFlag{
			Name:   "align",
			Usage:  `Align output by ws|root|pkg`,
			EnvVar: "BZL_TARGET_ALIGN",
			Value:  "root",
		},
	},
}

var grey = color.New(color.FgWhite).Add(color.Faint).SprintFunc()
var bold = color.New(color.FgWhite).Add(color.Bold).SprintFunc()

type target struct {
	label  string
	kind   string
	target *build.Target
	align  string
}

type ByLabel []*target

func (a ByLabel) Len() int      { return len(a) }
func (a ByLabel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLabel) Less(i, j int) bool {
	return strings.Compare(a[i].label, a[j].label) < 0
}

type ByKind []*target

func (a ByKind) Len() int      { return len(a) }
func (a ByKind) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByKind) Less(i, j int) bool {
	return strings.Compare(a[i].kind, a[j].kind) < 0
}

func execute(c *cli.Context) error {

	//
	// Prepare pattern
	//
	pattern := c.Args().First()
	if pattern == "" {
		pattern = ":*"
	}

	//
	// Perform 'bazel query'
	//
	query, err := bazelutil.New().Query(pattern)
	if err != nil {
		return err
	}

	align := c.String("align")

	color.NoColor = c.Bool("nocolor")

	//
	// Get results and pre-process into a different struct
	//
	queryTargets := query.GetTarget()

	targets := make([]*target, len(queryTargets))

	for i, queryTarget := range queryTargets {
		target := &target{
			target: queryTarget,
			align:  align,
		}
		switch *queryTarget.Type {
		case build.Target_SOURCE_FILE:
			target.kind = "source"
			target.label = queryTarget.SourceFile.GetName()
		case build.Target_GENERATED_FILE:
			target.kind = "generated"
			target.label = queryTarget.SourceFile.GetName()
		case build.Target_RULE:
			target.kind = queryTarget.Rule.GetName()
			target.label = queryTarget.Rule.GetRuleClass()
		case build.Target_PACKAGE_GROUP:
			target.kind = "package-group"
			target.label = queryTarget.PackageGroup.GetName()
		default:
			fmt.Printf("Skipping %+v\n", queryTarget)
		}
		targets[i] = target
	}

	//
	// Perform filtering, if requested
	//
	for _, snippet := range c.StringSlice("include") {
		filtered := make([]*target, 0)
		for _, t := range targets {
			if strings.Contains(t.kind, snippet) || strings.Contains(t.label, snippet) {
				filtered = append(filtered, t)
			}
		}
		targets = filtered
	}

	//
	// Perform sorting, if requested
	//
	for _, field := range c.StringSlice("sort") {
		switch field {
		case "kind":
			sort.Sort(ByKind(targets))
		case "label":
			sort.Sort(ByLabel(targets))
		default:
			return fmt.Errorf("Unknown sort field: %q", field)
		}
	}

	//fmt.Println("Targets:", len(query.GetTarget()))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight)

	for _, q := range targets {
		switch *q.target.Type {
		case build.Target_SOURCE_FILE:
			q.printSourceFile(w, q.target.SourceFile)
		case build.Target_GENERATED_FILE:
			q.printGeneratedFile(w, q.target.GeneratedFile)
		case build.Target_RULE:
			q.printRule(w, q.target.Rule)
		case build.Target_PACKAGE_GROUP:
			q.printPackageGroup(w, q.target.PackageGroup)
		default:
			fmt.Printf("Skipping %+v\n", q)
		}
	}
	w.Flush()
	return nil
}

func (t *target) printRule(w io.Writer, rule *build.Rule) {
	fmt.Fprintln(w, t.colorizeRuleClass(rule.GetRuleClass()), "\t", t.colorizeTarget(rule.GetName()))
}

func (t *target) printSourceFile(w io.Writer, file *build.SourceFile) {
	fmt.Fprintln(w, "source-file", "\t", t.colorizeTarget(*file.Name))
}

func (t *target) printGeneratedFile(w io.Writer, file *build.GeneratedFile) {
	fmt.Fprintln(w, "generated-file", "\t", t.colorizeTarget(*file.Name))
}

func (t *target) printPackageGroup(w io.Writer, group *build.PackageGroup) {
	fmt.Fprintln(w, "package-group", "\t", t.colorizeTarget(*group.Name))
}

func (t *target) colorizeRuleClass(ruleClass string) string {
	// t = strings.TrimSpace(t)
	// return color.GreenString(t)
	return ruleClass
}

func (g *target) colorizeTarget(t string) string {
	parts := strings.SplitN(t, ":", 2)
	selector := parts[0]
	target := parts[1]

	segments := strings.SplitN(selector, "//", 2)
	path := segments[1]
	ws := segments[0]
	if strings.HasPrefix(ws, "@") {
		ws = ws[1:]
	}

	s := color.BlueString(ws)
	if g.align == "root" {
		s += "\t"
	}
	s += grey("//") + bold(path)
	if g.align == "pkg" {
		s += "\t"
	}
	s += grey(":") + color.YellowString(target)

	if ws != "" {
		s = grey("@") + s
	}
	return s
}
