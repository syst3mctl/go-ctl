package metadata

import (
	"encoding/json"
	"fmt"
	"os"
)

// Option represents a selectable option with metadata
type Option struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	ImportPath   string              `json:"importPath"`
	Dependencies map[string][]string `json:"dependencies,omitempty"`
}

// ProjectOptions contains all available options for project generation
type ProjectOptions struct {
	GoVersions []string `json:"goVersions"`
	Http       []Option `json:"http"`
	Databases  []Option `json:"databases"`
	DbDrivers  []Option `json:"dbDrivers"`
	Features   []Option `json:"features"`
}

// ProjectConfig represents the user's selected configuration
type ProjectConfig struct {
	ProjectName    string   `json:"projectName"`
	GoVersion      string   `json:"goVersion"`
	HttpPackage    Option   `json:"httpPackage"`
	Database       Option   `json:"database"`
	DbDriver       Option   `json:"dbDriver"`
	Features       []Option `json:"features"`
	CustomPackages []string `json:"customPackages"`
}

// LoadOptions loads the project options from options.json file
func LoadOptions() (*ProjectOptions, error) {
	data, err := os.ReadFile("options.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read options.json: %w", err)
	}

	var options ProjectOptions
	if err := json.Unmarshal(data, &options); err != nil {
		return nil, fmt.Errorf("failed to parse options.json: %w", err)
	}

	return &options, nil
}

// FindOption finds an option by ID in a slice of options
func FindOption(options []Option, id string) Option {
	for _, option := range options {
		if option.ID == id {
			return option
		}
	}
	return Option{} // Return empty option if not found
}

// FindOptions finds multiple options by their IDs
func FindOptions(options []Option, ids []string) []Option {
	var result []Option
	for _, id := range ids {
		if option := FindOption(options, id); option.ID != "" {
			result = append(result, option)
		}
	}
	return result
}

// ValidateConfig validates the project configuration for compatibility
func ValidateConfig(config ProjectConfig) []string {
	var warnings []string

	// Check for incompatible combinations
	if config.Database.ID == "mongodb" && config.DbDriver.ID == "gorm" {
		warnings = append(warnings, "GORM does not support MongoDB. Consider using the MongoDB official driver instead.")
	}

	if config.Database.ID == "bigquery" && config.DbDriver.ID != "database-sql" {
		warnings = append(warnings, "BigQuery works best with the standard database/sql driver.")
	}

	// Validate project name
	if config.ProjectName == "" {
		warnings = append(warnings, "Project name cannot be empty.")
	}

	// Validate Go version
	if config.GoVersion == "" {
		warnings = append(warnings, "Go version must be selected.")
	}

	return warnings
}

// GetAllImports collects all import paths from the configuration
func (config ProjectConfig) GetAllImports() []string {
	var imports []string

	// Add HTTP package import
	if config.HttpPackage.ImportPath != "" {
		imports = append(imports, config.HttpPackage.ImportPath)
	}

	// Add database driver import and its dependencies
	if config.DbDriver.ImportPath != "" {
		imports = append(imports, config.DbDriver.ImportPath)

		// Add database-specific dependencies
		if deps, exists := config.DbDriver.Dependencies[config.Database.ID]; exists {
			imports = append(imports, deps...)
		}
	}

	// Add feature imports
	for _, feature := range config.Features {
		if feature.ImportPath != "" {
			imports = append(imports, feature.ImportPath)
		}
	}

	// Add custom packages
	imports = append(imports, config.CustomPackages...)

	return imports
}

// HasFeature checks if a specific feature is enabled in the configuration
func (config ProjectConfig) HasFeature(featureID string) bool {
	for _, feature := range config.Features {
		if feature.ID == featureID {
			return true
		}
	}
	return false
}
