package initialize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Verify the command is properly initialized
	assert.NotNil(t, c.flags)
	assert.NotEmpty(t, c.help)
	assert.Equal(t, "Initialize a configuration file with default values", c.Synopsis())
}

func TestInitDefaultYAMLFile(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	ui := cli.NewMockUi()
	c := New(ui)

	// Run init command with no arguments (should create aws-sso-config.yaml)
	exitCode := c.Run([]string{})
	assert.Equal(t, 0, exitCode)

	// Check that file was created
	_, err := os.Stat("aws-sso-config.yaml")
	assert.NoError(t, err)

	// Check file content
	content, err := os.ReadFile("aws-sso-config.yaml")
	require.NoError(t, err)
	assert.Contains(t, string(content), "sso_start_url:")
	assert.Contains(t, string(content), "AWS Config Tool Configuration")
	assert.Contains(t, string(content), "# SSO Configuration")
}

func TestInitJSONFile(t *testing.T) {
	testInitWithFormat(t, "json", "test-config.json", `"sso_start_url":`, `"_comment":`)
}

func TestInitTOMLFile(t *testing.T) {
	testInitWithFormat(t, "toml", "test-config.toml", `sso_start_url = `, "# AWS Config Tool Configuration")
}

// Helper function to reduce duplicate code
func testInitWithFormat(t *testing.T, format, filename, contentCheck1, contentCheck2 string) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	ui := cli.NewMockUi()
	c := New(ui)

	// Run init command with specified format
	exitCode := c.Run([]string{"-format=" + format, "-file=" + filename})
	assert.Equal(t, 0, exitCode)

	// Check that file was created
	_, err := os.Stat(filename)
	assert.NoError(t, err)

	// Check file content
	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	assert.Contains(t, string(content), contentCheck1)
	assert.Contains(t, string(content), contentCheck2)
}

func TestInitInvalidFormat(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	ui := cli.NewMockUi()
	c := New(ui)

	// Run init command with invalid format
	exitCode := c.Run([]string{"-format=xml"})
	assert.Equal(t, 1, exitCode)

	// Check error message
	assert.Contains(t, ui.ErrorWriter.String(), "format must be one of: yaml, json, toml")
}

func TestInitFileAlreadyExists(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	// Create an existing file
	err := os.WriteFile("existing.yaml", []byte("test"), 0600)
	require.NoError(t, err)

	ui := cli.NewMockUi()
	c := New(ui)

	// Try to create file that already exists
	exitCode := c.Run([]string{"-file=existing.yaml"})
	assert.Equal(t, 1, exitCode)

	// Check error message
	assert.Contains(t, ui.ErrorWriter.String(), "file existing.yaml already exists")
}

func TestInitCreateDirectories(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	ui := cli.NewMockUi()
	c := New(ui)

	// Create config in nested directory
	nestedPath := filepath.Join("config", "subdir", "aws-sso-config.yaml")
	exitCode := c.Run([]string{"-file=" + nestedPath})
	assert.Equal(t, 0, exitCode)

	// Check that directories were created
	_, err := os.Stat(nestedPath)
	assert.NoError(t, err)

	// Check that content is correct
	content, err := os.ReadFile(nestedPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "sso_start_url:")
}

func TestInitInvalidDirectoryPath(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	// Try to create in a path we can't create (assuming /root is not writable)
	exitCode := c.Run([]string{"-file=/root/cannot-create/config.yaml"})
	assert.Equal(t, 1, exitCode)

	// Should contain error about directory creation
	errorOutput := ui.ErrorWriter.String()
	assert.Contains(t, errorOutput, "Error creating directory")
}

func TestInitFormatExtensionMapping(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	tests := []struct {
		format   string
		expected string
	}{
		{"json", "aws-sso-config.json"},
		{"toml", "aws-sso-config.toml"},
		{"yaml", "aws-sso-config.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			ui := cli.NewMockUi()
			c := New(ui)

			exitCode := c.Run([]string{"-format=" + tt.format})
			assert.Equal(t, 0, exitCode)

			// Check that file was created with correct extension
			_, err := os.Stat(tt.expected)
			assert.NoError(t, err)

			// Clean up
			os.Remove(tt.expected)
		})
	}
}

func TestInitHelpOutput(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	help := c.Help()
	assert.Contains(t, help, "Usage: aws-sso-config init")
	assert.Contains(t, help, "-file=<path>")
	assert.Contains(t, help, "-format=<fmt>")
	assert.Contains(t, help, "Examples:")
}

func TestInitSynopsis(t *testing.T) {
	ui := cli.NewMockUi()
	c := New(ui)

	synopsis := c.Synopsis()
	assert.Equal(t, "Initialize a configuration file with default values", synopsis)
}

func TestInitOutputMessages(t *testing.T) {
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	ui := cli.NewMockUi()
	c := New(ui)

	exitCode := c.Run([]string{"-file=test-output.yaml"})
	assert.Equal(t, 0, exitCode)

	output := ui.OutputWriter.String()
	assert.Contains(t, output, "Successfully created configuration file: test-output.yaml")
	assert.Contains(t, output, "Edit the file to customize your settings")
	assert.Contains(t, output, "aws-sso-config generate -config=test-output.yaml")
}

func TestInitContentGeneration(t *testing.T) {
	tests := []struct {
		format            string
		expectedContent   []string
		unexpectedContent []string
	}{
		{
			format: "yaml",
			expectedContent: []string{
				"sso_start_url:",
				"# SSO Configuration",
				"backup_configs: true",
			},
			unexpectedContent: []string{
				`"sso_start_url":`,
				`sso_start_url =`,
			},
		},
		{
			format: "json",
			expectedContent: []string{
				`"sso_start_url":`,
				`"backup_configs": true`,
			},
			unexpectedContent: []string{
				"sso_start_url:",
				"sso_start_url =",
			},
		},
		{
			format: "toml",
			expectedContent: []string{
				"sso_start_url =",
				"# SSO Configuration",
				"backup_configs = true",
			},
			unexpectedContent: []string{
				"sso_start_url:",
				`"sso_start_url":`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			tempDir := t.TempDir()
			originalWd, _ := os.Getwd()
			os.Chdir(tempDir)
			defer os.Chdir(originalWd)

			ui := cli.NewMockUi()
			c := New(ui)

			filename := "test." + tt.format
			exitCode := c.Run([]string{"-format=" + tt.format, "-file=" + filename})
			assert.Equal(t, 0, exitCode)

			content, err := os.ReadFile(filename)
			require.NoError(t, err)
			contentStr := string(content)

			for _, expected := range tt.expectedContent {
				assert.Contains(t, contentStr, expected, "Expected content not found in %s format", tt.format)
			}

			for _, unexpected := range tt.unexpectedContent {
				assert.NotContains(t, contentStr, unexpected, "Unexpected content found in %s format", tt.format)
			}
		})
	}
}
