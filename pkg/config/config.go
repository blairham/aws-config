package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// SSO configuration
	SSOStartURL string `mapstructure:"sso_start_url" yaml:"sso_start_url" json:"sso_start_url" toml:"sso_start_url"`
	SSORegion   string `mapstructure:"sso_region" yaml:"sso_region" json:"sso_region" toml:"sso_region"`
	SSORole     string `mapstructure:"sso_role" yaml:"sso_role" json:"sso_role" toml:"sso_role"`

	// AWS configuration
	DefaultRegion string `mapstructure:"default_region" yaml:"default_region" json:"default_region" toml:"default_region"`
	ConfigFile    string `mapstructure:"config_file" yaml:"config_file" json:"config_file" toml:"config_file"`

	// Behavior settings
	BackupConfigs bool `mapstructure:"backup_configs" yaml:"backup_configs" json:"backup_configs" toml:"backup_configs"`
	DryRun        bool `mapstructure:"dry_run" yaml:"dry_run" json:"dry_run" toml:"dry_run"`
}

// Load loads configuration from file, environment variables, and defaults
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file search paths and name
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Default config file locations
		home, _ := homedir.Dir()
		v.SetConfigName("aws-sso-config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath(home)
		v.AddConfigPath(filepath.Join(home, ".config", "aws-sso-config"))
		v.AddConfigPath("/etc/aws-sso-config/")
	}

	// Environment variable configuration
	v.SetEnvPrefix("AWS_SSO_CONFIG")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Read config file (optional - don't fail if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Expand home directory in paths
	if strings.HasPrefix(config.ConfigFile, "~") {
		expanded, err := homedir.Expand(config.ConfigFile)
		if err == nil {
			config.ConfigFile = expanded
		}
	}

	return &config, nil
}

// Default returns a configuration with default values (for backward compatibility)
func Default() *Config {
	config, _ := Load("")
	return config
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	home, _ := homedir.Dir()

	v.SetDefault("sso_start_url", "https://your-sso-portal.awsapps.com/start")
	v.SetDefault("sso_region", "us-east-1")
	v.SetDefault("sso_role", "AdministratorAccess")
	v.SetDefault("default_region", "us-east-1")
	v.SetDefault("config_file", filepath.Join(home, ".aws", "config"))
	v.SetDefault("backup_configs", true)
	v.SetDefault("dry_run", false)
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.SSOStartURL == "" {
		return fmt.Errorf("SSO start URL is required")
	}
	if c.SSORegion == "" {
		return fmt.Errorf("SSO region is required")
	}
	if c.DefaultRegion == "" {
		return fmt.Errorf("default region is required")
	}
	if c.ConfigFile == "" {
		return fmt.Errorf("config file path is required")
	}
	return nil
}

// WriteExample creates an example configuration file
func WriteExample(path string) error {
	example := `# AWS Config Tool Configuration
# This file supports YAML, TOML, and JSON formats

# SSO Configuration
sso_start_url: "https://your-sso-portal.awsapps.com/start"
sso_region: "us-east-1"
sso_role: "AdministratorAccess"

# AWS Configuration
default_region: "us-east-1"
config_file: "~/.aws/config"

# Behavior Settings
backup_configs: true
dry_run: false

# Environment variables can also be used with AWS_CONFIG_ prefix:
# AWS_CONFIG_SSO_START_URL=https://your-sso-portal.awsapps.com/start
# AWS_CONFIG_SSO_REGION=us-east-1
# AWS_CONFIG_SSO_ROLE=AdministratorAccess
# AWS_CONFIG_DEFAULT_REGION=us-east-1
# AWS_CONFIG_CONFIG_FILE=~/.aws/config
# AWS_CONFIG_BACKUP_CONFIGS=true
# AWS_CONFIG_DRY_RUN=false
`

	return os.WriteFile(path, []byte(example), 0600)
}
