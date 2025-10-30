package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/syst3mctl/go-ctl/internal/cli/help"
	"github.com/syst3mctl/go-ctl/internal/cli/output"
	"github.com/syst3mctl/go-ctl/internal/cli/packages"
)

// NewPackageCommand creates the package management command
func NewPackageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Manage Go packages and dependencies",
		Long: `Search, validate, and manage Go packages for your projects.

This command provides utilities for discovering packages from pkg.go.dev,
validating package import paths, and getting package information.`,
		Aliases: []string{"pkg", "packages"},
	}

	// Add enhanced help
	help.AddEnhancedHelp(cmd)

	cmd.AddCommand(newPackageSearchCommand())
	cmd.AddCommand(newPackageInfoCommand())
	cmd.AddCommand(newPackageValidateCommand())
	cmd.AddCommand(newPackagePopularCommand())
	cmd.AddCommand(newPackageUpgradeCommand())

	return cmd
}

// newPackageSearchCommand creates the package search subcommand
func newPackageSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for Go packages",
		Long: `Search for Go packages using pkg.go.dev and other sources.

Examples:
  go-ctl package search http
  go-ctl package search web framework
  go-ctl package search --category=database postgres
  go-ctl package search --limit=5 --min-stars=100 logging`,
		Args: cobra.ExactArgs(1),
		RunE: runPackageSearch,
	}

	cmd.Flags().StringP("category", "c", "", "Search within category (web, database, testing, cli, logging, auth, validation, utils)")
	cmd.Flags().IntP("limit", "l", 10, "Maximum number of results")
	cmd.Flags().Int("min-stars", 0, "Minimum number of stars")
	cmd.Flags().StringP("license", "L", "", "Filter by license type")
	cmd.Flags().StringP("sort", "s", "relevance", "Sort by (relevance, stars, updated, name)")
	cmd.Flags().Bool("include-test", false, "Include test packages in results")

	return cmd
}

// newPackageInfoCommand creates the package info subcommand
func newPackageInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <import-path>",
		Short: "Get information about a specific package",
		Long: `Get detailed information about a Go package including version,
license, repository, and other metadata.

Examples:
  go-ctl package info github.com/gin-gonic/gin
  go-ctl package info gorm.io/gorm
  go-ctl package info github.com/stretchr/testify`,
		Args: cobra.ExactArgs(1),
		RunE: runPackageInfo,
	}

	return cmd
}

// newPackageValidateCommand creates the package validate subcommand
func newPackageValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate <import-path>...",
		Short: "Validate Go package import paths",
		Long: `Validate one or more Go package import paths to ensure they
are properly formatted and accessible.

Examples:
  go-ctl package validate github.com/gin-gonic/gin
  go-ctl package validate gorm.io/gorm github.com/stretchr/testify
  go-ctl package validate --file=packages.txt`,
		Args: cobra.MinimumNArgs(1),
		RunE: runPackageValidate,
	}

	cmd.Flags().StringP("file", "f", "", "Read import paths from file (one per line)")

	return cmd
}

// newPackagePopularCommand creates the package popular subcommand
func newPackagePopularCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "popular [category]",
		Short: "Show popular packages by category",
		Long: `Show popular Go packages organized by category.

Available categories:
  web        - HTTP frameworks and web utilities
  database   - Database drivers and ORMs
  testing    - Testing frameworks and utilities
  cli        - Command-line interface tools
  logging    - Logging libraries
  auth       - Authentication and authorization
  validation - Input validation and sanitization
  utils      - General utilities and helpers

Examples:
  go-ctl package popular
  go-ctl package popular web
  go-ctl package popular database`,
		Args: cobra.MaximumNArgs(1),
		RunE: runPackagePopular,
	}

	cmd.Flags().Bool("all", false, "Show all categories")
	cmd.Flags().Bool("detailed", false, "Show detailed information")

	return cmd
}

// newPackageUpgradeCommand creates the package upgrade subcommand
func newPackageUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [project-path]",
		Short: "Analyze and upgrade project dependencies",
		Long: `Analyze project dependencies for security issues, outdated packages,
and provide upgrade recommendations.

This command will:
  • Parse your go.mod file
  • Check for outdated dependencies
  • Identify security vulnerabilities
  • Suggest alternative packages
  • Provide upgrade recommendations

Examples:
  go-ctl package upgrade
  go-ctl package upgrade ./my-project
  go-ctl package upgrade --auto-apply
  go-ctl package upgrade --security-only`,
		Args: cobra.MaximumNArgs(1),
		RunE: runPackageUpgrade,
	}

	cmd.Flags().Bool("auto-apply", false, "Automatically apply safe upgrades")
	cmd.Flags().Bool("security-only", false, "Only show security-related upgrades")
	cmd.Flags().Bool("dry-run", false, "Show what would be changed without applying")
	cmd.Flags().StringP("output", "o", "", "Output format (table, json, yaml)")
	cmd.Flags().Int("max-risk", 3, "Maximum risk level to include (1-5, 5=critical)")

	return cmd
}

// runPackageSearch executes the package search command
func runPackageSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	searcher := packages.NewPackageSearcher()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get search options from flags
	category, _ := cmd.Flags().GetString("category")
	limit, _ := cmd.Flags().GetInt("limit")
	minStars, _ := cmd.Flags().GetInt("min-stars")
	license, _ := cmd.Flags().GetString("license")
	sortBy, _ := cmd.Flags().GetString("sort")
	includeTest, _ := cmd.Flags().GetBool("include-test")

	opts := packages.SearchOptions{
		Query:       query,
		Limit:       limit,
		MaxResults:  limit,
		MinStars:    minStars,
		License:     license,
		SortBy:      sortBy,
		IncludeTest: includeTest,
	}

	printInfo("Searching for packages matching '%s'...", query)

	var results []packages.Package

	if category != "" {
		// Search within specific category
		popular, popErr := searcher.SearchPopular(ctx, category)
		if popErr != nil {
			return fmt.Errorf("invalid category: %s", category)
		}

		// Filter by query
		for _, pkg := range popular {
			if strings.Contains(strings.ToLower(pkg.Name), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(pkg.Synopsis), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(pkg.ImportPath), strings.ToLower(query)) {
				results = append(results, pkg)
			}
		}
	} else {
		// General search
		searchResult, searchErr := searcher.Search(ctx, opts)
		if searchErr != nil {
			return fmt.Errorf("search failed: %w", searchErr)
		}
		results = searchResult.Packages
	}

	if len(results) == 0 {
		printWarning("No packages found matching '%s'", query)

		// Suggest alternative searches
		fmt.Printf("\n%s\n", color.CyanString("💡 Suggestions:"))
		fmt.Printf("  • Try broader search terms\n")
		fmt.Printf("  • Check spelling\n")
		fmt.Printf("  • Use category search: %s\n", color.YellowString("go-ctl package popular <category>"))
		return nil
	}

	// Display results
	PrintPackageList(results, fmt.Sprintf("Search Results for '%s'", query))

	if len(results) >= limit {
		fmt.Printf("\n%s Showing first %d results. Use %s to see more.\n",
			color.YellowString("ℹ️"), limit, color.CyanString("--limit=N"))
	}

	return nil
}

// runPackageInfo executes the package info command
func runPackageInfo(cmd *cobra.Command, args []string) error {
	importPath := args[0]

	searcher := packages.NewPackageSearcher()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	printInfo("Getting information for package '%s'...", importPath)

	pkg, err := searcher.GetPackageInfo(ctx, importPath)
	if err != nil {
		return fmt.Errorf("failed to get package info: %w", err)
	}

	// Display detailed package information
	fmt.Printf("\n%s\n", color.HiCyanString("Package Information"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	PrintPackage(*pkg)

	// Additional usage suggestions
	fmt.Printf("\n%s\n", color.HiGreenString("💡 Usage:"))
	fmt.Printf("  Add to go.mod: %s\n", color.YellowString(fmt.Sprintf("go get %s", pkg.ImportPath)))
	fmt.Printf("  Import in code: %s\n", color.YellowString(fmt.Sprintf("import \"%s\"", pkg.ImportPath)))

	return nil
}

// runPackageValidate executes the package validate command
func runPackageValidate(cmd *cobra.Command, args []string) error {
	importPaths := args

	// TODO: If file flag is specified, read from file
	file, _ := cmd.Flags().GetString("file")
	if file != "" {
		printWarning("File input not yet implemented. Using command line arguments.")
	}

	searcher := packages.NewPackageSearcher()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var validPaths []string
	var invalidPaths []string

	printInfo("Validating %d package(s)...", len(importPaths))

	for _, importPath := range importPaths {
		printVerbose("Validating %s", importPath)

		err := searcher.ValidatePackage(ctx, importPath)
		if err != nil {
			invalidPaths = append(invalidPaths, fmt.Sprintf("%s: %v", importPath, err))
		} else {
			validPaths = append(validPaths, importPath)
		}
	}

	// Display results
	fmt.Printf("\n%s\n", color.HiCyanString("Validation Results"))
	fmt.Printf("%s\n", strings.Repeat("=", 30))

	if len(validPaths) > 0 {
		fmt.Printf("\n%s Valid Packages (%d):\n", color.HiGreenString("✅"), len(validPaths))
		for _, path := range validPaths {
			fmt.Printf("  %s %s\n", color.GreenString("●"), path)
		}
	}

	if len(invalidPaths) > 0 {
		fmt.Printf("\n%s Invalid Packages (%d):\n", color.HiRedString("❌"), len(invalidPaths))
		for _, pathWithError := range invalidPaths {
			fmt.Printf("  %s %s\n", color.RedString("●"), pathWithError)
		}
	}

	// Summary
	fmt.Printf("\n%s\n", color.HiCyanString("Summary:"))
	fmt.Printf("  Valid: %s\n", color.GreenString("%d", len(validPaths)))
	fmt.Printf("  Invalid: %s\n", color.RedString("%d", len(invalidPaths)))
	fmt.Printf("  Total: %d\n", len(importPaths))

	if len(invalidPaths) > 0 {
		return fmt.Errorf("validation failed for %d package(s)", len(invalidPaths))
	}

	printSuccess("All packages are valid!")
	return nil
}

// runPackagePopular executes the package popular command
// PrintPackage prints package information in a formatted way
func PrintPackage(pkg packages.Package) {
	fmt.Printf("%s %s\n", color.HiGreenString("●"), color.HiWhiteString(pkg.Name))
	fmt.Printf("  %s: %s\n", color.CyanString("Import"), pkg.ImportPath)
	fmt.Printf("  %s: %s\n", color.CyanString("Description"), pkg.Synopsis)
	if pkg.Version != "" && pkg.Version != "latest" {
		fmt.Printf("  %s: %s\n", color.CyanString("Version"), pkg.Version)
	}
	if pkg.License != "" {
		fmt.Printf("  %s: %s\n", color.CyanString("License"), pkg.License)
	}
	if pkg.Stars > 0 {
		fmt.Printf("  %s: %d\n", color.CyanString("Stars"), pkg.Stars)
	}
}

// PrintPackageList prints a list of packages in a formatted way
func PrintPackageList(packages []packages.Package, title string) {
	if len(packages) == 0 {
		fmt.Printf("%s No packages found\n", color.YellowString("⚠️"))
		return
	}

	if title != "" {
		fmt.Printf("%s\n", color.HiCyanString(title))
		fmt.Printf("%s\n", strings.Repeat("=", len(title)))
	}

	for i, pkg := range packages {
		if i > 0 {
			fmt.Println()
		}
		PrintPackage(pkg)
	}
}

func runPackagePopular(cmd *cobra.Command, args []string) error {
	showAll, _ := cmd.Flags().GetBool("all")
	detailed, _ := cmd.Flags().GetBool("detailed")

	searcher := packages.NewPackageSearcher()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	categories := []string{"web", "database", "testing", "cli", "logging", "auth", "validation", "utils"}

	if len(args) > 0 {
		// Show specific category
		category := args[0]
		packages, err := searcher.SearchPopular(ctx, category)
		if err != nil {
			return fmt.Errorf("invalid category '%s'. Available: %s", category, strings.Join(categories, ", "))
		}

		title := fmt.Sprintf("Popular %s Packages", strings.Title(category))
		PrintPackageList(packages, title)
		return nil
	}

	if showAll {
		// Show all categories
		for i, category := range categories {
			if i > 0 {
				fmt.Println()
			}

			packages, err := searcher.SearchPopular(ctx, category)
			if err != nil {
				continue
			}

			if detailed {
				title := fmt.Sprintf("Popular %s Packages", strings.Title(category))
				PrintPackageList(packages, title)
			} else {
				fmt.Printf("%s %s\n", color.HiGreenString("●"), color.HiWhiteString(strings.Title(category)))
				for _, pkg := range packages {
					fmt.Printf("  %s - %s\n", color.CyanString(pkg.Name), pkg.Synopsis)
				}
			}
		}
		return nil
	}

	// Show category overview
	fmt.Printf("%s\n", color.HiCyanString("Popular Package Categories"))
	fmt.Printf("%s\n", strings.Repeat("=", 40))

	categoryDescriptions := map[string]string{
		"web":        "HTTP frameworks and web utilities",
		"database":   "Database drivers and ORMs",
		"testing":    "Testing frameworks and utilities",
		"cli":        "Command-line interface tools",
		"logging":    "Logging libraries",
		"auth":       "Authentication and authorization",
		"validation": "Input validation and sanitization",
		"utils":      "General utilities and helpers",
	}

	for _, category := range categories {
		fmt.Printf("\n%s %s\n", color.HiGreenString("●"), color.HiWhiteString(strings.Title(category)))
		if desc, exists := categoryDescriptions[category]; exists {
			fmt.Printf("  %s\n", desc)
		}
		fmt.Printf("  %s: %s\n", color.CyanString("View"), color.YellowString(fmt.Sprintf("go-ctl package popular %s", category)))
	}

	fmt.Printf("\n%s\n", color.HiCyanString("Usage:"))
	fmt.Printf("  • View category: %s\n", color.YellowString("go-ctl package popular <category>"))
	fmt.Printf("  • Show all: %s\n", color.YellowString("go-ctl package popular --all"))
	fmt.Printf("  • Detailed view: %s\n", color.YellowString("go-ctl package popular <category> --detailed"))

	return nil
}

// runPackageUpgrade executes the package upgrade command
func runPackageUpgrade(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	// Get flags
	autoApply, _ := cmd.Flags().GetBool("auto-apply")
	securityOnly, _ := cmd.Flags().GetBool("security-only")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	outputFormat, _ := cmd.Flags().GetString("output")
	maxRisk, _ := cmd.Flags().GetInt("max-risk")

	goModPath := fmt.Sprintf("%s/go.mod", projectPath)

	searcher := packages.NewPackageSearcher()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	printInfo("Analyzing dependencies in %s...", projectPath)

	// Analyze dependencies
	analyses, err := searcher.AnalyzeDependencies(ctx, goModPath)
	if err != nil {
		return fmt.Errorf("failed to analyze dependencies: %w", err)
	}

	if len(analyses) == 0 {
		printWarning("No dependencies found in go.mod")
		return nil
	}

	printSuccess("Found %d dependencies to analyze", len(analyses))

	// Generate recommendations
	recommendations := searcher.GenerateUpgradeRecommendations(ctx, analyses)

	// Filter recommendations based on flags
	if securityOnly {
		recommendations = filterSecurityRecommendations(recommendations)
	}

	recommendations = filterByRiskLevel(recommendations, maxRisk)

	if len(recommendations) == 0 {
		printSuccess("All dependencies are up to date and secure!")
		return nil
	}

	// Display results based on output format
	switch outputFormat {
	case "json":
		return displayUpgradeResultsJSON(recommendations)
	case "yaml":
		return displayUpgradeResultsYAML(recommendations)
	default:
		displayUpgradeResults(analyses, recommendations, securityOnly)
	}

	// Apply upgrades if requested
	if autoApply && !dryRun {
		printInfo("Applying automatic upgrades...")
		err := searcher.UpdateDependencies(ctx, goModPath, recommendations, true)
		if err != nil {
			return fmt.Errorf("failed to apply upgrades: %w", err)
		}
		printSuccess("Upgrades applied successfully!")
	} else if dryRun {
		printInfo("Dry run mode - no changes will be applied")
	} else {
		fmt.Printf("\n%s\n", color.HiCyanString("💡 To apply upgrades:"))
		fmt.Printf("  • Safe upgrades: %s\n", color.YellowString("go-ctl package upgrade --auto-apply"))
		fmt.Printf("  • Manual review: Review recommendations above and update go.mod manually\n")
	}

	return nil
}

// Helper functions for upgrade command

func filterSecurityRecommendations(recommendations []packages.UpgradeRecommendation) []packages.UpgradeRecommendation {
	var filtered []packages.UpgradeRecommendation
	for _, rec := range recommendations {
		if strings.Contains(strings.ToLower(rec.Reason), "security") ||
			strings.Contains(strings.ToLower(rec.Reason), "vulnerability") {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}

func filterByRiskLevel(recommendations []packages.UpgradeRecommendation, maxRisk int) []packages.UpgradeRecommendation {
	var filtered []packages.UpgradeRecommendation
	for _, rec := range recommendations {
		if rec.Priority <= maxRisk {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}

func displayUpgradeResults(analyses []packages.DependencyAnalysis, recommendations []packages.UpgradeRecommendation, securityOnly bool) {
	fmt.Printf("\n%s\n", color.HiCyanString("📦 Dependency Analysis Results"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Show summary
	var outdated, vulnerable, total int
	for _, analysis := range analyses {
		total++
		if analysis.IsOutdated {
			outdated++
		}
		if len(analysis.SecurityIssues) > 0 {
			vulnerable++
		}
	}

	fmt.Printf("\n%s\n", color.HiCyanString("📊 Summary:"))
	fmt.Printf("  Total dependencies: %d\n", total)
	fmt.Printf("  Outdated: %s\n", getCountColor(outdated, 0))
	fmt.Printf("  With security issues: %s\n", getCountColor(vulnerable, 0))
	fmt.Printf("  Recommendations: %s\n", getCountColor(len(recommendations), 0))

	if len(recommendations) == 0 {
		return
	}

	// Show recommendations
	fmt.Printf("\n%s\n", color.HiCyanString("🔧 Upgrade Recommendations:"))

	for i, rec := range recommendations {
		if i > 0 {
			fmt.Println()
		}

		priorityColor := getPriorityColor(rec.Priority)
		actionColor := getActionColor(rec.Action)

		fmt.Printf("%s %s\n", priorityColor("●"), color.HiWhiteString(rec.Package.Name))
		fmt.Printf("  %s: %s\n", color.CyanString("Package"), rec.Package.ImportPath)
		fmt.Printf("  %s: %s\n", color.CyanString("Action"), actionColor(rec.Action))
		if rec.NewVersion != "" {
			fmt.Printf("  %s: %s\n", color.CyanString("New Version"), rec.NewVersion)
		}
		if rec.Alternative != "" {
			fmt.Printf("  %s: %s\n", color.CyanString("Alternative"), rec.Alternative)
		}
		fmt.Printf("  %s: %s\n", color.CyanString("Reason"), rec.Reason)
		fmt.Printf("  %s: %s\n", color.CyanString("Priority"), priorityColor(getPriorityText(rec.Priority)))
	}
}

func displayUpgradeResultsJSON(recommendations []packages.UpgradeRecommendation) error {
	// TODO: Implement JSON output
	fmt.Println("JSON output not yet implemented")
	return nil
}

func displayUpgradeResultsYAML(recommendations []packages.UpgradeRecommendation) error {
	// TODO: Implement YAML output
	fmt.Println("YAML output not yet implemented")
	return nil
}

// outputPackageSearchJSON outputs package search results in JSON format
func outputPackageSearchJSON(results []packages.Package) error {
	type JSONPackageResult struct {
		Name        string `json:"name"`
		ImportPath  string `json:"import_path"`
		Description string `json:"description"`
		Stars       int    `json:"stars"`
		Version     string `json:"version"`
		Category    string `json:"category"`
		Updated     string `json:"last_updated"`
		License     string `json:"license"`
		URL         string `json:"url"`
	}

	var jsonResults []JSONPackageResult
	for _, result := range results {
		jsonResults = append(jsonResults, JSONPackageResult{
			Name:        result.Name,
			ImportPath:  result.ImportPath,
			Description: result.Synopsis,
			Stars:       result.Stars,
			Version:     result.Version,
			Category:    "", // Not available in Package struct
			Updated:     result.LastUpdated,
			License:     result.License,
			URL:         result.Repository,
		})
	}

	response := map[string]interface{}{
		"results": jsonResults,
		"count":   len(jsonResults),
		"query":   "search query", // This would need to be passed in
	}

	formatter := output.NewFormatter(output.FormatJSON, nil)
	return formatter.OutputResult(response)
}

func getCountColor(count, threshold int) string {
	if count == 0 {
		return color.GreenString("%d", count)
	} else if count <= threshold {
		return color.YellowString("%d", count)
	}
	return color.RedString("%d", count)
}

func getPriorityColor(priority int) func(...interface{}) string {
	switch {
	case priority >= 5:
		return func(a ...interface{}) string { return color.HiRedString("%v", a...) }
	case priority >= 3:
		return func(a ...interface{}) string { return color.HiYellowString("%v", a...) }
	case priority >= 2:
		return func(a ...interface{}) string { return color.YellowString("%v", a...) }
	default:
		return func(a ...interface{}) string { return color.GreenString("%v", a...) }
	}
}

func getActionColor(action string) func(...interface{}) string {
	switch action {
	case "update":
		return func(a ...interface{}) string { return color.GreenString("%v", a...) }
	case "replace":
		return func(a ...interface{}) string { return color.YellowString("%v", a...) }
	case "remove":
		return func(a ...interface{}) string { return color.RedString("%v", a...) }
	default:
		return func(a ...interface{}) string { return color.WhiteString("%v", a...) }
	}
}

func getPriorityText(priority int) string {
	switch {
	case priority >= 5:
		return "Critical"
	case priority >= 4:
		return "High"
	case priority >= 3:
		return "Medium"
	case priority >= 2:
		return "Low"
	default:
		return "Info"
	}
}
