package completion

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// SetupDynamicCompletion sets up dynamic completion for all commands
func SetupDynamicCompletion(rootCmd *cobra.Command) {
	// Load options for completion
	options, err := metadata.LoadOptions()
	if err != nil {
		// Fallback to static completion if options can't be loaded
		return
	}

	setupGenerateCompletion(rootCmd, options)
	setupPackageCompletion(rootCmd)
	setupTemplateCompletion(rootCmd)
	setupAnalyzeCompletion(rootCmd)
}

// setupGenerateCompletion sets up completion for the generate command
func setupGenerateCompletion(rootCmd *cobra.Command, options *metadata.ProjectOptions) {
	generateCmd := findCommand(rootCmd, "generate")
	if generateCmd == nil {
		return
	}

	// Go version completion
	generateCmd.RegisterFlagCompletionFunc("go-version", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for _, version := range options.GoVersions {
			if strings.HasPrefix(version, toComplete) {
				completions = append(completions, version)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// HTTP framework completion
	generateCmd.RegisterFlagCompletionFunc("http", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for _, framework := range options.Http {
			if strings.HasPrefix(framework.ID, toComplete) {
				description := framework.Name
				if framework.Description != "" {
					description += " - " + framework.Description
				}
				completions = append(completions, framework.ID+"\t"+description)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// Database completion
	generateCmd.RegisterFlagCompletionFunc("database", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for _, db := range options.Databases {
			if strings.HasPrefix(db.ID, toComplete) {
				description := db.Name
				if db.Description != "" {
					description += " - " + db.Description
				}
				completions = append(completions, db.ID+"\t"+description)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// Driver completion
	generateCmd.RegisterFlagCompletionFunc("driver", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for _, driver := range options.DbDrivers {
			if strings.HasPrefix(driver.ID, toComplete) {
				description := driver.Name
				if driver.Description != "" {
					description += " - " + driver.Description
				}
				completions = append(completions, driver.ID+"\t"+description)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// Features completion
	generateCmd.RegisterFlagCompletionFunc("features", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for _, feature := range options.Features {
			if strings.HasPrefix(feature.ID, toComplete) {
				description := feature.Name
				if feature.Description != "" {
					description += " - " + feature.Description
				}
				completions = append(completions, feature.ID+"\t"+description)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// Template completion
	generateCmd.RegisterFlagCompletionFunc("template", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		templates := []string{
			"minimal\tMinimal Go project with basic structure",
			"api\tREST API with database and authentication",
			"microservice\tMicroservice with gRPC and HTTP endpoints",
			"cli\tCommand-line application template",
			"worker\tBackground worker with queue processing",
			"web\tWeb application with HTML templates",
			"grpc\tgRPC service with protocol buffers",
		}

		var completions []string
		for _, template := range templates {
			if strings.HasPrefix(template, toComplete) {
				completions = append(completions, template)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// Output directory completion
	generateCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterDirs
	})
}

// setupPackageCompletion sets up completion for the package command
func setupPackageCompletion(rootCmd *cobra.Command) {
	packageCmd := findCommand(rootCmd, "package")
	if packageCmd == nil {
		return
	}

	// Package search subcommands
	searchCmd := findCommand(packageCmd, "search")
	if searchCmd != nil {
		// Category completion
		searchCmd.RegisterFlagCompletionFunc("category", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			categories := []string{
				"web\tWeb frameworks and HTTP libraries",
				"database\tDatabase drivers and ORMs",
				"testing\tTesting frameworks and utilities",
				"cli\tCommand-line interface libraries",
				"logging\tLogging and observability",
				"auth\tAuthentication and authorization",
				"validation\tInput validation libraries",
				"utils\tUtility libraries and helpers",
				"crypto\tCryptography and security",
				"json\tJSON processing libraries",
				"time\tTime and date utilities",
				"config\tConfiguration management",
				"cache\tCaching solutions",
				"queue\tMessage queues and job processing",
				"monitoring\tMonitoring and metrics",
			}

			var completions []string
			for _, category := range categories {
				if strings.HasPrefix(category, toComplete) {
					completions = append(completions, category)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		})
	}

	// Popular packages completion
	popularCmd := findCommand(packageCmd, "popular")
	if popularCmd != nil {
		popularCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			categories := []string{"web", "database", "testing", "cli", "logging", "auth", "validation", "utils"}
			var completions []string
			for _, category := range categories {
				if strings.HasPrefix(category, toComplete) {
					completions = append(completions, category)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		}
	}

	// Upgrade command completion
	upgradeCmd := findCommand(packageCmd, "upgrade")
	if upgradeCmd != nil {
		upgradeCmd.RegisterFlagCompletionFunc("focus", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			focuses := []string{
				"security\tFocus on security vulnerabilities",
				"outdated\tFocus on outdated dependencies",
				"breaking\tFocus on breaking changes",
				"performance\tFocus on performance improvements",
			}

			var completions []string
			for _, focus := range focuses {
				if strings.HasPrefix(focus, toComplete) {
					completions = append(completions, focus)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		})
	}
}

// setupTemplateCompletion sets up completion for the template command
func setupTemplateCompletion(rootCmd *cobra.Command) {
	templateCmd := findCommand(rootCmd, "template")
	if templateCmd == nil {
		return
	}

	// Template names completion for show command
	showCmd := findCommand(templateCmd, "show")
	if showCmd != nil {
		showCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			templates := []string{"minimal", "api", "microservice", "cli", "worker", "web", "grpc"}
			var completions []string
			for _, template := range templates {
				if strings.HasPrefix(template, toComplete) {
					completions = append(completions, template)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		}
	}

	// Template suggest completion
	suggestCmd := findCommand(templateCmd, "suggest")
	if suggestCmd != nil {
		suggestCmd.RegisterFlagCompletionFunc("use-case", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			useCases := []string{
				"api\tREST API service",
				"web\tWeb application",
				"cli\tCommand-line tool",
				"microservice\tMicroservice architecture",
				"worker\tBackground worker",
				"grpc\tgRPC service",
				"library\tReusable library",
				"desktop\tDesktop application",
				"game\tGame development",
				"iot\tIoT applications",
			}

			var completions []string
			for _, useCase := range useCases {
				if strings.HasPrefix(useCase, toComplete) {
					completions = append(completions, useCase)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		})

		suggestCmd.RegisterFlagCompletionFunc("requirements", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			requirements := []string{
				"database\tDatabase integration",
				"auth\tAuthentication system",
				"docker\tDocker containerization",
				"kubernetes\tKubernetes deployment",
				"testing\tTesting framework",
				"logging\tStructured logging",
				"monitoring\tMonitoring and metrics",
				"cache\tCaching layer",
				"queue\tMessage queue",
				"websocket\tWebSocket support",
				"grpc\tgRPC support",
				"rest\tREST API",
				"graphql\tGraphQL API",
				"spa\tSingle Page Application",
				"ssr\tServer-Side Rendering",
			}

			var completions []string
			for _, req := range requirements {
				if strings.HasPrefix(req, toComplete) {
					completions = append(completions, req)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		})
	}
}

// setupAnalyzeCompletion sets up completion for the analyze command
func setupAnalyzeCompletion(rootCmd *cobra.Command) {
	analyzeCmd := findCommand(rootCmd, "analyze")
	if analyzeCmd == nil {
		return
	}

	// Focus areas completion
	analyzeCmd.RegisterFlagCompletionFunc("focus", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		focuses := []string{
			"dependencies\tDependency analysis",
			"security\tSecurity analysis",
			"quality\tCode quality metrics",
			"architecture\tArchitecture patterns",
			"performance\tPerformance analysis",
			"testing\tTest coverage analysis",
		}

		var completions []string
		for _, focus := range focuses {
			if strings.HasPrefix(focus, toComplete) {
				completions = append(completions, focus)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// Output format completion
	analyzeCmd.RegisterFlagCompletionFunc("output-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		formats := []string{
			"text\tHuman-readable text output",
			"json\tJSON format for scripting",
			"yaml\tYAML format",
		}

		var completions []string
		for _, format := range formats {
			if strings.HasPrefix(format, toComplete) {
				completions = append(completions, format)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	})

	// Directory completion for project path
	analyzeCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterDirs
	}
}

// setupConfigCompletion sets up completion for the config command
func setupConfigCompletion(rootCmd *cobra.Command) {
	configCmd := findCommand(rootCmd, "config")
	if configCmd == nil {
		return
	}

	// Config file paths completion
	configCmd.RegisterFlagCompletionFunc("config", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"yaml", "yml", "json", "toml"}, cobra.ShellCompDirectiveFilterFileExt
	})
}

// findCommand recursively finds a command by name
func findCommand(parent *cobra.Command, name string) *cobra.Command {
	for _, cmd := range parent.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

// GetPackageSuggestions returns package suggestions based on partial input
func GetPackageSuggestions(partial string) []string {
	// This could be enhanced to query pkg.go.dev in the future
	popularPackages := []string{
		"github.com/gin-gonic/gin",
		"github.com/labstack/echo/v4",
		"github.com/gofiber/fiber/v2",
		"github.com/go-chi/chi/v5",
		"gorm.io/gorm",
		"gorm.io/driver/postgres",
		"gorm.io/driver/mysql",
		"github.com/jmoiron/sqlx",
		"go.mongodb.org/mongo-driver",
		"github.com/redis/go-redis/v9",
		"github.com/stretchr/testify",
		"github.com/spf13/cobra",
		"github.com/spf13/viper",
		"github.com/rs/zerolog",
		"github.com/sirupsen/logrus",
		"github.com/golang-jwt/jwt/v5",
		"github.com/google/uuid",
		"github.com/gorilla/websocket",
		"google.golang.org/grpc",
		"google.golang.org/protobuf",
	}

	var suggestions []string
	for _, pkg := range popularPackages {
		if strings.Contains(pkg, partial) {
			suggestions = append(suggestions, pkg)
		}
	}

	return suggestions
}

// GetTemplateSuggestions returns template suggestions based on use case
func GetTemplateSuggestions(useCase string) []string {
	templatesByUseCase := map[string][]string{
		"api":          {"api", "microservice"},
		"web":          {"web", "api"},
		"cli":          {"cli", "minimal"},
		"worker":       {"worker", "minimal"},
		"library":      {"minimal"},
		"microservice": {"microservice", "grpc"},
		"grpc":         {"grpc", "microservice"},
	}

	if templates, ok := templatesByUseCase[useCase]; ok {
		return templates
	}

	return []string{"minimal", "api", "microservice", "cli"}
}
