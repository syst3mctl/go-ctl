package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/syst3mctl/go-ctl/internal/cli/config"
	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// NewConfigCommand creates the config command
func NewConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage go-ctl configuration",
		Long: `Manage go-ctl configuration files and settings.

Configuration files allow you to define reusable project templates
and default settings for project generation.`,
	}

	cmd.AddCommand(newConfigInitCommand())
	cmd.AddCommand(newConfigValidateCommand())
	cmd.AddCommand(newConfigShowCommand())
	cmd.AddCommand(newConfigSetCommand())

	return cmd
}

// newConfigInitCommand creates the config init subcommand
func newConfigInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new configuration file",
		Long: `Create a new .go-ctl.yaml configuration file in the current directory
or your home directory with default settings.

Examples:
  # Create config in current directory
  go-ctl config init

  # Create config in home directory
  go-ctl config init --global`,
		RunE: runConfigInit,
	}

	cmd.Flags().Bool("global", false, "Create config in home directory")
	cmd.Flags().Bool("force", false, "Overwrite existing config file")

	return cmd
}

// newConfigValidateCommand creates the config validate subcommand
func newConfigValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate [config-file]",
		Short: "Validate a configuration file",
		Long: `Validate the syntax and content of a go-ctl configuration file.

Examples:
  # Validate default config file
  go-ctl config validate

  # Validate specific config file
  go-ctl config validate ./my-template.yaml`,
		Args: cobra.MaximumNArgs(1),
		RunE: runConfigValidate,
	}

	return cmd
}

// newConfigShowCommand creates the config show subcommand
func newConfigShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [config-file]",
		Short: "Show current configuration",
		Long: `Display the current configuration settings.

Examples:
  # Show current active configuration
  go-ctl config show

  # Show specific configuration file
  go-ctl config show ./my-template.yaml`,
		Args: cobra.MaximumNArgs(1),
		RunE: runConfigShow,
	}

	return cmd
}

// newConfigSetCommand creates the config set subcommand
func newConfigSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Long: `Set a configuration value in the default config file.

Examples:
  # Set default Go version
  go-ctl config set project.go_version 1.23

  # Set default HTTP framework
  go-ctl config set project.http_framework gin

  # Set default output directory
  go-ctl config set cli.default_output ./projects`,
		Args: cobra.ExactArgs(2),
		RunE: runConfigSet,
	}

	return cmd
}

// runConfigInit initializes a new configuration file
func runConfigInit(cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")
	force, _ := cmd.Flags().GetBool("force")

	var configPath string
	if global {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".go-ctl.yaml")
	} else {
		configPath = ".go-ctl.yaml"
	}

	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil && !force {
		return fmt.Errorf("configuration file already exists: %s (use --force to overwrite)", configPath)
	}

	// Create default configuration
	defaultConfig := config.DefaultCLIConfig()

	// Write configuration file
	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	printSuccess("Configuration file created: %s", configPath)
	printInfo("Edit the file to customize your default project settings")

	return nil
}

// runConfigValidate validates a configuration file
func runConfigValidate(cmd *cobra.Command, args []string) error {
	var configPath string
	if len(args) > 0 {
		configPath = args[0]
	} else {
		configPath = findConfigFile()
		if configPath == "" {
			return fmt.Errorf("no configuration file found")
		}
	}

	// Load and validate configuration
	cliConfig, err := config.LoadCLIConfig(configPath)
	if err != nil {
		printError("Configuration validation failed: %v", err)
		return err
	}

	// Load project options for validation
	options, err := metadata.LoadOptions()
	if err != nil {
		return fmt.Errorf("failed to load project options: %w", err)
	}

	// Validate project configuration if present
	if cliConfig.Project.ProjectName != "" {
		warnings := metadata.ValidateConfig(cliConfig.Project)
		if len(warnings) > 0 {
			printWarning("Configuration warnings:")
			for _, warning := range warnings {
				fmt.Printf("  • %s\n", warning)
			}
		}
	}

	// Validate option IDs
	if err := validateOptionIDs(cliConfig.Project, options); err != nil {
		printError("Invalid configuration: %v", err)
		return err
	}

	printSuccess("Configuration file is valid: %s", configPath)
	return nil
}

// runConfigShow displays current configuration
func runConfigShow(cmd *cobra.Command, args []string) error {
	var configPath string
	if len(args) > 0 {
		configPath = args[0]
	} else {
		configPath = findConfigFile()
		if configPath == "" {
			printInfo("No configuration file found")
			return nil
		}
	}

	// Load configuration
	cliConfig, err := config.LoadCLIConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Display configuration
	fmt.Printf("%s %s\n", color.HiCyanString("Configuration file:"), configPath)
	fmt.Printf("%s\n", color.HiCyanString("Settings:"))

	// CLI settings
	if cliConfig.CLI.DefaultOutput != "" {
		fmt.Printf("  %s: %s\n", color.CyanString("Default output"), cliConfig.CLI.DefaultOutput)
	}
	fmt.Printf("  %s: %t\n", color.CyanString("Interactive mode"), cliConfig.CLI.InteractiveMode)
	fmt.Printf("  %s: %t\n", color.CyanString("Color output"), cliConfig.CLI.ColorOutput)

	// Project defaults
	if cliConfig.Project.ProjectName != "" {
		fmt.Printf("  %s: %s\n", color.CyanString("Project template"), cliConfig.Project.ProjectName)
		fmt.Printf("  %s: %s\n", color.CyanString("Go version"), cliConfig.Project.GoVersion)
		if cliConfig.Project.HttpPackage.ID != "" {
			fmt.Printf("  %s: %s\n", color.CyanString("HTTP framework"), cliConfig.Project.HttpPackage.Name)
		}
		if len(cliConfig.Project.Databases) > 0 {
			fmt.Printf("  %s:\n", color.CyanString("Databases"))
			for _, db := range cliConfig.Project.Databases {
				fmt.Printf("    • %s (%s)\n", db.Database.Name, db.Driver.Name)
			}
		}
		if len(cliConfig.Project.Features) > 0 {
			fmt.Printf("  %s:\n", color.CyanString("Features"))
			for _, feature := range cliConfig.Project.Features {
				fmt.Printf("    • %s\n", feature.Name)
			}
		}
	}

	return nil
}

// runConfigSet sets a configuration value
func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	configPath := findConfigFile()
	if configPath == "" {
		return fmt.Errorf("no configuration file found. Run 'go-ctl config init' first")
	}

	// Load existing configuration
	v := viper.New()
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	// Set the value
	v.Set(key, value)

	// Write back to file
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	printSuccess("Configuration updated: %s = %s", key, value)
	return nil
}

// findConfigFile searches for configuration file in common locations
func findConfigFile() string {
	// Check current directory
	if _, err := os.Stat(".go-ctl.yaml"); err == nil {
		return ".go-ctl.yaml"
	}
	if _, err := os.Stat(".go-ctl.yml"); err == nil {
		return ".go-ctl.yml"
	}

	// Check home directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		yamlPath := filepath.Join(homeDir, ".go-ctl.yaml")
		if _, err := os.Stat(yamlPath); err == nil {
			return yamlPath
		}
		ymlPath := filepath.Join(homeDir, ".go-ctl.yml")
		if _, err := os.Stat(ymlPath); err == nil {
			return ymlPath
		}
	}

	return ""
}

// validateOptionIDs validates that all option IDs in configuration are valid
func validateOptionIDs(projectConfig metadata.ProjectConfig, options *metadata.ProjectOptions) error {
	// Validate Go version
	if projectConfig.GoVersion != "" {
		valid := false
		for _, version := range options.GoVersions {
			if version == projectConfig.GoVersion {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid Go version: %s", projectConfig.GoVersion)
		}
	}

	// Validate HTTP framework
	if projectConfig.HttpPackage.ID != "" {
		if metadata.FindOption(options.Http, projectConfig.HttpPackage.ID).ID == "" {
			return fmt.Errorf("invalid HTTP framework: %s", projectConfig.HttpPackage.ID)
		}
	}

	// Validate databases and drivers
	for _, db := range projectConfig.Databases {
		if metadata.FindOption(options.Databases, db.Database.ID).ID == "" {
			return fmt.Errorf("invalid database: %s", db.Database.ID)
		}
		if metadata.FindOption(options.DbDrivers, db.Driver.ID).ID == "" {
			return fmt.Errorf("invalid database driver: %s", db.Driver.ID)
		}
	}

	// Validate features
	for _, feature := range projectConfig.Features {
		if metadata.FindOption(options.Features, feature.ID).ID == "" {
			return fmt.Errorf("invalid feature: %s", feature.ID)
		}
	}

	return nil
}
