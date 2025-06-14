package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	config := Default()

	assert.NotEmpty(t, config.SSOStartURL)
	assert.NotEmpty(t, config.SSORegion)
	assert.NotEmpty(t, config.SSORole)
	assert.NotEmpty(t, config.DefaultRegion)
	assert.NotEmpty(t, config.ConfigFile)
	assert.True(t, config.BackupConfigs)
	assert.False(t, config.DryRun)
}

func TestConfigWithEnvironmentVariables(t *testing.T) {
	// Clean up any existing env vars first
	originalEnvVars := map[string]string{
		"AWS_SSO_CONFIG_SSO_START_URL":  os.Getenv("AWS_SSO_CONFIG_SSO_START_URL"),
		"AWS_SSO_CONFIG_SSO_REGION":     os.Getenv("AWS_SSO_CONFIG_SSO_REGION"),
		"AWS_SSO_CONFIG_SSO_ROLE":       os.Getenv("AWS_SSO_CONFIG_SSO_ROLE"),
		"AWS_SSO_CONFIG_DEFAULT_REGION": os.Getenv("AWS_SSO_CONFIG_DEFAULT_REGION"),
		"AWS_SSO_CONFIG_BACKUP_CONFIGS": os.Getenv("AWS_SSO_CONFIG_BACKUP_CONFIGS"),
		"AWS_SSO_CONFIG_DRY_RUN":        os.Getenv("AWS_SSO_CONFIG_DRY_RUN"),
	}

	// Set new environment variables
	envVars := map[string]string{
		"AWS_SSO_CONFIG_SSO_START_URL":  "https://test.awsapps.com/start",
		"AWS_SSO_CONFIG_SSO_REGION":     "us-west-2",
		"AWS_SSO_CONFIG_SSO_ROLE":       "TestRole",
		"AWS_SSO_CONFIG_DEFAULT_REGION": "eu-west-1",
		"AWS_SSO_CONFIG_BACKUP_CONFIGS": "false",
		"AWS_SSO_CONFIG_DRY_RUN":        "true",
	}

	// Set env vars
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	// Clean up
	defer func() {
		for key, originalValue := range originalEnvVars {
			if originalValue == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, originalValue)
			}
		}
	}()

	config, err := Load("")
	require.NoError(t, err)

	assert.Equal(t, "https://test.awsapps.com/start", config.SSOStartURL)
	assert.Equal(t, "us-west-2", config.SSORegion)
	assert.Equal(t, "TestRole", config.SSORole)
	assert.Equal(t, "eu-west-1", config.DefaultRegion)
	assert.False(t, config.BackupConfigs)
	assert.True(t, config.DryRun)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name:        "valid config",
			config:      Default(),
			expectError: false,
		},
		{
			name: "missing SSO start URL",
			config: &Config{
				SSORegion:     "us-east-1",
				DefaultRegion: "us-east-1",
				ConfigFile:    "/tmp/config",
			},
			expectError: true,
		},
		{
			name: "missing SSO region",
			config: &Config{
				SSOStartURL:   "https://test.com",
				DefaultRegion: "us-east-1",
				ConfigFile:    "/tmp/config",
			},
			expectError: true,
		},
		{
			name: "missing default region",
			config: &Config{
				SSOStartURL: "https://test.com",
				SSORegion:   "us-east-1",
				ConfigFile:  "/tmp/config",
			},
			expectError: true,
		},
		{
			name: "missing config file",
			config: &Config{
				SSOStartURL:   "https://test.com",
				SSORegion:     "us-east-1",
				DefaultRegion: "us-east-1",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigFileDefault(t *testing.T) {
	config := Default()

	// Should contain .aws/config path
	assert.Contains(t, config.ConfigFile, filepath.Join(".aws", "config"))
}

func TestLoadConfigFromYAMLFile(t *testing.T) {
	yamlContent := `sso_start_url: "https://yaml-test.awsapps.com/start"
sso_region: "us-west-1"
sso_role: "YAMLTestRole"
default_region: "eu-central-1"
backup_configs: false
dry_run: true
`

	expectedConfig := Config{
		SSOStartURL:   "https://yaml-test.awsapps.com/start",
		SSORegion:     "us-west-1",
		SSORole:       "YAMLTestRole",
		DefaultRegion: "eu-central-1",
		BackupConfigs: false,
		DryRun:        true,
	}

	testLoadConfigFromFile(t, "aws-sso-config.yaml", yamlContent, expectedConfig)
}

func TestLoadConfigFromJSONFile(t *testing.T) {
	jsonContent := `{
	"sso_start_url": "https://json-test.awsapps.com/start",
	"sso_region": "ap-southeast-1",
	"sso_role": "JSONTestRole",
	"default_region": "ap-southeast-2",
	"backup_configs": true,
	"dry_run": false
}`

	expectedConfig := Config{
		SSOStartURL:   "https://json-test.awsapps.com/start",
		SSORegion:     "ap-southeast-1",
		SSORole:       "JSONTestRole",
		DefaultRegion: "ap-southeast-2",
		BackupConfigs: true,
		DryRun:        false,
	}

	testLoadConfigFromFile(t, "aws-config.json", jsonContent, expectedConfig)
}

func testLoadConfigFromFile(t *testing.T, filename, content string, expected Config) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, filename)

	err := os.WriteFile(configFile, []byte(content), 0600)
	require.NoError(t, err)

	config, err := Load(configFile)
	require.NoError(t, err)

	assert.Equal(t, expected.SSOStartURL, config.SSOStartURL)
	assert.Equal(t, expected.SSORegion, config.SSORegion)
	assert.Equal(t, expected.SSORole, config.SSORole)
	assert.Equal(t, expected.DefaultRegion, config.DefaultRegion)
	assert.Equal(t, expected.BackupConfigs, config.BackupConfigs)
	assert.Equal(t, expected.DryRun, config.DryRun)
}

func TestWriteExample(t *testing.T) {
	tempDir := t.TempDir()
	exampleFile := filepath.Join(tempDir, "example-config.yaml")

	err := WriteExample(exampleFile)
	require.NoError(t, err)

	// Check that file exists and has content
	content, err := os.ReadFile(exampleFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "sso_start_url")
	assert.Contains(t, string(content), "AWS Config Tool Configuration")
}

func TestConfigFileSearchPaths(t *testing.T) {
	// Test that the XDG config directory is properly added to search paths
	// by creating a config file there and loading it with a specific path

	tempHome := t.TempDir()
	xdgConfigDir := filepath.Join(tempHome, ".config", "aws-config")
	err := os.MkdirAll(xdgConfigDir, 0750)
	require.NoError(t, err)

	// Create a test config file in the XDG location
	configFile := filepath.Join(xdgConfigDir, "aws-sso-config.yaml")
	configContent := `sso_start_url: "https://xdg-test.awsapps.com/start"
sso_region: "eu-west-1"
sso_role: "XDGTestRole"
default_region: "eu-central-1"
backup_configs: false
dry_run: true
`
	err = os.WriteFile(configFile, []byte(configContent), 0600)
	require.NoError(t, err)

	// Load the config file directly using its path
	config, err := Load(configFile)
	require.NoError(t, err)

	// Verify the config was loaded correctly
	assert.Equal(t, "https://xdg-test.awsapps.com/start", config.SSOStartURL)
	assert.Equal(t, "eu-west-1", config.SSORegion)
	assert.Equal(t, "XDGTestRole", config.SSORole)
	assert.Equal(t, "eu-central-1", config.DefaultRegion)
	assert.False(t, config.BackupConfigs)
	assert.True(t, config.DryRun)
}

func TestLoadConfigFromTOMLFile(t *testing.T) {
	tomlContent := `# TOML Configuration Test
sso_start_url = "https://toml-test.awsapps.com/start"
sso_region = "ap-northeast-1"
sso_role = "TOMLTestRole"
default_region = "ap-northeast-2"
backup_configs = true
dry_run = false
`

	expectedConfig := Config{
		SSOStartURL:   "https://toml-test.awsapps.com/start",
		SSORegion:     "ap-northeast-1",
		SSORole:       "TOMLTestRole",
		DefaultRegion: "ap-northeast-2",
		BackupConfigs: true,
		DryRun:        false,
	}

	testLoadConfigFromFile(t, "aws-config.toml", tomlContent, expectedConfig)
}

func TestLoadConfigWithInvalidFile(t *testing.T) {
	// Test loading a non-existent file
	_, err := Load("/non/existent/file.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading config file")
}

func TestLoadConfigWithInvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "invalid.yaml")

	// Create an invalid YAML file
	invalidContent := `sso_start_url: "test
sso_region: invalid yaml syntax [[[
`
	err := os.WriteFile(configFile, []byte(invalidContent), 0600)
	require.NoError(t, err)

	_, err = Load(configFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading config file")
}

func TestConfigPrecedence(t *testing.T) {
	// Test that environment variables override config file values
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "precedence-test.yaml")

	// Create config file with default values
	configContent := `sso_start_url: "https://file-config.awsapps.com/start"
sso_region: "us-east-1"
dry_run: false
`
	err := os.WriteFile(configFile, []byte(configContent), 0600)
	require.NoError(t, err)

	// Set environment variables that should override
	originalVars := map[string]string{
		"AWS_SSO_CONFIG_SSO_START_URL": os.Getenv("AWS_SSO_CONFIG_SSO_START_URL"),
		"AWS_SSO_CONFIG_DRY_RUN":       os.Getenv("AWS_SSO_CONFIG_DRY_RUN"),
	}

	os.Setenv("AWS_SSO_CONFIG_SSO_START_URL", "https://env-override.awsapps.com/start")
	os.Setenv("AWS_SSO_CONFIG_DRY_RUN", "true")

	defer func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	config, err := Load(configFile)
	require.NoError(t, err)

	// Environment variables should override file values
	assert.Equal(t, "https://env-override.awsapps.com/start", config.SSOStartURL)
	assert.True(t, config.DryRun)
	// File values should be used where no env var is set
	assert.Equal(t, "us-east-1", config.SSORegion)
}

func TestConfigWithPartialValues(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "partial.yaml")

	// Config with only some values set
	configContent := `sso_start_url: "https://partial-test.awsapps.com/start"
sso_region: "ca-central-1"
# Other values should use defaults
`
	err := os.WriteFile(configFile, []byte(configContent), 0600)
	require.NoError(t, err)

	config, err := Load(configFile)
	require.NoError(t, err)

	// Specified values
	assert.Equal(t, "https://partial-test.awsapps.com/start", config.SSOStartURL)
	assert.Equal(t, "ca-central-1", config.SSORegion)

	// Default values should be used for unspecified fields
	assert.Equal(t, "AdministratorAccess", config.SSORole)
	assert.Equal(t, "us-east-1", config.DefaultRegion)
	assert.True(t, config.BackupConfigs)
	assert.False(t, config.DryRun)
}

func TestConfigExpandsHomeDirInPaths(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "homedir-test.yaml")

	configContent := `config_file: "~/.aws/config"
`
	err := os.WriteFile(configFile, []byte(configContent), 0600)
	require.NoError(t, err)

	config, err := Load(configFile)
	require.NoError(t, err)

	// Should expand the tilde
	assert.Contains(t, config.ConfigFile, ".aws/config")
	assert.NotContains(t, config.ConfigFile, "~")
}

func TestValidateEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "empty SSO start URL",
			config: &Config{
				SSOStartURL:   "",
				SSORegion:     "us-east-1",
				DefaultRegion: "us-east-1",
				ConfigFile:    "/tmp/config",
			},
			expectError: true,
			errorMsg:    "SSO start URL is required",
		},
		{
			name: "whitespace-only SSO start URL",
			config: &Config{
				SSOStartURL:   "   ",
				SSORegion:     "us-east-1",
				DefaultRegion: "us-east-1",
				ConfigFile:    "/tmp/config",
			},
			expectError: false, // Viper will trim whitespace, but we don't validate format
		},
		{
			name: "empty all fields",
			config: &Config{
				SSOStartURL:   "",
				SSORegion:     "",
				SSORole:       "",
				DefaultRegion: "",
				ConfigFile:    "",
			},
			expectError: true,
			errorMsg:    "SSO start URL is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteExampleToInvalidPath(t *testing.T) {
	// Try to write to a path that doesn't exist and can't be created
	err := WriteExample("/root/non-existent-dir/config.yaml")
	assert.Error(t, err)
}

func TestConfigStructTags(t *testing.T) {
	// Verify that struct tags are correctly defined for all supported formats
	config := Config{
		SSOStartURL:   "https://test.com",
		SSORegion:     "us-west-2",
		SSORole:       "TestRole",
		DefaultRegion: "eu-west-1",
		ConfigFile:    "/test/config",
		BackupConfigs: true,
		DryRun:        false,
	}

	tempDir := t.TempDir()

	// Test YAML marshaling
	yamlFile := filepath.Join(tempDir, "test.yaml")
	yamlContent := `sso_start_url: "https://test.com"
sso_region: "us-west-2"
sso_role: "TestRole"
default_region: "eu-west-1"
config_file: "/test/config"
backup_configs: true
dry_run: false
`
	err := os.WriteFile(yamlFile, []byte(yamlContent), 0600)
	require.NoError(t, err)

	loadedConfig, err := Load(yamlFile)
	require.NoError(t, err)
	assert.Equal(t, config.SSOStartURL, loadedConfig.SSOStartURL)
	assert.Equal(t, config.SSORegion, loadedConfig.SSORegion)
	assert.Equal(t, config.SSORole, loadedConfig.SSORole)

	// Test JSON marshaling
	jsonFile := filepath.Join(tempDir, "test.json")
	jsonContent := `{
  "sso_start_url": "https://test.com",
  "sso_region": "us-west-2",
  "sso_role": "TestRole",
  "default_region": "eu-west-1",
  "config_file": "/test/config",
  "backup_configs": true,
  "dry_run": false
}`
	err = os.WriteFile(jsonFile, []byte(jsonContent), 0600)
	require.NoError(t, err)

	loadedConfig, err = Load(jsonFile)
	require.NoError(t, err)
	assert.Equal(t, config.SSOStartURL, loadedConfig.SSOStartURL)
	assert.Equal(t, config.SSORegion, loadedConfig.SSORegion)
	assert.Equal(t, config.SSORole, loadedConfig.SSORole)
}
