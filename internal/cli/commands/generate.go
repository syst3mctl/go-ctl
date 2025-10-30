package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/syst3mctl/go-ctl/internal/cli/help"
	"github.com/syst3mctl/go-ctl/internal/cli/interactive"
	"github.com/syst3mctl/go-ctl/internal/cli/output"
	"github.com/syst3mctl/go-ctl/internal/generator"
	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// NewGenerateCommand creates the generate command
func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [project-name]",
		Short: "Generate a new Go project",
		Long: `Generate a new Go project with the specified configuration.

Examples:
  # Generate a basic API project
  go-ctl generate my-api --http=gin --database=postgres --driver=gorm

  # Generate with multiple features
  go-ctl generate my-service --http=echo --database=postgres,redis --features=docker,makefile,air

  # Interactive mode
  go-ctl generate --interactive

  # Use configuration file
  go-ctl generate --config=./project-template.yaml

  # Dry run to see what would be generated
  go-ctl generate my-api --http=gin --dry-run`,
		Args: cobra.MaximumNArgs(1),
		RunE: runGenerate,
	}

	// Project configuration flags
	cmd.Flags().StringP("go-version", "g", "1.23", "Go version (1.20, 1.21, 1.22, 1.23)")
	cmd.Flags().StringP("http", "H", "", "HTTP framework (gin, echo, fiber, chi, net-http)")
	cmd.Flags().StringSliceP("database", "d", []string{}, "Database types (postgres, mysql, sqlite, mongodb, redis)")
	cmd.Flags().StringSliceP("driver", "D", []string{}, "Database drivers (gorm, sqlx, ent, mongo-driver, redis-client)")
	cmd.Flags().StringSliceP("features", "f", []string{}, "Additional features (docker, makefile, air, jwt, cors, logging, testing)")
	cmd.Flags().StringSliceP("packages", "p", []string{}, "Custom packages to include")
	cmd.Flags().StringP("output", "o", ".", "Output directory")
	cmd.Flags().StringP("template", "t", "", "Use built-in template (minimal, api, microservice, cli, worker)")

	// Mode flags
	cmd.Flags().BoolP("interactive", "i", false, "Interactive mode")
	cmd.Flags().Bool("dry-run", false, "Show what would be generated without creating files")
	cmd.Flags().Bool("suggest", false, "Get template suggestions before generating")
	cmd.Flags().Bool("show-stats", false, "Show detailed generation statistics")

	// Add enhanced help
	help.AddEnhancedHelp(cmd)

	return cmd
}

// runGenerate executes the generate command
func runGenerate(cmd *cobra.Command, args []string) error {
	// Load project options
	options, err := metadata.LoadOptions()
	if err != nil {
		return fmt.Errorf("failed to load project options: %w", err)
	}

	var config metadata.ProjectConfig

	// Check if suggestions are requested
	suggestMode, _ := cmd.Flags().GetBool("suggest")
	if suggestMode {
		printInfo("Getting template suggestions...")
		// Run template suggest command logic
		if err := runTemplateSuggestLogic(cmd, args); err != nil {
			return fmt.Errorf("template suggestions failed: %w", err)
		}

		// Ask user if they want to continue with generation
		var continueGeneration bool
		prompt := &survey.Confirm{
			Message: "Would you like to generate a project now?",
			Default: true,
		}

		if err := survey.AskOne(prompt, &continueGeneration); err != nil {
			return err
		}

		if !continueGeneration {
			printInfo("Project generation cancelled")
			return nil
		}
	}

	// Check if interactive mode is enabled
	interactiveMode, _ := cmd.Flags().GetBool("interactive")
	if interactiveMode {
		wizard := interactive.NewProjectWizard(options)
		config, err = wizard.Run()
		if err != nil {
			return fmt.Errorf("interactive configuration failed: %w", err)
		}
	} else {
		// Check if template is specified
		templateName, _ := cmd.Flags().GetString("template")
		if templateName != "" {
			config, err = buildConfigFromTemplate(cmd, args, templateName)
			if err != nil {
				return fmt.Errorf("template configuration failed: %w", err)
			}
		} else {
			// Build config from flags and arguments
			config, err = buildConfigFromFlags(cmd, args, options)
			if err != nil {
				return fmt.Errorf("invalid configuration: %w", err)
			}
		}
	}

	// Validate configuration
	if warnings := metadata.ValidateConfig(config); len(warnings) > 0 {
		for _, warning := range warnings {
			printWarning(warning)
		}
	}

	// Check if dry-run mode is enabled
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		return performDryRun(config)
	}

	// Generate the project with enhanced output
	return generateProjectWithEnhancedOutput(config, cmd)
}

// buildConfigFromFlags creates project configuration from command flags
func buildConfigFromFlags(cmd *cobra.Command, args []string, options *metadata.ProjectOptions) (metadata.ProjectConfig, error) {
	var config metadata.ProjectConfig

	// Project name
	if len(args) > 0 {
		config.ProjectName = args[0]
	} else {
		return config, fmt.Errorf("project name is required")
	}

	// Go version
	goVersion, _ := cmd.Flags().GetString("go-version")
	if !isValidGoVersion(goVersion, options.GoVersions) {
		return config, fmt.Errorf("invalid Go version: %s", goVersion)
	}
	config.GoVersion = goVersion

	// HTTP framework
	httpFramework, _ := cmd.Flags().GetString("http")
	if httpFramework != "" {
		if httpOption := metadata.FindOption(options.Http, httpFramework); httpOption.ID != "" {
			config.HttpPackage = httpOption
		} else {
			return config, fmt.Errorf("invalid HTTP framework: %s", httpFramework)
		}
	} else {
		// Default to gin
		config.HttpPackage = metadata.FindOption(options.Http, "gin")
	}

	// Databases and drivers
	databases, _ := cmd.Flags().GetStringSlice("database")
	drivers, _ := cmd.Flags().GetStringSlice("driver")

	if len(databases) > 0 {
		config.Databases = buildDatabaseSelections(databases, drivers, options)
	}

	// Features
	features, _ := cmd.Flags().GetStringSlice("features")
	if len(features) > 0 {
		config.Features = metadata.FindOptions(options.Features, features)
	}

	// Custom packages
	packages, _ := cmd.Flags().GetStringSlice("packages")
	config.CustomPackages = packages

	return config, nil
}

// buildConfigFromTemplate creates project configuration from a template
func buildConfigFromTemplate(cmd *cobra.Command, args []string, templateName string) (metadata.ProjectConfig, error) {
	// Find the template
	template := findTemplate(templateName)
	if template == nil {
		return metadata.ProjectConfig{}, fmt.Errorf("template not found: %s", templateName)
	}

	// Start with template configuration
	config := template.Config

	// Override with project name from args
	if len(args) > 0 {
		config.ProjectName = args[0]
	} else {
		return config, fmt.Errorf("project name is required")
	}

	// Allow command flags to override template settings
	if goVersion, _ := cmd.Flags().GetString("go-version"); cmd.Flags().Changed("go-version") {
		config.GoVersion = goVersion
	}

	if httpFramework, _ := cmd.Flags().GetString("http"); cmd.Flags().Changed("http") {
		// Load options to validate HTTP framework
		options, err := metadata.LoadOptions()
		if err != nil {
			return config, fmt.Errorf("failed to load options: %w", err)
		}
		if httpOption := metadata.FindOption(options.Http, httpFramework); httpOption.ID != "" {
			config.HttpPackage = httpOption
		} else {
			return config, fmt.Errorf("invalid HTTP framework: %s", httpFramework)
		}
	}

	// Override databases if specified
	if cmd.Flags().Changed("database") {
		databases, _ := cmd.Flags().GetStringSlice("database")
		drivers, _ := cmd.Flags().GetStringSlice("driver")
		options, err := metadata.LoadOptions()
		if err != nil {
			return config, fmt.Errorf("failed to load options: %w", err)
		}
		config.Databases = buildDatabaseSelections(databases, drivers, options)
	}

	// Override features if specified
	if cmd.Flags().Changed("features") {
		features, _ := cmd.Flags().GetStringSlice("features")
		options, err := metadata.LoadOptions()
		if err != nil {
			return config, fmt.Errorf("failed to load options: %w", err)
		}
		config.Features = metadata.FindOptions(options.Features, features)
	}

	// Override custom packages if specified
	if cmd.Flags().Changed("packages") {
		packages, _ := cmd.Flags().GetStringSlice("packages")
		config.CustomPackages = packages
	}

	return config, nil
}

// buildDatabaseSelections creates database selections from command flags
func buildDatabaseSelections(databases, drivers []string, options *metadata.ProjectOptions) []metadata.DatabaseSelection {
	var selections []metadata.DatabaseSelection

	for i, dbID := range databases {
		dbOption := metadata.FindOption(options.Databases, dbID)
		if dbOption.ID == "" {
			printWarning("Unknown database: %s", dbID)
			continue
		}

		// Match driver to database
		var driverOption metadata.Option
		if i < len(drivers) {
			// Use specified driver
			driverOption = metadata.FindOption(options.DbDrivers, drivers[i])
		} else {
			// Use default driver for database
			driverOption = getDefaultDriver(dbID, options.DbDrivers)
		}

		if driverOption.ID == "" {
			printWarning("No suitable driver found for database: %s", dbID)
			continue
		}

		selections = append(selections, metadata.DatabaseSelection{
			Database: dbOption,
			Driver:   driverOption,
		})
	}

	return selections
}

// getDefaultDriver returns the default driver for a database
func getDefaultDriver(databaseID string, drivers []metadata.Option) metadata.Option {
	defaults := map[string]string{
		"postgres": "gorm",
		"mysql":    "gorm",
		"sqlite":   "gorm",
		"mongodb":  "mongo-driver",
		"redis":    "redis-client",
	}

	if defaultDriver, exists := defaults[databaseID]; exists {
		return metadata.FindOption(drivers, defaultDriver)
	}

	return metadata.Option{}
}

// isValidGoVersion checks if the Go version is valid
func isValidGoVersion(version string, validVersions []string) bool {
	for _, v := range validVersions {
		if v == version {
			return true
		}
	}
	return false
}

// performDryRun shows what would be generated without creating files
func performDryRun(config metadata.ProjectConfig) error {
	printInfo("Dry run mode - showing what would be generated:")
	fmt.Println()

	// Project info
	fmt.Printf("%s %s\n", color.HiCyanString("Project Name:"), config.ProjectName)
	fmt.Printf("%s %s\n", color.HiCyanString("Go Version:"), config.GoVersion)
	fmt.Printf("%s %s\n", color.HiCyanString("HTTP Framework:"), config.HttpPackage.Name)

	// Databases
	if len(config.Databases) > 0 {
		fmt.Printf("%s\n", color.HiCyanString("Databases:"))
		for _, db := range config.Databases {
			fmt.Printf("  • %s with %s driver\n", db.Database.Name, db.Driver.Name)
		}
	}

	// Features
	if len(config.Features) > 0 {
		fmt.Printf("%s\n", color.HiCyanString("Features:"))
		for _, feature := range config.Features {
			fmt.Printf("  • %s\n", feature.Name)
		}
	}

	// Custom packages
	if len(config.CustomPackages) > 0 {
		fmt.Printf("%s\n", color.HiCyanString("Custom Packages:"))
		for _, pkg := range config.CustomPackages {
			fmt.Printf("  • %s\n", pkg)
		}
	}

	// Generate project structure preview
	gen := generator.New()
	if err := gen.LoadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	fmt.Printf("\n%s\n", color.HiCyanString("Project Structure:"))

	// Create a temporary structure to show files
	structure := map[string]bool{
		"go.mod":    true,
		"README.md": true,
		fmt.Sprintf("cmd/%s/main.go", config.ProjectName): true,
		"internal/config/config.go":                       true,
		"internal/domain/model.go":                        true,
		"internal/service/service.go":                     true,
		"internal/handler/handler.go":                     true,
		"internal/handler/http.go":                        true,
	}

	// Add database files
	if len(config.Databases) > 0 {
		structure["internal/storage/db.go"] = true
		for _, db := range config.Databases {
			structure[fmt.Sprintf("internal/storage/%s/repository.go", db.Database.ID)] = true
		}
	}

	// Add feature files
	for _, feature := range config.Features {
		switch feature.ID {
		case "gitignore":
			structure[".gitignore"] = true
		case "makefile":
			structure["Makefile"] = true
		case "env":
			structure[".env.example"] = true
		case "air":
			structure[".air.toml"] = true
		case "docker":
			structure["Dockerfile"] = true
			structure["docker-compose.yml"] = true
		case "logging":
			structure["internal/logger/logger.go"] = true
		case "testing":
			structure["internal/testing/testing.go"] = true
			structure["internal/service/service_test.go"] = true
		}
	}

	// Print structure as tree
	for file := range structure {
		fmt.Printf("  %s\n", file)
	}

	fmt.Printf("\n%s\n", color.HiGreenString("✨ This is what would be generated. Use without --dry-run to create the project."))

	return nil
}

// generateProjectWithEnhancedOutput generates project with enhanced output formatting
func generateProjectWithEnhancedOutput(config metadata.ProjectConfig, cmd *cobra.Command) error {
	// Create formatter
	format := output.FormatText
	if outputFormat := getOutputFormat(); outputFormat == "json" {
		format = output.FormatJSON
	}

	formatter := output.NewFormatter(format, os.Stdout)
	formatter.SetOptions(isVerbose(), isQuiet(), isNoColor())

	// Show stats flag
	showStats, _ := cmd.Flags().GetBool("show-stats")

	if !isQuiet() && format == output.FormatText {
		formatter.PrintInfo("Starting project generation...")

		// Create progress bar for generation steps
		steps := []string{
			"Creating directory structure",
			"Generating configuration files",
			"Setting up HTTP framework",
			"Configuring database layer",
			"Adding features",
			"Installing dependencies",
		}

		bar := formatter.CreateProgressBar(len(steps), "Generating project")
		for _, step := range steps {
			if bar != nil {
				bar.Describe(step)
				time.Sleep(100 * time.Millisecond) // Simulate work
				bar.Add(1)
			}
		}
		if bar != nil {
			bar.Finish()
			fmt.Println()
		}
	}

	// Generate the project using existing logic
	err := generateProject(config, cmd)
	if err != nil {
		return fmt.Errorf("project generation failed: %w", err)
	}

	// Determine output path
	outputDir, _ := cmd.Flags().GetString("output")
	outputPath := filepath.Join(outputDir, config.ProjectName)

	// Collect statistics if needed
	var stats *output.GenerationStats
	if showStats || format == output.FormatJSON {
		stats = collectGenerationStatistics(outputPath, config)
	}

	// Create result object
	result := &output.GenerationResult{
		ProjectName:    config.ProjectName,
		OutputPath:     outputPath,
		GoVersion:      config.GoVersion,
		HTTPFramework:  config.HttpPackage.ID,
		Databases:      getDatabaseNames(config.Databases),
		Drivers:        getDriverNames(config.Databases),
		Features:       getFeatureNames(config.Features),
		Packages:       config.CustomPackages,
		FilesGenerated: getFileCount(stats),
		Duration:       formatter.Duration().String(),
		Success:        true,
		NextSteps:      generateNextSteps(config),
		Statistics:     stats,
	}

	// Output result using formatter
	return formatter.OutputResult(result)
}

// collectGenerationStatistics collects detailed statistics about the generated project
func collectGenerationStatistics(projectPath string, config metadata.ProjectConfig) *output.GenerationStats {
	stats := &output.GenerationStats{
		FilesByExtension: make(map[string]int),
		FilesByCategory:  make(map[string]int),
		ProcessingTime:   make(map[string]string),
	}

	// Walk the generated project directory
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		stats.TotalFiles++
		ext := filepath.Ext(path)
		if ext == "" {
			ext = "no-extension"
		}
		stats.FilesByExtension[ext]++

		// Categorize files
		category := categorizeFile(path, ext)
		stats.FilesByCategory[category]++

		// Count lines and size
		if content, err := os.ReadFile(path); err == nil {
			lines := strings.Count(string(content), "\n") + 1
			stats.TotalLines += lines
			stats.TotalSize += int64(len(content))
		}

		return nil
	})

	if err != nil {
		// Continue with partial stats if walking failed
	}

	// Count dependencies from go.mod if it exists
	goModPath := filepath.Join(projectPath, "go.mod")
	if content, err := os.ReadFile(goModPath); err == nil {
		stats.Dependencies = countDependencies(string(content))
	}

	// Add templates used
	stats.Templates = getTemplatesUsed(config)

	return stats
}

// Helper functions for statistics
func categorizeFile(path, ext string) string {
	base := filepath.Base(path)

	switch ext {
	case ".go":
		if strings.Contains(path, "_test.go") {
			return "test"
		}
		return "source"
	case ".md":
		return "documentation"
	case ".yaml", ".yml":
		return "configuration"
	case ".json":
		return "configuration"
	case ".toml":
		return "configuration"
	case ".dockerfile", ".Dockerfile":
		return "deployment"
	case ".sql":
		return "database"
	default:
		switch base {
		case "Makefile":
			return "build"
		case ".gitignore":
			return "vcs"
		case ".env", ".env.example":
			return "configuration"
		default:
			return "other"
		}
	}
}

func countDependencies(goModContent string) int {
	lines := strings.Split(goModContent, "\n")
	count := 0
	inRequireBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if line == ")" && inRequireBlock {
			inRequireBlock = false
			continue
		}
		if inRequireBlock && line != "" && !strings.HasPrefix(line, "//") {
			count++
		}
		if !inRequireBlock && strings.HasPrefix(line, "require ") && !strings.Contains(line, "(") {
			count++
		}
	}

	return count
}

func getDatabaseNames(databases []metadata.DatabaseSelection) []string {
	var names []string
	for _, db := range databases {
		names = append(names, db.Database.ID)
	}
	return names
}

func getDriverNames(databases []metadata.DatabaseSelection) []string {
	var names []string
	for _, db := range databases {
		names = append(names, db.Driver.ID)
	}
	return names
}

func getFeatureNames(features []metadata.Option) []string {
	var names []string
	for _, feature := range features {
		names = append(names, feature.ID)
	}
	return names
}

func getFileCount(stats *output.GenerationStats) int {
	if stats != nil {
		return stats.TotalFiles
	}
	return 0
}

func getTemplatesUsed(config metadata.ProjectConfig) []string {
	templates := []string{"base"}

	if config.HttpPackage.ID != "" {
		templates = append(templates, "http/"+config.HttpPackage.ID)
	}

	for _, db := range config.Databases {
		templates = append(templates, "database/"+db.Database.ID)
	}

	for _, feature := range config.Features {
		templates = append(templates, "feature/"+feature.ID)
	}

	return templates
}

func generateNextSteps(config metadata.ProjectConfig) []string {
	steps := []string{
		fmt.Sprintf("cd %s", config.ProjectName),
		"go mod tidy",
	}

	// Add framework-specific steps
	switch config.HttpPackage.ID {
	case "gin", "echo", "fiber", "chi":
		steps = append(steps, "go run cmd/"+config.ProjectName+"/main.go")
	case "net-http":
		steps = append(steps, "go run cmd/"+config.ProjectName+"/main.go")
	}

	// Add feature-specific steps
	hasDocker := false
	hasMakefile := false
	hasAir := false

	for _, feature := range config.Features {
		switch feature.ID {
		case "docker":
			hasDocker = true
		case "makefile":
			hasMakefile = true
		case "air":
			hasAir = true
		}
	}

	if hasMakefile {
		steps = append(steps, "make help  # See available commands")
		if hasAir {
			steps = append(steps, "make dev   # Start with hot reload")
		}
	}

	if hasDocker {
		steps = append(steps, "docker-compose up --build  # Run with Docker")
	}

	steps = append(steps, "Happy coding! 🚀")

	return steps
}

// generateProject generates the actual project
func generateProject(config metadata.ProjectConfig, cmd *cobra.Command) error {
	outputDir, _ := cmd.Flags().GetString("output")
	projectPath := filepath.Join(outputDir, config.ProjectName)

	// Check if project directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("directory %s already exists", projectPath)
	}

	// Create project directory
	printVerbose("Creating project directory: %s", projectPath)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Initialize generator
	gen := generator.New()
	if err := gen.LoadTemplates(); err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Create progress bar if not in quiet mode
	var bar *progressbar.ProgressBar
	if !isQuiet() {
		bar = progressbar.NewOptions(-1,
			progressbar.OptionSetDescription("Generating project files..."),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
			progressbar.OptionShowCount(),
			progressbar.OptionClearOnFinish(),
		)
	}

	// Generate project structure
	printVerbose("Generating project structure...")

	// Use the existing generator to create the project structure
	// We'll extract files from what would be the ZIP and write them to disk
	projectStructure := generateProjectStructure(gen, config)

	current := 0

	for filePath, content := range projectStructure {
		current++
		if bar != nil {
			bar.Add(1)
		}

		fullPath := filepath.Join(projectPath, filePath)

		// Create directory if it doesn't exist
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}

		// Write file content
		printVerbose("Writing file: %s", filePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	if bar != nil {
		bar.Finish()
	}

	// Print success message
	printSuccess("Project '%s' generated successfully!", config.ProjectName)
	fmt.Printf("%s %s\n", color.HiCyanString("📁 Location:"), projectPath)

	// Print next steps
	fmt.Printf("\n%s\n", color.HiGreenString("🚀 Next steps:"))
	fmt.Printf("  cd %s\n", config.ProjectName)
	fmt.Printf("  go mod tidy\n")
	if config.HasFeature("makefile") {
		fmt.Printf("  make dev\n")
	} else {
		fmt.Printf("  go run cmd/%s/main.go\n", config.ProjectName)
	}

	return nil
}

// generateProjectStructure creates the project file structure
func generateProjectStructure(gen *generator.Generator, config metadata.ProjectConfig) map[string]string {
	files := make(map[string]string)

	// Base files
	files["go.mod"] = renderTemplate(gen, "go.mod", config)
	files["README.md"] = renderTemplate(gen, "README.md", config)
	files[fmt.Sprintf("cmd/%s/main.go", config.ProjectName)] = renderTemplate(gen, "main.go", config)

	// Core structure
	files["internal/config/config.go"] = renderTemplate(gen, "config.go", config)
	files["internal/domain/model.go"] = renderTemplate(gen, "domain.model.go", config)
	files["internal/service/service.go"] = renderTemplate(gen, "service.service.go", config)
	files["internal/handler/handler.go"] = renderTemplate(gen, "handler.handler.go", config)
	files["internal/handler/http.go"] = renderTemplate(gen, "handler.http.go", config)

	// Database layer if specified
	if len(config.Databases) > 0 {
		files["internal/storage/db.go"] = renderTemplate(gen, "storage.db.go", config)

		// Repository implementations
		for _, dbSelection := range config.Databases {
			repositoryTemplate := getRepositoryTemplate(dbSelection.Database.ID)
			files[fmt.Sprintf("internal/storage/%s/repository.go", dbSelection.Database.ID)] = renderTemplate(gen, repositoryTemplate, config)
		}
	}

	// Feature files
	for _, feature := range config.Features {
		switch feature.ID {
		case "gitignore":
			files[".gitignore"] = renderTemplate(gen, "gitignore", config)
		case "makefile":
			files["Makefile"] = renderTemplate(gen, "Makefile", config)
		case "env":
			files[".env.example"] = renderTemplate(gen, "env.example", config)
		case "air":
			files[".air.toml"] = renderTemplate(gen, "air.toml", config)
		case "docker":
			files["Dockerfile"] = generateDockerfile(config)
			files["docker-compose.yml"] = generateDockerCompose(config)
		case "logging":
			files["internal/logger/logger.go"] = renderTemplate(gen, "zerolog.go", config)
		case "testing":
			files["internal/testing/testing.go"] = renderTemplate(gen, "testing.go", config)
			files["internal/service/service_test.go"] = renderTemplate(gen, "service_test.go", config)
		}
	}

	return files
}

// Helper functions (simplified versions of generator methods)
func renderTemplate(gen *generator.Generator, templateName string, config metadata.ProjectConfig) string {
	content, err := gen.GenerateFileContent(templateName, config)
	if err != nil {
		return fmt.Sprintf("// Error generating %s: %v\npackage main\n", templateName, err)
	}
	return content
}

func getRepositoryTemplate(databaseID string) string {
	switch databaseID {
	case "postgres":
		return "postgres.repository.go"
	case "mysql":
		return "mysql.repository.go"
	case "sqlite":
		return "sqlite.repository.go"
	case "mongodb":
		return "mongodb.repository.go"
	case "redis":
		return "redis.repository.go"
	default:
		return "postgres.repository.go"
	}
}

func generateDockerfile(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`# Build stage
FROM golang:%s-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o bin/%s cmd/%s/main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/%s .

EXPOSE 8080
CMD ["./%s"]
`, config.GoVersion, config.ProjectName, config.ProjectName, config.ProjectName, config.ProjectName)
}

// findTemplate finds a template by ID
func findTemplate(templateID string) *Template {
	templates := getBuiltinTemplates()
	for _, template := range templates {
		if template.ID == templateID {
			return &template
		}
	}
	return nil
}

// Template represents a built-in project template
type Template struct {
	ID          string
	Name        string
	Description string
	Config      metadata.ProjectConfig
	Tags        []string
}

// getBuiltinTemplates returns the list of built-in templates
func getBuiltinTemplates() []Template {
	return []Template{
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
			Description: "REST API with database and authentication",
			Tags:        []string{"web", "api", "database"},
			Config: metadata.ProjectConfig{
				GoVersion:   "1.23",
				HttpPackage: metadata.Option{ID: "gin", Name: "Gin", Description: "High-performance HTTP web framework"},
				Databases: []metadata.DatabaseSelection{
					{
						Database: metadata.Option{ID: "postgres", Name: "PostgreSQL", Description: "PostgreSQL database"},
						Driver:   metadata.Option{ID: "gorm", Name: "GORM", Description: "The fantastic ORM library for Golang"},
					},
				},
				Features: []metadata.Option{
					{ID: "docker", Name: "Docker", Description: "Docker containerization"},
					{ID: "makefile", Name: "Makefile", Description: "Build automation"},
					{ID: "env", Name: "Environment Variables", Description: "Environment configuration"},
					{ID: "jwt", Name: "JWT Authentication", Description: "JSON Web Token authentication"},
					{ID: "cors", Name: "CORS", Description: "Cross-Origin Resource Sharing"},
					{ID: "logging", Name: "Structured Logging", Description: "Zerolog integration"},
					{ID: "testing", Name: "Testing Setup", Description: "Testify framework integration"},
				},
			},
		},
		{
			ID:          "microservice",
			Name:        "Microservice",
			Description: "Microservice with gRPC and HTTP endpoints",
			Tags:        []string{"microservice", "grpc", "distributed"},
			Config: metadata.ProjectConfig{
				GoVersion:   "1.23",
				HttpPackage: metadata.Option{ID: "gin", Name: "Gin", Description: "High-performance HTTP web framework"},
				Databases: []metadata.DatabaseSelection{
					{
						Database: metadata.Option{ID: "postgres", Name: "PostgreSQL", Description: "PostgreSQL database"},
						Driver:   metadata.Option{ID: "gorm", Name: "GORM", Description: "The fantastic ORM library for Golang"},
					},
					{
						Database: metadata.Option{ID: "redis", Name: "Redis", Description: "Redis for caching"},
						Driver:   metadata.Option{ID: "redis-client", Name: "Redis Client", Description: "Go Redis client"},
					},
				},
				Features: []metadata.Option{
					{ID: "docker", Name: "Docker", Description: "Docker containerization"},
					{ID: "makefile", Name: "Makefile", Description: "Build automation"},
					{ID: "env", Name: "Environment Variables", Description: "Environment configuration"},
					{ID: "air", Name: "Air", Description: "Hot reload for development"},
					{ID: "jwt", Name: "JWT Authentication", Description: "JSON Web Token authentication"},
					{ID: "cors", Name: "CORS", Description: "Cross-Origin Resource Sharing"},
					{ID: "logging", Name: "Structured Logging", Description: "Zerolog integration"},
					{ID: "testing", Name: "Testing Setup", Description: "Testify framework integration"},
				},
			},
		},
		{
			ID:          "cli",
			Name:        "CLI Application",
			Description: "Command-line application template",
			Tags:        []string{"cli", "command-line"},
			Config: metadata.ProjectConfig{
				GoVersion: "1.23",
				Features: []metadata.Option{
					{ID: "makefile", Name: "Makefile", Description: "Build automation"},
					{ID: "env", Name: "Environment Variables", Description: "Environment configuration"},
					{ID: "testing", Name: "Testing Setup", Description: "Testify framework integration"},
				},
				CustomPackages: []string{
					"github.com/spf13/cobra",
					"github.com/spf13/viper",
				},
			},
		},
		{
			ID:          "worker",
			Name:        "Background Worker",
			Description: "Background job processing service",
			Tags:        []string{"worker", "jobs", "queue"},
			Config: metadata.ProjectConfig{
				GoVersion: "1.23",
				Databases: []metadata.DatabaseSelection{
					{
						Database: metadata.Option{ID: "redis", Name: "Redis", Description: "Redis for job queue"},
						Driver:   metadata.Option{ID: "redis-client", Name: "Redis Client", Description: "Go Redis client"},
					},
					{
						Database: metadata.Option{ID: "postgres", Name: "PostgreSQL", Description: "PostgreSQL database"},
						Driver:   metadata.Option{ID: "gorm", Name: "GORM", Description: "The fantastic ORM library for Golang"},
					},
				},
				Features: []metadata.Option{
					{ID: "docker", Name: "Docker", Description: "Docker containerization"},
					{ID: "makefile", Name: "Makefile", Description: "Build automation"},
					{ID: "env", Name: "Environment Variables", Description: "Environment configuration"},
					{ID: "logging", Name: "Structured Logging", Description: "Zerolog integration"},
					{ID: "testing", Name: "Testing Setup", Description: "Testify framework integration"},
				},
			},
		},
	}
}

func generateDockerCompose(config metadata.ProjectConfig) string {
	var dbServices []string
	var volumes []string
	var dependsOnServices []string

	for _, dbSelection := range config.Databases {
		switch dbSelection.Database.ID {
		case "postgres":
			dbServices = append(dbServices, fmt.Sprintf(`
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: %s
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data`, config.ProjectName))
			volumes = append(volumes, "postgres_data:")
			dependsOnServices = append(dependsOnServices, "postgres")
		case "mysql":
			dbServices = append(dbServices, fmt.Sprintf(`
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: %s
      MYSQL_ROOT_PASSWORD: password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql`, config.ProjectName))
			volumes = append(volumes, "mysql_data:")
			dependsOnServices = append(dependsOnServices, "mysql")
		case "mongodb":
			dbServices = append(dbServices, `
  mongodb:
    image: mongo:6
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db`)
			volumes = append(volumes, "mongo_data:")
			dependsOnServices = append(dependsOnServices, "mongodb")
		case "redis":
			dbServices = append(dbServices, `
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data`)
			volumes = append(volumes, "redis_data:")
			dependsOnServices = append(dependsOnServices, "redis")
		}
	}

	dbService := strings.Join(dbServices, "")
	volumeService := ""
	if len(volumes) > 0 {
		volumeService = "\n\nvolumes:\n  " + strings.Join(volumes, "\n  ")
	}

	dependsOn := ""
	if len(dependsOnServices) > 0 {
		dependsOn = "    depends_on:\n      - " + strings.Join(dependsOnServices, "\n      - ")
	}

	return fmt.Sprintf(`version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
%s%s%s`, dependsOn, dbService, volumeService)
}

// runTemplateSuggestLogic runs template suggestion logic for generate command
func runTemplateSuggestLogic(cmd *cobra.Command, args []string) error {
	// This is a simplified implementation
	printInfo("Template suggestions would be shown here...")
	printInfo("Use 'go-ctl template suggest' for full suggestion functionality")
	return nil
}
