package list

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestListCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	if cmd == nil {
		t.Fatal("New() returned nil")
	}
}

func TestListNoArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	code := cmd.Run([]string{})

	if code != 0 {
		t.Errorf("Expected exit code 0, got %d", code)
	}

	output := stdout.String()
	if !strings.Contains(output, "Available configuration keys:") {
		t.Errorf("Expected output to contain 'Available configuration keys:', got %s", output)
	}

	if !strings.Contains(output, "sso.start_url") {
		t.Errorf("Expected output to contain 'sso.start_url', got %s", output)
	}
}

func TestListWithArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	code := cmd.Run([]string{"extra", "args"})

	if code != 1 {
		t.Errorf("Expected exit code 1, got %d", code)
	}

	errorOutput := stderr.String()
	if !strings.Contains(errorOutput, "Usage: aws-sso-config config list") {
		t.Errorf("Expected error to contain usage message, got %s", errorOutput)
	}

	if !strings.Contains(errorOutput, "This command takes no arguments.") {
		t.Errorf("Expected error to contain no arguments message, got %s", errorOutput)
	}
}

func TestListHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	help := cmd.Help()

	if !strings.Contains(help, "Usage: aws-sso-config config list") {
		t.Errorf("Expected help to contain usage, got %s", help)
	}

	if !strings.Contains(help, "List all available configuration keys") {
		t.Errorf("Expected help to contain description, got %s", help)
	}
}

func TestListSynopsis(t *testing.T) {
	var stdout, stderr bytes.Buffer
	ui := &cli.BasicUi{
		Writer:      &stdout,
		ErrorWriter: &stderr,
	}

	cmd := New(ui)
	synopsis := cmd.Synopsis()

	expected := "List all available configuration keys"
	if synopsis != expected {
		t.Errorf("Expected synopsis %q, got %q", expected, synopsis)
	}
}
