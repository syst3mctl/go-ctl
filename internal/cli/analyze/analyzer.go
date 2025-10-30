package analyze

import (
	"bufio"
	"context"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// ProjectAnalyzer analyzes Go projects for structure, dependencies, and patterns
type ProjectAnalyzer struct {
	projectPath string
	fileSet     *token.FileSet
}

// AnalysisResult contains the complete analysis of a project
type AnalysisResult struct {
	ProjectInfo   ProjectInfo       `json:"project_info"`
	Dependencies  []Dependency      `json:"dependencies"`
	Structure     ProjectStructure  `json:"structure"`
	Patterns      DetectedPatterns  `json:"patterns"`
	Metrics       ProjectMetrics    `json:"metrics"`
	Suggestions   []Suggestion      `json:"suggestions"`
	Issues        []Issue           `json:"issues"`
	Compatibility CompatibilityInfo `json:"compatibility"`
	GeneratedAt   time.Time         `json:"generated_at"`
}

// ProjectInfo contains basic project information
type ProjectInfo struct {
	Name         string            `json:"name"`
	ModulePath   string            `json:"module_path"`
	GoVersion    string            `json:"go_version"`
	ProjectType  string            `json:"project_type"` // web, cli, library, microservice, worker
	Description  string            `json:"description"`
	License      string            `json:"license"`
	Repository   string            `json:"repository"`
	LastModified time.Time         `json:"last_modified"`
	Size         int64             `json:"size_bytes"`
	FileCount    int               `json:"file_count"`
	LineCount    int               `json:"line_count"`
	TestCoverage float64           `json:"test_coverage"`
	Environment  map[string]string `json:"environment"`
}

// Dependency represents a project dependency
type Dependency struct {
	Name          string       `json:"name"`
	Version       string       `json:"version"`
	Type          string       `json:"type"`     // direct, indirect, test, build
	Category      string       `json:"category"` // web, database, testing, etc.
	License       string       `json:"license"`
	Vulnerability []string     `json:"vulnerability"`
	Usage         []string     `json:"usage"` // where it's imported
	Outdated      bool         `json:"outdated"`
	Alternative   []string     `json:"alternative"`
	Security      SecurityInfo `json:"security"`
}

// SecurityInfo contains security information about dependencies
type SecurityInfo struct {
	Score       int      `json:"score"` // 1-10
	Issues      []string `json:"issues"`
	LastUpdated string   `json:"last_updated"`
	Maintainers int      `json:"maintainers"`
	Downloads   int64    `json:"downloads"`
}

// ProjectStructure describes the project's directory structure
type ProjectStructure struct {
	RootPath      string          `json:"root_path"`
	Layout        string          `json:"layout"` // standard, flat, domain-driven, etc.
	Directories   []DirectoryInfo `json:"directories"`
	Files         []FileInfo      `json:"files"`
	ConfigFiles   []string        `json:"config_files"`
	Scripts       []string        `json:"scripts"`
	Documentation []string        `json:"documentation"`
}

// DirectoryInfo contains information about a directory
type DirectoryInfo struct {
	Path      string `json:"path"`
	Purpose   string `json:"purpose"` // cmd, internal, pkg, test, etc.
	FileCount int    `json:"file_count"`
	LineCount int    `json:"line_count"`
	TestCount int    `json:"test_count"`
}

// FileInfo contains information about a file
type FileInfo struct {
	Path         string    `json:"path"`
	Type         string    `json:"type"` // go, test, config, doc, script
	Size         int64     `json:"size"`
	LineCount    int       `json:"line_count"`
	Purpose      string    `json:"purpose"`
	LastModified time.Time `json:"last_modified"`
	Complexity   int       `json:"complexity"` // cyclomatic complexity
}

// DetectedPatterns contains detected architectural patterns
type DetectedPatterns struct {
	Architecture     string                 `json:"architecture"` // clean, hexagonal, layered, etc.
	Patterns         []ArchitecturalPattern `json:"patterns"`
	WebFramework     FrameworkInfo          `json:"web_framework"`
	Database         []DatabaseInfo         `json:"database"`
	Testing          TestingInfo            `json:"testing"`
	Logging          LoggingInfo            `json:"logging"`
	Configuration    ConfigurationInfo      `json:"configuration"`
	Containerization ContainerInfo          `json:"containerization"`
	CI               CIInfo                 `json:"ci_cd"`
}

// ArchitecturalPattern represents a detected pattern
type ArchitecturalPattern struct {
	Name        string   `json:"name"`
	Confidence  float64  `json:"confidence"` // 0-1
	Evidence    []string `json:"evidence"`
	Description string   `json:"description"`
}

// FrameworkInfo contains web framework information
type FrameworkInfo struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Router     string   `json:"router"`
	Middleware []string `json:"middleware"`
	Features   []string `json:"features"`
}

// DatabaseInfo contains database information
type DatabaseInfo struct {
	Type       string   `json:"type"`   // postgres, mysql, mongodb, etc.
	Driver     string   `json:"driver"` // gorm, sqlx, mongo-driver, etc.
	Migrations bool     `json:"migrations"`
	Seeds      bool     `json:"seeds"`
	Models     []string `json:"models"`
}

// TestingInfo contains testing information
type TestingInfo struct {
	Framework  string   `json:"framework"` // testify, ginkgo, etc.
	Coverage   float64  `json:"coverage"`
	TestFiles  int      `json:"test_files"`
	TestCount  int      `json:"test_count"`
	Benchmarks int      `json:"benchmarks"`
	Examples   int      `json:"examples"`
	Mocking    []string `json:"mocking"`
}

// LoggingInfo contains logging information
type LoggingInfo struct {
	Framework  string   `json:"framework"` // logrus, zap, zerolog, etc.
	Levels     []string `json:"levels"`
	Structured bool     `json:"structured"`
	Format     string   `json:"format"` // json, text, etc.
}

// ConfigurationInfo contains configuration information
type ConfigurationInfo struct {
	Method  string   `json:"method"` // viper, env, flags, etc.
	Files   []string `json:"files"`
	Formats []string `json:"formats"` // yaml, json, toml, etc.
	Sources []string `json:"sources"` // env, file, flag, etc.
}

// ContainerInfo contains containerization information
type ContainerInfo struct {
	Docker     bool     `json:"docker"`
	Compose    bool     `json:"compose"`
	Kubernetes bool     `json:"kubernetes"`
	Registry   string   `json:"registry"`
	Images     []string `json:"images"`
	Volumes    []string `json:"volumes"`
	Networks   []string `json:"networks"`
}

// CIInfo contains CI/CD information
type CIInfo struct {
	Platform   string   `json:"platform"` // github, gitlab, jenkins, etc.
	Workflows  []string `json:"workflows"`
	Triggers   []string `json:"triggers"`
	Stages     []string `json:"stages"`
	Deployment bool     `json:"deployment"`
}

// ProjectMetrics contains project metrics
type ProjectMetrics struct {
	Complexity      ComplexityMetrics  `json:"complexity"`
	Quality         QualityMetrics     `json:"quality"`
	Performance     PerformanceHints   `json:"performance"`
	Security        SecurityMetrics    `json:"security"`
	Maintainability MaintenanceMetrics `json:"maintainability"`
}

// ComplexityMetrics contains complexity measurements
type ComplexityMetrics struct {
	Cyclomatic   int `json:"cyclomatic"`
	Cognitive    int `json:"cognitive"`
	Depth        int `json:"depth"`
	FileSize     int `json:"file_size"`
	FunctionSize int `json:"function_size"`
	Score        int `json:"score"` // 1-10
}

// QualityMetrics contains quality measurements
type QualityMetrics struct {
	TestCoverage      float64 `json:"test_coverage"`
	Documentation     float64 `json:"documentation"`
	CodeDuplication   float64 `json:"code_duplication"`
	NamingConsistency float64 `json:"naming_consistency"`
	ErrorHandling     float64 `json:"error_handling"`
	Score             int     `json:"score"` // 1-10
}

// PerformanceHints contains performance analysis
type PerformanceHints struct {
	Bottlenecks   []string `json:"bottlenecks"`
	MemoryLeaks   []string `json:"memory_leaks"`
	Optimizations []string `json:"optimizations"`
	Goroutines    int      `json:"goroutines"`
	Channels      int      `json:"channels"`
	Score         int      `json:"score"` // 1-10
}

// SecurityMetrics contains security analysis
type SecurityMetrics struct {
	Vulnerabilities []string `json:"vulnerabilities"`
	Credentials     []string `json:"credentials"`
	Encryption      bool     `json:"encryption"`
	InputValidation bool     `json:"input_validation"`
	Logging         bool     `json:"logging"`
	Score           int      `json:"score"` // 1-10
}

// MaintenanceMetrics contains maintainability measurements
type MaintenanceMetrics struct {
	Dependencies   int      `json:"dependencies"`
	OutdatedDeps   int      `json:"outdated_deps"`
	TechnicalDebt  []string `json:"technical_debt"`
	RefactoringOps []string `json:"refactoring_opportunities"`
	Documentation  float64  `json:"documentation"`
	Score          int      `json:"score"` // 1-10
}

// Suggestion represents an improvement suggestion
type Suggestion struct {
	Type        string `json:"type"`     // dependency, structure, pattern, security, performance
	Severity    string `json:"severity"` // info, warning, error, critical
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`
	Before      string `json:"before,omitempty"`
	After       string `json:"after,omitempty"`
	Impact      string `json:"impact"`
	Effort      string `json:"effort"`   // low, medium, high
	Priority    int    `json:"priority"` // 1-10
}

// Issue represents a problem found in the project
type Issue struct {
	Type     string `json:"type"`     // error, warning, info
	Category string `json:"category"` // syntax, logic, security, performance
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Severity int    `json:"severity"` // 1-10
	Fixable  bool   `json:"fixable"`
}

// CompatibilityInfo contains compatibility information
type CompatibilityInfo struct {
	GoVersion      string             `json:"go_version"`
	OSSupport      []string           `json:"os_support"`
	ArchSupport    []string           `json:"arch_support"`
	Dependencies   []DependencyCompat `json:"dependencies"`
	Deployment     []DeploymentCompat `json:"deployment"`
	Upgradeability UpgradeInfo        `json:"upgradeability"`
}

// DependencyCompat contains dependency compatibility info
type DependencyCompat struct {
	Name            string   `json:"name"`
	CurrentVer      string   `json:"current_version"`
	LatestVer       string   `json:"latest_version"`
	Compatible      bool     `json:"compatible"`
	Issues          []string `json:"issues"`
	BreakingChanges []string `json:"breaking_changes"`
}

// DeploymentCompat contains deployment compatibility info
type DeploymentCompat struct {
	Platform     string   `json:"platform"`
	Supported    bool     `json:"supported"`
	Issues       []string `json:"issues"`
	Requirements []string `json:"requirements"`
}

// UpgradeInfo contains upgrade information
type UpgradeInfo struct {
	GoVersions  []string `json:"go_versions"`
	Blockers    []string `json:"blockers"`
	Effort      string   `json:"effort"` // low, medium, high
	Recommended string   `json:"recommended"`
}

// NewProjectAnalyzer creates a new project analyzer
func NewProjectAnalyzer(projectPath string) *ProjectAnalyzer {
	return &ProjectAnalyzer{
		projectPath: projectPath,
		fileSet:     token.NewFileSet(),
	}
}

// Analyze performs a comprehensive analysis of the project
func (pa *ProjectAnalyzer) Analyze(ctx context.Context) (*AnalysisResult, error) {
	if _, err := os.Stat(pa.projectPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("project path does not exist: %s", pa.projectPath)
	}

	result := &AnalysisResult{
		GeneratedAt: time.Now(),
	}

	// Analyze project info
	if err := pa.analyzeProjectInfo(ctx, &result.ProjectInfo); err != nil {
		return nil, fmt.Errorf("failed to analyze project info: %w", err)
	}

	// Analyze dependencies
	if err := pa.analyzeDependencies(ctx, &result.Dependencies); err != nil {
		return nil, fmt.Errorf("failed to analyze dependencies: %w", err)
	}

	// Analyze project structure
	if err := pa.analyzeStructure(ctx, &result.Structure); err != nil {
		return nil, fmt.Errorf("failed to analyze structure: %w", err)
	}

	// Detect patterns
	if err := pa.detectPatterns(ctx, &result.Patterns); err != nil {
		return nil, fmt.Errorf("failed to detect patterns: %w", err)
	}

	// Calculate metrics
	if err := pa.calculateMetrics(ctx, &result.Metrics); err != nil {
		return nil, fmt.Errorf("failed to calculate metrics: %w", err)
	}

	// Generate suggestions
	if err := pa.generateSuggestions(ctx, result, &result.Suggestions); err != nil {
		return nil, fmt.Errorf("failed to generate suggestions: %w", err)
	}

	// Find issues
	if err := pa.findIssues(ctx, &result.Issues); err != nil {
		return nil, fmt.Errorf("failed to find issues: %w", err)
	}

	// Check compatibility
	if err := pa.checkCompatibility(ctx, &result.Compatibility); err != nil {
		return nil, fmt.Errorf("failed to check compatibility: %w", err)
	}

	return result, nil
}

// analyzeProjectInfo extracts basic project information
func (pa *ProjectAnalyzer) analyzeProjectInfo(ctx context.Context, info *ProjectInfo) error {
	info.Name = filepath.Base(pa.projectPath)

	// Parse go.mod
	goModPath := filepath.Join(pa.projectPath, "go.mod")
	if content, err := os.ReadFile(goModPath); err == nil {
		pa.parseGoMod(string(content), info)
	}

	// Get repository info
	if repoInfo := pa.getRepositoryInfo(); repoInfo != "" {
		info.Repository = repoInfo
	}

	// Get project statistics
	stats, err := pa.getProjectStats()
	if err != nil {
		return err
	}
	info.Size = stats.size
	info.FileCount = stats.fileCount
	info.LineCount = stats.lineCount

	// Determine project type
	info.ProjectType = pa.detectProjectType()

	// Get description from README
	if desc := pa.getProjectDescription(); desc != "" {
		info.Description = desc
	}

	// Get license
	if license := pa.detectLicense(); license != "" {
		info.License = license
	}

	return nil
}

// analyzeDependencies analyzes project dependencies
func (pa *ProjectAnalyzer) analyzeDependencies(ctx context.Context, deps *[]Dependency) error {
	goModPath := filepath.Join(pa.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return nil // No go.mod file
	}

	lines := strings.Split(string(content), "\n")
	inRequireBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "require (" {
			inRequireBlock = true
			continue
		}

		if line == ")" && inRequireBlock {
			inRequireBlock = false
			continue
		}

		if inRequireBlock || strings.HasPrefix(line, "require ") {
			if dep := pa.parseDependencyLine(line); dep.Name != "" {
				*deps = append(*deps, dep)
			}
		}
	}

	// Analyze dependency usage
	for i := range *deps {
		pa.analyzeDependencyUsage(&(*deps)[i])
	}

	return nil
}

// analyzeStructure analyzes project structure
func (pa *ProjectAnalyzer) analyzeStructure(ctx context.Context, structure *ProjectStructure) error {
	structure.RootPath = pa.projectPath
	structure.Layout = pa.detectProjectLayout()

	// Walk through directories
	err := filepath.Walk(pa.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(pa.projectPath, path)

		if info.IsDir() {
			if pa.shouldSkipDirectory(relPath) {
				return filepath.SkipDir
			}

			dirInfo := pa.analyzeDirInfo(path, relPath)
			structure.Directories = append(structure.Directories, dirInfo)
		} else {
			fileInfo := pa.analyzeFileInfo(path, relPath, info)
			structure.Files = append(structure.Files, fileInfo)

			// Categorize special files
			if pa.isConfigFile(relPath) {
				structure.ConfigFiles = append(structure.ConfigFiles, relPath)
			}
			if pa.isScript(relPath) {
				structure.Scripts = append(structure.Scripts, relPath)
			}
			if pa.isDocumentation(relPath) {
				structure.Documentation = append(structure.Documentation, relPath)
			}
		}

		return nil
	})

	return err
}

// detectPatterns detects architectural patterns
func (pa *ProjectAnalyzer) detectPatterns(ctx context.Context, patterns *DetectedPatterns) error {
	// Detect architecture
	patterns.Architecture = pa.detectArchitecturalStyle()
	patterns.Patterns = pa.detectArchitecturalPatterns()

	// Analyze web framework
	patterns.WebFramework = pa.detectWebFramework()

	// Analyze databases
	patterns.Database = pa.detectDatabases()

	// Analyze testing
	patterns.Testing = pa.analyzeTestingSetup()

	// Analyze logging
	patterns.Logging = pa.analyzeLogging()

	// Analyze configuration
	patterns.Configuration = pa.analyzeConfiguration()

	// Analyze containerization
	patterns.Containerization = pa.analyzeContainerization()

	// Analyze CI/CD
	patterns.CI = pa.analyzeCICD()

	return nil
}

// calculateMetrics calculates various project metrics
func (pa *ProjectAnalyzer) calculateMetrics(ctx context.Context, metrics *ProjectMetrics) error {
	// Calculate complexity metrics
	metrics.Complexity = pa.calculateComplexityMetrics()

	// Calculate quality metrics
	metrics.Quality = pa.calculateQualityMetrics()

	// Analyze performance
	metrics.Performance = pa.analyzePerformance()

	// Analyze security
	metrics.Security = pa.analyzeSecurityMetrics()

	// Calculate maintainability
	metrics.Maintainability = pa.calculateMaintainabilityMetrics()

	return nil
}

// generateSuggestions generates improvement suggestions
func (pa *ProjectAnalyzer) generateSuggestions(ctx context.Context, result *AnalysisResult, suggestions *[]Suggestion) error {
	// Dependency suggestions
	*suggestions = append(*suggestions, pa.generateDependencySuggestions(result.Dependencies)...)

	// Structure suggestions
	*suggestions = append(*suggestions, pa.generateStructureSuggestions(result.Structure)...)

	// Pattern suggestions
	*suggestions = append(*suggestions, pa.generatePatternSuggestions(result.Patterns)...)

	// Performance suggestions
	*suggestions = append(*suggestions, pa.generatePerformanceSuggestions(result.Metrics.Performance)...)

	// Security suggestions
	*suggestions = append(*suggestions, pa.generateSecuritySuggestions(result.Metrics.Security)...)

	// Sort suggestions by priority
	sort.Slice(*suggestions, func(i, j int) bool {
		return (*suggestions)[i].Priority > (*suggestions)[j].Priority
	})

	return nil
}

// findIssues finds issues in the project
func (pa *ProjectAnalyzer) findIssues(ctx context.Context, issues *[]Issue) error {
	// Find syntax issues
	*issues = append(*issues, pa.findSyntaxIssues()...)

	// Find logic issues
	*issues = append(*issues, pa.findLogicIssues()...)

	// Find security issues
	*issues = append(*issues, pa.findSecurityIssues()...)

	// Find performance issues
	*issues = append(*issues, pa.findPerformanceIssues()...)

	return nil
}

// checkCompatibility checks compatibility information
func (pa *ProjectAnalyzer) checkCompatibility(ctx context.Context, compat *CompatibilityInfo) error {
	// Get Go version
	compat.GoVersion = pa.getGoVersion()

	// Check OS/Arch support
	compat.OSSupport = []string{"linux", "darwin", "windows"}
	compat.ArchSupport = []string{"amd64", "arm64"}

	// Check dependency compatibility
	compat.Dependencies = pa.checkDependencyCompatibility()

	// Check deployment compatibility
	compat.Deployment = pa.checkDeploymentCompatibility()

	// Analyze upgradeability
	compat.Upgradeability = pa.analyzeUpgradeability()

	return nil
}

// Helper methods for parsing and analysis
func (pa *ProjectAnalyzer) parseGoMod(content string, info *ProjectInfo) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			info.ModulePath = strings.TrimPrefix(line, "module ")
		} else if strings.HasPrefix(line, "go ") {
			info.GoVersion = strings.TrimPrefix(line, "go ")
		}
	}
}

type projectStats struct {
	size      int64
	fileCount int
	lineCount int
}

func (pa *ProjectAnalyzer) getProjectStats() (projectStats, error) {
	var stats projectStats

	err := filepath.Walk(pa.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			stats.size += info.Size()
			stats.fileCount++

			if strings.HasSuffix(path, ".go") {
				if lines := pa.countLines(path); lines > 0 {
					stats.lineCount += lines
				}
			}
		}

		return nil
	})

	return stats, err
}

func (pa *ProjectAnalyzer) countLines(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}

func (pa *ProjectAnalyzer) detectProjectType() string {
	// Check for web frameworks
	if pa.hasWebFramework() {
		return "web"
	}

	// Check for CLI patterns
	if pa.hasCLIPattern() {
		return "cli"
	}

	// Check for microservice patterns
	if pa.hasMicroservicePattern() {
		return "microservice"
	}

	// Check for worker patterns
	if pa.hasWorkerPattern() {
		return "worker"
	}

	// Default to library
	return "library"
}

func (pa *ProjectAnalyzer) hasWebFramework() bool {
	goModPath := filepath.Join(pa.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	webFrameworks := []string{
		"github.com/gin-gonic/gin",
		"github.com/labstack/echo",
		"github.com/gofiber/fiber",
		"github.com/gorilla/mux",
		"github.com/go-chi/chi",
	}

	contentStr := string(content)
	for _, framework := range webFrameworks {
		if strings.Contains(contentStr, framework) {
			return true
		}
	}

	return false
}

func (pa *ProjectAnalyzer) hasCLIPattern() bool {
	goModPath := filepath.Join(pa.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	cliLibs := []string{
		"github.com/spf13/cobra",
		"github.com/urfave/cli",
	}

	contentStr := string(content)
	for _, lib := range cliLibs {
		if strings.Contains(contentStr, lib) {
			return true
		}
	}

	return false
}

func (pa *ProjectAnalyzer) hasMicroservicePattern() bool {
	// Look for gRPC, service discovery, etc.
	goModPath := filepath.Join(pa.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	microserviceLibs := []string{
		"google.golang.org/grpc",
		"go.etcd.io/etcd",
		"github.com/hashicorp/consul",
	}

	contentStr := string(content)
	for _, lib := range microserviceLibs {
		if strings.Contains(contentStr, lib) {
			return true
		}
	}

	return false
}

func (pa *ProjectAnalyzer) hasWorkerPattern() bool {
	// Look for job queue libraries
	goModPath := filepath.Join(pa.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	workerLibs := []string{
		"github.com/go-redis/redis",
		"github.com/streadway/amqp",
		"github.com/nsqio/go-nsq",
	}

	contentStr := string(content)
	for _, lib := range workerLibs {
		if strings.Contains(contentStr, lib) {
			return true
		}
	}

	return false
}

func (pa *ProjectAnalyzer) getRepositoryInfo() string {
	// Try to get from .git/config
	gitConfigPath := filepath.Join(pa.projectPath, ".git", "config")
	if content, err := os.ReadFile(gitConfigPath); err == nil {
		re := regexp.MustCompile(`url = (.+)`)
		if match := re.FindStringSubmatch(string(content)); len(match) > 1 {
			return match[1]
		}
	}
	return ""
}

func (pa *ProjectAnalyzer) getProjectDescription() string {
	readmePath := filepath.Join(pa.projectPath, "README.md")
	if content, err := os.ReadFile(readmePath); err == nil {
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "# ") && i < len(lines)-1 {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine != "" && !strings.HasPrefix(nextLine, "#") {
					return nextLine
				}
			}
		}
	}
	return ""
}

func (pa *ProjectAnalyzer) detectLicense() string {
	licensePath := filepath.Join(pa.projectPath, "LICENSE")
	if content, err := os.ReadFile(licensePath); err == nil {
		contentStr := strings.ToUpper(string(content))
		if strings.Contains(contentStr, "MIT") {
			return "MIT"
		} else if strings.Contains(contentStr, "APACHE") {
			return "Apache-2.0"
		} else if strings.Contains(contentStr, "GPL") {
			return "GPL"
		} else if strings.Contains(contentStr, "BSD") {
			return "BSD"
		}
	}
	return ""
}

// Placeholder implementations for complex analysis methods
func (pa *ProjectAnalyzer) parseDependencyLine(line string) Dependency {
	// Simplified dependency parsing
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		name := strings.TrimPrefix(parts[0], "require ")
		version := parts[1]
		return Dependency{
			Name:     name,
			Version:  version,
			Type:     "direct",
			Category: pa.categorizeDependency(name),
		}
	}
	return Dependency{}
}

func (pa *ProjectAnalyzer) categorizeDependency(name string) string {
	if strings.Contains(name, "gin") || strings.Contains(name, "echo") || strings.Contains(name, "fiber") {
		return "web"
	}
	if strings.Contains(name, "gorm") || strings.Contains(name, "sqlx") || strings.Contains(name, "mongo") {
		return "database"
	}
	if strings.Contains(name, "testify") || strings.Contains(name, "ginkgo") {
		return "testing"
	}
	if strings.Contains(name, "cobra") || strings.Contains(name, "viper") {
		return "cli"
	}
	if strings.Contains(name, "logrus") || strings.Contains(name, "zap") || strings.Contains(name, "zerolog") {
		return "logging"
	}
	return "utils"
}

func (pa *ProjectAnalyzer) analyzeDependencyUsage(dep *Dependency) {
	// This would scan Go files for imports and usage
	// Simplified for now
	dep.Usage = []string{"main.go"}
}

func (pa *ProjectAnalyzer) shouldSkipDirectory(relPath string) bool {
	skipDirs := []string{".git", "vendor", "node_modules", ".DS_Store"}
	for _, skip := range skipDirs {
		if relPath == skip || strings.HasPrefix(relPath, skip+"/") {
			return true
		}
	}
	return false
}

func (pa *ProjectAnalyzer) detectProjectLayout() string {
	// Check for standard Go project layout
	if pa.hasDirectory("cmd") && pa.hasDirectory("internal") {
		return "standard"
	}
	if pa.hasDirectory("pkg") {
		return "library"
	}
	return "flat"
}

func (pa *ProjectAnalyzer) hasDirectory(name string) bool {
	dirPath := filepath.Join(pa.projectPath, name)
	if info, err := os.Stat(dirPath); err == nil && info.IsDir() {
		return true
	}
	return false
}

func (pa *ProjectAnalyzer) analyzeDirInfo(path, relPath string) DirectoryInfo {
	info := DirectoryInfo{
		Path:    relPath,
		Purpose: pa.getDirPurpose(relPath),
	}

	// Count files in directory
	if files, err := os.ReadDir(path); err == nil {
		for _, file := range files {
			if !file.IsDir() {
				info.FileCount++
				if strings.HasSuffix(file.Name(), ".go") {
					if strings.HasSuffix(file.Name(), "_test.go") {
						info.TestCount++
					}
					// Count lines
					filePath := filepath.Join(path, file.Name())
					info.LineCount += pa.countLines(filePath)
				}
			}
		}
	}

	return info
}

func (pa *ProjectAnalyzer) getDirPurpose(relPath string) string {
	switch {
	case relPath == "cmd" || strings.HasPrefix(relPath, "cmd/"):
		return "application"
	case relPath == "internal" || strings.HasPrefix(relPath, "internal/"):
		return "private"
	case relPath == "pkg" || strings.HasPrefix(relPath, "pkg/"):
		return "public"
	case relPath == "api" || strings.HasPrefix(relPath, "api/"):
		return "api"
	case strings.Contains(relPath, "test"):
		return "testing"
	case strings.Contains(relPath, "doc"):
		return "documentation"
	default:
		return "other"
	}
}

func (pa *ProjectAnalyzer) analyzeFileInfo(path, relPath string, info os.FileInfo) FileInfo {
	fileInfo := FileInfo{
		Path:         relPath,
		Size:         info.Size(),
		LastModified: info.ModTime(),
		Type:         pa.getFileType(relPath),
		Purpose:      pa.getFilePurpose(relPath),
	}

	if strings.HasSuffix(relPath, ".go") {
		fileInfo.LineCount = pa.countLines(path)
		fileInfo.Complexity = pa.calculateFileComplexity(path)
	}

	return fileInfo
}

func (pa *ProjectAnalyzer) getFileType(relPath string) string {
	switch {
	case strings.HasSuffix(relPath, ".go"):
		if strings.HasSuffix(relPath, "_test.go") {
			return "test"
		}
		return "go"
	case strings.HasSuffix(relPath, ".md"):
		return "doc"
	case strings.HasSuffix(relPath, ".yaml") || strings.HasSuffix(relPath, ".yml"):
		return "config"
	case strings.HasSuffix(relPath, ".json"):
		return "config"
	case relPath == "Dockerfile" || strings.HasPrefix(relPath, "docker-"):
		return "docker"
	case relPath == "Makefile":
		return "script"
	default:
		return "other"
	}
}

func (pa *ProjectAnalyzer) getFilePurpose(relPath string) string {
	switch filepath.Base(relPath) {
	case "main.go":
		return "entry"
	case "README.md":
		return "documentation"
	case "go.mod", "go.sum":
		return "module"
	case "Dockerfile":
		return "container"
	case "Makefile":
		return "build"
	default:
		return "code"
	}
}

func (pa *ProjectAnalyzer) calculateFileComplexity(path string) int {
	// Simplified cyclomatic complexity calculation
	content, err := os.ReadFile(path)
	if err != nil {
		return 0
	}

	complexity := 1 // Base complexity
	contentStr := string(content)

	// Count decision points
	keywords := []string{"if", "for", "switch", "case", "&&", "||"}
	for _, keyword := range keywords {
		complexity += strings.Count(contentStr, keyword)
	}

	return complexity
}

func (pa *ProjectAnalyzer) isConfigFile(relPath string) bool {
	configFiles := []string{".env", "config.yaml", "config.json", "app.toml"}
	fileName := filepath.Base(relPath)
	for _, cf := range configFiles {
		if fileName == cf {
			return true
		}
	}
	return strings.HasSuffix(relPath, ".yaml") || strings.HasSuffix(relPath, ".json") || strings.HasSuffix(relPath, ".toml")
}

func (pa *ProjectAnalyzer) isScript(relPath string) bool {
	return relPath == "Makefile" || strings.HasSuffix(relPath, ".sh") || strings.HasSuffix(relPath, ".bat")
}

func (pa *ProjectAnalyzer) isDocumentation(relPath string) bool {
	return strings.HasSuffix(relPath, ".md") || strings.HasSuffix(relPath, ".txt")
}

// Simplified implementations for remaining methods
func (pa *ProjectAnalyzer) detectArchitecturalStyle() string { return "clean" }
func (pa *ProjectAnalyzer) detectArchitecturalPatterns() []ArchitecturalPattern {
	return []ArchitecturalPattern{}
}
func (pa *ProjectAnalyzer) detectWebFramework() FrameworkInfo             { return FrameworkInfo{} }
func (pa *ProjectAnalyzer) detectDatabases() []DatabaseInfo               { return []DatabaseInfo{} }
func (pa *ProjectAnalyzer) analyzeTestingSetup() TestingInfo              { return TestingInfo{} }
func (pa *ProjectAnalyzer) analyzeLogging() LoggingInfo                   { return LoggingInfo{} }
func (pa *ProjectAnalyzer) analyzeConfiguration() ConfigurationInfo       { return ConfigurationInfo{} }
func (pa *ProjectAnalyzer) analyzeContainerization() ContainerInfo        { return ContainerInfo{} }
func (pa *ProjectAnalyzer) analyzeCICD() CIInfo                           { return CIInfo{} }
func (pa *ProjectAnalyzer) calculateComplexityMetrics() ComplexityMetrics { return ComplexityMetrics{} }
func (pa *ProjectAnalyzer) calculateQualityMetrics() QualityMetrics       { return QualityMetrics{} }
func (pa *ProjectAnalyzer) analyzePerformance() PerformanceHints          { return PerformanceHints{} }
func (pa *ProjectAnalyzer) analyzeSecurityMetrics() SecurityMetrics       { return SecurityMetrics{} }
func (pa *ProjectAnalyzer) calculateMaintainabilityMetrics() MaintenanceMetrics {
	return MaintenanceMetrics{}
}
func (pa *ProjectAnalyzer) generateDependencySuggestions(deps []Dependency) []Suggestion {
	return []Suggestion{}
}
func (pa *ProjectAnalyzer) generateStructureSuggestions(structure ProjectStructure) []Suggestion {
	return []Suggestion{}
}
func (pa *ProjectAnalyzer) generatePatternSuggestions(patterns DetectedPatterns) []Suggestion {
	return []Suggestion{}
}
func (pa *ProjectAnalyzer) generatePerformanceSuggestions(perf PerformanceHints) []Suggestion {
	return []Suggestion{}
}
func (pa *ProjectAnalyzer) generateSecuritySuggestions(security SecurityMetrics) []Suggestion {
	return []Suggestion{}
}
func (pa *ProjectAnalyzer) findSyntaxIssues() []Issue      { return []Issue{} }
func (pa *ProjectAnalyzer) findLogicIssues() []Issue       { return []Issue{} }
func (pa *ProjectAnalyzer) findSecurityIssues() []Issue    { return []Issue{} }
func (pa *ProjectAnalyzer) findPerformanceIssues() []Issue { return []Issue{} }
func (pa *ProjectAnalyzer) getGoVersion() string           { return "1.23" }
func (pa *ProjectAnalyzer) checkDependencyCompatibility() []DependencyCompat {
	return []DependencyCompat{}
}
func (pa *ProjectAnalyzer) checkDeploymentCompatibility() []DeploymentCompat {
	return []DeploymentCompat{}
}
func (pa *ProjectAnalyzer) analyzeUpgradeability() UpgradeInfo { return UpgradeInfo{} }

// ToGoCtlConfig converts analysis result to go-ctl project configuration
func (result *AnalysisResult) ToGoCtlConfig() *metadata.ProjectConfig {
	config := &metadata.ProjectConfig{
		ProjectName:    result.ProjectInfo.Name,
		GoVersion:      result.ProjectInfo.GoVersion,
		CustomPackages: []string{},
	}

	// Extract dependencies that could be useful
	for _, dep := range result.Dependencies {
		if dep.Category == "web" || dep.Category == "database" {
			config.CustomPackages = append(config.CustomPackages, dep.Name)
		}
	}

	return config
}
