# aws-sso-config

[![Go Report Card](https://goreportcard.com/badge/github.com/blairham/aws-sso-config)](https://goreportcard.com/report/github.com/blairham/aws-sso-config)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A command-line tool for managing AWS configuration files and SSO authentication.

## Features

- 🔧 **Automatic AWS Config Generation**: Generate AWS config files from SSO accounts
- 🚀 **Modern SSO Flow**: Browser-based authentication with automatic polling
- 🔄 **Profile Management**: Automatically detect and configure AWS profiles
- 🛡️ **Secure**: Uses AWS SDK v2 and follows security best practices
- 📦 **Easy Distribution**: Available as binary releases for multiple platforms

## Installation

### Download Binary

Download the latest release from the [releases page](https://github.com/blairham/aws-sso-config/releases).

### Build from Source

```bash
# Clone the repository
git clone https://github.com/blairham/aws-sso-config.git
cd aws-sso-config

# Build using Make
make build

# Or build using Go directly
go build -o aws-sso-config .
```

### Using Go Install

```bash
go install github.com/blairham/aws-sso-config@latest
```

## Usage

### Initialize Configuration

Create a configuration file with default settings:

```bash
# Create a YAML config file in current directory
aws-sso-config init

# Create a JSON config file
aws-sso-config init -format=json -file=my-config.json

# Create a TOML config file
aws-sso-config init -format=toml -file=my-config.toml
```

### Generate AWS Config

Generate an AWS config file with all accounts you have access to:

```bash
# Generate using default configuration
aws-sso-config generate

# Generate using a custom config file
aws-sso-config generate -config=my-config.yaml
```

Show differences before applying changes:

```bash
aws-sso-config generate --diff
```

### Run Commands with AWS Credentials

Execute commands with the appropriate AWS credentials automatically set:

```bash
aws-sso-config run aws s3 ls
aws-sso-config run terraform plan
```



## Configuration

aws-sso-config supports multiple configuration methods with the following precedence order (highest to lowest):

1. **Command-line flags** (e.g., `-config=my-config.yaml`)
2. **Configuration files** (YAML, JSON, or TOML)
3. **Environment variables** (with `AWS_CONFIG_` prefix)
4. **Default values**

### Configuration Files

Create a configuration file to customize aws-sso-config behavior:

```bash
# Create a YAML configuration file (default)
aws-sso-config init

# Create a JSON configuration file
aws-sso-config init -format=json -file=my-config.json

# Create a TOML configuration file  
aws-sso-config init -format=toml -file=my-config.toml
```

The configuration file will be created with these settings:

```yaml
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
```

Use your configuration file:

```bash
# Generate AWS config using your custom settings
aws-sso-config generate -config=my-config.yaml

# Show differences before applying
aws-sso-config generate -config=my-config.yaml -diff
```

### Configuration File Locations

aws-sso-config automatically searches for configuration files in:

1. Current directory (`./aws-sso-config.yaml`)
2. Home directory (`~/aws-sso-config.yaml`)
3. XDG config directory (`~/.config/aws-sso-config/aws-sso-config.yaml`)
4. System config directory (`/etc/aws-sso-config/aws-sso-config.yaml`)

Supported formats: `.yaml`, `.yml`, `.json`, `.toml`

### Environment Variables

You can override any configuration setting using environment variables with the `AWS_CONFIG_` prefix:

- `AWS_CONFIG_SSO_START_URL`: Your AWS SSO start URL
- `AWS_CONFIG_SSO_REGION`: AWS region for SSO (default: us-east-1)
- `AWS_CONFIG_SSO_ROLE`: SSO role name (default: AdministratorAccess)
- `AWS_CONFIG_DEFAULT_REGION`: Default AWS region (default: us-east-1)
- `AWS_CONFIG_CONFIG_FILE`: Path to AWS config file (default: ~/.aws/config)
- `AWS_CONFIG_BACKUP_CONFIGS`: Backup existing configs (default: true)
- `AWS_CONFIG_DRY_RUN`: Show changes without applying (default: false)

Example:

```bash
export AWS_CONFIG_SSO_START_URL="https://mycompany.awsapps.com/start"
export AWS_CONFIG_SSO_REGION="us-west-2"
export AWS_CONFIG_BACKUP_CONFIGS="false"

aws-sso-config generate
```

### Legacy Environment Variables

For compatibility, these environment variables are still supported:

- `AWS_PROFILE`: Override automatic profile detection

## Development

### Prerequisites

- Go 1.19+
- Make
- golangci-lint (for linting)
- goreleaser (for releases)

### Building

```bash
# Install dependencies
make deps

# Run all checks
make check

# Build for development
make build-dev

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Testing

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Generate coverage report (may fail on some tests)
make test-coverage
```

### Release Process

This project uses [GoReleaser](https://goreleaser.com) for automated releases:

1. **Create a new tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically**:
   - Builds binaries for multiple platforms (Linux, macOS, Windows)
   - Creates checksums and archives
   - Publishes the release on GitHub
   - Updates the changelog

3. **Manual release** (if needed):
   ```bash
   # Check GoReleaser configuration
   make goreleaser-check

   # Create a snapshot build for testing
   make snapshot

   # Create a full release (requires clean git state and proper tag)
   make release
   ```

### Pre-commit Hooks

Install pre-commit hooks for code quality:

```bash
pre-commit install
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linting (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

Please ensure your code follows the existing style and includes appropriate tests.

## Security

If you discover a security vulnerability, please send an email to [your-email]. All security vulnerabilities will be promptly addressed.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- AWS SDK for Go v2
- Mitchell Hashimoto's CLI library
- The Go community

---

**⚠️ Note**: This tool modifies your AWS configuration files. Please ensure you have backups before running.
