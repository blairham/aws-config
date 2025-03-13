package write

import (
	"github.com/blairham/aws-config/command/cli"
	"github.com/blairham/aws-config/command/flags"
)

type cmd struct {
	UI   cli.UI
	help string
}

const (
	synopsis = "Create or update a centralized config entry"
	help     = `Usage: aws-config config write [options] <configuration>

  Request a config entry to be created or updated. The configuration
  argument is either a file path or '-' to indicate that the config
  should be read from stdin. The data should be either in HCL or
  JSON form.

  Example (from file):

    $ aws-condif config write web.service.hcl

  Example (from stdin):

    $ aws-config config write -
`
)

func New(ui cli.UI) *cmd {
	c := &cmd{UI: ui}
	c.Init()
	return c
}

func (c *cmd) Init() {
	c.help = help
}

func (c *cmd) Run(args []string) int {
	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return flags.Usage(c.help, nil)
}
