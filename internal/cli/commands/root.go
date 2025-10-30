package commands

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syst3mctl/go-ctl/internal/cli/completion"
)

const version = "1.0.0"

// NewRootCommand creates and returns the root command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "go-ctl",
		Short: "A powerful Go project generator",
		Long: color.HiCyanString(`
🚀 go-ctl - Go Project Generator

go-ctl is a web-based and CLI Go project generator inspired by Spring Boot Initializr.
It provides an intuitive interface for developers to select project options and receive
a downloadable, ready-to-code project skeleton with clean architecture.

Features:
  • Multiple HTTP frameworks (Gin, Echo, Fiber, Chi, net/http)
  • Database integration (PostgreSQL, MySQL, SQLite, MongoDB, Redis)
  • Clean architecture project structure
  • Docker and development tool setup
  • Interactive project configuration
  • Template-based generation system
`),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringP("config", "C", "", "config file path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "quiet mode")

	// Version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s %s\n",
				color.HiGreenString("go-ctl"),
				color.HiYellowString("v%s", version))
		},
	})

	// Global flags for enhanced output
	rootCmd.PersistentFlags().StringP("output-format", "F", "text", "output format (text, json, yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")

	// Add subcommands
	rootCmd.AddCommand(NewGenerateCommand())
	rootCmd.AddCommand(NewConfigCommand())
	rootCmd.AddCommand(NewTemplateCommand())
	rootCmd.AddCommand(NewPackageCommand())
	rootCmd.AddCommand(NewAnalyzeCommand())
	rootCmd.AddCommand(NewCompletionCommand())
	rootCmd.AddCommand(NewDocsCommand())

	// Setup enhanced completion
	setupEnhancedCompletion(rootCmd)

	return rootCmd
}

// initConfig initializes configuration
func initConfig(cmd *cobra.Command) error {
	// Set up viper
	viper.SetEnvPrefix("GO_CTL")
	viper.AutomaticEnv()

	// Read config file if specified
	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		// Search for config in common locations
		viper.SetConfigName(".go-ctl")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")

		// Ignore error if config file is not found
		viper.ReadInConfig()
	}

	// Set up logging based on flags
	verbose, _ := cmd.Flags().GetBool("verbose")
	quiet, _ := cmd.Flags().GetBool("quiet")

	if verbose && quiet {
		return fmt.Errorf("cannot use both --verbose and --quiet flags")
	}

	// Store flags in viper for global access
	viper.Set("verbose", verbose)
	viper.Set("quiet", quiet)

	return nil
}

// isVerbose returns whether verbose output is enabled
func isVerbose() bool {
	return viper.GetBool("verbose")
}

// isQuiet returns whether quiet mode is enabled
func isQuiet() bool {
	return viper.GetBool("quiet")
}

// printInfo prints info message if not in quiet mode
func printInfo(format string, args ...interface{}) {
	if !isQuiet() {
		fmt.Printf(color.HiCyanString("ℹ ")+format+"\n", args...)
	}
}

// printSuccess prints success message if not in quiet mode
func printSuccess(format string, args ...interface{}) {
	if !isQuiet() {
		fmt.Printf(color.HiGreenString("✅ ")+format+"\n", args...)
	}
}

// printWarning prints warning message if not in quiet mode
func printWarning(format string, args ...interface{}) {
	if !isQuiet() {
		fmt.Printf(color.HiYellowString("⚠️  ")+format+"\n", args...)
	}
}

// printError prints error message to stderr
func printError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, color.HiRedString("❌ ")+format+"\n", args...)
}

// printVerbose prints verbose message if verbose mode is enabled
func printVerbose(format string, args ...interface{}) {
	if isVerbose() {
		fmt.Printf(color.HiBlackString("🔍 ")+format+"\n", args...)
	}
}

// setupEnhancedCompletion sets up enhanced completion for all commands
func setupEnhancedCompletion(rootCmd *cobra.Command) {
	// Setup dynamic completion for enhanced shell completion
	completion.SetupDynamicCompletion(rootCmd)
}

// getOutputFormat returns the configured output format
func getOutputFormat() string {
	return viper.GetString("output-format")
}

// isNoColor returns whether color output is disabled
func isNoColor() bool {
	return viper.GetBool("no-color")
}
