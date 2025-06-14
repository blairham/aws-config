package initialize

const (
	help = `
Usage: aws-sso-config init [options]

  Initialize a configuration file with default values.

Options:

  -file=<path>    Path to the configuration file to create.
                  Default: aws-sso-config.yaml in current directory.

  -format=<fmt>   Format for the configuration file (yaml, json, toml).
                  Default: yaml.

Examples:

  # Create a YAML config file in the current directory
  aws-sso-config init

  # Create a JSON config file at a specific location
  aws-sso-config init -file=config.json -format=json

  # Create a TOML config file
  aws-sso-config init -file=config.toml -format=toml
`
)
