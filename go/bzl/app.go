package bzl

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/bzl-io/bzl/bazel"
	"github.com/bzl-io/bzl/command/install"
	"github.com/bzl-io/bzl/command/release"
	"github.com/bzl-io/bzl/command/targets"
	"log"
)

// Will be replaced at link time to `git rev-parse HEAD`
var BUILD_SCM_REVISION = "0000000000000000000000000000000000000000"
var BUILD_SCM_DATE = "0000-00-00"

// App embeds an urfave/cli.App 
type App struct {
	*cli.App
}

// Create a new Application.
func New() *App {

	log.SetPrefix("((bzl)) ")
	
	// Create Cli inner app
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Usage = "A wrapper for the Bazel build tool"
	app.Version = fmt.Sprintf("%s (%s)", BUILD_SCM_REVISION, BUILD_SCM_DATE)

	// Global flags for bzl app
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name: "bazel",
			Usage: "Use this version(s) of bazel when running subcommand",
		},
	}
	
	// Add commands
	app.Commands = []cli.Command{
		*install.Command,
		*release.Command,
		*targets.Command,
	}

	instance := &App{
		App: app,
	}

	// Any command not found, just run bazel itself
	app.CommandNotFound = func(c *cli.Context, commandName string) {
		args := []string{ commandName }
		if len(c.GlobalStringSlice("bazel")) > 0 {
			args = append(args, c.Args().Tail()...)
			for _, version := range c.GlobalStringSlice("bazel") {
				err := bazel.SetVersion(version)
				if err != nil {
					log.Fatalf("Invalid bazel version %s, aborting.", version)
				}
				bazel.New().Invoke(args)				
			}
		} else {
			args = append(args, c.Args().Tail()...)
			bazel.New().Invoke(args)
		}
	}

	return instance
}
