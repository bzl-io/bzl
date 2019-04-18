package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/bzl-io/bzl/bazelutil"
	"github.com/bzl-io/bzl/command/install"
	"github.com/bzl-io/bzl/command/targets"
)

// Will be replaced at link time to `git rev-parse HEAD`
var BUILD_SCM_REVISION = "0000000000000000000000000000000000000000"
var BUILD_SCM_DATE = "0000-00-00"

// App embeds an urfave/cli.App
type App struct {
	*cli.App
}

// Create a new Application.
func NewApp() *App {

	log.SetPrefix("((bzl)) ")

	// Create Cli inner app
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Usage = "A wrapper for the Bazel build tool"
	app.Version = fmt.Sprintf("%s (%s)", BUILD_SCM_REVISION, BUILD_SCM_DATE)

	// Global flags for bzl app
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "bazel",
			Usage:  "Use this version(s) of bazel when running subcommand",
			EnvVar: "BAZEL_VERSION",
		},
	}

	// Add commands
	app.Commands = []cli.Command{
		*install.Command,
		*targets.Command,
	}

	instance := &App{
		App: app,
	}

	// Any command not found, just run bazel itself
	app.CommandNotFound = func(c *cli.Context, commandName string) {
		args := []string{commandName}
		if len(c.GlobalStringSlice("bazel")) > 0 {
			args = append(args, c.Args().Tail()...)
			for _, version := range c.GlobalStringSlice("bazel") {
				err := bazelutil.SetVersion(version)
				if err != nil {
					log.Fatalf("Invalid bazel version %s, aborting: %v", version, err)
				}
				err, exitCode := bazelutil.New().Invoke(args)
				if exitCode != 0 {
					log.Printf("bazel exited with exitCode %d: %v", exitCode, err)
					os.Exit(exitCode)
				}
			}
		} else {
			log.Println("BAZEL_VERSION not set, falling back to bazel on your PATH")
			args = append(args, c.Args().Tail()...)
			err, exitCode := bazelutil.New().Invoke(args)
			if exitCode != 0 {
				log.Printf("bazel exited with exitCode %d: %v", exitCode, err)
				os.Exit(exitCode)
			}
		}
	}

	return instance
}
