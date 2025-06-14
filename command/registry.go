package command

import (
	"fmt"

	mcli "github.com/mitchellh/cli"

	"github.com/blairham/aws-config/command/cli"
	"github.com/blairham/aws-config/command/config"
	configwrite "github.com/blairham/aws-config/command/config/write"
	"github.com/blairham/aws-config/command/generate"
	"github.com/blairham/aws-config/command/run"
)

// factory is a function that returns a new instance of a CLI-sub command.
type factory func(cli.UI) (cli.Command, error)

// entry is a struct that contains a command's name and a factory for that command.
type entry struct {
	name string
	fn   factory
}

func RegisteredCommands(ui cli.UI) map[string]mcli.CommandFactory {
	registry := map[string]mcli.CommandFactory{}
	registerCommands(ui, registry,
		// Add new commands here
		entry{"config", func(ui cli.UI) (cli.Command, error) { return config.New(ui), nil }},
		// entry{"config list", func(ui cli.UI) (cli.Command, error) { return configlist.New(ui), nil }},
		// entry{"config read", func(ui cli.UI) (cli.Command, error) { return configread.New(ui), nil }},
		entry{"config write", func(ui cli.UI) (cli.Command, error) { return configwrite.New(ui), nil }},
		entry{"generate", func(ui cli.UI) (cli.Command, error) { return generate.New(ui), nil }},
		entry{"run", func(ui cli.UI) (cli.Command, error) { return run.New(ui), nil }},
	)

	return registry
}

func registerCommands(ui cli.UI, m map[string]mcli.CommandFactory, cmdEntries ...entry) {
	for _, ent := range cmdEntries {
		thisFn := ent.fn
		if _, ok := m[ent.name]; ok {
			panic(fmt.Sprintf("duplicate command: %q", ent.name))
		}
		m[ent.name] = func() (mcli.Command, error) {
			return thisFn(ui)
		}
	}
}
