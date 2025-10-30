# go-ctl CLI Documentation

The `go-ctl` CLI is a powerful command-line tool for generating Go projects with clean architecture, inspired by Spring Boot Initializr.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Configuration](#configuration)
- [Templates](#templates)
- [Examples](#examples)
- [Shell Completion](#shell-completion)
- [Project Analysis](#project-analysis)
- [Package Management](#package-management)
- [Advanced Usage](#advanced-usage)

## Installation

### From Source

```bash
git clone https://github.com/syst3mctl/go-ctl.git
cd go-ctl
make build-cli
sudo cp bin/go-ctl /usr/local/bin/
```

### Using Go Install (Coming Soon)

```bash
go install github.com/syst3mctl/go-ctl/cmd/cli@latest
```

### Binary Releases (Coming Soon)

Download pre-compiled binaries from the [releases page](https://github.com/syst3mctl/go-ctl/releases).

## Quick Start

### Generate a Simple API Project

```bash
go-ctl generate my-api --http=gin --database=postgres --driver=gorm
```

### Use a Built-in Template

```bash
go-ctl generate my-project --template=api
```

### Interactive Mode

```bash
go-ctl generate --interactive
```

## Commands

### `go-ctl generate`

Generate a new Go project with specified configuration.

**Usage:**
```bash
go-ctl generate [project-name] [flags]
```

**Flags:**
- `-g, --go-version string`: Go version (1.20, 1.21, 1.22, 1.23) (default "1.23")
- `-H, --http string`: HTTP framework (gin, echo, fiber, chi, net-http)
- `-d, --database strings`: Database types (postgres, mysql, sqlite, mongodb, redis)
- `-D, --driver strings`: Database drivers (gorm, sqlx, ent, mongo-driver, redis-client)
- `-f, --features strings`: Additional features (docker, makefile, air, jwt, cors, logging, testing)
- `-p, --packages strings`: Custom packages to include
- `-o, --output string`: Output directory (default ".")
- `-t, --template string`: Use built-in template (minimal, api, microservice, cli, worker)
- `-i, --interactive`: Interactive mode
- `--dry-run`: Show what would be generated without creating files

**Examples:**
```bash
# Basic API with PostgreSQL
go-ctl generate my-api --http=gin --database=postgres --driver=gorm

# Multiple databases and features
go-ctl generate my-service \
  --http=echo \
  --database=postgres,redis \
  --driver=gorm,redis-client \
  --features=docker,makefile,air,jwt

# Using a template
go-ctl generate my-microservice --template=microservice

# Interactive mode
go-ctl generate --interactive

# Dry run to preview
go-ctl generate test-project --template=api --dry-run
```

### `go-ctl template`

Manage and explore built-in project templates.

#### `go-ctl template list`

List available templates.

```bash
go-ctl template list [--detailed]
```

#### `go-ctl template show <template-name>`

Show detailed information about a template.

```bash
go-ctl template show api
```

#### `go-ctl template preview <template-name>`

Preview the project structure for a template.

```bash
go-ctl template preview microservice --name=my-project
```

#### `go-ctl template create <template-id>`

Create a custom template.

```bash
# Create from existing project
go-ctl template create my-template --from-project ./my-existing-project

# Create empty template
go-ctl template create my-template --name="My Template" --description="Custom template"
```

#### `go-ctl template delete <template-id>`

Delete a custom template.

```bash
go-ctl template delete my-template --force
```

#### `go-ctl template export <template-id> <output-file>`

Export a template to file.

```bash
go-ctl template export my-template my-template.yaml
```

#### `go-ctl template import <template-file>`

Import a template from file.

```bash
go-ctl template import my-template.yaml
```

### `go-ctl package`

Manage Go packages and dependencies.

#### `go-ctl package search <query>`

Search for Go packages.

```bash
go-ctl package search http
go-ctl package search --category=database postgres
```

#### `go-ctl package popular [category]`

Show popular packages by category.

```bash
go-ctl package popular
go-ctl package popular web
```

#### `go-ctl package info <import-path>`

Get information about a specific package.

```bash
go-ctl package info github.com/gin-gonic/gin
```

#### `go-ctl package validate <import-path>...`

Validate Go package import paths.

```bash
go-ctl package validate github.com/gin-gonic/gin gorm.io/gorm
```

### `go-ctl analyze`

Analyze Go project structure and provide insights.

```bash
# Analyze current directory
go-ctl analyze

# Analyze specific project
go-ctl analyze ./my-project

# Generate detailed report
go-ctl analyze ./my-project --detailed

# Export analysis to file
go-ctl analyze ./my-project --output=analysis.json
```

### `go-ctl config`

Manage go-ctl configuration files.

#### `go-ctl config init`

Create a new configuration file.

```bash
go-ctl config init [--global] [--force]
```

#### `go-ctl config show`

Display current configuration.

```bash
go-ctl config show [config-file]
```

#### `go-ctl config validate`

Validate a configuration file.

```bash
go-ctl config validate [config-file]
```

#### `go-ctl config set <key> <value>`

Set a configuration value.

```bash
go-ctl config set project.go_version 1.23
go-ctl config set cli.default_output ./projects
```

### `go-ctl completion`

Generate shell completion scripts.

```bash
go-ctl completion [bash|zsh|fish|powershell]
```

### `go-ctl version`

Show version information.

```bash
go-ctl version
```

## Configuration

### Configuration File Locations

go-ctl looks for configuration files in the following order:

1. File specified with `--config` flag
2. `.go-ctl.yaml` in current directory
3. `.go-ctl.yaml` in home directory

### Configuration Format

Configuration files use YAML format:

```yaml
# CLI settings
cli:
  default_output: "./projects"
  interactive_mode: false
  color_output: true
  auto_update: false

# Default project settings
project:
  goversion: "1.23"
  httppackage:
    id: "gin"
    name: "Gin"
    description: "High-performance HTTP web framework"
    importpath: "github.com/gin-gonic/gin"
  
  # Default databases
  databases:
    - database:
        id: "postgres"
        name: "PostgreSQL"
        description: "PostgreSQL database"
      driver:
        id: "gorm"
        name: "GORM"
        description: "The fantastic ORM library for Golang"
        importpath: "gorm.io/gorm"
  
  # Default features
  features:
    - id: "docker"
      name: "Docker"
      description: "Docker containerization"
    - id: "makefile"
      name: "Makefile"
      description: "Build automation"
  
  # Default packages
  custompackages:
    - "github.com/google/uuid"
    - "github.com/stretchr/testify/assert"
```

### Environment Variables

Override configuration using environment variables with `GO_CTL_` prefix:

```bash
export GO_CTL_CLI_DEFAULT_OUTPUT="./my-projects"
export GO_CTL_PROJECT_GOVERSION="1.22"
```

## Templates

### Built-in Templates

#### minimal
- **Description**: Minimal Go project with basic structure
- **Use Case**: Simple applications, learning projects
- **Includes**: Basic project structure, net/http

#### api
- **Description**: REST API with database and authentication
- **Use Case**: Web APIs, backend services
- **Includes**: Gin, PostgreSQL, GORM, Docker, JWT, CORS, logging, testing

#### microservice
- **Description**: Microservice with full observability
- **Use Case**: Distributed systems, microservice architecture
- **Includes**: Gin, PostgreSQL, Redis, Docker, hot-reload, comprehensive monitoring

#### cli
- **Description**: Command-line application
- **Use Case**: CLI tools, automation scripts
- **Includes**: Cobra, Viper, testing framework

#### worker
- **Description**: Background job processing service
- **Use Case**: Job queues, background processing
- **Includes**: Redis for jobs, PostgreSQL for persistence, Docker

### Custom Templates

You can create custom templates from existing projects or define them manually:

```bash
# Create template from existing project
go-ctl template create my-api-template --from-project ./existing-api --name="My API Template"

# List all templates (built-in and custom)
go-ctl template list

# Export template for sharing
go-ctl template export my-api-template shared-template.yaml

# Import shared template
go-ctl template import shared-template.yaml
```

### Template Customization

Override template settings using flags:

```bash
# Use API template but change HTTP framework
go-ctl generate my-api --template=api --http=echo

# Use microservice template but add custom packages
go-ctl generate my-service --template=microservice --packages=github.com/my/package
```

## Examples

### Web API with PostgreSQL

```bash
go-ctl generate blog-api \
  --http=gin \
  --database=postgres \
  --driver=gorm \
  --features=docker,makefile,cors,jwt,logging
```

### Microservice with Multiple Databases

```bash
go-ctl generate user-service \
  --template=microservice \
  --database=postgres,redis \
  --packages=go.opentelemetry.io/otel
```

### CLI Application

```bash
go-ctl generate my-tool \
  --template=cli \
  --packages=github.com/fatih/color,github.com/AlecAivazis/survey/v2
```

### Background Worker

```bash
go-ctl generate job-processor \
  --template=worker \
  --features=docker,logging,testing
```

### Custom Configuration

```bash
# Create custom config
go-ctl config init --global

# Edit ~/.go-ctl.yaml to set defaults
# Then generate projects using defaults
go-ctl generate my-project
```

## Shell Completion

### Bash

```bash
# Install globally
go-ctl completion bash | sudo tee /etc/bash_completion.d/go-ctl

# Or for current user
go-ctl completion bash > ~/.bash_completions/go-ctl
echo 'source ~/.bash_completions/go-ctl' >> ~/.bashrc
```

### Zsh

```bash
# Install to fpath
go-ctl completion zsh > "${fpath[1]}/_go-ctl"

# Or add to .zshrc
echo 'source <(go-ctl completion zsh)' >> ~/.zshrc
```

### Fish

```bash
go-ctl completion fish > ~/.config/fish/completions/go-ctl.fish
```

### PowerShell

```powershell
go-ctl completion powershell > go-ctl.ps1
# Then source it in your PowerShell profile
```

## Advanced Usage

### Configuration Profiles

Create multiple configuration files for different project types:

```bash
# API projects
go-ctl config init --global
# Edit ~/.go-ctl.yaml for API defaults

# Create microservice config
cp ~/.go-ctl.yaml ~/.go-ctl-microservice.yaml
# Edit for microservice defaults

# Use specific config
go-ctl generate my-service --config=~/.go-ctl-microservice.yaml
```

### Scripting and Automation

```bash
#!/bin/bash
# Script to generate multiple related services

services=("user-service" "order-service" "notification-service")

for service in "${services[@]}"; do
  go-ctl generate "$service" \
    --template=microservice \
    --output="./services" \
    --quiet
done
```

### CI/CD Integration

```yaml
# .github/workflows/generate-service.yml
name: Generate Service
on:
  workflow_dispatch:
    inputs:
      service_name:
        description: 'Service name'
        required: true

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.23
      
      - name: Install go-ctl
        run: |
          git clone https://github.com/syst3mctl/go-ctl.git
          cd go-ctl && make build-cli
          sudo cp bin/go-ctl /usr/local/bin/
      
      - name: Generate Service
        run: |
          go-ctl generate ${{ github.event.inputs.service_name }} \
            --template=microservice \
            --quiet
```

### Custom Templates

```bash
# Create template from existing project
go-ctl template create my-template --from-project ./existing-service

# Create empty template interactively
go-ctl template create my-template --interactive

# Use custom template
go-ctl generate my-project --template=my-template

# Export template for sharing
go-ctl template export my-template my-template.yaml

# Import shared template
go-ctl template import my-template.yaml
```

## Project Analysis

Analyze existing Go projects to understand their structure and get improvement suggestions:

```bash
# Analyze current project
go-ctl analyze

# Analyze with detailed report
go-ctl analyze ./my-project --detailed

# Focus on specific areas
go-ctl analyze ./my-project --focus=dependencies,security

# Export analysis results
go-ctl analyze ./my-project --output=analysis.json
```

### Analysis Features

- **Project Overview**: Basic information, type detection, statistics
- **Dependencies**: Dependency analysis with categorization and security info
- **Structure**: Project layout analysis and best practices validation
- **Patterns**: Architectural pattern detection (frameworks, databases, etc.)
- **Metrics**: Code quality, complexity, security, and maintainability scores
- **Suggestions**: Actionable improvement recommendations
- **Issues**: Problem detection with severity levels
- **Compatibility**: Go version and platform compatibility analysis

## Package Management

Discover and validate Go packages for your projects:

```bash
# Search for packages
go-ctl package search http
go-ctl package search --category=web gin

# Browse popular packages
go-ctl package popular
go-ctl package popular database --detailed

# Get package information
go-ctl package info github.com/gin-gonic/gin

# Validate package paths
go-ctl package validate github.com/gin-gonic/gin gorm.io/gorm
```

### Package Categories

- **Web**: HTTP frameworks and web utilities
- **Database**: Database drivers and ORMs  
- **Testing**: Testing frameworks and utilities
- **CLI**: Command-line interface tools
- **Logging**: Logging libraries
- **Auth**: Authentication and authorization
- **Validation**: Input validation and sanitization
- **Utils**: General utilities and helpers

## Troubleshooting

### Common Issues

#### "Template not found" Error

```bash
# List available templates
go-ctl template list

# Check template name spelling
go-ctl template show api
```

#### "Invalid configuration" Error

```bash
# Validate your config file
go-ctl config validate

# Show current config
go-ctl config show
```

#### Permission Issues

```bash
# Check output directory permissions
ls -la ./

# Use different output directory
go-ctl generate my-project --output=/tmp/my-project
```

### Debug Mode

```bash
# Enable verbose output
go-ctl generate my-project --verbose

# Dry run to see what would be generated
go-ctl generate my-project --dry-run
```

### Getting Help

```bash
# General help
go-ctl --help

# Command-specific help
go-ctl generate --help
go-ctl template --help
go-ctl package --help
go-ctl analyze --help
```

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for information on contributing to go-ctl CLI.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.