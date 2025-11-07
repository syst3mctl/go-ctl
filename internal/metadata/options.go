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
	GoVersions []string        `json:"goVersions"`
	Http       []Option        `json:"http"`
	Databases  []Option         `json:"databases"`
	DbDrivers  []Option         `json:"dbDrivers"`
	Features   []Option         `json:"features"`
	Frontend   *FrontendOptions `json:"frontend,omitempty"`
}

// FrontendOptions contains all available front-end options
type FrontendOptions struct {
	Languages []Option `json:"languages"`
	BuildTools []Option `json:"buildTools"`
	Linters   []Option `json:"linters"`
	Features  []Option `json:"features"`
}

// DatabaseSelection represents a database with its selected driver
type DatabaseSelection struct {
	Database Option `json:"database"`
	Driver   Option `json:"driver"`
}

// ProjectConfig represents the user's selected configuration
type ProjectConfig struct {
	ProjectName    string              `json:"projectName"`
	ProjectType    string              `json:"projectType"` // "backend", "frontend", or "fullstack"
	GoVersion      string              `json:"goVersion"`
	HttpPackage    Option              `json:"httpPackage"`
	Databases      []DatabaseSelection `json:"databases"`
	Features       []Option            `json:"features"`
	CustomPackages []string            `json:"customPackages"`
	FrontendConfig *FrontendConfig     `json:"frontendConfig,omitempty"`
}

// FrontendConfig represents the front-end project configuration
type FrontendConfig struct {
	Language      Option   `json:"language"`      // JavaScript or TypeScript
	BuildTool     Option   `json:"buildTool"`    // Vite
	Linter        Option   `json:"linter"`        // ESLint
	Features      []Option `json:"features"`     // Prettier, React Router, Tailwind, etc.
	CustomPackages []string `json:"customPackages"` // npm packages
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

	// Validate project name
	if config.ProjectName == "" {
		warnings = append(warnings, "Project name cannot be empty.")
	}

	// Validate backend configuration
	if config.ProjectType == "backend" || config.ProjectType == "fullstack" {
		// Check for incompatible combinations in selected databases
		for _, dbSelection := range config.Databases {
			if dbSelection.Database.ID == "mongodb" && dbSelection.Driver.ID == "gorm" {
				warnings = append(warnings, "GORM does not support MongoDB. Consider using the MongoDB official driver instead.")
			}

			if dbSelection.Database.ID == "bigquery" && dbSelection.Driver.ID != "database-sql" {
				warnings = append(warnings, "BigQuery works best with the standard database/sql driver.")
			}
		}

		// Validate Go version
		if config.GoVersion == "" {
			warnings = append(warnings, "Go version must be selected.")
		}
	}

	// Validate frontend configuration
	if config.ProjectType == "frontend" || config.ProjectType == "fullstack" {
		if config.FrontendConfig == nil {
			warnings = append(warnings, "Frontend configuration is required for frontend projects.")
		} else {
			if config.FrontendConfig.Language.ID == "" {
				warnings = append(warnings, "Frontend language must be selected.")
			}
			if config.FrontendConfig.BuildTool.ID == "" {
				warnings = append(warnings, "Build tool must be selected.")
			}
		}
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

	// Add database driver imports and their dependencies
	for _, dbSelection := range config.Databases {
		if dbSelection.Driver.ImportPath != "" {
			imports = append(imports, dbSelection.Driver.ImportPath)

			// Add database-specific dependencies
			if deps, exists := dbSelection.Driver.Dependencies[dbSelection.Database.ID]; exists {
				imports = append(imports, deps...)
			}
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
