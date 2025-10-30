package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/syst3mctl/go-ctl/internal/metadata"
	"gopkg.in/yaml.v3"
)

// CLIConfig represents the CLI configuration structure
type CLIConfig struct {
	Project metadata.ProjectConfig `yaml:"project"`
	CLI     CLISettings            `yaml:"cli"`
}

// CLISettings contains CLI-specific settings
type CLISettings struct {
	DefaultOutput     string `yaml:"default_output"`
	InteractiveMode   bool   `yaml:"interactive_mode"`
	ColorOutput       bool   `yaml:"color_output"`
	AutoUpdate        bool   `yaml:"auto_update"`
	TemplateDirectory string `yaml:"template_directory"`
}

// DefaultCLIConfig returns a default CLI configuration
func DefaultCLIConfig() CLIConfig {
	return CLIConfig{
		CLI: CLISettings{
			DefaultOutput:     ".",
			InteractiveMode:   false,
			ColorOutput:       true,
			AutoUpdate:        false,
			TemplateDirectory: "",
		},
		Project: metadata.ProjectConfig{
			GoVersion: "1.23",
			HttpPackage: metadata.Option{
				ID:          "gin",
				Name:        "Gin",
				Description: "High-performance HTTP web framework",
				ImportPath:  "github.com/gin-gonic/gin",
			},
			Databases:      []metadata.DatabaseSelection{},
			Features:       []metadata.Option{},
			CustomPackages: []string{},
		},
	}
}

// LoadCLIConfig loads CLI configuration from a file
func LoadCLIConfig(configPath string) (CLIConfig, error) {
	var config CLIConfig

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, fmt.Errorf("configuration file does not exist: %s", configPath)
	}

	// Read file content
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	return config, nil
}

// SaveCLIConfig saves CLI configuration to a file
func SaveCLIConfig(config CLIConfig, configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// FindConfigFile searches for configuration file in common locations
func FindConfigFile() string {
	// Priority order: current directory, then home directory
	locations := []string{
		".go-ctl.yaml",
		".go-ctl.yml",
	}

	// Check current directory first
	for _, filename := range locations {
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}

	// Check home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		for _, filename := range locations {
			fullPath := filepath.Join(homeDir, filename)
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath
			}
		}
	}

	return ""
}

// MergeWithDefaults merges user config with default values
func MergeWithDefaults(userConfig CLIConfig) CLIConfig {
	defaultConfig := DefaultCLIConfig()

	// Merge CLI settings
	if userConfig.CLI.DefaultOutput == "" {
		userConfig.CLI.DefaultOutput = defaultConfig.CLI.DefaultOutput
	}

	// Merge project settings
	if userConfig.Project.GoVersion == "" {
		userConfig.Project.GoVersion = defaultConfig.Project.GoVersion
	}

	if userConfig.Project.HttpPackage.ID == "" {
		userConfig.Project.HttpPackage = defaultConfig.Project.HttpPackage
	}

	return userConfig
}

// ValidateConfig validates the CLI configuration
func ValidateConfig(config CLIConfig) error {
	// Validate CLI settings
	if config.CLI.DefaultOutput == "" {
		return fmt.Errorf("default_output cannot be empty")
	}

	// Validate default output directory exists or can be created
	if err := os.MkdirAll(config.CLI.DefaultOutput, 0755); err != nil {
		return fmt.Errorf("invalid default_output directory: %w", err)
	}

	// Validate template directory if specified
	if config.CLI.TemplateDirectory != "" {
		if _, err := os.Stat(config.CLI.TemplateDirectory); os.IsNotExist(err) {
			return fmt.Errorf("template_directory does not exist: %s", config.CLI.TemplateDirectory)
		}
	}

	// Validate project configuration if present
	if config.Project.ProjectName != "" {
		warnings := metadata.ValidateConfig(config.Project)
		if len(warnings) > 0 {
			// For now, we just return the first warning as an error
			// In a more sophisticated implementation, we might want to separate
			// warnings from errors
			return fmt.Errorf("project configuration warning: %s", warnings[0])
		}
	}

	return nil
}
