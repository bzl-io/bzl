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
		cli.BoolFlag{Name: "list"},
		cli.BoolFlag{Name: "force"},
		cli.BoolFlag{Name: "without_jdk"},
	},
	Action: execute,
}

func execute(c *cli.Context) error {
	version := c.Args().First()

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
	if c.Bool("list") {
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

	//fmt.Printf("INSTALL %s %s %s: %s\n", goos, goarch, version, *release.Assets[0].BrowserDownloadURL)
	// https://github.com/bazelbuild/bazel/releases/download/0.7.0/bazel-0.7.0-installer-darwin-x86_64.sh.sha256
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
		if c.Bool("force") {
			log.Printf(
				"%s: got sha256 '%s' but expected '%s' (from %s).  Proceeding with install due to --force\n",
				exe,
				actual,
				expected,
				sha)
		} else {
			return cli.NewExitError(
				fmt.Sprintf(
					"%s: got sha256 '%s' but expected '%s' (from %s)\n",
					exe,
					actual,
					expected,
					sha),
				2)
		}
	}

	log.Printf("Sha256 match %s: %s, proceeding with install\n", exe, sha)

	install(c, version, exe)
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
		err = dl(rc, out, int64(*asset.Size), *asset.Name)
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

	return dl(resp.Body, out, size, title)
}

func downloadUrl2(url string, out io.Writer, size int64, title string) error {
	f, err := download.Open(url, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	return dl(f, out, size, title)
}

// download accepts a reader, writer, and expected asset size, copies
// the bytes through and displays progress to the console.
func dl(reader io.Reader, writer io.Writer, size int64, title string) error {
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
