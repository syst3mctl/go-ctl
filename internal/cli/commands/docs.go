package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/syst3mctl/go-ctl/internal/cli/help"
)

// NewDocsCommand creates the docs command for generating documentation
func NewDocsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation for go-ctl",
		Long: `Generate various forms of documentation for go-ctl including man pages,
markdown documentation, and usage examples.

This command helps system administrators and users create documentation
for installation and reference purposes.`,
	}

	// Man page generation
	manCmd := &cobra.Command{
		Use:   "man [output-directory]",
		Short: "Generate man page",
		Long: `Generate a man page for go-ctl that can be installed on Unix-like systems.

Examples:
  # Generate man page to current directory
  go-ctl docs man

  # Generate man page to specific directory
  go-ctl docs man /usr/local/man/man1

  # Generate and install (requires sudo)
  sudo go-ctl docs man /usr/share/man/man1`,
		Args: cobra.MaximumNArgs(1),
		RunE: runGenerateManPage,
	}

	// Usage examples generation
	examplesCmd := &cobra.Command{
		Use:   "examples [output-file]",
		Short: "Generate usage examples",
		Long: `Generate comprehensive usage examples for go-ctl commands.

Examples:
  # Print examples to stdout
  go-ctl docs examples

  # Save examples to file
  go-ctl docs examples examples.md

  # Generate examples in different format
  go-ctl docs examples --format=text`,
		Args: cobra.MaximumNArgs(1),
		RunE: runGenerateExamples,
	}

	// Troubleshooting guide generation
	troubleshootCmd := &cobra.Command{
		Use:   "troubleshoot [output-file]",
		Short: "Generate troubleshooting guide",
		Long: `Generate a troubleshooting guide for common go-ctl issues.

Examples:
  # Print troubleshooting guide to stdout
  go-ctl docs troubleshoot

  # Save troubleshooting guide to file
  go-ctl docs troubleshoot troubleshoot.md`,
		Args: cobra.MaximumNArgs(1),
		RunE: runGenerateTroubleshooting,
	}

	// Add format flag for examples
	examplesCmd.Flags().StringP("format", "f", "markdown", "Output format (markdown, text)")

	// Add subcommands
	cmd.AddCommand(manCmd)
	cmd.AddCommand(examplesCmd)
	cmd.AddCommand(troubleshootCmd)

	return cmd
}

// runGenerateManPage generates a man page for go-ctl
func runGenerateManPage(cmd *cobra.Command, args []string) error {
	// Determine output directory
	outputDir := "."
	if len(args) > 0 {
		outputDir = args[0]
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate man page content
	rootCmd := cmd.Root()
	manContent := help.SetupManPage(rootCmd)

	// Write to file
	outputFile := filepath.Join(outputDir, "go-ctl.1")
	if err := os.WriteFile(outputFile, []byte(manContent), 0644); err != nil {
		return fmt.Errorf("failed to write man page: %w", err)
	}

	printSuccess("Man page generated: %s", outputFile)
	printInfo("To install the man page:")
	printInfo("  sudo cp %s /usr/share/man/man1/", outputFile)
	printInfo("  sudo mandb")
	printInfo("")
	printInfo("Then you can use: man go-ctl")

	return nil
}

// runGenerateExamples generates usage examples
func runGenerateExamples(cmd *cobra.Command, args []string) error {
	format, _ := cmd.Flags().GetString("format")

	// Generate examples content
	var content string
	switch format {
	case "markdown":
		content = generateMarkdownExamples()
	case "text":
		content = help.GetUsageExamples()
	default:
		return fmt.Errorf("unsupported format: %s (supported: markdown, text)", format)
	}

	// Output to file or stdout
	if len(args) > 0 {
		outputFile := args[0]
		if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write examples file: %w", err)
		}
		printSuccess("Examples generated: %s", outputFile)
	} else {
		fmt.Print(content)
	}

	return nil
}

// runGenerateTroubleshooting generates troubleshooting guide
func runGenerateTroubleshooting(cmd *cobra.Command, args []string) error {
	content := help.GetTroubleshootingGuide()

	// Output to file or stdout
	if len(args) > 0 {
		outputFile := args[0]

		// Convert to markdown format if outputting to file
		markdownContent := generateMarkdownTroubleshooting()

		if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
			return fmt.Errorf("failed to write troubleshooting file: %w", err)
		}
		printSuccess("Troubleshooting guide generated: %s", outputFile)
	} else {
		fmt.Print(content)
	}

	return nil
}

// generateMarkdownExamples generates examples in markdown format
func generateMarkdownExamples() string {
	return `# go-ctl Usage Examples

## Quick Start

### Interactive Project Generation
The easiest way to get started is with interactive mode:
` + "```bash" + `
go-ctl generate --interactive
` + "```" + `

### Simple API Project
Generate a basic API project with Gin and PostgreSQL:
` + "```bash" + `
go-ctl generate my-api --http=gin --database=postgres --driver=gorm
` + "```" + `

## Common Project Types

### REST API with Authentication
` + "```bash" + `
go-ctl generate auth-api --template=api --features=jwt,cors --database=postgres
` + "```" + `

### Microservice with gRPC and HTTP
` + "```bash" + `
go-ctl generate user-service --template=microservice --database=postgres --features=docker
` + "```" + `

### Command-Line Tool
` + "```bash" + `
go-ctl generate my-cli --template=cli --features=testing
` + "```" + `

### Background Worker
` + "```bash" + `
go-ctl generate job-worker --template=worker --database=redis --features=docker,monitoring
` + "```" + `

### Web Application
` + "```bash" + `
go-ctl generate my-webapp --template=web --database=postgres --features=sessions,csrf
` + "```" + `

## Advanced Configuration

### Using Configuration Files
Create a configuration file for reusable project templates:
` + "```yaml" + `
# .go-ctl.yaml
project:
  go_version: "1.23"
  http_framework: "gin"
  databases:
    - type: "postgres"
      driver: "gorm"
  features:
    - "docker"
    - "makefile"
    - "testing"
  custom_packages:
    - "github.com/google/uuid"
    - "github.com/stretchr/testify"
` + "```" + `

Then generate projects using the configuration:
` + "```bash" + `
go-ctl generate my-project --config=./.go-ctl.yaml
` + "```" + `

### Multiple Databases
` + "```bash" + `
go-ctl generate multi-db --http=echo \
  --database=postgres,redis,mongodb \
  --driver=gorm,redis-client,mongo-driver
` + "```" + `

### Custom Packages
` + "```bash" + `
go-ctl generate custom-api --http=gin \
  --packages=github.com/google/uuid,github.com/rs/zerolog \
  --features=logging,validation
` + "```" + `

## Project Analysis

### Full Project Analysis
` + "```bash" + `
go-ctl analyze --detailed --upgrade-check
` + "```" + `

### Security-Focused Analysis
` + "```bash" + `
go-ctl analyze --focus=security --output-format=json > security-report.json
` + "```" + `

### Dependency Analysis
` + "```bash" + `
go-ctl analyze --focus=dependencies --upgrade-check
` + "```" + `

## Package Management

### Search for Packages
` + "```bash" + `
go-ctl package search web --category=web --limit=5
` + "```" + `

### Popular Packages by Category
` + "```bash" + `
go-ctl package popular database
` + "```" + `

### Package Information
` + "```bash" + `
go-ctl package info github.com/gin-gonic/gin
` + "```" + `

### Upgrade Analysis
` + "```bash" + `
# Check for all upgrades
go-ctl package upgrade

# Security updates only
go-ctl package upgrade --security-only

# Preview changes
go-ctl package upgrade --dry-run

# Auto-apply safe updates
go-ctl package upgrade --auto-apply
` + "```" + `

## Template Discovery

### Get Template Suggestions
` + "```bash" + `
# Interactive questionnaire
go-ctl template suggest

# Specific use case
go-ctl template suggest --use-case=api

# With requirements
go-ctl template suggest --requirements=database,docker,auth

# Analyze existing project
go-ctl template suggest ./existing-project
` + "```" + `

### Explore Templates
` + "```bash" + `
# List all templates
go-ctl template list

# Show template details
go-ctl template show microservice

# Preview template structure
go-ctl template preview api --name=my-project
` + "```" + `

## Output Customization

### JSON Output for Scripting
` + "```bash" + `
go-ctl generate my-api --http=gin --output-format=json > generation-result.json
go-ctl analyze --output-format=json > analysis-result.json
` + "```" + `

### Quiet Mode
` + "```bash" + `
# Only output the project path
go-ctl generate my-api --http=gin --quiet

# Only output the analysis score
go-ctl analyze --quiet
` + "```" + `

### Verbose Mode
` + "```bash" + `
go-ctl generate my-api --http=gin --verbose
go-ctl analyze --verbose --detailed
` + "```" + `

## CI/CD Integration

### GitHub Actions Example
` + "```yaml" + `
name: Generate and Test Project
on: [push]
jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.23'
      - name: Install go-ctl
        run: go install github.com/syst3mctl/go-ctl/cmd/cli@latest
      - name: Generate project
        run: go-ctl generate test-project --http=gin --database=postgres --quiet
      - name: Test generated project
        run: |
          cd test-project
          go mod tidy
          go test ./...
` + "```" + `

### Automated Analysis
` + "```bash" + `
# In CI/CD pipeline
go-ctl analyze --focus=security --output-format=json | jq '.security.score < 8' && exit 1
` + "```" + `

## Development Workflow

### Development with Hot Reload
When generating projects with the Air feature:
` + "```bash" + `
go-ctl generate my-api --http=gin --features=air --database=postgres
cd my-api
make dev  # Uses Air for hot reload
` + "```" + `

### Docker Development
` + "```bash" + `
go-ctl generate my-service --template=microservice --features=docker
cd my-service
docker-compose up --build
` + "```" + `

### Testing Setup
` + "```bash" + `
go-ctl generate my-project --features=testing --database=postgres
cd my-project
make test
make coverage
` + "```" + `

For more examples and detailed tutorials, visit: https://github.com/syst3mctl/go-ctl/tree/main/examples
`
}

// generateMarkdownTroubleshooting generates troubleshooting guide in markdown format
func generateMarkdownTroubleshooting() string {
	return `# go-ctl Troubleshooting Guide

## Common Issues

### Installation Issues

**Problem:** "command not found: go-ctl"
**Solution:**
- Ensure go-ctl is installed: ` + "`go install github.com/syst3mctl/go-ctl/cmd/cli@latest`" + `
- Check that your GOPATH/bin is in your PATH
- Use the full path to the binary if needed

**Problem:** "permission denied" when running go-ctl
**Solution:**
- Ensure the binary has execute permissions: ` + "`chmod +x $(which go-ctl)`" + `
- Check file ownership and permissions

### Configuration Issues

**Problem:** "failed to load project options"
**Solution:**
- Ensure options.json exists in the go-ctl installation directory
- Check that the options.json file is valid JSON
- Reinstall go-ctl if the file is corrupted

**Problem:** "invalid configuration file"
**Solution:**
- Validate your YAML/JSON configuration syntax
- Use ` + "`go-ctl config validate`" + ` to check your configuration
- Check indentation in YAML files (use spaces, not tabs)

### Generation Issues

**Problem:** "incompatible database and driver combination"
**Solution:**
Use compatible combinations:
- PostgreSQL: ` + "`--database=postgres --driver=gorm`" + ` or ` + "`--driver=sqlx`" + `
- MySQL: ` + "`--database=mysql --driver=gorm`" + ` or ` + "`--driver=sqlx`" + `
- MongoDB: ` + "`--database=mongodb --driver=mongo-driver`" + `
- Redis: ` + "`--database=redis --driver=redis-client`" + `

**Problem:** "project generation failed"
**Solution:**
- Check write permissions in the output directory
- Ensure sufficient disk space
- Verify the output directory exists or can be created
- Use ` + "`--verbose`" + ` flag to see detailed error messages

**Problem:** "template not found"
**Solution:**
- Use ` + "`go-ctl template list`" + ` to see available templates
- Check template name spelling
- Ensure you're using the correct template ID (e.g., "api", not "API")

### Network Issues

**Problem:** "package search not working"
**Solution:**
- Check internet connection
- Verify access to pkg.go.dev (not blocked by firewall/proxy)
- Try again later if pkg.go.dev is temporarily unavailable
- Use ` + "`--verbose`" + ` to see network errors

**Problem:** "timeout errors during package operations"
**Solution:**
- Check network connectivity and proxy settings
- Increase timeout with environment variables if needed
- Use local package cache when available

### Analysis Issues

**Problem:** "analysis failed on large projects"
**Solution:**
- Use ` + "`--focus`" + ` flag to analyze specific areas
- Break analysis into smaller parts
- Ensure sufficient memory and processing power
- Check for circular dependencies that might cause infinite loops

**Problem:** "outdated vulnerability data"
**Solution:**
- Update go-ctl to the latest version
- Clear cache if available
- Check that security databases are accessible

## Debug Mode

Use the ` + "`--verbose`" + ` flag to enable detailed logging:

` + "```bash" + `
go-ctl generate my-project --http=gin --verbose
go-ctl analyze --verbose
go-ctl package search web --verbose
` + "```" + `

## Dry Run Mode

Use ` + "`--dry-run`" + ` to preview changes without making them:

` + "```bash" + `
go-ctl generate my-project --http=gin --dry-run
go-ctl package upgrade --dry-run
` + "```" + `

## Configuration Validation

Before using configuration files, validate them:

` + "```bash" + `
go-ctl config validate
go-ctl config validate --config=./my-config.yaml
` + "```" + `

## System Requirements

Ensure your system meets the requirements:
- Go 1.20 or later
- Sufficient disk space (at least 100MB free)
- Network access for package operations
- Write permissions in target directories

## Environment Variables

Set these environment variables for debugging:

` + "```bash" + `
export GO_CTL_VERBOSE=true          # Enable verbose logging
export GO_CTL_CONFIG_PATH=/path     # Override config path
export GO_CTL_CACHE_DIR=/path       # Override cache directory
` + "```" + `

## Getting Help

If you continue to experience issues:

1. **Check the documentation:** https://github.com/syst3mctl/go-ctl
2. **Search existing issues:** https://github.com/syst3mctl/go-ctl/issues
3. **Create a new issue:** Include the following information:
   - go-ctl version (` + "`go-ctl version`" + `)
   - Operating system and version
   - Go version (` + "`go version`" + `)
   - Complete command that failed
   - Full error output with ` + "`--verbose`" + ` flag
   - Configuration file contents (if applicable)

## Reporting Bugs

When reporting bugs, please include:
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- System information
- Log output with ` + "`--verbose`" + ` enabled

## Performance Issues

If go-ctl is running slowly:
- Use ` + "`--focus`" + ` flags to limit analysis scope
- Check available system resources
- Consider using ` + "`--quiet`" + ` mode for scripting
- Profile memory usage if generation fails

## Security Issues

For security-related issues:
- Do not open public GitHub issues
- Email security concerns to the maintainers
- Include proof-of-concept if applicable
- Allow reasonable time for response before disclosure
`
}
