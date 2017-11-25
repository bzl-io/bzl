package targets

import (
	"fmt"
	"os"
	"io"
	"text/tabwriter"
	"github.com/bzl-io/bzl/bazel"
	"github.com/urfave/cli"
	build "github.com/bzl-io/bzl/proto/build_go"
)

var Command = &cli.Command{
	Name:    "target",
	Aliases: []string{"targets"},
	Usage:   "Display available targets in the workspace",
	Action:  execute,
}

func execute(c *cli.Context) error {

	pattern := c.Args().First()
	if pattern == "" {
		pattern = "//..."
	} 	
	
	query, err := bazel.New().Query(pattern)
	if err != nil {
		return err
	}
	//fmt.Println("Targets:", len(query.GetTarget()))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	for _, target := range query.GetTarget() {
		switch *target.Type {
		case build.Target_SOURCE_FILE:
			printSourceFile(w, target.SourceFile)
		case build.Target_GENERATED_FILE:
			printGeneratedFile(w, target.GeneratedFile)
		case build.Target_RULE:
			printRule(w, target.Rule)
		default:
			fmt.Printf("Skipping %+v\n", target)
		}
	}
	w.Flush()
	return nil
}

func printRule(w io.Writer, rule *build.Rule) {
	//fmt.Fprintln(w, "rule\t", *rule.Name, "\t", *rule.RuleClass)
	fmt.Fprintln(w, *rule.RuleClass, "\trule\t", *rule.Name)
}

func printSourceFile(w io.Writer, file *build.SourceFile) {
	fmt.Fprintln(w, "source\tfile\t", *file.Name)
}

func printGeneratedFile(w io.Writer, file *build.GeneratedFile) {
	fmt.Fprintln(w, "generated\tfile\t", *file.Name)
}

