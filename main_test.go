package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRunHelp tests the Run function with help argument
func TestRunHelp(t *testing.T) {
	// Test with help command (actual mitchellh/cli behavior - help returns 127)
	exitCode := Run([]string{"help"})
	// Note: mitchellh/cli returns 127 for "help" subcommand
	assert.Equal(t, 127, exitCode, "Help command should return exit code 127 (mitchellh/cli behavior)")
}

// TestRunHelpFlag tests the Run function with help flag
func TestRunHelpFlag(t *testing.T) {
	// Test with --help flag
	exitCode := Run([]string{"--help"})
	assert.Equal(t, 0, exitCode, "Help flag should return exit code 0")
}

// TestRunWithInvalidCommand tests the Run function with an invalid command
func TestRunWithInvalidCommand(t *testing.T) {
	// Test with an invalid command
	exitCode := Run([]string{"invalidcommand"})
	assert.NotEqual(t, 0, exitCode, "Invalid command should return non-zero exit code")
}

// TestRunNoArguments tests the Run function with no arguments
func TestRunNoArguments(t *testing.T) {
	// Test with no arguments (shows help and returns 127)
	exitCode := Run([]string{})
	assert.Equal(t, 127, exitCode, "No arguments should show help and return exit code 127")
}

// TestRunVersion tests the Run function with version flag
func TestRunVersion(t *testing.T) {
	// Test with version flag
	exitCode := Run([]string{"-version"})
	assert.Equal(t, 0, exitCode, "Version flag should return exit code 0")
}

// TestMainFunction tests the main function indirectly since we can't test os.Exit
func TestMainFunction(t *testing.T) {
	// Save original os.Exit and os.Args
	oldExit := osExit
	oldArgs := os.Args
	defer func() {
		osExit = oldExit
		os.Args = oldArgs
	}()

	var exitCode int
	// Override osExit with mock
	osExit = func(code int) {
		exitCode = code
		// Don't actually exit
	}

	// Set test arguments for a successful command
	os.Args = []string{"aws-sso-config", "-version"}

	// Call main which should call our mocked osExit with exit code 0
	main()

	assert.Equal(t, 0, exitCode, "Main should call os.Exit with 0 for version command")
}

// TestMainWithDifferentCommands tests main function with different commands
func TestMainWithDifferentCommands(t *testing.T) {
	// Save original os.Exit and os.Args
	oldExit := osExit
	oldArgs := os.Args
	defer func() {
		osExit = oldExit
		os.Args = oldArgs
	}()

	var exitCode int
	// Override osExit with mock
	osExit = func(code int) {
		exitCode = code
		// Don't actually exit
	}

	// Test with help command which should exit with code 127 (mitchellh/cli behavior)
	os.Args = []string{"aws-sso-config", "help"}
	main()
	assert.Equal(t, 127, exitCode, "Main should call os.Exit with 127 for help command")

	// Test with invalid command which should exit with non-zero code
	os.Args = []string{"aws-sso-config", "nonexistentcommand"}
	main()
	assert.NotEqual(t, 0, exitCode, "Main should call os.Exit with non-zero code for invalid command")
}
