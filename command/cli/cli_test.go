package cli_test

import (
	"bytes"
	"testing"

	mcli "github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"

	"github.com/blairham/aws-config/command/cli"
)

func TestBasicUI(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	ui := &cli.BasicUI{
		BasicUi: mcli.BasicUi{
			Writer:      &stdout,
			ErrorWriter: &stderr,
		},
	}

	ui.Output("Hello, world!")
	ui.Error("Oops, something went wrong.")

	assert.Equal(t, "Hello, world!\n", stdout.String())
	assert.Equal(t, "Oops, something went wrong.\n", stderr.String())
}
