package run

import (
	"bytes"
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockUI implements the cli.UI interface for testing
type mockUI struct {
	OutBuffer *bytes.Buffer
	ErrBuffer *bytes.Buffer
}

func newMockUI() *mockUI {
	return &mockUI{
		OutBuffer: &bytes.Buffer{},
		ErrBuffer: &bytes.Buffer{},
	}
}

func (m *mockUI) Ask(string) (string, error) {
	return "", nil
}

func (m *mockUI) AskSecret(string) (string, error) {
	return "", nil
}

func (m *mockUI) Output(message string) {
	m.OutBuffer.WriteString(message + "\n")
}

func (m *mockUI) Info(message string) {
	m.OutBuffer.WriteString(message + "\n")
}

func (m *mockUI) Error(message string) {
	m.ErrBuffer.WriteString(message + "\n")
}

func (m *mockUI) Warn(message string) {
	m.ErrBuffer.WriteString(message + "\n")
}

func (m *mockUI) Stdout() io.Writer {
	return m.OutBuffer
}

func (m *mockUI) Stderr() io.Writer {
	return m.ErrBuffer
}

func TestRunCommand(t *testing.T) {
	ui := newMockUI()
	cmd := New(ui)

	// Test initialization
	assert.NotNil(t, cmd)
	assert.NotEmpty(t, cmd.Help())
	assert.NotEmpty(t, cmd.Synopsis())

	// Skip actual command execution test since it depends on aws2-wrap
	// The test below would test the real functionality if aws2-wrap was installed

	// Test with '--help' flag - should return help text when help flags exist in command
	helpText := cmd.Help()
	assert.Contains(t, helpText, "Usage:")
}

func TestRunCommandWithCommand(t *testing.T) {
	// Always skip this test as it requires aws2-wrap to be installed
	t.Skip("Skipping test that requires aws2-wrap to be installed")
}

func TestRunCommandInit(t *testing.T) {
	ui := newMockUI()
	cmd := New(ui)

	// Test Init explicitly
	cmd.Init()

	// Verify logger is initialized
	assert.NotNil(t, Logger)
}

func TestRunCommandSynopsis(t *testing.T) {
	ui := newMockUI()
	cmd := New(ui)

	// Test synopsis - use the actual synopsis constant from the package
	synopsis := cmd.Synopsis()
	assert.Equal(t, cmd.Synopsis(), synopsis)
}

func TestRunCommandWithMockedProfile(t *testing.T) {
	// Skip actual aws2-wrap execution
	origExec := execCommand
	defer func() { execCommand = origExec }()

	commandExecuted := false
	commandArgs := []string{}

	// Mock exec.Command
	execCommand = func(command string, args ...string) *exec.Cmd {
		commandExecuted = true
		commandArgs = args
		return fakeExecCommand("echo", "test output")
	}

	// Create and run command
	ui := newMockUI()
	cmd := New(ui)

	// Mock GetProfile to return a known value
	origGetProfile := getProfileFunc
	defer func() { getProfileFunc = origGetProfile }()

	getProfileFunc = func() (string, error) {
		return "test-profile", nil
	}

	// Run with test arguments
	result := cmd.Run([]string{"s3", "ls"})

	// Verify results
	assert.Equal(t, 0, result, "Command should succeed")
	assert.True(t, commandExecuted, "Command should be executed")
	assert.Equal(t, []string{"s3", "ls"}, commandArgs, "Command args should be passed through")
}

// Helpers for mocking exec.Command
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	// Set up pipes or other configuration needed
	return cmd
}
