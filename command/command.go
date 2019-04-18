package command

import (
	"github.com/urfave/cli"
)

type Command interface {
	Execute(cli *cli.Context) error
}
