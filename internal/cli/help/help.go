package help

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// AddEnhancedHelp adds enhanced help information to commands
func AddEnhancedHelp(cmd *cobra.Command) {
	switch cmd.Name() {
	case "generate":
		setupGenerateHelp(cmd)
	case "template":
		setupTemplateHelp(cmd)
	case "package":
		setupPackageHelp(cmd)
	case "analyze":
		setupAnalyzeHelp(cmd)
	case "config":
		setupConfigHelp(cmd)
	}
}

// setupGenerateHelp adds comprehensive help for the generate command
func setupGenerateHelp(cmd *cobra.Command) {
	cmd.Long = color.HiCyanString(`Generate a new Go project with clean architecture and best practices.

`) + color.HiYellowString(`DESCRIPTION:
`) + `  The generate command creates a new Go project with your specified configuration.
  It supports multiple HTTP frameworks, databases, and additional features to help
  you get started quickly with a production-ready codebase.

` + color.HiYellowString(`EXAMPLES:
`) + `  ` + color.GreenString(`# Generate a minimal API project`) + `
  go-ctl generate my-api --http=gin --database=postgres --driver=gorm

  ` + color.GreenString(`# Generate with multiple databases`) + `
  go-ctl generate my-service --http=echo --database=postgres,redis --driver=gorm,redis-client

  ` + color.GreenString(`# Generate with additional features`) + `
  go-ctl generate my-app --http=fiber --database=mysql --features=docker,makefile,air,jwt

  ` + color.GreenString(`# Interactive mode for step-by-step configuration`) + `
  go-ctl generate --interactive

  ` + color.GreenString(`# Use a predefined template`) + `
  go-ctl generate my-microservice --template=microservice

  ` + color.GreenString(`# Dry run to see what would be generated`) + `
  go-ctl generate my-test --http=gin --dry-run

  ` + color.GreenString(`# Generate with smart template suggestions`) + `
  go-ctl generate --suggest --use-case=api

  ` + color.GreenString(`# Use configuration file`) + `
  go-ctl generate --config=./my-project-config.yaml

  ` + color.GreenString(`# Generate with custom packages`) + `
  go-ctl generate my-api --http=gin --packages=github.com/google/uuid,github.com/stretchr/testify

  ` + color.GreenString(`# Specify output directory`) + `
  go-ctl generate my-project --output=./projects --http=chi

` + color.HiYellowString(`SUPPORTED OPTIONS:
`) + `  ` + color.CyanString(`HTTP Frameworks:`) + ` gin, echo, fiber, chi, net-http
  ` + color.CyanString(`Databases:`) + ` postgres, mysql, sqlite, mongodb, redis, bigquery
  ` + color.CyanString(`Drivers:`) + ` gorm, sqlx, ent, mongo-driver, redis-client
  ` + color.CyanString(`Features:`) + ` docker, makefile, air, jwt, cors, logging, testing, env

` + color.HiYellowString(`TEMPLATES:
`) + `  ` + color.MagentaString(`minimal`) + `      - Basic Go project structure
  ` + color.MagentaString(`api`) + `          - REST API with database integration
  ` + color.MagentaString(`microservice`) + ` - Full microservice with gRPC and HTTP
  ` + color.MagentaString(`cli`) + `          - Command-line application
  ` + color.MagentaString(`worker`) + `       - Background worker with job processing
  ` + color.MagentaString(`web`) + `          - Web application with HTML templates
  ` + color.MagentaString(`grpc`) + `         - gRPC service with protocol buffers

` + color.HiYellowString(`OUTPUT FORMATS:
`) + `  Use ` + color.CyanString(`--output-format=json`) + ` for machine-readable output
  Use ` + color.CyanString(`--quiet`) + ` for minimal output (just the project path)

` + color.HiYellowString(`CONFIGURATION:
`) + `  Configuration files support YAML, JSON, and TOML formats.
  Example configuration file (.go-ctl.yaml):

  ` + color.WhiteString(`project:
    name: "my-api"
    go_version: "1.23"
    http_framework: "gin"
    databases:
      - type: "postgres"
        driver: "gorm"
    features:
      - "docker"
      - "makefile"
    custom_packages:
      - "github.com/google/uuid"`) + `

` + color.HiYellowString(`MORE INFO:`) + `
  Documentation: https://github.com/syst3mctl/go-ctl
  Examples: https://github.com/syst3mctl/go-ctl/tree/main/examples
  Templates: https://github.com/syst3mctl/go-ctl/tree/main/templates`
}

// setupTemplateHelp adds comprehensive help for the template command
func setupTemplateHelp(cmd *cobra.Command) {
	cmd.Long = color.HiCyanString(`Manage and explore project templates.

`) + color.HiYellowString(`DESCRIPTION:
`) + `  The template command helps you discover, explore, and get suggestions for
  project templates. Templates provide pre-configured project structures
  for common use cases and architectural patterns.

` + color.HiYellowString(`EXAMPLES:
`) + `  ` + color.GreenString(`# List all available templates`) + `
  go-ctl template list

  ` + color.GreenString(`# Show detailed information about a template`) + `
  go-ctl template show api

  ` + color.GreenString(`# Preview a template with your project name`) + `
  go-ctl template preview api --name=my-project

  ` + color.GreenString(`# Get smart template suggestions`) + `
  go-ctl template suggest

  ` + color.GreenString(`# Get suggestions for specific use case`) + `
  go-ctl template suggest --use-case=microservice

  ` + color.GreenString(`# Get suggestions with requirements`) + `
  go-ctl template suggest --requirements=database,docker,auth

  ` + color.GreenString(`# Analyze existing project for template suggestions`) + `
  go-ctl template suggest ./existing-project

` + color.HiYellowString(`AVAILABLE TEMPLATES:
`) + `  ` + color.MagentaString(`minimal`) + `      - Minimal Go project with basic structure
                Good for: Libraries, simple tools, learning
                Includes: go.mod, basic main.go, README

  ` + color.MagentaString(`api`) + `          - REST API with database integration
                Good for: REST APIs, web services, CRUD applications
                Includes: HTTP server, database layer, clean architecture

  ` + color.MagentaString(`microservice`) + ` - Full microservice with gRPC and HTTP
                Good for: Distributed systems, scalable services
                Includes: gRPC + HTTP servers, service discovery, monitoring

  ` + color.MagentaString(`cli`) + `          - Command-line application
                Good for: CLI tools, automation scripts, utilities
                Includes: Cobra framework, subcommands, configuration

  ` + color.MagentaString(`worker`) + `       - Background worker with job processing
                Good for: Background jobs, data processing, queues
                Includes: Job queue, worker pools, monitoring

  ` + color.MagentaString(`web`) + `          - Web application with HTML templates
                Good for: Server-rendered web apps, admin panels
                Includes: Template engine, static files, sessions

  ` + color.MagentaString(`grpc`) + `         - gRPC service with protocol buffers
                Good for: High-performance services, microservices
                Includes: Protocol buffers, gRPC server, client generation

` + color.HiYellowString(`TEMPLATE SELECTION:
`) + `  The suggest command uses intelligent algorithms to recommend templates
  based on your requirements:

  ` + color.CyanString(`Use Cases:`) + ` api, web, cli, microservice, worker, library, grpc
  ` + color.CyanString(`Requirements:`) + ` database, auth, docker, kubernetes, testing, monitoring

  Suggestions include confidence scores and explanations to help you choose.

` + color.HiYellowString(`MORE INFO:`) + `
  Template documentation: https://github.com/syst3mctl/go-ctl/tree/main/templates
  Custom templates: https://github.com/syst3mctl/go-ctl/blob/main/docs/custom-templates.md`
}

// setupPackageHelp adds comprehensive help for the package command
func setupPackageHelp(cmd *cobra.Command) {
	cmd.Long = color.HiCyanString(`Search, discover, and manage Go packages and dependencies.

`) + color.HiYellowString(`DESCRIPTION:
`) + `  The package command helps you discover Go packages, analyze dependencies,
  and manage upgrades. It integrates with pkg.go.dev for package discovery
  and provides intelligent upgrade recommendations.

` + color.HiYellowString(`EXAMPLES:
`) + `  ` + color.GreenString(`# Search for web framework packages`) + `
  go-ctl package search web

  ` + color.GreenString(`# Search with category filter`) + `
  go-ctl package search --category=web --limit=5

  ` + color.GreenString(`# Get popular packages by category`) + `
  go-ctl package popular database

  ` + color.GreenString(`# Get detailed package information`) + `
  go-ctl package info github.com/gin-gonic/gin

  ` + color.GreenString(`# Validate package compatibility`) + `
  go-ctl package validate github.com/gin-gonic/gin gorm.io/gorm

  ` + color.GreenString(`# Analyze current project dependencies`) + `
  go-ctl package upgrade

  ` + color.GreenString(`# Show only security-related upgrades`) + `
  go-ctl package upgrade --security-only

  ` + color.GreenString(`# Preview upgrade changes`) + `
  go-ctl package upgrade --dry-run

  ` + color.GreenString(`# Apply safe upgrades automatically`) + `
  go-ctl package upgrade --auto-apply

  ` + color.GreenString(`# Focus on specific types of issues`) + `
  go-ctl package upgrade --focus=security,outdated

` + color.HiYellowString(`PACKAGE CATEGORIES:
`) + `  ` + color.CyanString(`web`) + `          - Web frameworks (Gin, Echo, Fiber)
  ` + color.CyanString(`database`) + `     - Database drivers and ORMs (GORM, sqlx)
  ` + color.CyanString(`testing`) + `      - Testing frameworks (Testify, Ginkgo)
  ` + color.CyanString(`cli`) + `          - CLI libraries (Cobra, Viper)
  ` + color.CyanString(`logging`) + `      - Logging libraries (Logrus, Zap)
  ` + color.CyanString(`auth`) + `         - Authentication (JWT, OAuth2)
  ` + color.CyanString(`validation`) + `   - Input validation libraries
  ` + color.CyanString(`utils`) + `        - Utility libraries and helpers
  ` + color.CyanString(`crypto`) + `       - Cryptography and security
  ` + color.CyanString(`monitoring`) + `   - Metrics and observability

` + color.HiYellowString(`UPGRADE ANALYSIS:
`) + `  The upgrade command provides intelligent dependency analysis:

  ` + color.CyanString(`Security Analysis:`) + ` Identifies vulnerabilities with CVSS scores
  ` + color.CyanString(`Version Compatibility:`) + ` Checks for breaking changes
  ` + color.CyanString(`Alternative Suggestions:`) + ` Recommends better packages
  ` + color.CyanString(`Risk Assessment:`) + ` Rates upgrade risk (low/medium/high)

  Recommendations are prioritized by:
  • Security vulnerabilities (Priority 5 - Critical)
  • Major version updates available (Priority 4)
  • Deprecated packages (Priority 3)
  • Minor updates (Priority 2)
  • Patch updates (Priority 1)

` + color.HiYellowString(`OUTPUT FORMATS:
`) + `  Use ` + color.CyanString(`--output-format=json`) + ` for machine-readable output
  Perfect for CI/CD integration and automation scripts

` + color.HiYellowString(`MORE INFO:`) + `
  Package discovery: https://pkg.go.dev
  Security advisories: https://pkg.go.dev/vuln/
  Go modules guide: https://go.dev/blog/using-go-modules`
}

// setupAnalyzeHelp adds comprehensive help for the analyze command
func setupAnalyzeHelp(cmd *cobra.Command) {
	cmd.Long = color.HiCyanString(`Analyze Go projects for architecture, quality, and security insights.

`) + color.HiYellowString(`DESCRIPTION:
`) + `  The analyze command provides comprehensive analysis of Go projects,
  including architecture patterns, code quality metrics, security analysis,
  and actionable recommendations for improvement.

` + color.HiYellowString(`EXAMPLES:
`) + `  ` + color.GreenString(`# Analyze current directory`) + `
  go-ctl analyze

  ` + color.GreenString(`# Analyze specific project`) + `
  go-ctl analyze ./my-project

  ` + color.GreenString(`# Include upgrade analysis`) + `
  go-ctl analyze --upgrade-check

  ` + color.GreenString(`# Detailed analysis with verbose output`) + `
  go-ctl analyze --detailed

  ` + color.GreenString(`# Focus on specific areas`) + `
  go-ctl analyze --focus=security,dependencies

  ` + color.GreenString(`# Export analysis results`) + `
  go-ctl analyze --output=analysis-report.json --output-format=json

  ` + color.GreenString(`# Quick quality check`) + `
  go-ctl analyze --focus=quality --quiet

` + color.HiYellowString(`ANALYSIS AREAS:
`) + `  ` + color.CyanString(`Architecture:`) + `
    • Detects architectural patterns (Clean Architecture, MVC, etc.)
    • Analyzes layer separation and dependencies
    • Measures complexity and maintainability
    • Suggests structural improvements

  ` + color.CyanString(`Dependencies:`) + `
    • Maps dependency relationships
    • Identifies outdated packages
    • Detects security vulnerabilities
    • Suggests upgrade paths

  ` + color.CyanString(`Quality Metrics:`) + `
    • Code complexity analysis
    • Test coverage assessment
    • Code duplication detection
    • Documentation coverage

  ` + color.CyanString(`Security Analysis:`) + `
    • Vulnerability scanning
    • Dependency security audit
    • Security best practices check
    • Risk assessment and scoring

  ` + color.CyanString(`Performance:`) + `
    • Performance anti-patterns detection
    • Resource usage analysis
    • Optimization recommendations

` + color.HiYellowString(`FOCUS OPTIONS:
`) + `  ` + color.MagentaString(`dependencies`) + ` - Dependency analysis and upgrade suggestions
  ` + color.MagentaString(`security`) + `     - Security vulnerabilities and best practices
  ` + color.MagentaString(`quality`) + `      - Code quality metrics and improvements
  ` + color.MagentaString(`architecture`) + ` - Architectural patterns and structure
  ` + color.MagentaString(`performance`) + `  - Performance analysis and optimization
  ` + color.MagentaString(`testing`) + `      - Test coverage and testing practices

` + color.HiYellowString(`SCORING SYSTEM:
`) + `  Projects receive scores from 1-10 in each area:
  ` + color.GreenString(`9-10`) + ` - Excellent (following all best practices)
  ` + color.GreenString(`7-8`) + `  - Good (minor improvements needed)
  ` + color.YellowString(`5-6`) + `  - Fair (several areas for improvement)
  ` + color.RedString(`3-4`) + `  - Poor (significant issues present)
  ` + color.RedString(`1-2`) + `  - Critical (major problems requiring attention)

` + color.HiYellowString(`RECOMMENDATIONS:
`) + `  Analysis includes prioritized recommendations:
  • Priority 1-2: Nice-to-have improvements
  • Priority 3: Recommended improvements
  • Priority 4: Important issues to address
  • Priority 5: Critical issues requiring immediate attention

` + color.HiYellowString(`OUTPUT FORMATS:
`) + `  ` + color.CyanString(`text`) + ` - Human-readable report with colors and formatting
  ` + color.CyanString(`json`) + ` - Machine-readable format for CI/CD integration
  ` + color.CyanString(`yaml`) + ` - YAML format for configuration processing

` + color.HiYellowString(`MORE INFO:`) + `
  Best practices: https://github.com/syst3mctl/go-ctl/blob/main/docs/best-practices.md
  Security guide: https://github.com/syst3mctl/go-ctl/blob/main/docs/security.md`
}

// setupConfigHelp adds comprehensive help for the config command
func setupConfigHelp(cmd *cobra.Command) {
	cmd.Long = color.HiCyanString(`Manage go-ctl configuration files and settings.

`) + color.HiYellowString(`DESCRIPTION:
`) + `  The config command helps you create, validate, and manage configuration
  files for go-ctl. Configuration files allow you to save project templates
  and reuse them across multiple generations.

` + color.HiYellowString(`EXAMPLES:
`) + `  ` + color.GreenString(`# Initialize a new configuration file`) + `
  go-ctl config init

  ` + color.GreenString(`# Initialize with template`) + `
  go-ctl config init --template=api

  ` + color.GreenString(`# Validate configuration file`) + `
  go-ctl config validate

  ` + color.GreenString(`# Validate specific config file`) + `
  go-ctl config validate --config=./my-config.yaml

  ` + color.GreenString(`# Show current configuration`) + `
  go-ctl config show

  ` + color.GreenString(`# Show configuration with defaults`) + `
  go-ctl config show --include-defaults

  ` + color.GreenString(`# Set configuration value`) + `
  go-ctl config set project.go_version 1.23

  ` + color.GreenString(`# Get configuration value`) + `
  go-ctl config get project.http_framework

` + color.HiYellowString(`CONFIGURATION FILES:
`) + `  go-ctl searches for configuration files in this order:
  1. File specified by ` + color.CyanString(`--config`) + ` flag
  2. ` + color.CyanString(`.go-ctl.yaml`) + ` in current directory
  3. ` + color.CyanString(`.go-ctl.yaml`) + ` in home directory
  4. Global system configuration

  Supported formats: YAML, JSON, TOML

` + color.HiYellowString(`CONFIGURATION STRUCTURE:
`) + `  ` + color.WhiteString(`# .go-ctl.yaml
project:
  name: "my-api"                    # Project name template
  go_version: "1.23"                # Default Go version
  http_framework: "gin"             # Preferred HTTP framework
  databases:                        # Database configurations
    - type: "postgres"
      driver: "gorm"
  features:                         # Default features
    - "docker"
    - "makefile"
    - "air"
  custom_packages:                  # Custom packages to include
    - "github.com/google/uuid"
    - "github.com/stretchr/testify"

cli:                                # CLI-specific settings
  default_output: "./projects"      # Default output directory
  interactive_mode: false           # Enable interactive mode by default
  color_output: true                # Enable colored output
  auto_update: true                 # Auto-check for updates

templates:                          # Custom template configurations
  my_api_template:
    http_framework: "gin"
    databases: ["postgres"]
    features: ["docker", "testing"]`) + `

` + color.HiYellowString(`ENVIRONMENT VARIABLES:
`) + `  Configuration can be overridden with environment variables:
  ` + color.CyanString(`GO_CTL_PROJECT_NAME`) + `           - Project name
  ` + color.CyanString(`GO_CTL_PROJECT_GO_VERSION`) + `     - Go version
  ` + color.CyanString(`GO_CTL_PROJECT_HTTP_FRAMEWORK`) + ` - HTTP framework
  ` + color.CyanString(`GO_CTL_CLI_COLOR_OUTPUT`) + `       - Enable/disable colors

` + color.HiYellowString(`VALIDATION:
`) + `  Configuration validation checks:
  • Required fields are present
  • Values are within allowed options
  • Dependencies are compatible
  • Custom packages are valid Go module paths
  • Template references exist

` + color.HiYellowString(`MORE INFO:`) + `
  Configuration reference: https://github.com/syst3mctl/go-ctl/blob/main/docs/configuration.md
  Examples: https://github.com/syst3mctl/go-ctl/tree/main/examples/configs`
}

// GetUsageExamples returns usage examples for different scenarios
func GetUsageExamples() string {
	return color.HiCyanString(`🚀 go-ctl Usage Examples

`) + color.HiYellowString(`QUICK START:
`) + color.GreenString(`  # Interactive project generation
  go-ctl generate --interactive

  # Simple API project
  go-ctl generate my-api --http=gin --database=postgres

`) + color.HiYellowString(`COMMON SCENARIOS:
`) + color.GreenString(`  # REST API with authentication
  go-ctl generate auth-api --template=api --features=jwt,cors

  # Microservice with gRPC and HTTP
  go-ctl generate user-service --template=microservice --database=postgres

  # CLI tool
  go-ctl generate my-cli --template=cli --features=testing

  # Background worker
  go-ctl generate job-worker --template=worker --database=redis

`) + color.HiYellowString(`ADVANCED USAGE:
`) + color.GreenString(`  # Custom configuration
  go-ctl generate --config=./project-template.yaml

  # Multiple databases
  go-ctl generate multi-db --http=echo --database=postgres,redis,mongodb

  # With custom packages
  go-ctl generate custom-api --http=gin --packages=github.com/google/uuid

`) + color.HiYellowString(`PROJECT ANALYSIS:
`) + color.GreenString(`  # Full project analysis
  go-ctl analyze --detailed --upgrade-check

  # Security focus
  go-ctl analyze --focus=security --output-format=json

`) + color.HiYellowString(`PACKAGE MANAGEMENT:
`) + color.GreenString(`  # Find packages
  go-ctl package search logging --category=logging

  # Check for upgrades
  go-ctl package upgrade --security-only

`) + color.HiYellowString(`TEMPLATE DISCOVERY:
`) + color.GreenString(`  # Get suggestions
  go-ctl template suggest --use-case=api --requirements=auth,docker

  # Explore templates
  go-ctl template show microservice`) + `

` + color.HiYellowString(`For more examples and tutorials:`) + `
https://github.com/syst3mctl/go-ctl/tree/main/examples`
}

// GetTroubleshootingGuide returns common troubleshooting information
func GetTroubleshootingGuide() string {
	return color.HiCyanString(`🔧 Troubleshooting Guide

`) + color.HiYellowString(`COMMON ISSUES:

`) + color.RedString(`Problem: `) + `"command not found: go-ctl"
` + color.GreenString(`Solution: `) + `Ensure go-ctl is in your PATH or use full path to binary

` + color.RedString(`Problem: `) + `"failed to load project options"
` + color.GreenString(`Solution: `) + `Check that options.json exists and is valid JSON

` + color.RedString(`Problem: `) + `"incompatible database and driver combination"
` + color.GreenString(`Solution: `) + `Use compatible combinations (e.g., postgres+gorm, mongodb+mongo-driver)

` + color.RedString(`Problem: `) + `"project generation failed"
` + color.GreenString(`Solution: `) + `Check write permissions and disk space in output directory

` + color.RedString(`Problem: `) + `"package search not working"
` + color.GreenString(`Solution: `) + `Check internet connection for pkg.go.dev access

` + color.HiYellowString(`DEBUG MODE:
`) + `Use ` + color.CyanString(`--verbose`) + ` flag for detailed logging
Use ` + color.CyanString(`--dry-run`) + ` to preview changes without modification

` + color.HiYellowString(`GETTING HELP:
`) + `• GitHub Issues: https://github.com/syst3mctl/go-ctl/issues
• Documentation: https://github.com/syst3mctl/go-ctl
• Examples: https://github.com/syst3mctl/go-ctl/tree/main/examples`
}

// SetupManPage generates man page content
func SetupManPage(rootCmd *cobra.Command) string {
	var sb strings.Builder

	// Man page header
	sb.WriteString(fmt.Sprintf(".TH GO-CTL 1 \"%s\" \"go-ctl\" \"User Commands\"\n", "2024-01-01"))
	sb.WriteString(".SH NAME\n")
	sb.WriteString("go-ctl \\- Go project generator with clean architecture\n")
	sb.WriteString(".SH SYNOPSIS\n")
	sb.WriteString(".B go-ctl\n")
	sb.WriteString("[\\fIOPTIONS\\fR] \\fICOMMAND\\fR [\\fIARGS\\fR]\n")
	sb.WriteString(".SH DESCRIPTION\n")
	sb.WriteString("go-ctl is a powerful Go project generator inspired by Spring Boot Initializr.\n")
	sb.WriteString("It provides an intuitive interface for developers to select project options\n")
	sb.WriteString("and receive a downloadable, ready-to-code project skeleton with clean architecture.\n")

	// Commands section
	sb.WriteString(".SH COMMANDS\n")
	for _, cmd := range rootCmd.Commands() {
		if cmd.Hidden {
			continue
		}
		sb.WriteString(fmt.Sprintf(".TP\n.B %s\n%s\n", cmd.Name(), cmd.Short))
	}

	// Options section
	sb.WriteString(".SH OPTIONS\n")
	sb.WriteString(".TP\n.B \\-\\-config \\fIFILE\\fR\nSpecify configuration file path\n")
	sb.WriteString(".TP\n.B \\-\\-verbose, \\-v\nEnable verbose output\n")
	sb.WriteString(".TP\n.B \\-\\-quiet, \\-q\nEnable quiet mode\n")
	sb.WriteString(".TP\n.B \\-\\-help, \\-h\nShow help information\n")

	// Files section
	sb.WriteString(".SH FILES\n")
	sb.WriteString(".TP\n.B ~/.go-ctl.yaml\nUser configuration file\n")
	sb.WriteString(".TP\n.B .go-ctl.yaml\nProject-specific configuration file\n")

	// Examples section
	sb.WriteString(".SH EXAMPLES\n")
	sb.WriteString("Generate a basic API project:\n")
	sb.WriteString(".IP\n")
	sb.WriteString("go-ctl generate my-api --http=gin --database=postgres\n")
	sb.WriteString(".PP\n")
	sb.WriteString("Interactive project generation:\n")
	sb.WriteString(".IP\n")
	sb.WriteString("go-ctl generate --interactive\n")

	// See also section
	sb.WriteString(".SH SEE ALSO\n")
	sb.WriteString("Project documentation: https://github.com/syst3mctl/go-ctl\n")

	return sb.String()
}
