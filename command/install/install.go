package install

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/bzl-io/bzl/bazelutil"
	"github.com/bzl-io/bzl/gh"
	"github.com/google/go-github/github"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/ioprogress"
	"github.com/urfave/cli"
	"golang.org/x/net/context"

	humanize "github.com/dustin/go-humanize"
	download "github.com/joeybloggs/go-download"
)

var Command = &cli.Command{
	Name:  "install",
	Usage: "Install a bazel release (or list release assets)",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "assets",
			Usage: "List assets for a particular release (example: bazel install 0.24.1 --assets)",
		},
		cli.BoolFlag{
			Name:  "force",
			Usage: "Force install",
		},
		cli.BoolFlag{Name: "without_jdk"},
		cli.StringFlag{
			Name:  "bazel_source_dir",
			Usage: "Location of the bazelbuild/bazel git repository (if building from source)",
			Value: "$HOME/.cache/bzl/github.com/bazelbuild/bazel",
		},
	},
	Action: func(c *cli.Context) error {
		if err := execute(c); err != nil {
			return cli.NewExitError(fmt.Sprintf("Install aborted: %v", err), 1)
		}
		return nil
	},
}

func execute(c *cli.Context) error {
	version := c.Args().First()

	//
	// Check if already installed
	//
	if version != "" && !c.Bool("assets") {
		homeDir, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf("Failed to get home directory: %v", err)
		}
		releaseDir := path.Join(homeDir, ".cache", "bzl", "release", version)
		if FileExists(releaseDir) && !c.Bool("force") && version != "snapshot" {
			return fmt.Errorf("%v is already installed (use --force to re-install it)", version)
		}
	}

	//
	// if the version looks like a sha1 value, build from source instead.
	//
	if len(version) == 40 || version == "snapshot" {
		return installFromSource(c, version)
	}

	//
	// List available releases
	//
	listOptions := &github.ListOptions{}
	client := gh.Client()
	releases, _, err := client.Repositories.ListReleases(context.Background(), "bazelbuild", "bazel", listOptions)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to get release list: %v", err), 1)
	}

	if version != "" {
		var match *github.RepositoryRelease
		for _, release := range releases {
			if version == *release.TagName {
				match = release
				break
			}
		}
		if match != nil {
			return processRelease(c, match)
		} else {
			return cli.NewExitError(fmt.Sprintf("Release not found: '%s', should be one of %v\n", version, strings.Join(getReleaseTagNames(releases), ", ")), 1)
		}

	} else {
		return listReleases(c, releases)
	}
}

func getReleaseTagNames(releases []*github.RepositoryRelease) []string {
	tags := make([]string, len(releases))
	for i, release := range releases {
		tags[i] = *release.TagName
	}
	return tags
}

func listReleases(c *cli.Context, releases []*github.RepositoryRelease) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, release := range releases {
		fmt.Fprintln(w, *release.TagName, "\t", (*release.PublishedAt).Format("Mon Jan 02 2006"))
	}
	w.Flush()
	return nil
}

func processRelease(c *cli.Context, release *github.RepositoryRelease) error {
	command := "install"
	if c.Bool("assets") {
		command = "assets"
	}

	switch command {
	case "install":
		return installRelease(c, release)
	case "assets":
		return listReleaseAssets(c, release)
	default:
		return cli.NewExitError(fmt.Sprintf("Unknown subcommand '%s' for release '%s')\n", command, *release.TagName), 1)
	}
	return nil
}

func listReleaseAssets(c *cli.Context, release *github.RepositoryRelease) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, asset := range release.Assets {
		if filepath.Ext(*asset.Name) != ".sha256" && filepath.Ext(*asset.Name) != ".sig" {
			fmt.Fprintln(w, *asset.Name, "\t", humanize.Bytes(uint64(*asset.Size)))
		}
	}
	w.Flush()
	return nil
}

func installRelease(c *cli.Context, release *github.RepositoryRelease) error {
	version := *release.TagName
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	if goarch == "amd64" {
		goarch = "x86_64" // is this right?
	}

	baseDir := "/tmp"
	homeDir, err := homedir.Dir()
	if err == nil {
		baseDir = path.Join(homeDir, ".cache", "bzl", "install")
		os.MkdirAll(baseDir, os.ModePerm)
	}

	//fmt.Printf("INSTALL %s %s %s: %s\n", goos, goarch, version, *release.Assets[0].BrowserDownloadURL)

	installer := fmt.Sprintf("bazel-%s-installer-%s-%s.sh", version, goos, goarch)
	if c.Bool("without_jdk") {
		installer = fmt.Sprintf("bazel-%s-without-jdk-installer-%s-%s.sh", version, goos, goarch)
	}

	assets := makeAssetMap(release)

	exeAsset := assets[installer]
	if *exeAsset.Name == "" {
		return cli.NewExitError(fmt.Sprintf("Can't install from release %s (no installer found)", version), 1)
	}
	shaAsset := assets[installer+".sha256"]
	if *shaAsset.Name == "" {
		return cli.NewExitError(fmt.Sprintf("Can't install from release %s (no installer sha256 file)", version), 1)
	}
	sigAsset := assets[installer+".sig"]
	if *sigAsset.Name == "" {
		return cli.NewExitError(fmt.Sprintf("Can't install from release %s (no installer gpg sig)", version), 1)
	}

	sha, err := getOrDownloadAsset(baseDir, &shaAsset)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Can't install from release: %s", err), 1)
	}
	// sig, err := getOrDownloadAsset(baseDir, sigAsset)
	// if err != nil {
	// 	return err
	// }

	content, err := ioutil.ReadFile(sha)
	if err != nil {
		return err
	}

	expected := strings.Fields(string(content))

	// gpg, err := ioutil.ReadFile(sig)
	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("gpg: %s\n", gpg)

	exe, err := getOrDownloadAsset(baseDir, &exeAsset)
	if err != nil {
		return err
	}

	actual, err := GetFileSha256(exe)
	if err != nil {
		return err
	}

	if actual != expected[0] {
		return fmt.Errorf(
			"%s: got sha256 '%s' but expected '%s' (from %s)\n",
			exe,
			actual,
			expected,
			sha,
		)
	}

	log.Printf("Sha256 match %s: %s, proceeding with install...\n", exe, sha)

	if err := install(c, version, exe); err != nil {
		return err
	}

	log.Printf("Install was successful ('export BAZEL_VERSION=%s' to use it)", version)

	return nil
}

func makeAssetMap(release *github.RepositoryRelease) map[string]github.ReleaseAsset {
	m := make(map[string]github.ReleaseAsset)
	for _, a := range release.Assets {
		m[*a.Name] = a
	}
	return m
}

func GetFileSha256(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func getOrDownloadAsset(baseDir string, asset *github.ReleaseAsset) (string, error) {
	filename := path.Join(baseDir, *asset.Name)
	if FileExists(filename) {
		return filename, nil
	}
	return downloadAsset(baseDir, asset)
}

func FileExists(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()
	return true
}

func downloadAsset(baseDir string, asset *github.ReleaseAsset) (string, error) {
	filename := path.Join(baseDir, *asset.Name)
	rc, redirect, err := gh.Client().Repositories.DownloadReleaseAsset(context.Background(), "bazelbuild", "bazel", *asset.ID)
	if err != nil {
		return filename, err
	}

	out, err := os.Create(filename)
	if err != nil {
		return filename, err
	}
	defer out.Close()

	if redirect != "" {
		err = downloadUrl(redirect, out, int64(*asset.Size), *asset.Name)
	} else {
		defer rc.Close()
		err = Download(rc, out, int64(*asset.Size), *asset.Name)
	}

	return filename, nil
}

func downloadUrl(url string, out io.Writer, size int64, title string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	return Download(resp.Body, out, size, title)
}

func downloadUrl2(url string, out io.Writer, size int64, title string) error {
	f, err := download.Open(url, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	return Download(f, out, size, title)
}

// download accepts a reader, writer, and expected asset size, copies
// the bytes through and displays progress to the console.
func Download(reader io.Reader, writer io.Writer, size int64, title string) error {
	// Attempt 1:
	// bar := pb.New(size).SetUnits(pb.U_BYTES)
	// bar.Prefix(title)
	// bar.Start()
	// reader = bar.NewProxyReader(reader)

	// Attempt 2: mpb
	//progress := mpb.New()
	//bar := progress.AddBar(size, mpb.BarTrim())
	//reader = bar.ProxyReader(reader)

	// Attempt 3: ioprogress
	reader = &ioprogress.Reader{
		Reader: reader,
		Size:   size,
	}

	_, err := io.Copy(writer, reader)
	if err != nil {
		return err
	}

	//bar.Finish()
	//progress.Stop()

	return nil
}

func install(c *cli.Context, version, filename string) error {
	cmdName := "bash"

	prefix := ""
	homeDir, err := homedir.Dir()
	if err == nil {
		prefix = path.Join(homeDir, ".cache", "bzl", "release", version)
		os.MkdirAll(prefix, os.ModePerm)
	}

	cmdArgs := []string{
		filename,
	}

	if prefix == "" {
		cmdArgs = append(cmdArgs, "--user")
	} else {
		cmdArgs = append(cmdArgs, "--prefix="+prefix)
	}

	log.Printf("Installing %v", cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installFromSource(c *cli.Context, version string) error {

	// make sure bazel is checked out at the approprite location
	bazelDir := os.ExpandEnv(c.String("bazel_source_dir"))

	// If the directory does not exist, create it and clone bazel.
	if !FileExists(bazelDir) {
		parentDir := filepath.Dir(bazelDir)
		if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
			return fmt.Errorf("Failed to prepare bazel source directory: %v", err)
		}

		// Execute a git clone
		cloneArgs := []string{
			"clone",
			"https://github.com/bazelbuild/bazel.git",
		}
		cmd := exec.Command("git", cloneArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = parentDir
		cmd.Env = os.Environ()

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Failed to clone bazel: %v", err)
		}

	}

	//
	// Make sure we have the correct version checked out.
	//
	if version != "snapshot" {
		// Pull most recent stuff
		cmd := exec.Command("git", []string{
			"fetch",
			// version,
		}...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = bazelDir
		cmd.Env = os.Environ()

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Failed to fetch %q: %v", version, err)
		}

		// Checkout to the requested commit
		cmd = exec.Command("git", []string{
			"checkout",
			version,
		}...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = bazelDir
		cmd.Env = os.Environ()

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("Failed to checkout %q: %v", version, err)
		}
	}

	//
	// And run the build...
	//
	if err, exitCode := bazelutil.New().Invoke([]string{
		"build",
		"//src:bazel",
	}, bazelDir); err != nil {
		return fmt.Errorf("Failed to build bazel: %v (exit %d)", err, exitCode)
	}

	//
	// Make the release dir
	//
	var releaseDir string
	homeDir, err := homedir.Dir()
	if err != nil {
		return fmt.Errorf("Failed to read home dir: %v", err)
	}
	releaseDir = filepath.Join(homeDir, ".cache", "bzl", "release", version, "bin")
	if err := os.MkdirAll(releaseDir, os.ModePerm); err != nil {
		return fmt.Errorf("Failed to prepare bazel release directory: %v", err)
	}

	//
	// Copy the binary
	//
	srcFile := filepath.Join(bazelDir, "bazel-bin", "src", "bazel")
	dstFile := filepath.Join(releaseDir, "bazel")
	if err := bazelutil.CopyFile(srcFile, dstFile); err != nil {
		return fmt.Errorf("Failed to copy bazel binary to release directory: %v", err)
	}

	//
	// Set executable permisssions
	//
	os.Chmod(dstFile, 0755)

	log.Printf("Build+Install was successful ('export BAZEL_VERSION=%s' to use it)", version)

	return nil
}
