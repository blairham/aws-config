package list

import (
	"github.com/mitchellh/cli"

	"github.com/blairham/aws-sso-config/command/config/shared"
)

type cmd struct {
	UI cli.Ui
}

func New(ui cli.Ui) *cmd {
	return &cmd{UI: ui}
}

func (c *cmd) Run(args []string) int {
	if len(args) != 0 {
		c.UI.Error("Usage: aws-sso-config config list")
		c.UI.Error("")
		c.UI.Error("This command takes no arguments.")
		return 1
	}

	c.UI.Output("Available configuration keys:")
	shared.OutputAvailableKeys(c.UI)
	return 0
}

func (c *cmd) Help() string {
	return `Usage: aws-sso-config config list

  List all available configuration keys and their descriptions.

Examples:
  # List all configuration keys
  aws-sso-config config list
`
}

func (c *cmd) Synopsis() string {
	return "List all available configuration keys"
}
