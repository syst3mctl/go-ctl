package interactive

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// ProjectWizard handles interactive project configuration
type ProjectWizard struct {
	options *metadata.ProjectOptions
	config  metadata.ProjectConfig
}

// NewProjectWizard creates a new project wizard
func NewProjectWizard(options *metadata.ProjectOptions) *ProjectWizard {
	return &ProjectWizard{
		options: options,
		config:  metadata.ProjectConfig{},
	}
}

// Run executes the interactive wizard
func (w *ProjectWizard) Run() (metadata.ProjectConfig, error) {
	fmt.Printf("%s\n", color.HiCyanString("✨ Welcome to go-ctl project generator!"))
	fmt.Printf("%s\n\n", color.HiBlackString("This wizard will guide you through creating a new Go project."))

	// Step 1: Project name
	if err := w.askProjectName(); err != nil {
		return w.config, err
	}

	// Step 2: Go version
	if err := w.askGoVersion(); err != nil {
		return w.config, err
	}

	// Step 3: HTTP framework
	if err := w.askHttpFramework(); err != nil {
		return w.config, err
	}

	// Step 4: Database selection
	if err := w.askDatabases(); err != nil {
		return w.config, err
	}

	// Step 5: Additional features
	if err := w.askFeatures(); err != nil {
		return w.config, err
	}

	// Step 6: Custom packages
	if err := w.askCustomPackages(); err != nil {
		return w.config, err
	}

	// Step 7: Configuration summary
	if err := w.showSummary(); err != nil {
		return w.config, err
	}

	return w.config, nil
}

// askProjectName prompts for project name
func (w *ProjectWizard) askProjectName() error {
	var projectName string
	prompt := &survey.Input{
		Message: "Project name:",
		Help:    "Enter a name for your Go project (e.g., my-awesome-api)",
	}

	if err := survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	// Clean project name (remove spaces, convert to lowercase, etc.)
	projectName = strings.ReplaceAll(strings.TrimSpace(projectName), " ", "-")
	projectName = strings.ToLower(projectName)

	w.config.ProjectName = projectName
	return nil
}

// askGoVersion prompts for Go version selection
func (w *ProjectWizard) askGoVersion() error {
	var selectedVersion string
	prompt := &survey.Select{
		Message: "Go version:",
		Options: w.options.GoVersions,
		Default: w.options.GoVersions[0], // Latest version as default
		Help:    "Select the Go version for your project",
	}

	if err := survey.AskOne(prompt, &selectedVersion); err != nil {
		return err
	}

	w.config.GoVersion = selectedVersion
	return nil
}

// askHttpFramework prompts for HTTP framework selection
func (w *ProjectWizard) askHttpFramework() error {
	var options []string
	var descriptions []string

	for _, http := range w.options.Http {
		options = append(options, fmt.Sprintf("%s - %s", http.Name, http.Description))
		descriptions = append(descriptions, http.ID)
	}

	var selected string
	prompt := &survey.Select{
		Message: "HTTP framework:",
		Options: options,
		Default: options[0], // Gin as default
		Help:    "Choose the HTTP framework for your web API",
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return err
	}

	// Find the selected option
	for i, option := range options {
		if option == selected {
			w.config.HttpPackage = w.options.Http[i]
			break
		}
	}

	return nil
}

// askDatabases prompts for database and driver selection
func (w *ProjectWizard) askDatabases() error {
	// First ask if they want to use databases
	var useDatabase bool
	prompt := &survey.Confirm{
		Message: "Do you want to add database support?",
		Default: true,
		Help:    "Select yes to configure database connections for your project",
	}

	if err := survey.AskOne(prompt, &useDatabase); err != nil {
		return err
	}

	if !useDatabase {
		return nil
	}

	// Select databases
	var databaseOptions []string
	for _, db := range w.options.Databases {
		databaseOptions = append(databaseOptions, fmt.Sprintf("%s - %s", db.Name, db.Description))
	}

	var selectedDatabases []string
	dbPrompt := &survey.MultiSelect{
		Message: "Select databases:",
		Options: databaseOptions,
		Help:    "You can select multiple databases for your project",
	}

	if err := survey.AskOne(dbPrompt, &selectedDatabases); err != nil {
		return err
	}

	// For each selected database, ask for driver
	for _, selectedDb := range selectedDatabases {
		var dbOption metadata.Option
		for i, option := range databaseOptions {
			if option == selectedDb {
				dbOption = w.options.Databases[i]
				break
			}
		}

		// Find compatible drivers for this database
		var compatibleDrivers []metadata.Option
		for _, driver := range w.options.DbDrivers {
			if w.isDriverCompatible(dbOption.ID, driver) {
				compatibleDrivers = append(compatibleDrivers, driver)
			}
		}

		if len(compatibleDrivers) == 0 {
			continue // Skip if no compatible drivers
		}

		var driverOptions []string
		for _, driver := range compatibleDrivers {
			driverOptions = append(driverOptions, fmt.Sprintf("%s - %s", driver.Name, driver.Description))
		}

		var selectedDriver string
		driverPrompt := &survey.Select{
			Message: fmt.Sprintf("Select driver for %s:", dbOption.Name),
			Options: driverOptions,
			Default: driverOptions[0],
			Help:    fmt.Sprintf("Choose the driver/ORM for %s database", dbOption.Name),
		}

		if err := survey.AskOne(driverPrompt, &selectedDriver); err != nil {
			return err
		}

		// Find the selected driver
		var driverOption metadata.Option
		for i, option := range driverOptions {
			if option == selectedDriver {
				driverOption = compatibleDrivers[i]
				break
			}
		}

		// Add to configuration
		w.config.Databases = append(w.config.Databases, metadata.DatabaseSelection{
			Database: dbOption,
			Driver:   driverOption,
		})
	}

	return nil
}

// askFeatures prompts for additional features
func (w *ProjectWizard) askFeatures() error {
	var addFeatures bool
	prompt := &survey.Confirm{
		Message: "Do you want to add additional features?",
		Default: true,
		Help:    "Additional features like Docker, Makefile, hot-reload, testing, etc.",
	}

	if err := survey.AskOne(prompt, &addFeatures); err != nil {
		return err
	}

	if !addFeatures {
		return nil
	}

	var featureOptions []string
	for _, feature := range w.options.Features {
		featureOptions = append(featureOptions, fmt.Sprintf("%s - %s", feature.Name, feature.Description))
	}

	var selectedFeatures []string
	featurePrompt := &survey.MultiSelect{
		Message: "Select additional features:",
		Options: featureOptions,
		Help:    "Select features to enhance your project setup",
	}

	if err := survey.AskOne(featurePrompt, &selectedFeatures); err != nil {
		return err
	}

	// Convert selected features to options
	for _, selectedFeature := range selectedFeatures {
		for i, option := range featureOptions {
			if option == selectedFeature {
				w.config.Features = append(w.config.Features, w.options.Features[i])
				break
			}
		}
	}

	return nil
}

// askCustomPackages prompts for custom Go packages
func (w *ProjectWizard) askCustomPackages() error {
	var addPackages bool
	prompt := &survey.Confirm{
		Message: "Do you want to add custom Go packages?",
		Default: false,
		Help:    "Add additional Go packages that will be included in go.mod",
	}

	if err := survey.AskOne(prompt, &addPackages); err != nil {
		return err
	}

	if !addPackages {
		return nil
	}

	for {
		var packageName string
		packagePrompt := &survey.Input{
			Message: "Enter package import path (or press Enter to finish):",
			Help:    "e.g., github.com/google/uuid, github.com/stretchr/testify/assert",
		}

		if err := survey.AskOne(packagePrompt, &packageName); err != nil {
			return err
		}

		packageName = strings.TrimSpace(packageName)
		if packageName == "" {
			break // User pressed Enter without input
		}

		// Basic validation
		if !strings.Contains(packageName, "/") {
			fmt.Printf("%s Invalid package path. Please use full import path (e.g., github.com/user/package)\n",
				color.YellowString("⚠️"))
			continue
		}

		w.config.CustomPackages = append(w.config.CustomPackages, packageName)
		fmt.Printf("%s Added package: %s\n", color.GreenString("✅"), packageName)
	}

	return nil
}

// showSummary shows configuration summary and asks for confirmation
func (w *ProjectWizard) showSummary() error {
	fmt.Printf("\n%s\n", color.HiCyanString("📋 Configuration Summary"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Project details
	fmt.Printf("%s %s\n", color.HiCyanString("Project Name:"), w.config.ProjectName)
	fmt.Printf("%s %s\n", color.HiCyanString("Go Version:"), w.config.GoVersion)
	fmt.Printf("%s %s\n", color.HiCyanString("HTTP Framework:"), w.config.HttpPackage.Name)

	// Databases
	if len(w.config.Databases) > 0 {
		fmt.Printf("%s\n", color.HiCyanString("Databases:"))
		for _, db := range w.config.Databases {
			fmt.Printf("  • %s with %s driver\n", db.Database.Name, db.Driver.Name)
		}
	} else {
		fmt.Printf("%s No databases selected\n", color.HiCyanString("Databases:"))
	}

	// Features
	if len(w.config.Features) > 0 {
		fmt.Printf("%s\n", color.HiCyanString("Features:"))
		for _, feature := range w.config.Features {
			fmt.Printf("  • %s\n", feature.Name)
		}
	} else {
		fmt.Printf("%s No additional features selected\n", color.HiCyanString("Features:"))
	}

	// Custom packages
	if len(w.config.CustomPackages) > 0 {
		fmt.Printf("%s\n", color.HiCyanString("Custom Packages:"))
		for _, pkg := range w.config.CustomPackages {
			fmt.Printf("  • %s\n", pkg)
		}
	}

	// Validation warnings
	if warnings := metadata.ValidateConfig(w.config); len(warnings) > 0 {
		fmt.Printf("\n%s\n", color.YellowString("⚠️  Configuration Warnings:"))
		for _, warning := range warnings {
			fmt.Printf("  • %s\n", warning)
		}
	}

	fmt.Printf("\n%s\n", strings.Repeat("=", 50))

	// Confirmation
	var confirmed bool
	confirmPrompt := &survey.Confirm{
		Message: "Generate project with this configuration?",
		Default: true,
	}

	if err := survey.AskOne(confirmPrompt, &confirmed); err != nil {
		return err
	}

	if !confirmed {
		return fmt.Errorf("project generation cancelled by user")
	}

	return nil
}

// isDriverCompatible checks if a driver is compatible with a database
func (w *ProjectWizard) isDriverCompatible(databaseID string, driver metadata.Option) bool {
	// Define compatibility matrix
	compatibility := map[string][]string{
		"postgres": {"gorm", "sqlx", "database-sql"},
		"mysql":    {"gorm", "sqlx", "database-sql"},
		"sqlite":   {"gorm", "sqlx", "database-sql"},
		"mongodb":  {"mongo-driver"},
		"redis":    {"redis-client"},
		"bigquery": {"database-sql"},
	}

	if compatibleDrivers, exists := compatibility[databaseID]; exists {
		for _, compatibleDriver := range compatibleDrivers {
			if driver.ID == compatibleDriver {
				return true
			}
		}
	}

	return false
}
