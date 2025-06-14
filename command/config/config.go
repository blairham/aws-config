package config

import (
	"fmt"

	"github.com/mitchellh/cli"

	"github.com/blairham/aws-sso-config/command/config/edit"
	"github.com/blairham/aws-sso-config/command/config/get"
	"github.com/blairham/aws-sso-config/command/config/list"
	"github.com/blairham/aws-sso-config/command/config/set"
	"github.com/blairham/aws-sso-config/command/config/unset"
)

type cmd struct {
	UI cli.Ui
}

func New(ui cli.Ui) *cmd {
	return &cmd{UI: ui}
}

func (c *cmd) Run(args []string) int {
	if len(args) == 0 {
		c.UI.Error("Usage: aws-sso-config config <subcommand>")
		c.UI.Error("")
		c.UI.Error("Available subcommands:")
		c.UI.Error("  get <key>             Get a configuration value")
		c.UI.Error("  set <key> <value>     Set a configuration value")
		c.UI.Error("  unset <key>           Reset a configuration value to its default")
		c.UI.Error("  list                  List all available configuration keys")
		c.UI.Error("  edit [config-file]    Open configuration file in an editor")
		return 1
	}

	subcommand := args[0]
	subArgs := args[1:]

	switch subcommand {
	case "get":
		getCmd := get.New(c.UI)
		return getCmd.Run(subArgs)
	case "set":
		setCmd := set.New(c.UI)
		return setCmd.Run(subArgs)
	case "unset":
		unsetCmd := unset.New(c.UI)
		return unsetCmd.Run(subArgs)
	case "list":
		listCmd := list.New(c.UI)
		return listCmd.Run(subArgs)
	case "edit":
		editCmd := edit.New(c.UI)
		return editCmd.Run(subArgs)
	default:
		c.UI.Error(fmt.Sprintf("Unknown subcommand: %s", subcommand))
		c.UI.Error("")
		c.UI.Error("Available subcommands:")
		c.UI.Error("  get <key>             Get a configuration value")
		c.UI.Error("  set <key> <value>     Set a configuration value")
		c.UI.Error("  unset <key>           Reset a configuration value to its default")
		c.UI.Error("  list                  List all available configuration keys")
		c.UI.Error("  edit [config-file]    Open configuration file in an editor")
		return 1
	}
}

func (c *cmd) Help() string {
	return `Usage: aws-sso-config config <subcommand>

  Manage configuration settings for aws-sso-config.

Subcommands:
  get <key>            Get a configuration value
  set <key> <value>    Set a configuration value
  unset <key>          Reset a configuration value to its default
  list                 List all available configuration keys
  edit [config-file]   Open configuration file in an editor

Available configuration keys:
  sso.start_url        Your AWS SSO start URL
  sso.region          AWS region for SSO (e.g., us-east-1)
  sso.role            SSO role name (e.g., AdministratorAccess)
  aws.default_region  Default AWS region for profiles
  aws.config_file     Path to AWS config file

Examples:
  # Get the SSO start URL
  aws-sso-config config get sso.start_url

  # Set the SSO start URL (no quotes needed)
  aws-sso-config config set sso.start_url https://mycompany.awsapps.com/start

  # Reset SSO start URL to default
  aws-sso-config config unset sso.start_url

  # Set the default region
  aws-sso-config config set aws.default_region us-west-2

  # Reset default region to default
  aws-sso-config config unset aws.default_region

  # Set the AWS config file path
  aws-sso-config config set aws.config_file ~/.aws/config

  # List all available configuration keys
  aws-sso-config config list

  # Edit the configuration file
  aws-sso-config config edit

  # Edit a specific configuration file
  aws-sso-config config edit /path/to/config
`
}

func (c *cmd) Synopsis() string {
	return "Read and write configuration values"
}
