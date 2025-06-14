package config

import (
	"github.com/mitchellh/cli"

	"github.com/blairham/aws-config/command/flags"
)

type cmd struct {
	UI cli.Ui
}

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	return c
}

func (c *cmd) Run(args []string) int {
	return cli.RunResultHelp
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return flags.Usage(help, nil)
}
