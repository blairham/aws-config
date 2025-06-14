package initialize

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"

	"github.com/blairham/aws-sso-config/command/flags"
)

const (
	formatYAML = "yaml"
	formatJSON = "json"
	formatTOML = "toml"
)

type cmd struct {
	UI    cli.Ui
	flags *flag.FlagSet
	help  string

	file   string
	format string
}

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.Init()
	return c
}

func (c *cmd) Init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.StringVar(&c.file, "file", "aws-sso-config.yaml", "Path to the configuration file to create")
	c.flags.StringVar(&c.format, "format", formatYAML, "Format for the configuration file (yaml, json, toml)")

	c.help = flags.Usage(help, c.flags)
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	// Validate format
	format := strings.ToLower(c.format)
	if format != formatYAML && format != formatJSON && format != formatTOML {
		c.UI.Error("Error: format must be one of: yaml, json, toml")
		return 1
	}

	// Determine file extension based on format if not specified
	if c.file == "aws-sso-config.yaml" && format != formatYAML {
		c.file = fmt.Sprintf("aws-sso-config.%s", format)
	}

	// Check if file already exists
	if _, err := os.Stat(c.file); err == nil {
		c.UI.Error(fmt.Sprintf("Error: file %s already exists", c.file))
		return 1
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(c.file)
	if dir != "." {
		if err := os.MkdirAll(dir, 0750); err != nil {
			c.UI.Error(fmt.Sprintf("Error creating directory %s: %v", dir, err))
			return 1
		}
	}

	// Generate config content based on format
	var content string
	var err error

	switch format {
	case formatYAML:
		content = generateYAMLExample()
	case formatJSON:
		content = generateJSONExample()
	case formatTOML:
		content = generateTOMLExample()
	default:
		content = generateYAMLExample()
	}

	// Write the file
	if err = os.WriteFile(c.file, []byte(content), 0600); err != nil {
		c.UI.Error(fmt.Sprintf("Error writing config file: %v", err))
		return 1
	}

	c.UI.Output(fmt.Sprintf("Successfully created configuration file: %s", c.file))
	c.UI.Output("Edit the file to customize your settings, then use:")
	c.UI.Output(fmt.Sprintf("  aws-sso-config generate -config=%s", c.file))

	return 0
}

func (c *cmd) Help() string {
	return c.help
}

func (c *cmd) Synopsis() string {
	return "Initialize a configuration file with default values"
}

func generateYAMLExample() string {
	return `# AWS Config Tool Configuration
# This file supports YAML, TOML, and JSON formats

# SSO Configuration
sso_start_url: "https://your-sso-portal.awsapps.com/start"
sso_region: "us-east-1"
sso_role: "AdministratorAccess"

# AWS Configuration
default_region: "us-east-1"
config_file: "~/.aws/config"

# Behavior Settings
backup_configs: true
dry_run: false

# Environment variables can also be used with AWS_CONFIG_ prefix:
# AWS_CONFIG_SSO_START_URL=https://your-sso-portal.awsapps.com/start
# AWS_CONFIG_SSO_REGION=us-east-1
# AWS_CONFIG_SSO_ROLE=AdministratorAccess
# AWS_CONFIG_DEFAULT_REGION=us-east-1
# AWS_CONFIG_CONFIG_FILE=~/.aws/config
# AWS_CONFIG_BACKUP_CONFIGS=true
# AWS_CONFIG_DRY_RUN=false
`
}

func generateJSONExample() string {
	return `{
  "_comment": "AWS Config Tool Configuration - supports YAML, TOML, and JSON formats",
  "sso_start_url": "https://your-sso-portal.awsapps.com/start",
  "sso_region": "us-east-1",
  "sso_role": "AdministratorAccess",
  "default_region": "us-east-1",
  "config_file": "~/.aws/config",
  "backup_configs": true,
  "dry_run": false
}
`
}

func generateTOMLExample() string {
	return `# AWS Config Tool Configuration
# This file supports YAML, TOML, and JSON formats

# SSO Configuration
sso_start_url = "https://your-sso-portal.awsapps.com/start"
sso_region = "us-east-1"
sso_role = "AdministratorAccess"

# AWS Configuration
default_region = "us-east-1"
config_file = "~/.aws/config"

# Behavior Settings
backup_configs = true
dry_run = false

# Environment variables can also be used with AWS_CONFIG_ prefix:
# AWS_CONFIG_SSO_START_URL=https://your-sso-portal.awsapps.com/start
# AWS_CONFIG_SSO_REGION=us-east-1
# AWS_CONFIG_SSO_ROLE=AdministratorAccess
# AWS_CONFIG_DEFAULT_REGION=us-east-1
# AWS_CONFIG_CONFIG_FILE=~/.aws/config
# AWS_CONFIG_BACKUP_CONFIGS=true
# AWS_CONFIG_DRY_RUN=false
`
}
