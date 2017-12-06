package bzl

import (
	"github.com/urfave/cli"
	"github.com/bzl-io/bzl/bazel"
	"github.com/bzl-io/bzl/plugin"
	"github.com/bzl-io/bzl/command/install"
	"github.com/bzl-io/bzl/command/release"
	"github.com/bzl-io/bzl/command/targets"
	"log"
)

// Will be replaced at link time to `git rev-parse HEAD`
var BUILD_SCM_REVISION = "0000000000000000000000000000000000000000"

// App embeds an urfave/cli.App 
type App struct {
	*cli.App
	PluginManager *plugin.Manager
}

// Create a new Application.
func New() *App {

	log.SetPrefix("((bzl)) ")
	
	// Create Cli inner app
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Usage = "A candy wrapper for the Bazel build tool"
	app.Version = BUILD_SCM_REVISION

	// Global flags for bzl app
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name: "bazel",
			Usage: "Use this version(s) of bazel when running subcommand",
		},
		cli.StringFlag{
			Name: "plugin_server_address",
			Usage: "Use this address to get plugins",
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
		PluginManager: nil,
	}

	// Any command not found, just run bazel itself
	app.CommandNotFound = func(c *cli.Context, commandName string) {
		if tryPlugin(instance, c, commandName) {
			return
		}
		
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

func tryPlugin(app *App, c *cli.Context, commandName string) bool {
	log.Println("Looking for plugin: " + commandName)
	
	if app.PluginManager == nil {
		manager, err := plugin.NewManager(c.GlobalString("plugin_server_address"))
		if err == nil {
			log.Println("Created new plugin manager: " + commandName)
			app.PluginManager = manager
		} else {
			log.Fatalf("Unable to create plugin manager: %v\n", err)
		}
	}

	
	if app.PluginManager == nil {
		return false
	}


	if !app.PluginManager.HasPlugin(commandName) {
		log.Println("No plugin named", commandName)
		return false
	}

	cmd, err := app.PluginManager.GetCommandPlugin(commandName)
	if err != nil {
		return false
	}

	err = cmd.Execute(c)
	return true
}
