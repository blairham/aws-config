package config

const synopsis = "Interact with aws-config's Centralized Configurations"
const help = `
Usage: aws-config config <subcommand> [options] [args]

    This command has subcommands for interacting with aws-config's Centralized
    Configuration system. Here are some simple examples, and more detailed
    examples are available in the subcommands or the documentation.

    Write a config:

    $ aws-config config write web.serviceconf.hcl

    Read a config:

    $ aws-config config read -kind service-defaults -name web

    List all configs for a type:

    $ aws-config config list -kind service-defaults

    Delete a config:

    $ aws config delete -kind service-defaults -name web

    For more examples, ask for subcommand help or view the documentation.
`
