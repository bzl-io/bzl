package bzl

import (
	"os"
	"os/exec"
	"github.com/urfave/cli"
	"github.com/bzl-io/bzl/command/install"
	"github.com/bzl-io/bzl/command/targets"
)

type App struct {
	app *cli.App
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
		*targets.Command,
	}

	// Any command not found, just run bazel itself
	app.CommandNotFound = func(c *cli.Context, command string) {
		cmdName := "bazel"
		cmdArgs := append([]string{
			command,
		}, c.Args().Tail()...)
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = ""
		cmd.Run() 
	}

	return &App{
		app: app,
	}
}

// Run the application with the given arguments.
func (bzl *App) Run(args []string) error {
	return bzl.app.Run(args)
}


