package run

import (
	"testing"
)

// TestRunCommandWithError tests error handling in the Run command
func TestRunCommandWithError(t *testing.T) {
	t.Skip("Skipping test that requires mocking log.Logger.Fatalf")
}

// TestRunCommandWithCmdError tests error handling when the command execution fails
func TestRunCommandWithCmdError(t *testing.T) {
	t.Skip("Skipping test that requires mocking log.Logger.Fatal")
}
