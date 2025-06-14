package main

import (
	"fmt"
	"os"

	mcli "github.com/mitchellh/cli"

	"github.com/blairham/aws-config/command"
	"github.com/blairham/aws-config/command/cli"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {
	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	cmds := command.RegisteredCommands(ui)
	var names []string
	for c := range cmds {
		names = append(names, c)
	}

	cliInstance := &mcli.CLI{
		Name: "aws-config",
		// Version:                    version.GetVersion().FullVersionNumber(true),
		Args:                       args,
		Commands:                   cmds,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
		HelpFunc:                   mcli.FilteredHelpFunc(names, mcli.BasicHelpFunc("aws-config")),
		HelpWriter:                 os.Stdout,
	}

	exitCode, err := cliInstance.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %v\n", err)
		return 1
	}

	return exitCode
}
