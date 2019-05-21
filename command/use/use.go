package use

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	// "github.com/bazelbuild/buildtools/build"
	"github.com/bzl-io/bzl/command/install"
	"github.com/bzl-io/bzl/gh"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

var Command = &cli.Command{
	Name:  "use",
	Usage: "Output a workspace rule for a given github repository",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "rule",
			Usage:  "Name of the rule to output",
			Value:  "http_archive",
			EnvVar: "BZL_USE_RULE_OUTPUT",
		},
		cli.StringFlag{
			Name:   "type",
			Usage:  "Type of asset to download (zip, tar)",
			Value:  "tar",
			EnvVar: "BZL_USE_ASSET_TYPE",
		},
		cli.StringFlag{
			Name:   "history",
			Usage:  "Query recent commit log for the given branch (example: --history=master)",
			EnvVar: "BZL_USE_COMMIT_HISTORY",
		},
	},
	Action: func(c *cli.Context) error {
		err := execute(c)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("use failed: %v", err), 1)
		}
		return nil
	},
}

func execute(c *cli.Context) error {

	// Yes, it's one function.

	arg1 := c.Args().First()
	if arg1 == "" {
		return fmt.Errorf("usage: bazel use [OWNER/]REPO TAG (example: bazel use rules_go 0.18.3)")
	}

	// Get the first argument, either a single string like 'rules_go' or with an
	// owner qualifier 'bazelbuild/rules_go'.
	parts := strings.SplitN(arg1, "/", 2)

	// 2nd argument should be the tag name desired.  If not present, just list
	// the releases.
	tag := c.Args().Get(1)
	archiveType := c.String("type")

	var owner string
	var repo string

	if len(parts) == 1 {
		// case like 'bazel use rules_go'
		owner = "bazelbuild"
		repo = parts[0]
	} else if len(parts) == 2 && len(parts[1]) == 40 {
		// case like 'bazel use rules_go ef7cca8857f2f3a905b86c264737d0043c6a766e'
		owner = "bazelbuild"
		repo = parts[0]
	} else if len(parts) == 2 {
		// case like 'bazel use bazelbuild/rules_go ef7cca8857f2f3a905b86c264737d0043c6a766e'
		owner = parts[0]
		repo = parts[1]
	} else {
		return fmt.Errorf("want OWNER/REPO, got %q", c.Args().First())
	}

	//
	// List the github releases
	//
	client := gh.Client()

	var commits []*github.RepositoryCommit
	var err error

	if c.String("history") != "" {
		log.Printf("Listing recent commit history for %s/%s...", owner, repo)
		commits, _, err = client.Repositories.ListCommits(context.Background(), owner, repo, &github.CommitsListOptions{
			SHA: c.String("history"),
		})
		if err != nil {
			return fmt.Errorf("Failed to list commits for ref %q: %v", c.String("history"), err)
		}

		if tag == "" {
			listCommits(commits)
			return fmt.Errorf("please specify a commit ID (example: %s)", commits[0].GetSHA())
		}
	} else if len(tag) != 40 && !strings.HasPrefix(tag, "refs/") {

		releases, _, err := client.Repositories.ListReleases(context.Background(), owner, repo, nil)
		if err != nil {
			return fmt.Errorf("Failed to list releases: %v", err)
		}

		if len(releases) == 0 {
			return fmt.Errorf("No releases found for %s/%s", owner, repo)
		}

		var release *github.RepositoryRelease

		if tag == "" {
			listReleases(releases)
			return fmt.Errorf("please specify release tag (example: %s)", releases[0].GetTagName())
		}

		// Try and match desired release
		//
		for _, r := range releases {
			if r.GetTagName() == tag {
				release = r
				break
			}
		}

		if release == nil {
			return fmt.Errorf("release %q not found in %s/%s", tag, owner, repo)
		}

		tag = release.GetTagName()
	}

	if strings.HasPrefix(tag, "refs/heads") {
		parts = strings.Split(tag, "/")
		tag = parts[len(parts)-1]
	}

	//
	// Normalize the archiveType
	//
	switch archiveType {
	case "tar":
		fallthrough
	case "tgz":
		archiveType = "tar.gz"
	default:
		return fmt.Errorf("Unknown --type=%q", archiveType)
	}

	// Prep the archive url
	//
	url := fmt.Sprintf("https://github.com/%s/%s/archive/%s.%s", owner, repo, tag, archiveType)

	// Prep to download to a temp file
	//
	tmpFile, err := ioutil.TempFile("", tag)
	if err != nil {
		return fmt.Errorf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Get the url (following redirects)
	httpClient := &http.Client{}

	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		dump, err := httputil.DumpResponse(resp, true)
		return fmt.Errorf("%s\n\n%s\n%v", url, dump, err)
	}

	err = install.Download(resp.Body, tmpFile, resp.ContentLength, url)
	if err != nil {
		return fmt.Errorf("Failed to download %s: %v", url, err)
	}

	// Calc a sha256 for the downloaded file
	//
	sha256, err := install.GetFileSha256(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("Failed calculate sha256 %s: %v", tmpFile.Name(), err)
	}

	// Prep the workspace name, defaulting to a canonical github form
	wsName := strings.ToLower(fmt.Sprintf("com_github_%s_%s", owner, repo))
	wsName = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(wsName, "_")

	// Fetch the WORKSPACE file...
	wsFile, err := client.Repositories.DownloadContents(context.Background(), owner, repo, "WORKSPACE", &github.RepositoryContentGetOptions{
		Ref: tag,
	})

	// Ignore errors, only parse it on success
	if err == nil {
		defer wsFile.Close()

		r := regexp.MustCompile(`\s*workspace\s*\(\s*name\s*=\s*['"](?P<name>[-_.a-zA-Z0-9]+)['"]\s*\)\s*`)

		// Parse the WORKSPACE file for the workspace name.  Yes, using a regexp
		// is suboptimal.  However, I could not link buildifier due to failure
		// to GoCompile (goyacc errors?)
		scanner := bufio.NewScanner(wsFile)
		for scanner.Scan() {
			line := scanner.Text()
			match := r.FindAllSubmatchIndex([]byte(line), 1)
			// Example: for 'workspace(name = "io_bazel_rules_closure")', expect
			// [[0 42 18 40]], 18 being start index and 40 being end index
			// (exclusive).
			if len(match) == 1 && len(match[0]) == 4 {
				pair := match[0]
				wsName = line[pair[2]:pair[3]]
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("Failed to scan workspace: %v", err)
		}
	}

	//
	// Success, now print the rule.
	//

	rule := c.String("rule")
	switch rule {
	case "http_archive":
		printHttpArchive(wsName, owner, repo, tag, sha256, archiveType)
	case "go_repository":
		printGoRepository(wsName, owner, repo, tag, sha256)
	default:
		return fmt.Errorf("Unknown --rule=%q", rule)
	}

	return nil
}

func listReleases(releases []*github.RepositoryRelease) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, release := range releases {
		fmt.Fprintln(w, *release.TagName, "\t", (*release.PublishedAt).Format("Mon Jan 02 2006"))
	}
	w.Flush()
	return nil
}

func listCommits(commits []*github.RepositoryCommit) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, commit := range commits {
		date := (*commit.Commit.Author.Date).Format("Mon Jan 02 2006")
		msg := scanFirstLine(commit.Commit.GetMessage())
		fmt.Fprintln(w, *commit.SHA, "\t", date, "\t", msg)
	}
	w.Flush()
	return nil
}

func scanFirstLine(in string) string {
	scanner := bufio.NewScanner(strings.NewReader(in))
	for scanner.Scan() {
		return scanner.Text()
	}
	return in
}

func printHttpArchive(wsName, owner, repo, tag, sha256, archiveType string) {
	stripVersion := tag
	if strings.HasPrefix(tag, "v") {
		stripVersion = tag[1:]
	}
	stripPrefix := fmt.Sprintf("%s-%s", repo, stripVersion)
	fmt.Printf(`
http_archive(
    name = %q,
    sha256 = %q,
    strip_prefix = %q,
    urls = ["https://github.com/%s/%s/archive/%s.%s"],
)

`, wsName, sha256, stripPrefix, owner, repo, tag, archiveType)
}

func printGoRepository(wsName, owner, repo, tag, sha256 string) {
	attr := "tag"
	if len(tag) == 40 {
		attr = "commit"
	}
	fmt.Printf(`
go_repository(
    name = %q,
    importpath = "github.com/%s/%s",
    %s = %q,
    sha256 = %q,
)

`, wsName, owner, repo, attr, tag, sha256)
}
