# SecureVault CLI

A secure password and secrets management CLI tool that stores encrypted credentials locally and integrates with cloud secret managers.

![Build Status](https://img.shields.io/github/actions/workflow/status/username/securevault-cli/ci.yml?branch=main)
![Go Version](https://img.shields.io/github/go-mod/go-version/username/securevault-cli)
![License](https://img.shields.io/github/license/username/securevault-cli)
![Release](https://img.shields.io/github/v/release/username/securevault-cli)

## Overview

SecureVault CLI helps developers, DevOps engineers, and system administrators
securely manage secrets from the command line. It uses AES-256 encryption for
local storage and seamlessly integrates with popular cloud secret managers like
AWS Secrets Manager, Azure Key Vault, and HashiCorp Vault.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [Contributing](#contributing)
- [Testing](#testing)
- [License](#license)

## Features

- **Secure Local Storage** - AES-256 encryption for credentials stored locally
- **Cloud Integration** - Connect to AWS Secrets Manager, Azure Key Vault, and
  HashiCorp Vault
- **Import/Export** - Easily migrate secrets between environments
- **Cross-Platform** - Works on Linux, macOS, and Windows
- **Shell Completion** - Built-in support for bash, zsh, and fish shells
- **Audit Logging** - Track all secret access and modifications
- **Git-Safe** - Automatically excluded from version control

## Installation

### Prerequisites

- Go 1.21+ (only for building from source)
- No runtime dependencies required

### Quick Install

**macOS/Linux (Homebrew):**

```bash
brew install username/tap/securevault-cli
```

**All Platforms (Binary Release):**

Download the latest release from
[GitHub Releases](https://github.com/username/securevault-cli/releases) and
add to your PATH.

### From Source

```bash
# Clone the repository
git clone https://github.com/username/securevault-cli.git
cd securevault-cli

# Build the binary
go build -o securevault cmd/securevault/main.go

# Move to your PATH
sudo mv securevault /usr/local/bin/
```

### Verification

```bash
securevault --version
# Expected output: securevault version 1.0.0
```

### Shell Completion (Optional)

```bash
# Bash
securevault completion bash > /etc/bash_completion.d/securevault

# Zsh
securevault completion zsh > "${fpath[1]}/_securevault"

# Fish
securevault completion fish > ~/.config/fish/completions/securevault.fish
```

## Quick Start

Create and retrieve your first secret in under 60 seconds:

```bash
# Initialize the vault
securevault init

# Set a master password when prompted
# Master password: ********

# Store a secret
securevault set my-api-key "sk-1234567890"

# Retrieve the secret
securevault get my-api-key
# Output: sk-1234567890
```

## Usage

### Basic Usage

```bash
# Store a secret
securevault set <key> <value>

# Retrieve a secret
securevault get <key>

# List all secret keys
securevault list

# Delete a secret
securevault delete <key>

# Update a secret
securevault update <key> <new-value>
```

### Common Scenarios

#### Store Database Credentials

```bash
securevault set db-password "my-secure-password"
securevault set db-username "admin"
securevault set db-host "localhost:5432"
```

#### Sync with Cloud Provider

```bash
# Configure AWS Secrets Manager backend
securevault config set backend aws

# Push local secrets to AWS
securevault sync push

# Pull secrets from AWS
securevault sync pull
```

#### Use in Scripts

```bash
#!/bin/bash
API_KEY=$(securevault get api-key)
curl -H "Authorization: Bearer $API_KEY" https://api.example.com/data
```

#### Rotate Credentials

```bash
# Rotate a credential and update it everywhere
securevault rotate my-api-key --update-cloud
```

For more examples, see the [examples/](./examples) directory.

## Configuration

### Environment Variables

| Variable                      | Description                          | Required | Default   |
| ----------------------------- | ------------------------------------ | -------- | --------- |
| `SECUREVAULT_MASTER_PASSWORD` | Master password for encryption       | No       | Prompted  |
| `SECUREVAULT_BACKEND`         | Backend provider (local, aws, azure) | No       | `local`   |
| `SECUREVAULT_CONFIG_PATH`     | Custom config file path              | No       | See below |

### Configuration File

SecureVault looks for configuration at `~/.securevault/config.yaml`:

```yaml
# Backend configuration
backend: local # Options: local, aws, azure, vault

# AWS Secrets Manager configuration
aws:
  region: us-east-1
  profile: default

# Azure Key Vault configuration
azure:
  vault_name: my-keyvault
  tenant_id: your-tenant-id

# HashiCorp Vault configuration
vault:
  address: https://vault.example.com:8200
  token: your-vault-token

# Local storage settings
local:
  database_path: ~/.securevault/secrets.db
  encryption: aes-256

# Audit logging
audit:
  enabled: true
  log_path: ~/.securevault/audit.log
```

## Examples

### Example 1: Manage Secrets for Multiple Environments

```bash
# Store production secrets
securevault set --env=prod db-password "prod-password"

# Store staging secrets
securevault set --env=staging db-password "staging-password"

# Retrieve environment-specific secret
securevault get --env=prod db-password
```

### Example 2: Import Secrets from File

```bash
# Import from JSON file
securevault import secrets.json

# secrets.json format:
# {
#   "api-key": "sk-123",
#   "db-password": "pass123"
# }
```

### Example 3: Integration with CI/CD

```bash
# In your CI/CD pipeline
export SECUREVAULT_MASTER_PASSWORD="${CI_MASTER_PASSWORD}"
export API_KEY=$(securevault get api-key)
./deploy.sh
```

## Contributing

We welcome contributions! Here's how to get started:

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and add tests
4. Run the test suite: `go test ./...`
5. Commit your changes: `git commit -m 'Add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/username/securevault-cli.git
cd securevault-cli

# Install dependencies
go mod download

# Run in development mode
go run cmd/securevault/main.go

# Run tests
go test ./...

# Run linter
golangci-lint run
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...

# Run specific test
go test -run TestVaultOperations ./pkg/vault
```

## License

This project is licensed under the MIT License - see the
[LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [AWS SDK for Go](https://github.com/aws/aws-sdk-go) - AWS integration
- [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go) - Azure
  integration
