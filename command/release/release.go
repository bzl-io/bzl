package release

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bzl-io/bzl/bazelutil"
	"github.com/bzl-io/bzl/gh"
	"github.com/davecgh/go-spew/spew"
	"github.com/golang/sync/errgroup"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"golang.org/x/net/context"

	bes "github.com/bzl-io/bzl/proto/bes"
)

var Command = &cli.Command{
	Name: "release",
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "platform",
			Usage: "Name of the @io_bazel_rules_go//go/toolchain:PLATFORM to cross-compile to",
			Value: &cli.StringSlice{
				"linux_amd64",
				"windows_amd64",
				"darwin_amd64",
			},
		},
		cli.StringSliceFlag{
			Name:  "platform_name",
			Usage: "A string mapping of the form PLATFORM=NAME, such as 'windows_amd64=windows-x64_64'",
			Value: &cli.StringSlice{
				"windows_amd64=windows-x64_64",
				"darwin_amd64=darwin-x64_64",
				"linux_amd64=linux-x64_64",
			},
		},
		cli.StringFlag{
			Name:  "asset_dir",
			Usage: "Name of directory where built platform-specific assets should be assembled",
			Value: ".assets",
		},
		cli.StringFlag{
			Name:  "owner",
			Usage: "Name of github owner to publish release",
			Value: "bzl-io",
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "Name of github repo to publish release",
			Value: "bzl",
		},
		cli.StringFlag{
			Name:  "tag",
			Usage: "Tag name for the release such as '1.0.0' (required to publish)",
		},
		cli.StringFlag{
			Name:  "notes",
			Usage: "Release notes filename (a path to markdown file)",
			Value: "RELEASE.md",
		},
		cli.StringFlag{
			Name:  "commit",
			Usage: "Commit ID for the release (required when publishing a release).",
		},
		cli.BoolFlag{
			Name:  "dry_run",
			Usage: "Build assets, but don't actually create a release",
		},
	},
	Usage: "Build target golang binaries for (multiple) platform(s) and publish a release to GitHub",
	Action: func(c *cli.Context) error {
		err := execute(c)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("release failed: %v", err), 1)
		}
		return nil
	},
}

func execute(c *cli.Context) error {
	platforms := c.StringSlice("platform")
	target := ""

	if len(platforms) > 0 {
		target = c.Args().First()
		if target == "" {
			return fmt.Errorf("The 'release' command requires a build target in combination with platforms")
		}
	}

	username := os.Getenv("BZL_GH_USERNAME")
	password := os.Getenv("BZL_GH_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("need environment variables BZL_GH_USERNAME={your-username}, BZL_GH_PASSWORD={personal-access-token}")
	}

	allFiles := make([]string, 0)

	for _, platform := range platforms {
		args := []string{
			"build",
			"--platforms", fmt.Sprintf("@io_bazel_rules_go//go/toolchain:%s", platform),
			target,
		}

		events, err := bazelutil.New().InvokeWithEvents(args)
		if err != nil {
			return err
		}

		completed := bazelutil.FirstTargetComplete(events)
		if completed == nil || !completed.Success {
			return fmt.Errorf("The invocation failed to complete: %s", args)
		}

		files, err := handleTargetCompleted(c, platform, completed)
		if err != nil {
			return err
		}

		allFiles = append(allFiles, files...)
	}

	if c.String("tag") != "" && len(allFiles) > 0 {
		release, err := uploadRelease(c, allFiles)
		if err != nil {
			return err
		}
		fmt.Printf("Release successful: %s\n", release.GetTagName())
	}

	return nil
}

func handleTargetCompleted(c *cli.Context, platform string, completed *bes.TargetComplete) ([]string, error) {
	importantFiles := completed.ImportantOutput
	copies := make([]string, 0)
	for _, file := range importantFiles {
		copy, err := copyFileToPlatformDir(c, c.String("asset_dir"), platform, file)
		if err != nil {
			return nil, fmt.Errorf("Error while relocating output file %s: %v", file.GetUri(), err)
		}
		copies = append(copies, copy)
	}
	return copies, nil
}

func copyFileToPlatformDir(c *cli.Context, assetDir string, platform string, file *bes.File) (string, error) {
	uri := file.GetUri()

	if !strings.HasPrefix(uri, "file://") {
		return "", errors.New("Copying non-file URIs is not implemented: " + uri)
	}

	filename := strings.TrimPrefix(uri, "file://")
	basename := path.Base(filename)

	if assetDir == "" {
		assetDir = "dist"
	}

	platformName := getPreferredPlatformName(c, platform)
	platformDir := assetDir
	err := os.MkdirAll(platformDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	name := basename
	if c.String("tag") != "" {
		name += "-" + c.String("tag")
	}
	name += "-" + platformName
	platformFile := path.Join(platformDir, name)
	err = bazelutil.CopyFile(filename, platformFile)
	if err != nil {
		return "", err
	}

	fmt.Printf("Staged %s for '%s' to %s\n", file.GetName(), platform, platformFile)
	return platformFile, nil
}

func processBuildEvent(event *bes.BuildEvent) error {
	switch x := event.Payload.(type) {
	case *bes.BuildEvent_Progress:
		//fmt.Printf("Progress: %+v\n\n", event)
	case *bes.BuildEvent_Aborted:
		fmt.Printf("Aborted: %+v\n\n", event)
	case *bes.BuildEvent_LoadingFailed:
		fmt.Printf("LoadingFailed: %+v\n\n", event)
	case *bes.BuildEvent_AnalysisFailed:
		fmt.Printf("AnalysisFailed: %+v\n\n", event)
	case *bes.BuildEvent_Started:
		fmt.Printf("Started: %+v\n\n", event)
	case *bes.BuildEvent_CommandLine:
		fmt.Printf("CommandLine: %+v\n\n", event)
	case *bes.BuildEvent_OptionsParsed:
		fmt.Printf("OptionsParsed: %+v\n\n", event)
	case *bes.BuildEvent_WorkspaceStatus:
		fmt.Printf("WorkspaceStatus: %+v\n\n", event)
	case *bes.BuildEvent_Configuration:
		fmt.Printf("Configuration: %+v\n\n", event)
	case *bes.BuildEvent_Expanded:
		fmt.Printf("Expanded: %+v\n\n", event)
	case *bes.BuildEvent_Configured:
		fmt.Printf("Configured: %+v\n\n", event)
	case *bes.BuildEvent_Action:
		fmt.Printf("Action: %+v\n\n", event)
	case *bes.BuildEvent_NamedSetOfFiles:
		fmt.Printf("NamedSetOfFiles: %+v\n\n", event)
	case *bes.BuildEvent_Completed:
		fmt.Printf("Completed: %+v\n\n", event)
	case *bes.BuildEvent_TestResult:
		fmt.Printf("TestResult: %+v\n\n", event)
	case *bes.BuildEvent_TestSummary:
		fmt.Printf("TestSummary: %+v\n\n", event)
	case *bes.BuildEvent_Finished:
		fmt.Printf("Finished: %+v\n\n", event)
	case nil:
	default:
		return fmt.Errorf("BuildEvent.Payload has unexpected type %T", x)
	}
	return nil
}

func uploadRelease(c *cli.Context, files []string) (*github.RepositoryRelease, error) {
	owner := c.String("owner")
	if owner == "" {
		return nil, errors.New("--owner is required when publishing a release")
	}
	repo := c.String("repo")
	if repo == "" {
		return nil, errors.New("--repo is required when publishing a release")
	}
	tag := c.String("tag")
	if tag == "" {
		return nil, errors.New("--tag is required when publishing a release")
	}
	commit := c.String("commit")
	if commit == "" {
		return nil, errors.New("--commit is required when publishing a release")
	}
	notes, err := getReleaseNotes(c.String("notes"))
	if err != nil {
		return nil, err
	}
	fmt.Println("Uploading assets for release", tag, "...")

	client := gh.Client()

	req := &github.RepositoryRelease{
		TagName:         &tag,
		TargetCommitish: &commit,
		Body:            &notes,
	}

	release, err := createRelease(c, client, req, files)
	return release, err
}

// Read the given filename into a string.  Return err if any io error
// occured.
func getReleaseNotes(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func createRelease(c *cli.Context, client *github.Client, req *github.RepositoryRelease, files []string) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	if c.Bool("dry_run") {
		return nil, fmt.Errorf("Create release stopped early (dry_run is ON)")
	}
	owner := c.String("owner")
	if owner == "" {
		return nil, fmt.Errorf("Create release: 'owner' is required")
	}
	repo := c.String("repo")
	if repo == "" {
		return nil, fmt.Errorf("Create release: 'repo' is required")
	}
	log.Printf("Creating release for %s/%s...: %+v", owner, repo, req)

	release, res, err := client.Repositories.CreateRelease(ctx, owner, repo, req)
	if err != nil {
		if release != nil {
			spew.Dump(res, err)
			if res.StatusCode == 404 {
				fmt.Fprintf(os.Stderr, "Github responded with 404 for %s/%s; this may represent an authentication error.  Confirm that the env vars BZL_GH_USERNAME and BZL_GH_PASSWORD are set with PUSH access to this repository (https://developer.github.com/v3/troubleshooting/).\n", c.String("owner"), c.String("repo"))
			}
			if res.StatusCode == 422 {
				fmt.Fprintf(os.Stderr, "Github responded with 422 (validation Failed).  This can occur for multiple reasons (The release already exists? The --commit actually exists at the remote repository?)\n")
			}
		}
		return nil, fmt.Errorf("Create release failed: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Create release failed (invalid status %s): %v", res.Status, err)
	}

	if *release.ID <= 0 {
		return nil, fmt.Errorf("Create release failed to assign a valid ID: %v", release)
	}

	err = uploadAssets(c, client, ctx, *release.ID, files, 5)
	if err != nil {
		return nil, fmt.Errorf("Upload release assets failed: %v", err)
	}
	return release, nil
}

func uploadAssets(c *cli.Context, client *github.Client, ctx context.Context, releaseID int, localAssets []string, parallel int) error {
	start := time.Now()

	defer func() {
		fmt.Printf("UploadAssets: time: %d ms\n", int(time.Since(start).Seconds()*1000))
	}()

	eg, ctx := errgroup.WithContext(ctx)

	semaphore := make(chan struct{}, parallel)

	for _, localAsset := range localAssets {
		localAsset := localAsset
		eg.Go(func() error {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()

			fmt.Fprintf(os.Stdout, "--> Uploading: %15s\n", filepath.Base(localAsset))
			_, err := uploadAsset(c, client, ctx, releaseID, localAsset)
			if err != nil {
				return errors.Wrapf(err,
					"failed to upload asset: %s", localAsset)
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "one of the goroutines failed")
	}

	return nil
}

func uploadAsset(c *cli.Context, client *github.Client, ctx context.Context, releaseID int, filename string) (*github.ReleaseAsset, error) {

	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get abs path")
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}

	opts := &github.UploadOptions{
		// Use base name by default
		Name: filepath.Base(filename),
	}

	asset, res, err := client.Repositories.UploadReleaseAsset(ctx, c.String("owner"), c.String("repo"), releaseID, opts, f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to upload release asset: %s", filename)
	}

	switch res.StatusCode {
	case http.StatusCreated:
		return asset, nil
	case 422:
		return nil, errors.Errorf(
			"upload release asset: invalid status code: %s",
			"422 (this is probably because the asset already uploaded)")
	default:
		return nil, errors.Errorf(
			"upload release asset: invalid status code: %s", res.Status)
	}
}

func getPreferredPlatformName(c *cli.Context, platform string) string {
	// Get all entryies like 'windows_amd64=windows-x86_64'
	entries := c.StringSlice("platform_name")
	// Create a mapping from platform to preferred name
	names := make(map[string]string)
	// Fill it with entries like 'windows_amd4' -> 'windows-x86_64'
	for _, entry := range entries {
		parts := strings.Split(entry, "=")
		if len(parts) == 2 {
			names[parts[0]] = parts[1]
		} else {
			fmt.Println("Malformed platform name mapping:", entry)
		}
	}
	// Look it up
	name := names[platform]
	if name != "" {
		return name
	} else {
		return platform
	}
}
