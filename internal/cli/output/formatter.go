package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatYAML OutputFormat = "yaml"
)

// Formatter handles different output formats and styling
type Formatter struct {
	format    OutputFormat
	writer    io.Writer
	verbose   bool
	quiet     bool
	noColor   bool
	startTime time.Time
}

// NewFormatter creates a new output formatter
func NewFormatter(format OutputFormat, writer io.Writer) *Formatter {
	if writer == nil {
		writer = os.Stdout
	}
	return &Formatter{
		format:    format,
		writer:    writer,
		startTime: time.Now(),
	}
}

// SetOptions configures formatter options
func (f *Formatter) SetOptions(verbose, quiet, noColor bool) {
	f.verbose = verbose
	f.quiet = quiet
	f.noColor = noColor
}

// GenerationResult represents the result of a project generation
type GenerationResult struct {
	ProjectName    string                 `json:"project_name"`
	OutputPath     string                 `json:"output_path"`
	GoVersion      string                 `json:"go_version"`
	HTTPFramework  string                 `json:"http_framework,omitempty"`
	Databases      []string               `json:"databases,omitempty"`
	Drivers        []string               `json:"drivers,omitempty"`
	Features       []string               `json:"features,omitempty"`
	Packages       []string               `json:"packages,omitempty"`
	FilesGenerated int                    `json:"files_generated"`
	Duration       string                 `json:"duration"`
	Success        bool                   `json:"success"`
	Message        string                 `json:"message,omitempty"`
	NextSteps      []string               `json:"next_steps,omitempty"`
	Statistics     *GenerationStats       `json:"statistics,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// GenerationStats provides detailed statistics about the generation process
type GenerationStats struct {
	TotalFiles       int               `json:"total_files"`
	FilesByExtension map[string]int    `json:"files_by_extension"`
	FilesByCategory  map[string]int    `json:"files_by_category"`
	TotalLines       int               `json:"total_lines"`
	TotalSize        int64             `json:"total_size_bytes"`
	Dependencies     int               `json:"dependencies"`
	Templates        []string          `json:"templates_used"`
	ProcessingTime   map[string]string `json:"processing_time"`
}

// AnalysisResult represents the result of project analysis
type AnalysisResult struct {
	ProjectPath     string           `json:"project_path"`
	GoVersion       string           `json:"go_version"`
	ModulePath      string           `json:"module_path"`
	Dependencies    []DependencyInfo `json:"dependencies"`
	Architecture    ArchitectureInfo `json:"architecture"`
	Quality         QualityMetrics   `json:"quality"`
	Security        SecurityAnalysis `json:"security"`
	Recommendations []Recommendation `json:"recommendations"`
	Summary         string           `json:"summary"`
	Score           float64          `json:"score"`
	Timestamp       time.Time        `json:"timestamp"`
}

// DependencyInfo represents information about a dependency
type DependencyInfo struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Latest      string    `json:"latest,omitempty"`
	Category    string    `json:"category"`
	Security    string    `json:"security"`
	Outdated    bool      `json:"outdated"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
}

// ArchitectureInfo provides architectural analysis
type ArchitectureInfo struct {
	Pattern         string   `json:"pattern"`
	Layers          []string `json:"layers"`
	Complexity      string   `json:"complexity"`
	Maintainability float64  `json:"maintainability"`
}

// QualityMetrics provides code quality metrics
type QualityMetrics struct {
	Coverage        float64 `json:"coverage"`
	Complexity      float64 `json:"complexity"`
	Maintainability float64 `json:"maintainability"`
	TestFiles       int     `json:"test_files"`
	CodeLines       int     `json:"code_lines"`
	TestLines       int     `json:"test_lines"`
}

// SecurityAnalysis provides security analysis results
type SecurityAnalysis struct {
	Issues      []SecurityIssue `json:"issues"`
	Score       float64         `json:"score"`
	Severity    string          `json:"severity"`
	LastScanned time.Time       `json:"last_scanned"`
}

// SecurityIssue represents a security issue
type SecurityIssue struct {
	Package     string  `json:"package"`
	Severity    string  `json:"severity"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
	Solution    string  `json:"solution,omitempty"`
}

// Recommendation represents a recommendation
type Recommendation struct {
	Priority    int    `json:"priority"`
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action,omitempty"`
}

// OutputResult outputs a result in the configured format
func (f *Formatter) OutputResult(result interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.outputJSON(result)
	case FormatText:
		return f.outputText(result)
	default:
		return f.outputText(result)
	}
}

// outputJSON outputs result in JSON format
func (f *Formatter) outputJSON(result interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputText outputs result in human-readable format
func (f *Formatter) outputText(result interface{}) error {
	switch r := result.(type) {
	case *GenerationResult:
		return f.outputGenerationResult(r)
	case *AnalysisResult:
		return f.outputAnalysisResult(r)
	default:
		return fmt.Errorf("unsupported result type for text output")
	}
}

// outputGenerationResult outputs generation result in text format
func (f *Formatter) outputGenerationResult(result *GenerationResult) error {
	if f.quiet && !result.Success {
		fmt.Fprintf(f.writer, "Error: %s\n", result.Message)
		return nil
	}

	if f.quiet && result.Success {
		fmt.Fprintf(f.writer, "%s\n", result.OutputPath)
		return nil
	}

	// Header
	f.printHeader("🎉 Project Generation Complete!")

	// Project info
	f.printSection("Project Information")
	f.printKeyValue("Name", result.ProjectName)
	f.printKeyValue("Location", result.OutputPath)
	f.printKeyValue("Go Version", result.GoVersion)

	if result.HTTPFramework != "" {
		f.printKeyValue("HTTP Framework", result.HTTPFramework)
	}

	if len(result.Databases) > 0 {
		f.printKeyValue("Databases", strings.Join(result.Databases, ", "))
	}

	if len(result.Drivers) > 0 {
		f.printKeyValue("Drivers", strings.Join(result.Drivers, ", "))
	}

	if len(result.Features) > 0 {
		f.printKeyValue("Features", strings.Join(result.Features, ", "))
	}

	// Statistics
	if result.Statistics != nil {
		f.printSection("Generation Statistics")
		f.printKeyValue("Files Generated", fmt.Sprintf("%d", result.Statistics.TotalFiles))
		f.printKeyValue("Total Lines", fmt.Sprintf("%d", result.Statistics.TotalLines))
		f.printKeyValue("Dependencies", fmt.Sprintf("%d", result.Statistics.Dependencies))
		f.printKeyValue("Duration", result.Duration)

		if f.verbose && len(result.Statistics.FilesByExtension) > 0 {
			f.printSubSection("Files by Extension")
			for ext, count := range result.Statistics.FilesByExtension {
				f.printKeyValue(fmt.Sprintf("  %s", ext), fmt.Sprintf("%d", count))
			}
		}
	}

	// Next steps
	if len(result.NextSteps) > 0 {
		f.printSection("Next Steps")
		for i, step := range result.NextSteps {
			f.printListItem(fmt.Sprintf("%d. %s", i+1, step))
		}
	}

	f.printSeparator()
	return nil
}

// outputAnalysisResult outputs analysis result in text format
func (f *Formatter) outputAnalysisResult(result *AnalysisResult) error {
	if f.quiet {
		fmt.Fprintf(f.writer, "Score: %.1f/10\n", result.Score)
		return nil
	}

	// Header
	f.printHeader("📊 Project Analysis Results")

	// Project info
	f.printSection("Project Information")
	f.printKeyValue("Path", result.ProjectPath)
	f.printKeyValue("Module", result.ModulePath)
	f.printKeyValue("Go Version", result.GoVersion)
	f.printKeyValue("Overall Score", fmt.Sprintf("%.1f/10", result.Score))

	// Architecture
	f.printSection("Architecture")
	f.printKeyValue("Pattern", result.Architecture.Pattern)
	f.printKeyValue("Complexity", result.Architecture.Complexity)
	f.printKeyValue("Maintainability", fmt.Sprintf("%.1f/10", result.Architecture.Maintainability))

	// Quality metrics
	f.printSection("Quality Metrics")
	f.printKeyValue("Code Lines", fmt.Sprintf("%d", result.Quality.CodeLines))
	f.printKeyValue("Test Lines", fmt.Sprintf("%d", result.Quality.TestLines))
	f.printKeyValue("Test Coverage", fmt.Sprintf("%.1f%%", result.Quality.Coverage))

	// Security
	if len(result.Security.Issues) > 0 {
		f.printSection("Security Issues")
		for _, issue := range result.Security.Issues {
			severity := f.colorBySeverity(issue.Severity, issue.Severity)
			f.printKeyValue(issue.Package, fmt.Sprintf("%s - %s", severity, issue.Description))
		}
	}

	// Dependencies
	if len(result.Dependencies) > 0 && f.verbose {
		f.printSection("Dependencies")
		for _, dep := range result.Dependencies {
			status := "✅"
			if dep.Outdated {
				status = "⚠️"
			}
			f.printKeyValue(fmt.Sprintf("  %s %s", status, dep.Name), dep.Version)
		}
	}

	// Recommendations
	if len(result.Recommendations) > 0 {
		f.printSection("Recommendations")
		for _, rec := range result.Recommendations {
			priority := f.colorByPriority(rec.Priority, fmt.Sprintf("P%d", rec.Priority))
			f.printListItem(fmt.Sprintf("%s %s: %s", priority, rec.Title, rec.Description))
		}
	}

	f.printSeparator()
	return nil
}

// Helper methods for formatting

func (f *Formatter) printHeader(text string) {
	if f.noColor {
		fmt.Fprintf(f.writer, "\n%s\n%s\n\n", text, strings.Repeat("=", len(text)))
	} else {
		fmt.Fprintf(f.writer, "\n%s\n%s\n\n",
			color.HiCyanString(text),
			color.HiBlackString(strings.Repeat("=", len(text))))
	}
}

func (f *Formatter) printSection(title string) {
	if f.noColor {
		fmt.Fprintf(f.writer, "%s:\n", title)
	} else {
		fmt.Fprintf(f.writer, "%s:\n", color.HiYellowString(title))
	}
}

func (f *Formatter) printSubSection(title string) {
	if f.noColor {
		fmt.Fprintf(f.writer, "  %s:\n", title)
	} else {
		fmt.Fprintf(f.writer, "  %s:\n", color.YellowString(title))
	}
}

func (f *Formatter) printKeyValue(key, value string) {
	if f.noColor {
		fmt.Fprintf(f.writer, "  %-20s %s\n", key+":", value)
	} else {
		fmt.Fprintf(f.writer, "  %-20s %s\n",
			color.HiBlackString(key+":"),
			color.WhiteString(value))
	}
}

func (f *Formatter) printListItem(text string) {
	if f.noColor {
		fmt.Fprintf(f.writer, "  • %s\n", text)
	} else {
		fmt.Fprintf(f.writer, "  %s %s\n",
			color.HiBlackString("•"),
			color.WhiteString(text))
	}
}

func (f *Formatter) printSeparator() {
	if f.noColor {
		fmt.Fprintf(f.writer, "\n%s\n", strings.Repeat("-", 60))
	} else {
		fmt.Fprintf(f.writer, "\n%s\n", color.HiBlackString(strings.Repeat("-", 60)))
	}
}

func (f *Formatter) colorBySeverity(severity, text string) string {
	if f.noColor {
		return text
	}

	switch strings.ToLower(severity) {
	case "critical", "high":
		return color.HiRedString(text)
	case "medium", "moderate":
		return color.HiYellowString(text)
	case "low":
		return color.HiGreenString(text)
	default:
		return text
	}
}

func (f *Formatter) colorByPriority(priority int, text string) string {
	if f.noColor {
		return text
	}

	switch priority {
	case 5:
		return color.HiRedString(text)
	case 4:
		return color.RedString(text)
	case 3:
		return color.YellowString(text)
	case 2:
		return color.GreenString(text)
	case 1:
		return color.HiGreenString(text)
	default:
		return text
	}
}

// CreateProgressBar creates a styled progress bar
func (f *Formatter) CreateProgressBar(max int, description string) *progressbar.ProgressBar {
	if f.quiet || f.format == FormatJSON {
		return nil
	}

	bar := progressbar.NewOptions(max,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(f.writer),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)

	if !f.noColor {
		bar = progressbar.NewOptions(max,
			progressbar.OptionSetDescription(color.HiCyanString(description)),
			progressbar.OptionSetWriter(f.writer),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(50),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowIts(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        color.GreenString("█"),
				SaucerHead:    color.HiGreenString("█"),
				SaucerPadding: " ",
				BarStart:      color.HiBlackString("|"),
				BarEnd:        color.HiBlackString("|"),
			}),
		)
	}

	return bar
}

// PrintInfo prints an info message
func (f *Formatter) PrintInfo(message string) {
	if f.quiet || f.format == FormatJSON {
		return
	}

	if f.noColor {
		fmt.Fprintf(f.writer, "INFO: %s\n", message)
	} else {
		fmt.Fprintf(f.writer, "%s %s\n", color.HiCyanString("ℹ"), message)
	}
}

// PrintSuccess prints a success message
func (f *Formatter) PrintSuccess(message string) {
	if f.quiet || f.format == FormatJSON {
		return
	}

	if f.noColor {
		fmt.Fprintf(f.writer, "SUCCESS: %s\n", message)
	} else {
		fmt.Fprintf(f.writer, "%s %s\n", color.HiGreenString("✅"), message)
	}
}

// PrintWarning prints a warning message
func (f *Formatter) PrintWarning(message string) {
	if f.quiet || f.format == FormatJSON {
		return
	}

	if f.noColor {
		fmt.Fprintf(f.writer, "WARNING: %s\n", message)
	} else {
		fmt.Fprintf(f.writer, "%s %s\n", color.HiYellowString("⚠️"), message)
	}
}

// PrintError prints an error message
func (f *Formatter) PrintError(message string) {
	if f.format == FormatJSON {
		return
	}

	if f.noColor {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
	} else {
		fmt.Fprintf(os.Stderr, "%s %s\n", color.HiRedString("❌"), message)
	}
}

// PrintVerbose prints a verbose message
func (f *Formatter) PrintVerbose(message string) {
	if !f.verbose || f.quiet || f.format == FormatJSON {
		return
	}

	if f.noColor {
		fmt.Fprintf(f.writer, "VERBOSE: %s\n", message)
	} else {
		fmt.Fprintf(f.writer, "%s %s\n", color.HiBlackString("🔍"), message)
	}
}

// Duration returns the time elapsed since formatter creation
func (f *Formatter) Duration() time.Duration {
	return time.Since(f.startTime)
}
