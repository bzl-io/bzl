package bzl

import (
	"github.com/urfave/cli"
	"github.com/bzl-io/bzl/bazel"
	"github.com/bzl-io/bzl/command/install"
	"github.com/bzl-io/bzl/command/release"
	"github.com/bzl-io/bzl/command/targets"
)

// App embeds an urfave/cli.App 
type App struct {
	*cli.App
}

// Create a new Application.
func New() *App {
	
	// Create Cli inner app
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Usage = "A candy wrapper for the Bazel build tool"
	app.Version = "0.1.0"

	// Add commands
	app.Commands = []cli.Command{
		*install.Command,
		*release.Command,
		*targets.Command,
	}

	// Any command not found, just run bazel itself
	app.CommandNotFound = func(c *cli.Context, command string) {
		args := append([]string{
			command,
		}, c.Args().Tail()...)
		bazel.New().Invoke(args)
	}

	return &App{
		app,
	}
}


