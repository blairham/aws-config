package generate

const synopsis = "Generate AWS config file"
const help = `
Usage: aws-sso-config generate [options]

  This command will auto-generate an AWS config file
  with all accounts you have access to.

Options:

  -diff             Enable diff output to see changes before writing.

  -config=<path>    Path to configuration file. If not specified,
                    uses environment variables and defaults.

Examples:

  # Generate using environment variables and defaults
  aws-sso-config generate

  # Generate using a custom config file
  aws-sso-config generate -config=my-config.yaml

  # Show diff before writing changes
  aws-sso-config generate -diff -config=my-config.yaml
`
