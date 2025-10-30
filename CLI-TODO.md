# CLI Implementation Roadmap for go-ctl

## Overview

This document outlines the implementation plan for adding a Command Line Interface (CLI) to the `go-ctl` project generator. The CLI will provide the same functionality as the web interface but through terminal commands, enabling automation, CI/CD integration, and developer-friendly scripting capabilities.

## 🎯 Project Goals

### Primary Objectives
- **Feature Parity**: CLI should offer all features available in the web interface
- **Automation Support**: Enable scripting and CI/CD pipeline integration
- **Developer Experience**: Intuitive commands with comprehensive help and validation
- **Cross-Platform**: Work seamlessly on Linux, macOS, and Windows
- **Configuration Files**: Support for project templates and reusable configurations

### Secondary Objectives
- **Interactive Mode**: Step-by-step project configuration wizard
- **Template Management**: Create, share, and manage custom project templates
- **Plugin System**: Extensible architecture for custom generators
- **Shell Completion**: Bash, Zsh, Fish, and PowerShell completion support

## 🏗️ Architecture Design

### CLI Framework
**Recommendation: Cobra CLI** (`github.com/spf13/cobra`)
- Industry standard for Go CLI applications
- Excellent subcommand support
- Built-in help generation
- Shell completion support
- Works well with Viper for configuration

### Project Structure
```
go-ctl/
├── cmd/
│   ├── cli/                    # CLI entry point
│   │   └── main.go
│   └── server/                 # Web server (existing)
│       └── main.go
├── internal/
│   ├── cli/                    # CLI-specific code
│   │   ├── commands/           # Command implementations
│   │   ├── config/            # CLI configuration
│   │   ├── interactive/       # Interactive wizard
│   │   └── templates/         # Template management
│   ├── generator/             # Shared generator (existing)
│   └── metadata/              # Shared metadata (existing)
└── pkg/
    └── cli/                   # Public CLI interfaces
```

## 📋 Implementation Tasks

### Phase 1: Core CLI Foundation

#### 1.1 Project Setup
- [ ] **Create CLI entry point** (`cmd/cli/main.go`)
- [ ] **Set up Cobra framework** with basic commands structure
- [ ] **Configure Viper** for configuration file support
- [ ] **Add CLI-specific go.mod dependencies**
- [ ] **Update Makefile** with CLI build targets

#### 1.2 Basic Commands Structure
- [ ] **Root command** (`go-ctl`)
  ```bash
  go-ctl --version
  go-ctl --help
  ```

- [ ] **Generate command** (`go-ctl generate`)
  ```bash
  go-ctl generate [project-name] [flags]
  ```

- [ ] **Config commands** (`go-ctl config`)
  ```bash
  go-ctl config init
  go-ctl config validate
  go-ctl config show
  ```

- [ ] **Template commands** (`go-ctl template`)
  ```bash
  go-ctl template list
  go-ctl template show [name]
  ```

#### 1.3 Core Flags and Options
- [ ] **Project configuration flags**
  ```bash
  --name string          Project name
  --go-version string    Go version (1.20, 1.21, 1.22, 1.23)
  --http string          HTTP framework (gin, echo, fiber, chi, net-http)
  --database strings     Database types (postgres, mysql, sqlite, mongodb, redis)
  --driver strings       Database drivers (gorm, sqlx, ent, mongo-driver, redis-client)
  --features strings     Additional features (docker, makefile, air, jwt, cors, logging)
  --packages strings     Custom packages to include
  --output string        Output directory (default: current directory)
  ```

- [ ] **Global flags**
  ```bash
  --config string        Config file path
  --verbose, -v         Verbose output
  --quiet, -q           Quiet mode
  --dry-run            Show what would be generated without creating files
  ```

### Phase 2: Command Implementations

#### 2.1 Generate Command Implementation
- [ ] **Basic project generation**
  ```bash
  go-ctl generate my-api --http=gin --database=postgres --driver=gorm
  ```

- [ ] **Configuration file support**
  ```bash
  go-ctl generate --config=./project-template.yaml
  ```

- [ ] **Interactive mode**
  ```bash
  go-ctl generate --interactive
  ```

- [ ] **Dry run functionality**
  ```bash
  go-ctl generate my-api --dry-run
  ```

#### 2.2 Configuration Management
- [ ] **Configuration file formats** (YAML, JSON, TOML)
  ```yaml
  # .go-ctl.yaml
  project:
    name: "my-api"
    go_version: "1.23"
    http_framework: "gin"
    databases:
      - type: "postgres"
        driver: "gorm"
    features:
      - "docker"
      - "makefile"
      - "air"
    custom_packages:
      - "github.com/google/uuid"
      - "github.com/stretchr/testify"
  ```

- [ ] **Config validation and schema**
- [ ] **Config file discovery** (project root, home directory, global)
- [ ] **Environment variable overrides**

#### 2.3 Template Management System
- [ ] **Built-in template listing**
  ```bash
  go-ctl template list
  # Output:
  # Built-in Templates:
  #   minimal     - Minimal Go project with basic structure
  #   api         - REST API with database and authentication
  #   microservice - Microservice with gRPC and HTTP endpoints
  #   cli         - Command-line application template
  ```

- [ ] **Template details and preview**
  ```bash
  go-ctl template show api
  go-ctl template preview api --name=my-project
  ```

- [ ] **Custom template support**
  ```bash
  go-ctl template create my-template --from-project ./existing-project
  go-ctl template install https://github.com/user/go-template.git
  ```

### Phase 3: Advanced Features

#### 3.1 Interactive Mode
- [ ] **Step-by-step wizard**
  ```
  $ go-ctl generate --interactive
  
  ✨ Welcome to go-ctl project generator!
  
  ? Project name: my-awesome-api
  ? Go version: (Use arrow keys)
    ❯ 1.23 (recommended)
      1.22
      1.21
      1.20
  
  ? HTTP framework: (Use arrow keys)
    ❯ Gin (high-performance)
      Echo (minimalist)
      Fiber (express-inspired)
      Chi (lightweight router)
      net/http (standard library)
  ```

- [ ] **Smart defaults and recommendations**
- [ ] **Conditional prompts** based on previous selections
- [ ] **Configuration summary and confirmation**

#### 3.2 Package Management Integration
- [ ] **Package search integration**
  ```bash
  go-ctl package search logging
  go-ctl package add github.com/rs/zerolog
  go-ctl package info github.com/gin-gonic/gin
  ```

- [ ] **Dependency analysis and suggestions**
- [ ] **Version compatibility checking**
- [ ] **Popular package recommendations**

#### 3.3 Project Analysis and Upgrades
- [ ] **Existing project analysis**
  ```bash
  go-ctl analyze ./existing-project
  go-ctl upgrade ./existing-project --to-version=1.23
  ```

- [ ] **Migration suggestions**
- [ ] **Dependency updates**
- [ ] **Best practice recommendations**

### Phase 4: Developer Experience Enhancements

#### 4.1 Output and Formatting
- [ ] **Colored output with progress indicators**
  ```
  ✅ Generating project structure...
  ✅ Creating go.mod file...
  ✅ Setting up HTTP server...
  ✅ Configuring database layer...
  ✅ Adding optional features...
  
  🎉 Project 'my-api' generated successfully!
  📁 Location: ./my-api
  🚀 Next steps:
      cd my-api
      make dev
  ```

- [ ] **JSON output for scripting**
  ```bash
  go-ctl generate my-api --output-format=json
  ```

- [ ] **Summary statistics and insights**

#### 4.2 Shell Completion
- [ ] **Bash completion**
  ```bash
  go-ctl completion bash > /etc/bash_completion.d/go-ctl
  ```

- [ ] **Zsh completion**
  ```bash
  go-ctl completion zsh > "${fpath[1]}/_go-ctl"
  ```

- [ ] **Fish and PowerShell support**

#### 4.3 Help and Documentation
- [ ] **Comprehensive help system**
- [ ] **Usage examples and tutorials**
- [ ] **Man page generation**
- [ ] **Online documentation links**

### Phase 5: Advanced Integration Features

#### 5.1 CI/CD Integration
- [ ] **GitHub Actions workflow generation**
  ```bash
  go-ctl generate my-api --ci=github-actions
  ```

- [ ] **GitLab CI, Jenkins, and other CI systems**
- [ ] **Deployment configurations** (Kubernetes, Docker Compose)
- [ ] **Infrastructure as Code** integration (Terraform, Pulumi)

#### 5.2 Plugin System
- [ ] **Plugin architecture design**
- [ ] **Plugin discovery and installation**
  ```bash
  go-ctl plugin install go-ctl-k8s
  go-ctl plugin list
  ```

- [ ] **Custom generator plugins**
- [ ] **Template plugin system**

#### 5.3 Cloud Provider Integration
- [ ] **AWS integration** (SAM, CDK, CloudFormation)
- [ ] **Google Cloud** (Cloud Functions, App Engine)
- [ ] **Azure integration** (Functions, Container Apps)
- [ ] **Database service configurations**

## 🔧 Technical Implementation Details

### Command Structure
```go
// cmd/cli/main.go
package main

import (
    "github.com/spf13/cobra"
    "github.com/syst3mctl/go-ctl/internal/cli/commands"
)

func main() {
    rootCmd := commands.NewRootCommand()
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### Configuration System
```go
// internal/cli/config/config.go
type CLIConfig struct {
    Project ProjectConfig `yaml:"project"`
    CLI     CLISettings   `yaml:"cli"`
}

type CLISettings struct {
    DefaultOutput     string `yaml:"default_output"`
    InteractiveMode   bool   `yaml:"interactive_mode"`
    ColorOutput      bool   `yaml:"color_output"`
    AutoUpdate       bool   `yaml:"auto_update"`
}
```

### Interactive Mode Implementation
```go
// internal/cli/interactive/wizard.go
type ProjectWizard struct {
    prompter survey.Prompter
    config   *metadata.ProjectConfig
}

func (w *ProjectWizard) Run() error {
    // Step-by-step configuration wizard
    return w.collectProjectDetails()
}
```

## 📦 Dependencies

### Required Packages
```go
// CLI framework
github.com/spf13/cobra v1.8.0
github.com/spf13/viper v1.18.0

// Interactive prompts
github.com/AlecAivazis/survey/v2 v2.3.7

// Output formatting
github.com/fatih/color v1.16.0
github.com/schollz/progressbar/v3 v3.14.1

// Validation
github.com/go-playground/validator/v10 v10.16.0

// Template engine (if needed)
github.com/Masterminds/sprig/v3 v3.2.3
```

## 🧪 Testing Strategy

### Unit Tests
- [ ] **Command logic testing**
- [ ] **Configuration parsing tests**
- [ ] **Flag validation tests**
- [ ] **Template generation tests**

### Integration Tests
- [ ] **End-to-end CLI workflow tests**
- [ ] **Configuration file integration**
- [ ] **Template system integration**
- [ ] **Cross-platform compatibility tests**

### CLI-Specific Tests
```go
func TestGenerateCommand(t *testing.T) {
    cmd := commands.NewGenerateCommand()
    cmd.SetArgs([]string{
        "test-project",
        "--http=gin",
        "--database=postgres",
        "--driver=gorm",
    })
    
    err := cmd.Execute()
    assert.NoError(t, err)
    
    // Verify project structure
    assert.DirExists(t, "test-project")
    assert.FileExists(t, "test-project/go.mod")
    assert.FileExists(t, "test-project/cmd/test-project/main.go")
}
```

## 📖 Documentation Plan

### User Documentation
- [ ] **Installation guide**
- [ ] **Quick start tutorial**
- [ ] **Command reference**
- [ ] **Configuration file documentation**
- [ ] **Template creation guide**
- [ ] **Examples and use cases**

### Developer Documentation
- [ ] **Architecture overview**
- [ ] **Plugin development guide**
- [ ] **Contributing guidelines**
- [ ] **API documentation**

## 🚀 Release Strategy

### Version 1.0.0 (CLI Beta)
- ✅ Core generate command
- ✅ Basic configuration support
- ✅ Template system
- ✅ Interactive mode

### Version 1.1.0 (Enhanced Features)
- ✅ Package management integration
- ✅ Shell completion
- ✅ Advanced formatting

### Version 1.2.0 (Enterprise Features)
- ✅ Plugin system
- ✅ CI/CD integration
- ✅ Cloud provider support

## 📋 Development Checklist

### Phase 1 Setup
- [ ] Create CLI project structure
- [ ] Set up Cobra and Viper
- [ ] Implement basic root command
- [ ] Add version and help commands
- [ ] Update build system (Makefile)

### Phase 2 Core Commands
- [ ] Implement generate command
- [ ] Add configuration file support
- [ ] Create template management
- [ ] Add validation and error handling

### Phase 3 User Experience
- [ ] Implement interactive mode
- [ ] Add progress indicators
- [ ] Create shell completion
- [ ] Improve error messages and help

### Phase 4 Advanced Features
- [ ] Package search integration
- [ ] Plugin architecture
- [ ] CI/CD template generation
- [ ] Cloud integration

### Phase 5 Polish and Release
- [ ] Comprehensive testing
- [ ] Documentation completion
- [ ] Performance optimization
- [ ] Release preparation

## 🔍 Success Metrics

- **Adoption**: CLI download/usage metrics
- **User Experience**: Time to first successful project generation
- **Feature Coverage**: Parity with web interface features
- **Community**: User-contributed templates and plugins
- **Integration**: Usage in CI/CD pipelines and automation

---

**Next Steps**: Start with Phase 1 implementation, focusing on core CLI structure and basic project generation functionality.