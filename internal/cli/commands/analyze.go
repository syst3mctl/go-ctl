package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/syst3mctl/go-ctl/internal/cli/analyze"
	"github.com/syst3mctl/go-ctl/internal/cli/output"
	"github.com/syst3mctl/go-ctl/internal/cli/packages"
)

// NewAnalyzeCommand creates the analyze command
func NewAnalyzeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze [project-path]",
		Short: "Analyze Go project structure and provide insights",
		Long: `Analyze a Go project to understand its structure, dependencies, patterns,
and provide suggestions for improvements.

This command performs comprehensive analysis including:
  • Project structure and layout
  • Dependencies and their usage
  • Architectural patterns detection
  • Code quality metrics
  • Security analysis
  • Performance hints
  • Compatibility checks
  • Improvement suggestions

Examples:
  # Analyze current directory
  go-ctl analyze

  # Analyze specific project
  go-ctl analyze ./my-project

  # Generate detailed report
  go-ctl analyze ./my-project --detailed

  # Export analysis to file
  go-ctl analyze ./my-project --output=analysis.json

  # Focus on specific aspects
  go-ctl analyze --focus=security,dependencies

  # Include upgrade analysis
  go-ctl analyze --upgrade-check

  # JSON output for CI/CD
  go-ctl analyze --output-format=json
  go-ctl analyze ./my-project --focus=dependencies,security`,
		Args: cobra.MaximumNArgs(1),
		RunE: runAnalyze,
	}

	cmd.Flags().Bool("detailed", false, "Show detailed analysis report")
	cmd.Flags().StringP("output", "o", "", "Export analysis to file (JSON format)")
	cmd.Flags().StringSlice("focus", []string{}, "Focus on specific areas (dependencies, structure, patterns, metrics, security, compatibility)")
	cmd.Flags().Bool("suggestions", true, "Include improvement suggestions")
	cmd.Flags().Bool("issues", true, "Include issue detection")
	cmd.Flags().Int("timeout", 60, "Analysis timeout in seconds")
	cmd.Flags().Bool("json", false, "Output in JSON format")
	cmd.Flags().Bool("upgrade-check", false, "Include dependency upgrade analysis")

	return cmd
}

// runAnalyze executes the analyze command
func runAnalyze(cmd *cobra.Command, args []string) error {
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	// Get command flags
	detailed, _ := cmd.Flags().GetBool("detailed")
	outputFile, _ := cmd.Flags().GetString("output")
	focus, _ := cmd.Flags().GetStringSlice("focus")
	includeSuggestions, _ := cmd.Flags().GetBool("suggestions")
	includeIssues, _ := cmd.Flags().GetBool("issues")
	timeoutSecs, _ := cmd.Flags().GetInt("timeout")
	jsonOutput, _ := cmd.Flags().GetBool("json")
	upgradeCheck, _ := cmd.Flags().GetBool("upgrade-check")

	// Validate project path
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", projectPath)
	}

	printInfo("Analyzing project: %s", projectPath)

	// Create analyzer
	analyzer := analyze.NewProjectAnalyzer(projectPath)

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	// Perform analysis
	result, err := analyzer.Analyze(ctx)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Add upgrade analysis if requested
	var upgradeRecommendations []packages.UpgradeRecommendation
	if upgradeCheck {
		printInfo("Performing dependency upgrade analysis...")
		searcher := packages.NewPackageSearcher()
		goModPath := filepath.Join(projectPath, "go.mod")

		if _, err := os.Stat(goModPath); err == nil {
			analyses, err := searcher.AnalyzeDependencies(ctx, goModPath)
			if err == nil {
				upgradeRecommendations = searcher.GenerateUpgradeRecommendations(ctx, analyses)
				printSuccess("Found %d upgrade recommendations", len(upgradeRecommendations))
			} else {
				printWarning("Failed to analyze dependencies: %v", err)
			}
		}
	}

	// Filter results based on focus areas
	if len(focus) > 0 {
		filterAnalysisResults(result, focus)
	}

	// Export to file if requested
	if outputFile != "" {
		if err := exportAnalysis(result, outputFile); err != nil {
			return fmt.Errorf("failed to export analysis: %w", err)
		}
		printSuccess("Analysis exported to %s", outputFile)
	}

	// Display results
	if jsonOutput {
		return displayJSONResults(result)
	}

	return displayAnalysisResults(result, detailed, includeSuggestions, includeIssues, upgradeRecommendations)
}

// displayAnalysisResults displays the analysis results in a formatted way
func displayAnalysisResults(result *analyze.AnalysisResult, detailed, includeSuggestions, includeIssues bool, upgradeRecs []packages.UpgradeRecommendation) error {
	// Project Overview
	displayProjectOverview(&result.ProjectInfo)

	// Dependencies Analysis
	if len(result.Dependencies) > 0 {
		displayDependencies(result.Dependencies, detailed)
	}

	// Project Structure
	displayProjectStructure(&result.Structure, detailed)

	// Detected Patterns
	displayDetectedPatterns(&result.Patterns, detailed)

	// Metrics
	displayMetrics(&result.Metrics, detailed)

	// Suggestions
	if includeSuggestions && len(result.Suggestions) > 0 {
		displaySuggestions(result.Suggestions, detailed)
	}

	// Issues
	if includeIssues && len(result.Issues) > 0 {
		displayIssues(result.Issues, detailed)
	}

	// Upgrade Recommendations
	if len(upgradeRecs) > 0 {
		displayUpgradeRecommendations(upgradeRecs, detailed)
	}

	// Compatibility
	displayCompatibility(&result.Compatibility, detailed)

	// Summary
	displaySummary(result)

	return nil
}

// displayProjectOverview displays basic project information
func displayProjectOverview(info *analyze.ProjectInfo) {
	fmt.Printf("\n%s\n", color.HiCyanString("📋 Project Overview"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	fmt.Printf("%s %s\n", color.CyanString("Name:"), info.Name)
	if info.ModulePath != "" {
		fmt.Printf("%s %s\n", color.CyanString("Module:"), info.ModulePath)
	}
	fmt.Printf("%s %s\n", color.CyanString("Go Version:"), info.GoVersion)
	fmt.Printf("%s %s\n", color.CyanString("Type:"), strings.Title(info.ProjectType))

	if info.Description != "" {
		fmt.Printf("%s %s\n", color.CyanString("Description:"), info.Description)
	}

	if info.License != "" {
		fmt.Printf("%s %s\n", color.CyanString("License:"), info.License)
	}

	fmt.Printf("%s %d files, %d lines, %.1f MB\n",
		color.CyanString("Size:"),
		info.FileCount,
		info.LineCount,
		float64(info.Size)/1024/1024)
}

// displayDependencies displays dependency analysis
func displayDependencies(deps []analyze.Dependency, detailed bool) {
	fmt.Printf("\n%s\n", color.HiCyanString("📦 Dependencies"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Group by category
	categories := make(map[string][]analyze.Dependency)
	for _, dep := range deps {
		categories[dep.Category] = append(categories[dep.Category], dep)
	}

	for category, categoryDeps := range categories {
		if len(categoryDeps) == 0 {
			continue
		}

		fmt.Printf("\n%s %s (%d)\n",
			color.HiGreenString("●"),
			color.HiWhiteString(strings.Title(category)),
			len(categoryDeps))

		for _, dep := range categoryDeps {
			fmt.Printf("  %s %s", color.GreenString("•"), dep.Name)
			if dep.Version != "" {
				fmt.Printf(" %s", color.YellowString(dep.Version))
			}
			if detailed && dep.License != "" {
				fmt.Printf(" (%s)", dep.License)
			}
			if dep.Outdated {
				fmt.Printf(" %s", color.RedString("(outdated)"))
			}
			fmt.Println()
		}
	}

	fmt.Printf("\n%s Total: %d dependencies\n", color.CyanString("📊"), len(deps))
}

// displayProjectStructure displays project structure analysis
func displayProjectStructure(structure *analyze.ProjectStructure, detailed bool) {
	fmt.Printf("\n%s\n", color.HiCyanString("🏗️ Project Structure"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	fmt.Printf("%s %s\n", color.CyanString("Layout:"), strings.Title(structure.Layout))
	fmt.Printf("%s %d directories, %d files\n",
		color.CyanString("Structure:"),
		len(structure.Directories),
		len(structure.Files))

	if detailed {
		// Show key directories
		fmt.Printf("\n%s\n", color.HiWhiteString("Key Directories:"))
		for _, dir := range structure.Directories {
			if dir.Purpose != "other" {
				fmt.Printf("  %s %-15s - %s (%d files)\n",
					color.GreenString("•"),
					dir.Path,
					dir.Purpose,
					dir.FileCount)
			}
		}
	}

	// Show special files
	if len(structure.ConfigFiles) > 0 {
		fmt.Printf("\n%s %s\n",
			color.CyanString("Config Files:"),
			strings.Join(structure.ConfigFiles, ", "))
	}

	if len(structure.Scripts) > 0 {
		fmt.Printf("%s %s\n",
			color.CyanString("Scripts:"),
			strings.Join(structure.Scripts, ", "))
	}
}

// displayDetectedPatterns displays architectural patterns
func displayDetectedPatterns(patterns *analyze.DetectedPatterns, detailed bool) {
	fmt.Printf("\n%s\n", color.HiCyanString("🎯 Detected Patterns"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	fmt.Printf("%s %s\n", color.CyanString("Architecture:"), strings.Title(patterns.Architecture))

	// Web Framework
	if patterns.WebFramework.Name != "" {
		fmt.Printf("%s %s", color.CyanString("Web Framework:"), patterns.WebFramework.Name)
		if patterns.WebFramework.Version != "" {
			fmt.Printf(" %s", color.YellowString(patterns.WebFramework.Version))
		}
		fmt.Println()
	}

	// Databases
	if len(patterns.Database) > 0 {
		fmt.Printf("%s", color.CyanString("Databases:"))
		for _, db := range patterns.Database {
			fmt.Printf(" %s", db.Type)
			if db.Driver != "" {
				fmt.Printf("(%s)", db.Driver)
			}
		}
		fmt.Println()
	}

	// Testing
	if patterns.Testing.Framework != "" {
		fmt.Printf("%s %s", color.CyanString("Testing:"), patterns.Testing.Framework)
		if patterns.Testing.Coverage > 0 {
			fmt.Printf(" (%.1f%% coverage)", patterns.Testing.Coverage)
		}
		fmt.Println()
	}

	// Logging
	if patterns.Logging.Framework != "" {
		fmt.Printf("%s %s\n", color.CyanString("Logging:"), patterns.Logging.Framework)
	}

	// Containerization
	if patterns.Containerization.Docker {
		fmt.Printf("%s Docker", color.CyanString("Container:"))
		if patterns.Containerization.Compose {
			fmt.Printf(" + Compose")
		}
		if patterns.Containerization.Kubernetes {
			fmt.Printf(" + Kubernetes")
		}
		fmt.Println()
	}

	// CI/CD
	if patterns.CI.Platform != "" {
		fmt.Printf("%s %s\n", color.CyanString("CI/CD:"), patterns.CI.Platform)
	}
}

// displayMetrics displays project metrics
func displayMetrics(metrics *analyze.ProjectMetrics, detailed bool) {
	fmt.Printf("\n%s\n", color.HiCyanString("📊 Project Metrics"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Complexity
	complexity := getScoreColor(metrics.Complexity.Score)
	fmt.Printf("%s %s (%d/10)\n",
		color.CyanString("Complexity:"),
		complexity,
		metrics.Complexity.Score)

	// Quality
	quality := getScoreColor(metrics.Quality.Score)
	fmt.Printf("%s %s (%d/10)\n",
		color.CyanString("Quality:"),
		quality,
		metrics.Quality.Score)

	// Performance
	performance := getScoreColor(metrics.Performance.Score)
	fmt.Printf("%s %s (%d/10)\n",
		color.CyanString("Performance:"),
		performance,
		metrics.Performance.Score)

	// Security
	security := getScoreColor(metrics.Security.Score)
	fmt.Printf("%s %s (%d/10)\n",
		color.CyanString("Security:"),
		security,
		metrics.Security.Score)

	// Maintainability
	maintainability := getScoreColor(metrics.Maintainability.Score)
	fmt.Printf("%s %s (%d/10)\n",
		color.CyanString("Maintainability:"),
		maintainability,
		metrics.Maintainability.Score)

	if detailed {
		// Show detailed metrics
		fmt.Printf("\n%s\n", color.HiWhiteString("Detailed Metrics:"))

		if metrics.Quality.TestCoverage > 0 {
			fmt.Printf("  Test Coverage: %.1f%%\n", metrics.Quality.TestCoverage)
		}

		if metrics.Maintainability.Dependencies > 0 {
			fmt.Printf("  Dependencies: %d", metrics.Maintainability.Dependencies)
			if metrics.Maintainability.OutdatedDeps > 0 {
				fmt.Printf(" (%d outdated)", metrics.Maintainability.OutdatedDeps)
			}
			fmt.Println()
		}
	}
}

// displaySuggestions displays improvement suggestions
func displaySuggestions(suggestions []analyze.Suggestion, detailed bool) {
	fmt.Printf("\n%s\n", color.HiCyanString("💡 Suggestions"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Group by severity
	critical := []analyze.Suggestion{}
	warnings := []analyze.Suggestion{}
	info := []analyze.Suggestion{}

	for _, suggestion := range suggestions {
		switch suggestion.Severity {
		case "critical":
			critical = append(critical, suggestion)
		case "warning":
			warnings = append(warnings, suggestion)
		default:
			info = append(info, suggestion)
		}
	}

	// Display critical suggestions
	if len(critical) > 0 {
		fmt.Printf("\n%s Critical Issues (%d)\n", color.HiRedString("🔴"), len(critical))
		for _, s := range critical {
			fmt.Printf("  %s %s\n", color.RedString("•"), s.Title)
			if detailed {
				fmt.Printf("    %s\n", s.Description)
				fmt.Printf("    %s: %s\n", color.CyanString("Action"), s.Action)
			}
		}
	}

	// Display warnings
	if len(warnings) > 0 {
		fmt.Printf("\n%s Warnings (%d)\n", color.HiYellowString("🟡"), len(warnings))
		for _, s := range warnings {
			fmt.Printf("  %s %s\n", color.YellowString("•"), s.Title)
			if detailed {
				fmt.Printf("    %s\n", s.Description)
			}
		}
	}

	// Display info suggestions
	if len(info) > 0 && detailed {
		fmt.Printf("\n%s Recommendations (%d)\n", color.HiBlueString("🔵"), len(info))
		for _, s := range info {
			fmt.Printf("  %s %s\n", color.BlueString("•"), s.Title)
			fmt.Printf("    %s\n", s.Description)
		}
	}
}

// displayIssues displays detected issues
func displayIssues(issues []analyze.Issue, detailed bool) {
	fmt.Printf("\n%s\n", color.HiCyanString("⚠️ Issues"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Group by type
	errors := []analyze.Issue{}
	warnings := []analyze.Issue{}

	for _, issue := range issues {
		if issue.Type == "error" {
			errors = append(errors, issue)
		} else {
			warnings = append(warnings, issue)
		}
	}

	if len(errors) > 0 {
		fmt.Printf("\n%s Errors (%d)\n", color.HiRedString("❌"), len(errors))
		for _, issue := range errors {
			fmt.Printf("  %s %s:%d - %s\n",
				color.RedString("•"),
				issue.File,
				issue.Line,
				issue.Message)
		}
	}

	if len(warnings) > 0 {
		fmt.Printf("\n%s Warnings (%d)\n", color.HiYellowString("⚠️"), len(warnings))
		for _, issue := range warnings {
			fmt.Printf("  %s %s:%d - %s\n",
				color.YellowString("•"),
				issue.File,
				issue.Line,
				issue.Message)
		}
	}
}

// displayCompatibility displays compatibility information
func displayCompatibility(compat *analyze.CompatibilityInfo, detailed bool) {
	fmt.Printf("\n%s\n", color.HiCyanString("🔄 Compatibility"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	fmt.Printf("%s %s\n", color.CyanString("Go Version:"), compat.GoVersion)

	if len(compat.OSSupport) > 0 {
		fmt.Printf("%s %s\n",
			color.CyanString("OS Support:"),
			strings.Join(compat.OSSupport, ", "))
	}

	if len(compat.ArchSupport) > 0 {
		fmt.Printf("%s %s\n",
			color.CyanString("Architecture:"),
			strings.Join(compat.ArchSupport, ", "))
	}

	// Upgrade recommendations
	if compat.Upgradeability.Recommended != "" {
		fmt.Printf("%s %s\n",
			color.CyanString("Recommended Go:"),
			compat.Upgradeability.Recommended)
	}
}

// displaySummary displays analysis summary
func displaySummary(result *analyze.AnalysisResult) {
	fmt.Printf("\n%s\n", color.HiCyanString("📈 Analysis Summary"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Overall health score (simplified calculation)
	healthScore := (result.Metrics.Complexity.Score +
		result.Metrics.Quality.Score +
		result.Metrics.Performance.Score +
		result.Metrics.Security.Score +
		result.Metrics.Maintainability.Score) / 5

	healthColor := getScoreColor(healthScore)
	fmt.Printf("%s %s (%d/10)\n",
		color.CyanString("Overall Health:"),
		healthColor,
		healthScore)

	// Key statistics
	fmt.Printf("%s %d dependencies, %d files analyzed\n",
		color.CyanString("Analyzed:"),
		len(result.Dependencies),
		result.ProjectInfo.FileCount)

	if len(result.Suggestions) > 0 {
		fmt.Printf("%s %d suggestions for improvement\n",
			color.CyanString("Suggestions:"),
			len(result.Suggestions))
	}

	if len(result.Issues) > 0 {
		fmt.Printf("%s %d issues found\n",
			color.CyanString("Issues:"),
			len(result.Issues))
	}

	fmt.Printf("\n%s Analysis completed in %v\n",
		color.CyanString("⏱️"),
		time.Since(result.GeneratedAt).Truncate(time.Millisecond))
}

// displayJSONResults outputs results in JSON format
func displayJSONResults(result *analyze.AnalysisResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

// convertToAnalysisResult converts internal analysis result to output format
func convertToAnalysisResult(result *analyze.AnalysisResult, projectPath string) *output.AnalysisResult {
	analysisResult := &output.AnalysisResult{
		ProjectPath:  projectPath,
		GoVersion:    result.ProjectInfo.GoVersion,
		ModulePath:   result.ProjectInfo.ModulePath,
		Dependencies: make([]output.DependencyInfo, len(result.Dependencies)),
		Architecture: output.ArchitectureInfo{
			Pattern:         result.Patterns.Architecture,
			Layers:          []string{"cmd", "internal", "pkg"}, // Default layers
			Complexity:      "Medium",                           // Default value since not directly available
			Maintainability: 8.0,                                // Default value since not directly available
		},
		Quality: output.QualityMetrics{
			Coverage:        result.ProjectInfo.TestCoverage,
			Complexity:      float64(result.Metrics.Complexity.Cyclomatic),
			Maintainability: float64(result.Metrics.Quality.Score),
			TestFiles:       0, // Default value since Frameworks field doesn't exist
			CodeLines:       result.ProjectInfo.LineCount,
			TestLines:       0, // Not directly available
		},
		Security: output.SecurityAnalysis{
			Issues:      make([]output.SecurityIssue, 0),
			Score:       float64(result.Metrics.Security.Score),
			Severity:    "Low",
			LastScanned: time.Now(),
		},
		Recommendations: make([]output.Recommendation, len(result.Suggestions)),
		Summary:         fmt.Sprintf("Analysis of %s project", result.ProjectInfo.ProjectType),
		Score:           calculateOverallScore(result),
		Timestamp:       result.GeneratedAt,
	}

	// Convert dependencies
	for i, dep := range result.Dependencies {
		analysisResult.Dependencies[i] = output.DependencyInfo{
			Name:        dep.Name,
			Version:     dep.Version,
			Latest:      "", // Not directly available
			Category:    dep.Type,
			Security:    "Unknown", // Default value
			Outdated:    dep.Outdated,
			LastUpdated: time.Now(), // Default value
		}
	}

	// Convert suggestions to recommendations
	for i, suggestion := range result.Suggestions {
		priority := 3 // Default priority
		priority = suggestion.Priority
		if priority > 10 {
			priority = 10
		}
		if priority < 1 {
			priority = 1
		}

		analysisResult.Recommendations[i] = output.Recommendation{
			Priority:    priority,
			Category:    suggestion.Type,
			Title:       suggestion.Title,
			Description: suggestion.Description,
			Action:      suggestion.Action,
		}
	}

	return analysisResult
}

// calculateOverallScore calculates an overall score based on various metrics
func calculateOverallScore(result *analyze.AnalysisResult) float64 {
	score := 5.0 // Base score

	// Add points for test coverage
	if result.ProjectInfo.TestCoverage > 80 {
		score += 2.0
	} else if result.ProjectInfo.TestCoverage > 60 {
		score += 1.0
	}

	// Add points for good structure
	if len(result.Structure.Directories) >= 3 {
		score += 1.0
	}

	// Subtract points for issues
	highPriorityIssues := 0
	for _, issue := range result.Issues {
		if issue.Severity >= 8 { // High severity is 8-10
			highPriorityIssues++
		}
	}
	score -= float64(highPriorityIssues) * 0.5

	// Ensure score is between 0 and 10
	if score > 10 {
		score = 10
	}
	if score < 0 {
		score = 0
	}

	return score
}

// outputAnalysisResultsEnhanced outputs analysis results using enhanced formatter
func outputAnalysisResultsEnhanced(result *output.AnalysisResult, outputFile string, formatter *output.Formatter) error {
	if outputFile != "" {
		// Write to file
		var content []byte
		var err error

		if strings.HasSuffix(outputFile, ".json") {
			content, err = json.MarshalIndent(result, "", "  ")
		} else {
			// For non-JSON files, use JSON as default
			content, err = json.MarshalIndent(result, "", "  ")
		}

		if err != nil {
			return fmt.Errorf("failed to marshal results: %w", err)
		}

		if err := os.WriteFile(outputFile, content, 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		if !isQuiet() {
			formatter.PrintSuccess(fmt.Sprintf("Analysis results saved to: %s", outputFile))
		}
		return nil
	}

	// Output using formatter
	return formatter.OutputResult(result)
}

// exportAnalysis exports analysis results to a file
func exportAnalysis(result *analyze.AnalysisResult, filename string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// filterAnalysisResults filters results based on focus areas
func filterAnalysisResults(result *analyze.AnalysisResult, focus []string) {
	focusMap := make(map[string]bool)
	for _, f := range focus {
		focusMap[strings.ToLower(f)] = true
	}

	// Clear sections not in focus
	if !focusMap["dependencies"] {
		result.Dependencies = nil
	}
	if !focusMap["structure"] {
		result.Structure = analyze.ProjectStructure{}
	}
	if !focusMap["patterns"] {
		result.Patterns = analyze.DetectedPatterns{}
	}
	if !focusMap["metrics"] {
		result.Metrics = analyze.ProjectMetrics{}
	}
	if !focusMap["security"] {
		result.Issues = filterIssuesByCategory(result.Issues, "security")
		result.Suggestions = filterSuggestionsByType(result.Suggestions, "security")
	}
	if !focusMap["compatibility"] {
		result.Compatibility = analyze.CompatibilityInfo{}
	}
}

// filterIssuesByCategory filters issues by category
func filterIssuesByCategory(issues []analyze.Issue, category string) []analyze.Issue {
	var filtered []analyze.Issue
	for _, issue := range issues {
		if issue.Category == category {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// filterSuggestionsByType filters suggestions by type
func filterSuggestionsByType(suggestions []analyze.Suggestion, suggestionType string) []analyze.Suggestion {
	var filtered []analyze.Suggestion
	for _, suggestion := range suggestions {
		if suggestion.Type == suggestionType {
			filtered = append(filtered, suggestion)
		}
	}
	return filtered
}

// displayUpgradeRecommendations displays dependency upgrade recommendations
func displayUpgradeRecommendations(recommendations []packages.UpgradeRecommendation, detailed bool) {
	if len(recommendations) == 0 {
		return
	}

	fmt.Printf("\n%s\n", color.HiCyanString("🔄 Dependency Upgrade Recommendations"))
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	// Group recommendations by priority
	criticalRecs := []packages.UpgradeRecommendation{}
	highRecs := []packages.UpgradeRecommendation{}
	mediumRecs := []packages.UpgradeRecommendation{}
	lowRecs := []packages.UpgradeRecommendation{}

	for _, rec := range recommendations {
		switch {
		case rec.Priority >= 5:
			criticalRecs = append(criticalRecs, rec)
		case rec.Priority >= 4:
			highRecs = append(highRecs, rec)
		case rec.Priority >= 3:
			mediumRecs = append(mediumRecs, rec)
		default:
			lowRecs = append(lowRecs, rec)
		}
	}

	// Display by priority
	if len(criticalRecs) > 0 {
		fmt.Printf("\n%s Critical Priority (%d)\n", color.HiRedString("🚨"), len(criticalRecs))
		for _, rec := range criticalRecs {
			displaySingleUpgradeRecommendation(rec, detailed)
		}
	}

	if len(highRecs) > 0 {
		fmt.Printf("\n%s High Priority (%d)\n", color.HiYellowString("⚠️"), len(highRecs))
		for _, rec := range highRecs {
			displaySingleUpgradeRecommendation(rec, detailed)
		}
	}

	if len(mediumRecs) > 0 {
		fmt.Printf("\n%s Medium Priority (%d)\n", color.YellowString("💡"), len(mediumRecs))
		for _, rec := range mediumRecs {
			displaySingleUpgradeRecommendation(rec, detailed)
		}
	}

	if len(lowRecs) > 0 {
		fmt.Printf("\n%s Low Priority (%d)\n", color.GreenString("ℹ️"), len(lowRecs))
		for _, rec := range lowRecs {
			displaySingleUpgradeRecommendation(rec, detailed)
		}
	}

	// Summary
	fmt.Printf("\n%s\n", color.HiCyanString("Summary:"))
	fmt.Printf("  Total recommendations: %d\n", len(recommendations))
	if len(criticalRecs) > 0 {
		fmt.Printf("  Critical: %s (immediate attention required)\n", color.HiRedString("%d", len(criticalRecs)))
	}
	if len(highRecs) > 0 {
		fmt.Printf("  High: %s (should be addressed soon)\n", color.HiYellowString("%d", len(highRecs)))
	}
	if len(mediumRecs) > 0 {
		fmt.Printf("  Medium: %s (consider for next update cycle)\n", color.YellowString("%d", len(mediumRecs)))
	}
	if len(lowRecs) > 0 {
		fmt.Printf("  Low: %s (optional improvements)\n", color.GreenString("%d", len(lowRecs)))
	}

	fmt.Printf("\n%s Use %s to apply upgrades\n",
		color.CyanString("💡"), color.YellowString("go-ctl package upgrade"))
}

// displaySingleUpgradeRecommendation displays a single upgrade recommendation
func displaySingleUpgradeRecommendation(rec packages.UpgradeRecommendation, detailed bool) {
	priorityIcon := getPriorityIcon(rec.Priority)
	actionColor := getUpgradeActionColor(rec.Action)

	fmt.Printf("  %s %s\n", priorityIcon, color.HiWhiteString(rec.Package.Name))
	fmt.Printf("    %s: %s\n", color.CyanString("Action"), actionColor(rec.Action))

	if rec.NewVersion != "" {
		fmt.Printf("    %s: %s\n", color.CyanString("Version"), rec.NewVersion)
	}

	if rec.Alternative != "" {
		fmt.Printf("    %s: %s\n", color.CyanString("Alternative"), rec.Alternative)
	}

	fmt.Printf("    %s: %s\n", color.CyanString("Reason"), rec.Reason)

	if detailed {
		fmt.Printf("    %s: %s\n", color.CyanString("Package"), rec.Package.ImportPath)
		if rec.Package.Synopsis != "" {
			fmt.Printf("    %s: %s\n", color.CyanString("Description"), rec.Package.Synopsis)
		}
	}
}

// getPriorityIcon returns icon for priority level
func getPriorityIcon(priority int) string {
	switch {
	case priority >= 5:
		return color.HiRedString("🔴")
	case priority >= 4:
		return color.HiYellowString("🟡")
	case priority >= 3:
		return color.YellowString("🟠")
	case priority >= 2:
		return color.GreenString("🟢")
	default:
		return color.HiBlackString("⚪")
	}
}

// getUpgradeActionColor returns color function for action type
func getUpgradeActionColor(action string) func(...interface{}) string {
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

// getScoreColor returns colored score representation
func getScoreColor(score int) string {
	switch {
	case score >= 8:
		return color.HiGreenString("Excellent")
	case score >= 6:
		return color.GreenString("Good")
	case score >= 4:
		return color.YellowString("Fair")
	case score >= 2:
		return color.RedString("Poor")
	default:
		return color.HiRedString("Critical")
	}
}
