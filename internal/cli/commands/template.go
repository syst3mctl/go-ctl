package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/syst3mctl/go-ctl/internal/cli/help"
	"github.com/syst3mctl/go-ctl/internal/cli/output"
	"github.com/syst3mctl/go-ctl/internal/cli/templates"
	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// NewTemplateCommand creates the template command
func NewTemplateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Manage project templates with smart suggestions",
		Long: `Manage and explore built-in project templates with intelligent suggestions.

Templates provide pre-configured project setups for common use cases
like APIs, microservices, CLI applications, and more. The smart suggestion
system helps you find the perfect template based on your requirements.`,
	}

	// Add enhanced help
	help.AddEnhancedHelp(cmd)

	cmd.AddCommand(newTemplateListCommand())
	cmd.AddCommand(newTemplateShowCommand())
	cmd.AddCommand(newTemplatePreviewCommand())
	cmd.AddCommand(newTemplateCreateCommand())
	cmd.AddCommand(newTemplateDeleteCommand())
	cmd.AddCommand(newTemplateExportCommand())
	cmd.AddCommand(newTemplateImportCommand())
	cmd.AddCommand(newTemplateSuggestCommand())

	return cmd
}

// newTemplateListCommand creates the template list subcommand
func newTemplateListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates with enhanced formatting",
		Long: `List all available built-in project templates with detailed information.

Each template includes a pre-configured set of options for specific project
types and use cases. Use --output-format=json for machine-readable output.

Examples:
  go-ctl template list                    # Standard output
  go-ctl template list --output-format=json  # JSON output for scripts
  go-ctl template list --detailed        # Show additional template info`,
		RunE: runTemplateList,
	}

	cmd.Flags().Bool("detailed", false, "show detailed template information")
	cmd.Flags().StringP("output-format", "f", "text", "output format (text, json)")

	return cmd
}

// newTemplateShowCommand creates the template show subcommand
func newTemplateShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <template-name>",
		Short: "Show template details",
		Long: `Show detailed information about a specific template
including its configuration and included features.

Examples:
  go-ctl template show api
  go-ctl template show microservice
  go-ctl template show minimal`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplateShow,
	}

	return cmd
}

// newTemplatePreviewCommand creates the template preview subcommand
func newTemplatePreviewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview <template-name>",
		Short: "Preview template project structure",
		Long: `Preview the project structure that would be generated
using a specific template.

Examples:
  go-ctl template preview api --name=my-api
  go-ctl template preview microservice --name=my-service`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplatePreview,
	}

	cmd.Flags().String("name", "example-project", "Project name for preview")

	return cmd
}

// newTemplateCreateCommand creates the template create subcommand
func newTemplateCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <template-id>",
		Short: "Create a custom template",
		Long: `Create a custom template from an existing project or interactively.

Examples:
  # Create from existing project
  go-ctl template create my-template --from-project ./my-existing-project --name="My Template"

  # Create empty template (interactive)
  go-ctl template create my-template --interactive

  # Create with specific configuration
  go-ctl template create my-template --name="My Template" --description="Custom template"`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplateCreate,
	}

	cmd.Flags().String("from-project", "", "Create template from existing project directory")
	cmd.Flags().String("name", "", "Template display name")
	cmd.Flags().String("description", "", "Template description")
	cmd.Flags().String("author", "", "Template author")
	cmd.Flags().StringSlice("tags", []string{}, "Template tags")
	cmd.Flags().BoolP("interactive", "i", false, "Interactive template creation")
	cmd.Flags().Bool("force", false, "Overwrite existing template")

	return cmd
}

// newTemplateDeleteCommand creates the template delete subcommand
func newTemplateDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <template-id>",
		Short: "Delete a custom template",
		Long: `Delete a custom template by ID.

Examples:
  go-ctl template delete my-template
  go-ctl template delete my-template --force`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplateDelete,
	}

	cmd.Flags().Bool("force", false, "Skip confirmation prompt")

	return cmd
}

// newTemplateExportCommand creates the template export subcommand
func newTemplateExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <template-id> <output-file>",
		Short: "Export a template to file",
		Long: `Export a template to a YAML or JSON file for sharing.

Examples:
  go-ctl template export my-template my-template.yaml
  go-ctl template export my-template my-template.json`,
		Args: cobra.ExactArgs(2),
		RunE: runTemplateExport,
	}

	return cmd
}

// newTemplateImportCommand creates the template import subcommand
func newTemplateImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <template-file>",
		Short: "Import a template from file",
		Long: `Import a template from a YAML or JSON file.

Examples:
  go-ctl template import my-template.yaml
  go-ctl template import my-template.json
  go-ctl template import https://example.com/template.yaml`,
		Args: cobra.ExactArgs(1),
		RunE: runTemplateImport,
	}

	cmd.Flags().Bool("force", false, "Overwrite existing template")

	return cmd
}

// Note: Template definitions are imported from generate.go to avoid duplication

// runTemplateList lists available templates
func runTemplateList(cmd *cobra.Command, args []string) error {
	detailed, _ := cmd.Flags().GetBool("detailed")
	outputFormat, _ := cmd.Flags().GetString("output-format")
	builtinTemplates := getBuiltinTemplatesForTemplateCmd()

	// Get custom templates
	tm := templates.NewTemplateManager()
	customTemplates, err := tm.ListTemplates()
	if err != nil && !isQuiet() {
		printWarning("Failed to load custom templates: %v", err)
	}

	if len(builtinTemplates) == 0 && len(customTemplates) == 0 {
		if !isQuiet() {
			printInfo("No templates available")
		}
		return nil
	}

	// Handle JSON output format
	if outputFormat == "json" {
		return outputTemplateListJSON(builtinTemplates, customTemplates, detailed)
	}

	// Create formatter for text output
	formatter := output.NewFormatter(output.FormatText, os.Stdout)
	formatter.SetOptions(isVerbose(), isQuiet(), isNoColor())

	// Display built-in templates
	if len(builtinTemplates) > 0 {
		if !isQuiet() {
			formatter.PrintInfo("Built-in Templates")
			fmt.Printf("%s\n", strings.Repeat("=", 50))
		}

		for _, template := range builtinTemplates {
			if detailed {
				if isNoColor() {
					fmt.Printf("\n● %s\n", template.Name)
				} else {
					fmt.Printf("\n%s %s\n", color.HiGreenString("●"), color.HiWhiteString(template.Name))
				}

				printTemplateDetails(template, formatter)
			} else {
				if isNoColor() {
					fmt.Printf("  ● %-12s - %s\n", template.ID, template.Description)
				} else {
					fmt.Printf("  %s %-12s - %s\n",
						color.HiGreenString("●"),
						color.HiWhiteString(template.ID),
						template.Description)
				}
			}
		}
	}

	// Display custom templates
	if len(customTemplates) > 0 {
		if len(builtinTemplates) > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", color.HiCyanString("Custom Templates:"))
		fmt.Printf("%s\n", strings.Repeat("=", 50))

		for _, template := range customTemplates {
			if detailed {
				fmt.Printf("\n%s %s\n", color.HiMagentaString("●"), color.HiWhiteString(template.Name))
				fmt.Printf("  %s: %s\n", color.CyanString("ID"), template.ID)
				fmt.Printf("  %s: %s\n", color.CyanString("Description"), template.Description)
				if template.Author != "" {
					fmt.Printf("  %s: %s\n", color.CyanString("Author"), template.Author)
				}
				if len(template.Tags) > 0 {
					fmt.Printf("  %s: %s\n", color.CyanString("Tags"), strings.Join(template.Tags, ", "))
				}
				fmt.Printf("  %s: %s\n", color.CyanString("Go Version"), template.Config.GoVersion)
				if template.Config.HttpPackage.ID != "" {
					fmt.Printf("  %s: %s\n", color.CyanString("HTTP Framework"), template.Config.HttpPackage.Name)
				}
				fmt.Printf("  %s: %s\n", color.CyanString("Created"), template.CreatedAt.Format("2006-01-02"))
			} else {
				fmt.Printf("  %s %-12s - %s\n",
					color.HiMagentaString("●"),
					color.HiWhiteString(template.ID),
					template.Description)
			}
		}
	}

	if !detailed {
		fmt.Printf("\nUse '%s' to see detailed information about a template\n",
			color.CyanString("go-ctl template show <template-name>"))
	}

	return nil
}

// outputTemplateListJSON outputs template list in JSON format
func outputTemplateListJSON(builtinTemplates []BuiltinTemplate, customTemplates []*templates.CustomTemplate, detailed bool) error {
	type JSONTemplateInfo struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Type        string   `json:"type"`
		Tags        []string `json:"tags,omitempty"`
		GoVersion   string   `json:"go_version,omitempty"`
		HTTP        string   `json:"http_framework,omitempty"`
		Databases   []string `json:"databases,omitempty"`
		Features    []string `json:"features,omitempty"`
		UseCase     string   `json:"use_case,omitempty"`
	}

	var allTemplates []JSONTemplateInfo

	// Add built-in templates
	for _, template := range builtinTemplates {
		jsonTemplate := JSONTemplateInfo{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			Type:        "builtin",
			Tags:        template.Tags,
		}

		if detailed && template.Config.GoVersion != "" {
			jsonTemplate.GoVersion = template.Config.GoVersion
		}

		if detailed && template.Config.HttpPackage.ID != "" {
			jsonTemplate.HTTP = template.Config.HttpPackage.Name
		}

		if detailed && len(template.Config.Databases) > 0 {
			for _, db := range template.Config.Databases {
				jsonTemplate.Databases = append(jsonTemplate.Databases, db.Database.Name)
			}
		}

		if detailed && len(template.Config.Features) > 0 {
			for _, feature := range template.Config.Features {
				jsonTemplate.Features = append(jsonTemplate.Features, feature.Name)
			}
		}

		allTemplates = append(allTemplates, jsonTemplate)
	}

	// Add custom templates
	for _, template := range customTemplates {
		jsonTemplate := JSONTemplateInfo{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			Type:        "custom",
		}

		allTemplates = append(allTemplates, jsonTemplate)
	}

	response := map[string]interface{}{
		"templates": allTemplates,
		"count":     len(allTemplates),
		"builtin":   len(builtinTemplates),
		"custom":    len(customTemplates),
	}

	formatter := output.NewFormatter(output.FormatJSON, os.Stdout)
	return formatter.OutputResult(response)
}

// printTemplateDetails prints detailed template information
func printTemplateDetails(template BuiltinTemplate, formatter *output.Formatter) {
	if isNoColor() {
		fmt.Printf("  ID: %s\n", template.ID)
		fmt.Printf("  Description: %s\n", template.Description)
		if len(template.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(template.Tags, ", "))
		}
		fmt.Printf("  Go Version: %s\n", template.Config.GoVersion)
		if template.Config.HttpPackage.ID != "" {
			fmt.Printf("  HTTP Framework: %s\n", template.Config.HttpPackage.Name)
		}
		if len(template.Config.Databases) > 0 {
			fmt.Printf("  Databases: ")
			var dbNames []string
			for _, db := range template.Config.Databases {
				dbNames = append(dbNames, db.Database.Name)
			}
			fmt.Printf("%s\n", strings.Join(dbNames, ", "))
		}
		if len(template.Config.Features) > 0 {
			fmt.Printf("  Features: %d features\n", len(template.Config.Features))
		}
	} else {
		fmt.Printf("  %s: %s\n", color.CyanString("ID"), template.ID)
		fmt.Printf("  %s: %s\n", color.CyanString("Description"), template.Description)
		if len(template.Tags) > 0 {
			fmt.Printf("  %s: %s\n", color.CyanString("Tags"), strings.Join(template.Tags, ", "))
		}
		fmt.Printf("  %s: %s\n", color.CyanString("Go Version"), template.Config.GoVersion)
		if template.Config.HttpPackage.ID != "" {
			fmt.Printf("  %s: %s\n", color.CyanString("HTTP Framework"), template.Config.HttpPackage.Name)
		}
		if len(template.Config.Databases) > 0 {
			fmt.Printf("  %s: ", color.CyanString("Databases"))
			var dbNames []string
			for _, db := range template.Config.Databases {
				dbNames = append(dbNames, db.Database.Name)
			}
			fmt.Printf("%s\n", strings.Join(dbNames, ", "))
		}
		if len(template.Config.Features) > 0 {
			fmt.Printf("  %s: %d features\n", color.CyanString("Features"), len(template.Config.Features))
		}
	}
}

// BuiltinTemplate represents a project template with configuration
type BuiltinTemplate struct {
	ID          string
	Name        string
	Description string
	Config      metadata.ProjectConfig
	Tags        []string
}

// getBuiltinTemplatesForTemplateCmd returns built-in templates for template command
func getBuiltinTemplatesForTemplateCmd() []BuiltinTemplate {
	return []BuiltinTemplate{
		{
			ID:          "minimal",
			Name:        "Minimal",
			Description: "Minimal Go project with basic structure",
			Tags:        []string{"basic", "starter"},
			Config: metadata.ProjectConfig{
				GoVersion:   "1.23",
				HttpPackage: metadata.Option{ID: "net-http", Name: "net/http", Description: "Standard library HTTP"},
			},
		},
		{
			ID:          "api",
			Name:        "REST API",
			Description: "REST API with database integration and clean architecture",
			Tags:        []string{"api", "web", "database"},
			Config: metadata.ProjectConfig{
				GoVersion:   "1.23",
				HttpPackage: metadata.Option{ID: "gin", Name: "Gin", Description: "High-performance HTTP web framework"},
				Databases: []metadata.DatabaseSelection{
					{
						Database: metadata.Option{ID: "postgres", Name: "PostgreSQL"},
						Driver:   metadata.Option{ID: "gorm", Name: "GORM"},
					},
				},
				Features: []metadata.Option{
					{ID: "docker", Name: "Docker"},
					{ID: "makefile", Name: "Makefile"},
				},
			},
		},
		{
			ID:          "microservice",
			Name:        "Microservice",
			Description: "Full microservice with gRPC and HTTP endpoints",
			Tags:        []string{"microservice", "grpc", "api"},
			Config: metadata.ProjectConfig{
				GoVersion:   "1.23",
				HttpPackage: metadata.Option{ID: "gin", Name: "Gin"},
				Databases: []metadata.DatabaseSelection{
					{
						Database: metadata.Option{ID: "postgres", Name: "PostgreSQL"},
						Driver:   metadata.Option{ID: "gorm", Name: "GORM"},
					},
				},
				Features: []metadata.Option{
					{ID: "docker", Name: "Docker"},
					{ID: "grpc", Name: "gRPC"},
					{ID: "monitoring", Name: "Monitoring"},
				},
			},
		},
		{
			ID:          "cli",
			Name:        "CLI Application",
			Description: "Command-line application with Cobra framework",
			Tags:        []string{"cli", "tool"},
			Config: metadata.ProjectConfig{
				GoVersion: "1.23",
				Features: []metadata.Option{
					{ID: "cobra", Name: "Cobra CLI"},
					{ID: "testing", Name: "Testing"},
				},
			},
		},
		{
			ID:          "worker",
			Name:        "Background Worker",
			Description: "Background worker with job processing and queues",
			Tags:        []string{"worker", "queue", "background"},
			Config: metadata.ProjectConfig{
				GoVersion: "1.23",
				Databases: []metadata.DatabaseSelection{
					{
						Database: metadata.Option{ID: "redis", Name: "Redis"},
						Driver:   metadata.Option{ID: "redis-client", Name: "Redis Client"},
					},
				},
				Features: []metadata.Option{
					{ID: "docker", Name: "Docker"},
					{ID: "monitoring", Name: "Monitoring"},
				},
			},
		},
		{
			ID:          "web",
			Name:        "Web Application",
			Description: "Web application with HTML templates and sessions",
			Tags:        []string{"web", "html", "templates"},
			Config: metadata.ProjectConfig{
				GoVersion:   "1.23",
				HttpPackage: metadata.Option{ID: "gin", Name: "Gin"},
				Features: []metadata.Option{
					{ID: "sessions", Name: "Sessions"},
					{ID: "static", Name: "Static Files"},
				},
			},
		},
		{
			ID:          "grpc",
			Name:        "gRPC Service",
			Description: "gRPC service with protocol buffers",
			Tags:        []string{"grpc", "protobuf", "service"},
			Config: metadata.ProjectConfig{
				GoVersion: "1.23",
				Features: []metadata.Option{
					{ID: "grpc", Name: "gRPC"},
					{ID: "protobuf", Name: "Protocol Buffers"},
					{ID: "docker", Name: "Docker"},
				},
			},
		},
	}
}

// runTemplateShow shows template details
func runTemplateShow(cmd *cobra.Command, args []string) error {
	templateID := args[0]

	// Try built-in templates first
	builtinTemplates := getBuiltinTemplates()
	var selectedBuiltinTemplate *Template
	for _, template := range builtinTemplates {
		if template.ID == templateID {
			selectedBuiltinTemplate = &template
			break
		}
	}

	// Try custom templates
	tm := templates.NewTemplateManager()
	var selectedCustomTemplate *templates.CustomTemplate
	if tm.TemplateExists(templateID) {
		var err error
		selectedCustomTemplate, err = tm.LoadTemplate(templateID)
		if err != nil {
			return fmt.Errorf("failed to load custom template: %w", err)
		}
	}

	if selectedBuiltinTemplate == nil && selectedCustomTemplate == nil {
		return fmt.Errorf("template not found: %s", templateID)
	}

	if selectedBuiltinTemplate != nil {
		// Display built-in template

		// Display built-in template information
		fmt.Printf("%s %s %s\n", color.HiCyanString("Template:"), selectedBuiltinTemplate.Name, color.HiBlackString("(built-in)"))
		fmt.Printf("%s %s\n", color.HiCyanString("ID:"), selectedBuiltinTemplate.ID)
		fmt.Printf("%s %s\n", color.HiCyanString("Description:"), selectedBuiltinTemplate.Description)

		if len(selectedBuiltinTemplate.Tags) > 0 {
			fmt.Printf("%s %s\n", color.HiCyanString("Tags:"), strings.Join(selectedBuiltinTemplate.Tags, ", "))
		}

		fmt.Printf("\n%s\n", color.HiCyanString("Configuration:"))
		fmt.Printf("  %s: %s\n", color.CyanString("Go Version"), selectedBuiltinTemplate.Config.GoVersion)

		if selectedBuiltinTemplate.Config.HttpPackage.ID != "" {
			fmt.Printf("  %s: %s\n", color.CyanString("HTTP Framework"), selectedBuiltinTemplate.Config.HttpPackage.Name)
		}

		if len(selectedBuiltinTemplate.Config.Databases) > 0 {
			fmt.Printf("  %s:\n", color.CyanString("Databases"))
			for _, db := range selectedBuiltinTemplate.Config.Databases {
				fmt.Printf("    • %s with %s driver\n", db.Database.Name, db.Driver.Name)
			}
		}

		if len(selectedBuiltinTemplate.Config.Features) > 0 {
			fmt.Printf("  %s:\n", color.CyanString("Features"))
			for _, feature := range selectedBuiltinTemplate.Config.Features {
				fmt.Printf("    • %s\n", feature.Name)
			}
		}

		if len(selectedBuiltinTemplate.Config.CustomPackages) > 0 {
			fmt.Printf("  %s:\n", color.CyanString("Custom Packages"))
			for _, pkg := range selectedBuiltinTemplate.Config.CustomPackages {
				fmt.Printf("    • %s\n", pkg)
			}
		}

		fmt.Printf("\n%s\n", color.HiGreenString("Usage:"))
		fmt.Printf("  go-ctl generate my-project --template=%s\n", selectedBuiltinTemplate.ID)
		fmt.Printf("  go-ctl template preview %s --name=my-project\n", selectedBuiltinTemplate.ID)
	} else {
		// Display custom template information
		fmt.Printf("%s %s %s\n", color.HiCyanString("Template:"), selectedCustomTemplate.Name, color.HiMagentaString("(custom)"))
		fmt.Printf("%s %s\n", color.HiCyanString("ID:"), selectedCustomTemplate.ID)
		fmt.Printf("%s %s\n", color.HiCyanString("Description:"), selectedCustomTemplate.Description)

		if selectedCustomTemplate.Author != "" {
			fmt.Printf("%s %s\n", color.HiCyanString("Author:"), selectedCustomTemplate.Author)
		}

		if selectedCustomTemplate.Version != "" {
			fmt.Printf("%s %s\n", color.HiCyanString("Version:"), selectedCustomTemplate.Version)
		}

		if len(selectedCustomTemplate.Tags) > 0 {
			fmt.Printf("%s %s\n", color.HiCyanString("Tags:"), strings.Join(selectedCustomTemplate.Tags, ", "))
		}

		fmt.Printf("%s %s\n", color.HiCyanString("Created:"), selectedCustomTemplate.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("%s %s\n", color.HiCyanString("Updated:"), selectedCustomTemplate.UpdatedAt.Format("2006-01-02 15:04:05"))

		fmt.Printf("\n%s\n", color.HiCyanString("Configuration:"))
		fmt.Printf("  %s: %s\n", color.CyanString("Go Version"), selectedCustomTemplate.Config.GoVersion)

		if selectedCustomTemplate.Config.HttpPackage.ID != "" {
			fmt.Printf("  %s: %s\n", color.CyanString("HTTP Framework"), selectedCustomTemplate.Config.HttpPackage.Name)
		}

		if len(selectedCustomTemplate.Config.Databases) > 0 {
			fmt.Printf("  %s:\n", color.CyanString("Databases"))
			for _, db := range selectedCustomTemplate.Config.Databases {
				fmt.Printf("    • %s with %s driver\n", db.Database.Name, db.Driver.Name)
			}
		}

		if len(selectedCustomTemplate.Config.Features) > 0 {
			fmt.Printf("  %s:\n", color.CyanString("Features"))
			for _, feature := range selectedCustomTemplate.Config.Features {
				fmt.Printf("    • %s\n", feature.Name)
			}
		}

		if len(selectedCustomTemplate.Config.CustomPackages) > 0 {
			fmt.Printf("  %s:\n", color.CyanString("Custom Packages"))
			for _, pkg := range selectedCustomTemplate.Config.CustomPackages {
				fmt.Printf("    • %s\n", pkg)
			}
		}

		if selectedCustomTemplate.Metadata.SourceProject != "" {
			fmt.Printf("\n%s\n", color.HiCyanString("Source:"))
			fmt.Printf("  %s: %s\n", color.CyanString("Project"), selectedCustomTemplate.Metadata.SourceProject)
		}

		fmt.Printf("\n%s\n", color.HiGreenString("Usage:"))
		fmt.Printf("  go-ctl generate my-project --template=%s\n", selectedCustomTemplate.ID)
		fmt.Printf("  go-ctl template preview %s --name=my-project\n", selectedCustomTemplate.ID)
		fmt.Printf("  go-ctl template export %s my-template.yaml\n", selectedCustomTemplate.ID)
		fmt.Printf("  go-ctl template delete %s\n", selectedCustomTemplate.ID)
	}

	return nil
}

// runTemplatePreview previews template project structure
func runTemplatePreview(cmd *cobra.Command, args []string) error {
	templateID := args[0]
	projectName, _ := cmd.Flags().GetString("name")

	templates := getBuiltinTemplates()

	var selectedTemplate *Template
	for _, template := range templates {
		if template.ID == templateID {
			selectedTemplate = &template
			break
		}
	}

	if selectedTemplate == nil {
		return fmt.Errorf("template not found: %s", templateID)
	}

	// Create config with project name
	config := selectedTemplate.Config
	config.ProjectName = projectName

	fmt.Printf("%s %s\n", color.HiCyanString("Template:"), selectedTemplate.Name)
	fmt.Printf("%s %s\n", color.HiCyanString("Project Name:"), projectName)
	fmt.Printf("\n%s\n", color.HiCyanString("Project Structure Preview:"))

	// Generate basic structure preview
	structure := []string{
		"go.mod",
		"README.md",
		fmt.Sprintf("cmd/%s/main.go", projectName),
		"internal/config/config.go",
		"internal/domain/model.go",
		"internal/service/service.go",
		"internal/handler/handler.go",
	}

	// Add HTTP-specific files
	if config.HttpPackage.ID != "" {
		structure = append(structure, "internal/handler/http.go")
	}

	// Add database files
	if len(config.Databases) > 0 {
		structure = append(structure, "internal/storage/db.go")
		for _, db := range config.Databases {
			structure = append(structure, fmt.Sprintf("internal/storage/%s/repository.go", db.Database.ID))
		}
	}

	// Add feature files
	for _, feature := range config.Features {
		switch feature.ID {
		case "gitignore":
			structure = append(structure, ".gitignore")
		case "makefile":
			structure = append(structure, "Makefile")
		case "env":
			structure = append(structure, ".env.example")
		case "air":
			structure = append(structure, ".air.toml")
		case "docker":
			structure = append(structure, "Dockerfile", "docker-compose.yml")
		case "logging":
			structure = append(structure, "internal/logger/logger.go")
		case "testing":
			structure = append(structure, "internal/testing/testing.go", "internal/service/service_test.go")
		}
	}

	// Sort and display structure
	sort.Strings(structure)
	for _, file := range structure {
		fmt.Printf("  %s\n", file)
	}

	fmt.Printf("\n%s\n", color.HiGreenString("✨ Generate this project with:"))
	fmt.Printf("  go-ctl generate %s --template=%s\n", projectName, templateID)

	return nil
}

// runTemplateCreate creates a new custom template
func runTemplateCreate(cmd *cobra.Command, args []string) error {
	templateID := args[0]

	fromProject, _ := cmd.Flags().GetString("from-project")
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")
	author, _ := cmd.Flags().GetString("author")
	tags, _ := cmd.Flags().GetStringSlice("tags")
	interactive, _ := cmd.Flags().GetBool("interactive")
	force, _ := cmd.Flags().GetBool("force")

	tm := templates.NewTemplateManager()

	// Check if template already exists
	if tm.TemplateExists(templateID) && !force {
		return fmt.Errorf("template '%s' already exists. Use --force to overwrite", templateID)
	}

	if fromProject != "" {
		// Create template from existing project
		if name == "" {
			name = fmt.Sprintf("Template from %s", filepath.Base(fromProject))
		}

		printInfo("Creating template from project: %s", fromProject)

		if err := tm.CreateFromProject(fromProject, templateID, name); err != nil {
			return fmt.Errorf("failed to create template: %w", err)
		}

		printSuccess("Template '%s' created from project", templateID)
		return nil
	}

	if interactive {
		return fmt.Errorf("interactive template creation not yet implemented")
	}

	// Create basic template
	if name == "" {
		name = templateID
	}

	template := &templates.CustomTemplate{
		ID:          templateID,
		Name:        name,
		Description: description,
		Author:      author,
		Version:     "1.0.0",
		Tags:        tags,
		Config: metadata.ProjectConfig{
			GoVersion: "1.23",
		},
	}

	if err := tm.CreateTemplate(template); err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	printSuccess("Template '%s' created", templateID)
	printInfo("Edit the template: %s", filepath.Join(tm.GetTemplatesDir(), templateID))

	return nil
}

// runTemplateDelete deletes a custom template
func runTemplateDelete(cmd *cobra.Command, args []string) error {
	templateID := args[0]
	force, _ := cmd.Flags().GetBool("force")

	tm := templates.NewTemplateManager()

	if !tm.TemplateExists(templateID) {
		return fmt.Errorf("template '%s' not found", templateID)
	}

	if !force {
		// TODO: Add confirmation prompt
		printWarning("Use --force to confirm deletion")
		return fmt.Errorf("deletion cancelled")
	}

	if err := tm.DeleteTemplate(templateID); err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	printSuccess("Template '%s' deleted", templateID)
	return nil
}

// runTemplateExport exports a template to file
func runTemplateExport(cmd *cobra.Command, args []string) error {
	templateID := args[0]
	outputFile := args[1]

	tm := templates.NewTemplateManager()

	if !tm.TemplateExists(templateID) {
		return fmt.Errorf("template '%s' not found", templateID)
	}

	printInfo("Exporting template '%s' to %s", templateID, outputFile)

	if err := tm.ExportTemplate(templateID, outputFile); err != nil {
		return fmt.Errorf("failed to export template: %w", err)
	}

	printSuccess("Template exported to %s", outputFile)
	return nil
}

// runTemplateImport imports a template from file
func runTemplateImport(cmd *cobra.Command, args []string) error {
	templateFile := args[0]
	force, _ := cmd.Flags().GetBool("force")

	// TODO: Handle URL imports
	if strings.HasPrefix(templateFile, "http") {
		return fmt.Errorf("URL imports not yet implemented")
	}

	tm := templates.NewTemplateManager()

	printInfo("Importing template from %s", templateFile)

	if err := tm.ImportTemplate(templateFile); err != nil {
		if strings.Contains(err.Error(), "already exists") && !force {
			return fmt.Errorf("%w. Use --force to overwrite", err)
		}
		return fmt.Errorf("failed to import template: %w", err)
	}

	printSuccess("Template imported successfully")
	return nil
}

// newTemplateSuggestCommand creates the template suggest subcommand
func newTemplateSuggestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suggest [project-path]",
		Short: "Get intelligent template suggestions based on project analysis",
		Long: `Analyze an existing project or answer questions to get personalized
template recommendations that best fit your use case.

This command will:
  • Analyze your current project (if provided)
  • Ask intelligent questions about your requirements
  • Suggest the most suitable templates
  • Explain why each template was recommended

Examples:
  # Get suggestions for current directory
  go-ctl template suggest

  # Analyze specific project for suggestions
  go-ctl template suggest ./my-existing-project

  # Interactive questionnaire mode
  go-ctl template suggest --interactive

  # Get suggestions for specific use case
  go-ctl template suggest --use-case=api`,
		Args: cobra.MaximumNArgs(1),
		RunE: runTemplateSuggest,
	}

	cmd.Flags().Bool("interactive", false, "Use interactive questionnaire")
	cmd.Flags().StringP("use-case", "u", "", "Specify use case (api, cli, web, microservice, worker)")
	cmd.Flags().StringSlice("requirements", []string{}, "Specify requirements (database, auth, docker, etc.)")
	cmd.Flags().Int("max-suggestions", 3, "Maximum number of template suggestions")
	cmd.Flags().Bool("explain", true, "Include explanations for suggestions")

	return cmd
}

// runTemplateSuggest executes the template suggest command
func runTemplateSuggest(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	interactive, _ := cmd.Flags().GetBool("interactive")
	useCase, _ := cmd.Flags().GetString("use-case")
	requirements, _ := cmd.Flags().GetStringSlice("requirements")
	maxSuggestions, _ := cmd.Flags().GetInt("max-suggestions")
	explain, _ := cmd.Flags().GetBool("explain")

	fmt.Printf("%s\n", color.HiCyanString("🎯 Smart Template Suggestions"))
	fmt.Printf("%s\n", strings.Repeat("=", 40))

	var projectAnalysis *ProjectRequirements
	var err error

	// Analyze existing project if path exists and contains go.mod
	if goModPath := projectPath + "/go.mod"; fileExists(goModPath) {
		printInfo("Analyzing existing project at %s...", projectPath)
		projectAnalysis, err = analyzeExistingProject(projectPath)
		if err != nil {
			printWarning("Failed to analyze project: %v", err)
		} else {
			printSuccess("Project analysis completed")
		}
	}

	// Interactive mode or missing information
	if interactive || (projectAnalysis == nil && useCase == "") {
		printInfo("Starting interactive questionnaire...")
		projectAnalysis, err = conductInteractiveQuestionnaire(projectAnalysis)
		if err != nil {
			return fmt.Errorf("questionnaire failed: %w", err)
		}
	}

	// Use command line inputs if provided
	if useCase != "" {
		if projectAnalysis == nil {
			projectAnalysis = &ProjectRequirements{}
		}
		projectAnalysis.UseCase = useCase
	}
	if len(requirements) > 0 {
		if projectAnalysis == nil {
			projectAnalysis = &ProjectRequirements{}
		}
		projectAnalysis.Requirements = append(projectAnalysis.Requirements, requirements...)
	}

	// Generate suggestions
	suggestions := generateTemplateSuggestions(projectAnalysis, maxSuggestions)

	if len(suggestions) == 0 {
		printWarning("No suitable templates found for your requirements")
		return nil
	}

	// Display suggestions
	displayTemplateSuggestions(suggestions, explain)

	return nil
}

// ProjectRequirements represents analyzed project requirements
type ProjectRequirements struct {
	UseCase          string   `json:"use_case"`
	Requirements     []string `json:"requirements"`
	ExistingPatterns []string `json:"existing_patterns"`
	Technologies     []string `json:"technologies"`
	DatabaseTypes    []string `json:"database_types"`
	Architecture     string   `json:"architecture"`
	Scale            string   `json:"scale"`
	Team             string   `json:"team"`
}

// TemplateSuggestion represents a template recommendation
type TemplateSuggestion struct {
	Template   Template `json:"template"`
	Score      float64  `json:"score"`
	Reasons    []string `json:"reasons"`
	Pros       []string `json:"pros"`
	Cons       []string `json:"cons"`
	Confidence string   `json:"confidence"`
}

// analyzeExistingProject analyzes an existing project to understand requirements
func analyzeExistingProject(projectPath string) (*ProjectRequirements, error) {
	req := &ProjectRequirements{
		Requirements:     []string{},
		ExistingPatterns: []string{},
		Technologies:     []string{},
		DatabaseTypes:    []string{},
	}

	// Check for common patterns in the project
	if hasWebFramework(projectPath) {
		req.UseCase = "api"
		req.ExistingPatterns = append(req.ExistingPatterns, "web-api")
	}

	if hasCLIPattern(projectPath) {
		req.UseCase = "cli"
		req.ExistingPatterns = append(req.ExistingPatterns, "command-line")
	}

	// Check for database usage
	if hasGormUsage(projectPath) {
		req.Requirements = append(req.Requirements, "database")
		req.Technologies = append(req.Technologies, "gorm")
		req.DatabaseTypes = append(req.DatabaseTypes, "sql")
	}

	if hasMongoUsage(projectPath) {
		req.Requirements = append(req.Requirements, "database")
		req.Technologies = append(req.Technologies, "mongodb")
		req.DatabaseTypes = append(req.DatabaseTypes, "nosql")
	}

	// Check for Docker
	if fileExists(projectPath + "/Dockerfile") {
		req.Requirements = append(req.Requirements, "docker")
	}

	// Check for testing
	if hasTestingFramework(projectPath) {
		req.Requirements = append(req.Requirements, "testing")
	}

	return req, nil
}

// conductInteractiveQuestionnaire conducts interactive questionnaire
func conductInteractiveQuestionnaire(existing *ProjectRequirements) (*ProjectRequirements, error) {
	req := existing
	if req == nil {
		req = &ProjectRequirements{}
	}

	fmt.Printf("\n%s\n", color.HiCyanString("📝 Project Requirements Questionnaire"))

	// Use case question
	if req.UseCase == "" {
		var useCase string
		prompt := &survey.Select{
			Message: "What type of project are you building?",
			Options: []string{
				"REST API - Web API with HTTP endpoints",
				"CLI Tool - Command-line application",
				"Web Application - Full-stack web app",
				"Microservice - Service in distributed architecture",
				"Worker/Job - Background processing service",
				"Library/Package - Reusable Go package",
			},
		}

		if err := survey.AskOne(prompt, &useCase); err != nil {
			return nil, err
		}

		// Extract use case from selection
		if strings.Contains(useCase, "REST API") {
			req.UseCase = "api"
		} else if strings.Contains(useCase, "CLI Tool") {
			req.UseCase = "cli"
		} else if strings.Contains(useCase, "Web Application") {
			req.UseCase = "web"
		} else if strings.Contains(useCase, "Microservice") {
			req.UseCase = "microservice"
		} else if strings.Contains(useCase, "Worker") {
			req.UseCase = "worker"
		} else {
			req.UseCase = "library"
		}
	}

	// Requirements questions
	var requirements []string
	reqPrompt := &survey.MultiSelect{
		Message: "Select the features you need:",
		Options: []string{
			"Database integration",
			"Authentication/Authorization",
			"Docker containerization",
			"Testing framework",
			"Logging system",
			"Configuration management",
			"Hot reload (development)",
			"API documentation",
			"Metrics/Monitoring",
			"Message queues",
		},
	}

	if err := survey.AskOne(reqPrompt, &requirements); err != nil {
		return nil, err
	}

	// Convert selections to requirement tags
	for _, selection := range requirements {
		if strings.Contains(selection, "Database") {
			req.Requirements = append(req.Requirements, "database")
		}
		if strings.Contains(selection, "Authentication") {
			req.Requirements = append(req.Requirements, "auth")
		}
		if strings.Contains(selection, "Docker") {
			req.Requirements = append(req.Requirements, "docker")
		}
		if strings.Contains(selection, "Testing") {
			req.Requirements = append(req.Requirements, "testing")
		}
		if strings.Contains(selection, "Logging") {
			req.Requirements = append(req.Requirements, "logging")
		}
		if strings.Contains(selection, "Configuration") {
			req.Requirements = append(req.Requirements, "config")
		}
		if strings.Contains(selection, "Hot reload") {
			req.Requirements = append(req.Requirements, "air")
		}
		if strings.Contains(selection, "documentation") {
			req.Requirements = append(req.Requirements, "docs")
		}
		if strings.Contains(selection, "Metrics") {
			req.Requirements = append(req.Requirements, "metrics")
		}
		if strings.Contains(selection, "Message queues") {
			req.Requirements = append(req.Requirements, "messaging")
		}
	}

	return req, nil
}

// generateTemplateSuggestions generates template suggestions based on requirements
func generateTemplateSuggestions(req *ProjectRequirements, maxSuggestions int) []TemplateSuggestion {
	if req == nil {
		return []TemplateSuggestion{}
	}

	allTemplates := getBuiltinTemplates()
	var suggestions []TemplateSuggestion

	for _, template := range allTemplates {
		score := calculateTemplateScore(template, req)
		if score > 0.3 { // Only suggest if reasonably relevant
			suggestion := TemplateSuggestion{
				Template:   template,
				Score:      score,
				Reasons:    generateReasons(template, req),
				Pros:       generatePros(template, req),
				Cons:       generateCons(template, req),
				Confidence: getConfidenceLevel(score),
			}
			suggestions = append(suggestions, suggestion)
		}
	}

	// Sort by score
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Score > suggestions[j].Score
	})

	// Limit results
	if len(suggestions) > maxSuggestions {
		suggestions = suggestions[:maxSuggestions]
	}

	return suggestions
}

// calculateTemplateScore calculates how well a template matches requirements
func calculateTemplateScore(template Template, req *ProjectRequirements) float64 {
	score := 0.0

	// Use case matching (most important)
	if matchesUseCase(template, req.UseCase) {
		score += 0.4
	}

	// Requirements matching
	reqScore := 0.0
	for _, requirement := range req.Requirements {
		if hasRequirement(template, requirement) {
			reqScore += 0.1
		}
	}
	score += reqScore

	// Technology matching
	techScore := 0.0
	for _, tech := range req.Technologies {
		if hasTechnology(template, tech) {
			techScore += 0.05
		}
	}
	score += techScore

	// Penalize over-complexity
	if len(template.Config.Features) > len(req.Requirements)*2 {
		score *= 0.8
	}

	return score
}

// displayTemplateSuggestions displays template suggestions
func displayTemplateSuggestions(suggestions []TemplateSuggestion, explain bool) {
	fmt.Printf("\n%s\n", color.HiGreenString("🎯 Template Recommendations"))
	fmt.Printf("%s\n", strings.Repeat("=", 40))

	for i, suggestion := range suggestions {
		rank := i + 1
		confidence := suggestion.Confidence

		fmt.Printf("\n%s %s %s\n",
			color.HiCyanString("#%d", rank),
			color.HiWhiteString(suggestion.Template.Name),
			getConfidenceColor(confidence)(confidence))

		fmt.Printf("  %s: %s\n", color.CyanString("Description"), suggestion.Template.Description)
		fmt.Printf("  %s: %.0f%%\n", color.CyanString("Match Score"), suggestion.Score*100)

		if explain {
			if len(suggestion.Reasons) > 0 {
				fmt.Printf("  %s:\n", color.GreenString("Why this template"))
				for _, reason := range suggestion.Reasons {
					fmt.Printf("    • %s\n", reason)
				}
			}

			if len(suggestion.Pros) > 0 {
				fmt.Printf("  %s:\n", color.GreenString("Advantages"))
				for _, pro := range suggestion.Pros {
					fmt.Printf("    ✓ %s\n", pro)
				}
			}

			if len(suggestion.Cons) > 0 {
				fmt.Printf("  %s:\n", color.YellowString("Considerations"))
				for _, con := range suggestion.Cons {
					fmt.Printf("    ! %s\n", con)
				}
			}
		}

		// Usage example
		fmt.Printf("  %s: %s\n", color.CyanString("Usage"),
			color.YellowString("go-ctl generate my-project --template=%s", suggestion.Template.ID))
	}

	fmt.Printf("\n%s\n", color.HiCyanString("💡 Next Steps:"))
	fmt.Printf("  1. Choose a template: %s\n", color.YellowString("go-ctl generate my-project --template=<template-id>"))
	fmt.Printf("  2. Preview template: %s\n", color.YellowString("go-ctl template preview <template-id>"))
	fmt.Printf("  3. Interactive setup: %s\n", color.YellowString("go-ctl generate --interactive"))
}

// Helper functions for template suggestion

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hasWebFramework(projectPath string) bool {
	// Simple check for common web framework imports
	// In practice, you'd use go/parser to analyze imports
	return fileExists(projectPath + "/main.go") // Simplified
}

func hasCLIPattern(projectPath string) bool {
	// Check for CLI patterns
	return fileExists(projectPath + "/cmd") // Simplified
}

func hasGormUsage(projectPath string) bool {
	// Check for GORM usage
	return true // Simplified - would check imports
}

func hasMongoUsage(projectPath string) bool {
	// Check for MongoDB usage
	return false // Simplified
}

func hasTestingFramework(projectPath string) bool {
	// Check for test files
	return fileExists(projectPath + "/*_test.go") // Simplified
}

func matchesUseCase(template Template, useCase string) bool {
	return strings.Contains(strings.ToLower(template.ID), useCase) ||
		strings.Contains(strings.ToLower(template.Description), useCase)
}

func hasRequirement(template Template, requirement string) bool {
	for _, tag := range template.Tags {
		if strings.Contains(strings.ToLower(tag), requirement) {
			return true
		}
	}
	return false
}

func hasTechnology(template Template, tech string) bool {
	// Check if template includes specific technology
	return strings.Contains(strings.ToLower(template.Description), tech)
}

func generateReasons(template Template, req *ProjectRequirements) []string {
	var reasons []string

	if matchesUseCase(template, req.UseCase) {
		reasons = append(reasons, fmt.Sprintf("Designed for %s applications", req.UseCase))
	}

	matchingReqs := 0
	for _, requirement := range req.Requirements {
		if hasRequirement(template, requirement) {
			matchingReqs++
		}
	}

	if matchingReqs > 0 {
		reasons = append(reasons, fmt.Sprintf("Includes %d of your required features", matchingReqs))
	}

	return reasons
}

func generatePros(template Template, req *ProjectRequirements) []string {
	return []string{
		"Production-ready project structure",
		"Includes best practices and patterns",
		"Well-documented and tested",
	}
}

func generateCons(template Template, req *ProjectRequirements) []string {
	var cons []string

	if len(template.Config.Features) > len(req.Requirements)*2 {
		cons = append(cons, "May include more features than needed")
	}

	return cons
}

func getConfidenceLevel(score float64) string {
	switch {
	case score >= 0.8:
		return "Excellent Match"
	case score >= 0.6:
		return "Good Match"
	case score >= 0.4:
		return "Fair Match"
	default:
		return "Partial Match"
	}
}

func getConfidenceColor(confidence string) func(...interface{}) string {
	switch confidence {
	case "Excellent Match":
		return func(a ...interface{}) string { return color.HiGreenString("%v", a...) }
	case "Good Match":
		return func(a ...interface{}) string { return color.GreenString("%v", a...) }
	case "Fair Match":
		return func(a ...interface{}) string { return color.YellowString("%v", a...) }
	default:
		return func(a ...interface{}) string { return color.HiBlackString("%v", a...) }
	}
}

// Note: findTemplate is imported from generate.go to avoid duplication
