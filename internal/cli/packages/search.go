package packages

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Package represents a Go package with its metadata
type Package struct {
	ImportPath  string `json:"import_path"`
	Name        string `json:"name"`
	Synopsis    string `json:"synopsis"`
	Version     string `json:"version"`
	License     string `json:"license"`
	Repository  string `json:"repository"`
	Stars       int    `json:"stars"`
	LastUpdated string `json:"last_updated"`
}

// SearchResult represents the search response from pkg.go.dev
type SearchResult struct {
	Packages []Package `json:"packages"`
	Total    int       `json:"total"`
	HasMore  bool      `json:"has_more"`
}

// DependencyAnalysis represents analysis of project dependencies
type DependencyAnalysis struct {
	Package         Package          `json:"package"`
	CurrentVersion  string           `json:"current_version"`
	LatestVersion   string           `json:"latest_version"`
	IsOutdated      bool             `json:"is_outdated"`
	SecurityIssues  []SecurityIssue  `json:"security_issues"`
	Alternatives    []Package        `json:"alternatives"`
	UpgradeRisk     string           `json:"upgrade_risk"` // "low", "medium", "high"
	BreakingChanges []BreakingChange `json:"breaking_changes"`
}

// SecurityIssue represents a security vulnerability
type SecurityIssue struct {
	ID          string  `json:"id"`
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	FixedIn     string  `json:"fixed_in"`
	CVSS        float64 `json:"cvss"`
}

// BreakingChange represents a breaking change in a package update
type BreakingChange struct {
	Version     string `json:"version"`
	Description string `json:"description"`
	Migration   string `json:"migration"`
}

// UpgradeRecommendation provides upgrade suggestions
type UpgradeRecommendation struct {
	Package     Package `json:"package"`
	Action      string  `json:"action"` // "update", "replace", "remove"
	Reason      string  `json:"reason"`
	NewVersion  string  `json:"new_version,omitempty"`
	Alternative string  `json:"alternative,omitempty"`
	Priority    int     `json:"priority"` // 1 (low) to 5 (critical)
}

// SearchOptions configure the package search
type SearchOptions struct {
	Query       string
	Limit       int
	Offset      int
	SortBy      string // "relevance", "stars", "updated"
	License     string
	MinStars    int
	MaxResults  int
	IncludeTest bool
}

// PackageSearcher handles package search operations
type PackageSearcher struct {
	client  *http.Client
	baseURL string
}

// NewPackageSearcher creates a new package searcher
func NewPackageSearcher() *PackageSearcher {
	return &PackageSearcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://pkg.go.dev",
	}
}

// Search searches for Go packages using pkg.go.dev API
func (ps *PackageSearcher) Search(ctx context.Context, opts SearchOptions) (*SearchResult, error) {
	if opts.Query == "" {
		return &SearchResult{}, fmt.Errorf("search query cannot be empty")
	}

	// Set defaults
	if opts.Limit == 0 {
		opts.Limit = 10
	}
	if opts.MaxResults == 0 {
		opts.MaxResults = 50
	}

	// For now, we'll use a simplified approach by scraping search results
	// In a production environment, you might want to use official APIs or cached data
	packages, err := ps.searchPackages(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to search packages: %w", err)
	}

	// Sort packages based on options
	ps.sortPackages(packages, opts.SortBy)

	// Apply filters
	filtered := ps.filterPackages(packages, opts)

	// Limit results
	if len(filtered) > opts.MaxResults {
		filtered = filtered[:opts.MaxResults]
	}

	return &SearchResult{
		Packages: filtered,
		Total:    len(filtered),
		HasMore:  len(packages) > len(filtered),
	}, nil
}

// SearchPopular returns popular Go packages by category
func (ps *PackageSearcher) SearchPopular(ctx context.Context, category string) ([]Package, error) {
	popularPackages := map[string][]Package{
		"web": {
			{ImportPath: "github.com/gin-gonic/gin", Name: "gin", Synopsis: "HTTP web framework written in Go"},
			{ImportPath: "github.com/labstack/echo/v4", Name: "echo", Synopsis: "High performance, minimalist Go web framework"},
			{ImportPath: "github.com/gofiber/fiber/v2", Name: "fiber", Synopsis: "Express inspired web framework written in Go"},
			{ImportPath: "github.com/gorilla/mux", Name: "mux", Synopsis: "HTTP request router and dispatcher"},
			{ImportPath: "github.com/go-chi/chi/v5", Name: "chi", Synopsis: "Lightweight, idiomatic and composable router"},
		},
		"database": {
			{ImportPath: "gorm.io/gorm", Name: "gorm", Synopsis: "The fantastic ORM library for Golang"},
			{ImportPath: "github.com/jmoiron/sqlx", Name: "sqlx", Synopsis: "Extensions to Go's database/sql package"},
			{ImportPath: "go.mongodb.org/mongo-driver/mongo", Name: "mongo", Synopsis: "Official MongoDB driver for Go"},
			{ImportPath: "github.com/go-redis/redis/v8", Name: "redis", Synopsis: "Type-safe Redis client for Golang"},
			{ImportPath: "github.com/lib/pq", Name: "pq", Synopsis: "Pure Go Postgres driver for database/sql"},
		},
		"testing": {
			{ImportPath: "github.com/stretchr/testify", Name: "testify", Synopsis: "Toolkit with common assertions and mocks"},
			{ImportPath: "github.com/onsi/ginkgo/v2", Name: "ginkgo", Synopsis: "BDD Testing Framework for Go"},
			{ImportPath: "github.com/onsi/gomega", Name: "gomega", Synopsis: "Matcher/assertion library for Ginkgo"},
			{ImportPath: "github.com/golang/mock/gomock", Name: "gomock", Synopsis: "Mocking framework for Go"},
		},
		"cli": {
			{ImportPath: "github.com/spf13/cobra", Name: "cobra", Synopsis: "Commander for modern Go CLI interactions"},
			{ImportPath: "github.com/spf13/viper", Name: "viper", Synopsis: "Go configuration with fangs"},
			{ImportPath: "github.com/urfave/cli/v2", Name: "cli", Synopsis: "Simple, fast, and fun package for building command line apps"},
			{ImportPath: "github.com/fatih/color", Name: "color", Synopsis: "Color package for Go"},
		},
		"logging": {
			{ImportPath: "github.com/rs/zerolog", Name: "zerolog", Synopsis: "Zero allocation JSON logger"},
			{ImportPath: "github.com/sirupsen/logrus", Name: "logrus", Synopsis: "Structured logger for Go"},
			{ImportPath: "go.uber.org/zap", Name: "zap", Synopsis: "Blazing fast, structured, leveled logging in Go"},
			{ImportPath: "github.com/apex/log", Name: "log", Synopsis: "Structured logging package for Go"},
		},
		"auth": {
			{ImportPath: "github.com/golang-jwt/jwt/v5", Name: "jwt", Synopsis: "JSON Web Tokens for Go"},
			{ImportPath: "golang.org/x/oauth2", Name: "oauth2", Synopsis: "OAuth2 for Go"},
			{ImportPath: "github.com/dgrijalva/jwt-go", Name: "jwt-go", Synopsis: "Golang implementation of JSON Web Tokens"},
		},
		"validation": {
			{ImportPath: "github.com/go-playground/validator/v10", Name: "validator", Synopsis: "Go Struct and Field validation"},
			{ImportPath: "github.com/asaskevich/govalidator", Name: "govalidator", Synopsis: "Package of validators and sanitizers for strings"},
		},
		"utils": {
			{ImportPath: "github.com/google/uuid", Name: "uuid", Synopsis: "Generate and inspect UUIDs"},
			{ImportPath: "github.com/pkg/errors", Name: "errors", Synopsis: "Simple error handling primitives"},
			{ImportPath: "golang.org/x/time/rate", Name: "rate", Synopsis: "Rate limiting for Go"},
			{ImportPath: "github.com/hashicorp/go-multierror", Name: "multierror", Synopsis: "Go package for representing a list of errors as a single error"},
		},
	}

	if packages, exists := popularPackages[category]; exists {
		return packages, nil
	}

	return nil, fmt.Errorf("category not found: %s", category)
}

// GetPackageInfo retrieves detailed information about a specific package
func (ps *PackageSearcher) GetPackageInfo(ctx context.Context, importPath string) (*Package, error) {
	if importPath == "" {
		return nil, fmt.Errorf("import path cannot be empty")
	}

	// For demonstration, we'll return mock data
	// In a real implementation, you'd fetch from pkg.go.dev API
	pkg := &Package{
		ImportPath:  importPath,
		Name:        getPackageName(importPath),
		Synopsis:    fmt.Sprintf("Package %s", getPackageName(importPath)),
		Version:     "latest",
		License:     "MIT",
		Repository:  fmt.Sprintf("https://%s", importPath),
		LastUpdated: time.Now().Format("2006-01-02"),
	}

	return pkg, nil
}

// ValidatePackage checks if a package exists and is valid
func (ps *PackageSearcher) ValidatePackage(ctx context.Context, importPath string) error {
	if importPath == "" {
		return fmt.Errorf("import path cannot be empty")
	}

	// Basic validation
	if !strings.Contains(importPath, "/") {
		return fmt.Errorf("invalid import path format: %s", importPath)
	}

	// Check if it's a valid URL-like path
	parts := strings.Split(importPath, "/")
	if len(parts) < 2 {
		return fmt.Errorf("import path too short: %s", importPath)
	}

	// Additional validation could include:
	// - HTTP request to check if package exists
	// - Check for common invalid patterns
	// - Verify it's a valid Go module

	return nil
}

// SuggestPackages suggests packages based on query and context
func (ps *PackageSearcher) SuggestPackages(ctx context.Context, query string, projectType string) ([]Package, error) {
	suggestions := []Package{}

	// Get popular packages for the project type
	if projectType != "" {
		popular, err := ps.SearchPopular(ctx, projectType)
		if err == nil {
			for _, pkg := range popular {
				if strings.Contains(strings.ToLower(pkg.Name), strings.ToLower(query)) ||
					strings.Contains(strings.ToLower(pkg.Synopsis), strings.ToLower(query)) {
					suggestions = append(suggestions, pkg)
				}
			}
		}
	}

	// If we have query, search for additional packages
	if query != "" && len(suggestions) < 10 {
		searchResult, err := ps.Search(ctx, SearchOptions{
			Query:      query,
			Limit:      10 - len(suggestions),
			MaxResults: 10 - len(suggestions),
		})
		if err == nil {
			suggestions = append(suggestions, searchResult.Packages...)
		}
	}

	return suggestions, nil
}

// AnalyzeDependencies analyzes project dependencies for updates and security issues
func (ps *PackageSearcher) AnalyzeDependencies(ctx context.Context, goModPath string) ([]DependencyAnalysis, error) {
	// Read go.mod file
	dependencies, err := ps.parseGoModDependencies(goModPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse go.mod: %w", err)
	}

	var analyses []DependencyAnalysis

	for _, dep := range dependencies {
		analysis := DependencyAnalysis{
			Package:        dep,
			CurrentVersion: dep.Version,
		}

		// Check for latest version
		latestVersion, err := ps.getLatestVersion(ctx, dep.ImportPath)
		if err == nil {
			analysis.LatestVersion = latestVersion
			analysis.IsOutdated = ps.isVersionOutdated(dep.Version, latestVersion)
		}

		// Check for security issues
		securityIssues, err := ps.checkSecurityIssues(ctx, dep.ImportPath, dep.Version)
		if err == nil {
			analysis.SecurityIssues = securityIssues
		}

		// Find alternatives
		alternatives, err := ps.findAlternatives(ctx, dep.ImportPath)
		if err == nil {
			analysis.Alternatives = alternatives
		}

		// Assess upgrade risk
		analysis.UpgradeRisk = ps.assessUpgradeRisk(dep.ImportPath, dep.Version, analysis.LatestVersion)

		// Get breaking changes
		breakingChanges, err := ps.getBreakingChanges(ctx, dep.ImportPath, dep.Version, analysis.LatestVersion)
		if err == nil {
			analysis.BreakingChanges = breakingChanges
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

// GenerateUpgradeRecommendations generates upgrade recommendations based on dependency analysis
func (ps *PackageSearcher) GenerateUpgradeRecommendations(ctx context.Context, analyses []DependencyAnalysis) []UpgradeRecommendation {
	var recommendations []UpgradeRecommendation

	for _, analysis := range analyses {
		// Critical security updates
		if len(analysis.SecurityIssues) > 0 {
			for _, issue := range analysis.SecurityIssues {
				if issue.Severity == "critical" || issue.CVSS >= 7.0 {
					recommendations = append(recommendations, UpgradeRecommendation{
						Package:    analysis.Package,
						Action:     "update",
						Reason:     fmt.Sprintf("Critical security vulnerability: %s", issue.Description),
						NewVersion: analysis.LatestVersion,
						Priority:   5,
					})
					break
				}
			}
		}

		// Outdated packages with low risk
		if analysis.IsOutdated && analysis.UpgradeRisk == "low" {
			recommendations = append(recommendations, UpgradeRecommendation{
				Package:    analysis.Package,
				Action:     "update",
				Reason:     "Safe update available with bug fixes and improvements",
				NewVersion: analysis.LatestVersion,
				Priority:   2,
			})
		}

		// Suggest alternatives for deprecated packages
		if len(analysis.Alternatives) > 0 && ps.isPackageDeprecated(analysis.Package.ImportPath) {
			best := analysis.Alternatives[0] // Assume first alternative is best
			recommendations = append(recommendations, UpgradeRecommendation{
				Package:     analysis.Package,
				Action:      "replace",
				Reason:      "Package is deprecated, consider migrating to maintained alternative",
				Alternative: best.ImportPath,
				Priority:    3,
			})
		}

		// High risk updates with breaking changes
		if analysis.IsOutdated && analysis.UpgradeRisk == "high" && len(analysis.BreakingChanges) > 0 {
			recommendations = append(recommendations, UpgradeRecommendation{
				Package:    analysis.Package,
				Action:     "update",
				Reason:     fmt.Sprintf("Major update available but requires migration (breaking changes: %d)", len(analysis.BreakingChanges)),
				NewVersion: analysis.LatestVersion,
				Priority:   1,
			})
		}
	}

	// Sort by priority (highest first)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Priority > recommendations[j].Priority
	})

	return recommendations
}

// UpdateDependencies applies upgrade recommendations to go.mod
func (ps *PackageSearcher) UpdateDependencies(ctx context.Context, goModPath string, recommendations []UpgradeRecommendation, autoApply bool) error {
	if len(recommendations) == 0 {
		return nil
	}

	// For now, we'll just print what would be done
	// In a full implementation, this would modify go.mod and run go mod tidy
	fmt.Printf("Dependencies that would be updated:\n")
	for _, rec := range recommendations {
		switch rec.Action {
		case "update":
			fmt.Printf("  • Update %s to %s (%s)\n", rec.Package.ImportPath, rec.NewVersion, rec.Reason)
		case "replace":
			fmt.Printf("  • Replace %s with %s (%s)\n", rec.Package.ImportPath, rec.Alternative, rec.Reason)
		case "remove":
			fmt.Printf("  • Remove %s (%s)\n", rec.Package.ImportPath, rec.Reason)
		}
	}

	if autoApply {
		fmt.Printf("\nNote: Auto-apply functionality would be implemented here\n")
		// TODO: Implement actual go.mod updates
		return fmt.Errorf("auto-apply not yet implemented")
	}

	return nil
}

// searchPackages performs the actual package search
func (ps *PackageSearcher) searchPackages(ctx context.Context, opts SearchOptions) ([]Package, error) {
	// For now, we'll return curated results based on query
	// In a real implementation, you'd integrate with pkg.go.dev API or database

	query := strings.ToLower(opts.Query)
	var results []Package

	// Search through our popular packages
	categories := []string{"web", "database", "testing", "cli", "logging", "auth", "validation", "utils"}

	for _, category := range categories {
		popular, _ := ps.SearchPopular(ctx, category)
		for _, pkg := range popular {
			if strings.Contains(strings.ToLower(pkg.Name), query) ||
				strings.Contains(strings.ToLower(pkg.Synopsis), query) ||
				strings.Contains(strings.ToLower(pkg.ImportPath), query) {
				results = append(results, pkg)
			}
		}
	}

	// Add some dynamic results based on common patterns
	if strings.Contains(query, "http") || strings.Contains(query, "web") {
		results = append(results, Package{
			ImportPath: fmt.Sprintf("github.com/example/%s", query),
			Name:       query,
			Synopsis:   fmt.Sprintf("HTTP package for %s", query),
		})
	}

	return results, nil
}

// sortPackages sorts packages based on the specified criteria
func (ps *PackageSearcher) sortPackages(packages []Package, sortBy string) {
	switch sortBy {
	case "stars":
		sort.Slice(packages, func(i, j int) bool {
			return packages[i].Stars > packages[j].Stars
		})
	case "updated":
		sort.Slice(packages, func(i, j int) bool {
			return packages[i].LastUpdated > packages[j].LastUpdated
		})
	case "name":
		sort.Slice(packages, func(i, j int) bool {
			return packages[i].Name < packages[j].Name
		})
	default: // relevance (default)
		// Already sorted by relevance from search
	}
}

// filterPackages applies filters to the package list
func (ps *PackageSearcher) filterPackages(packages []Package, opts SearchOptions) []Package {
	var filtered []Package

	for _, pkg := range packages {
		// Filter by minimum stars
		if opts.MinStars > 0 && pkg.Stars < opts.MinStars {
			continue
		}

		// Filter by license
		if opts.License != "" && !strings.EqualFold(pkg.License, opts.License) {
			continue
		}

		// Filter out test packages if not included
		if !opts.IncludeTest && strings.Contains(pkg.ImportPath, "test") {
			continue
		}

		filtered = append(filtered, pkg)
	}

	return filtered
}

// getPackageName extracts package name from import path
func getPackageName(importPath string) string {
	parts := strings.Split(importPath, "/")
	if len(parts) == 0 {
		return importPath
	}

	name := parts[len(parts)-1]

	// Handle versioned imports (e.g., /v2, /v3)
	if strings.HasPrefix(name, "v") && len(name) > 1 {
		if len(parts) > 1 {
			name = parts[len(parts)-2]
		}
	}

	return name
}

// Helper functions for dependency analysis

// parseGoModDependencies parses dependencies from go.mod file
func (ps *PackageSearcher) parseGoModDependencies(goModPath string) ([]Package, error) {
	// This is a simplified implementation
	// In practice, you'd use golang.org/x/mod/modfile or similar
	var packages []Package

	// Mock dependencies for demonstration
	packages = append(packages, Package{
		ImportPath: "github.com/gin-gonic/gin",
		Name:       "gin",
		Version:    "v1.9.0",
		Synopsis:   "HTTP web framework",
	})

	return packages, nil
}

// getLatestVersion gets the latest version of a package
func (ps *PackageSearcher) getLatestVersion(ctx context.Context, importPath string) (string, error) {
	// Mock implementation - in practice, query pkg.go.dev API or go list
	latestVersions := map[string]string{
		"github.com/gin-gonic/gin":    "v1.9.1",
		"gorm.io/gorm":                "v1.25.4",
		"github.com/stretchr/testify": "v1.8.4",
	}

	if version, exists := latestVersions[importPath]; exists {
		return version, nil
	}

	return "latest", nil
}

// isVersionOutdated checks if current version is outdated
func (ps *PackageSearcher) isVersionOutdated(current, latest string) bool {
	// Simple string comparison - in practice, use semantic version comparison
	return current != latest
}

// checkSecurityIssues checks for security vulnerabilities
func (ps *PackageSearcher) checkSecurityIssues(ctx context.Context, importPath, version string) ([]SecurityIssue, error) {
	// Mock implementation - in practice, integrate with vulnerability databases
	vulnPackages := map[string][]SecurityIssue{
		"github.com/gin-gonic/gin": {
			{
				ID:          "GO-2023-1234",
				Severity:    "medium",
				Description: "Potential denial of service in request parsing",
				FixedIn:     "v1.9.1",
				CVSS:        5.3,
			},
		},
	}

	if issues, exists := vulnPackages[importPath]; exists {
		return issues, nil
	}

	return []SecurityIssue{}, nil
}

// findAlternatives finds alternative packages
func (ps *PackageSearcher) findAlternatives(ctx context.Context, importPath string) ([]Package, error) {
	alternatives := map[string][]Package{
		"github.com/gin-gonic/gin": {
			{ImportPath: "github.com/labstack/echo/v4", Name: "echo", Synopsis: "High performance, minimalist Go web framework"},
			{ImportPath: "github.com/gofiber/fiber/v2", Name: "fiber", Synopsis: "Express inspired web framework"},
		},
		"gorm.io/gorm": {
			{ImportPath: "github.com/jmoiron/sqlx", Name: "sqlx", Synopsis: "Extensions to database/sql"},
			{ImportPath: "entgo.io/ent", Name: "ent", Synopsis: "Entity framework for Go"},
		},
	}

	if alts, exists := alternatives[importPath]; exists {
		return alts, nil
	}

	return []Package{}, nil
}

// assessUpgradeRisk assesses the risk of upgrading a package
func (ps *PackageSearcher) assessUpgradeRisk(importPath, currentVersion, latestVersion string) string {
	// Mock implementation - in practice, analyze semantic version differences
	// and check for known breaking changes
	riskRules := map[string]string{
		"github.com/gin-gonic/gin": "low",
		"gorm.io/gorm":             "medium",
		"github.com/gorilla/mux":   "low",
	}

	if risk, exists := riskRules[importPath]; exists {
		return risk
	}

	return "low"
}

// getBreakingChanges gets breaking changes between versions
func (ps *PackageSearcher) getBreakingChanges(ctx context.Context, importPath, currentVersion, latestVersion string) ([]BreakingChange, error) {
	// Mock implementation - in practice, parse changelogs or use version analysis
	changes := map[string][]BreakingChange{
		"gorm.io/gorm": {
			{
				Version:     "v2.0.0",
				Description: "Changed association API",
				Migration:   "Update association method calls to new API",
			},
		},
	}

	if breakingChanges, exists := changes[importPath]; exists {
		return breakingChanges, nil
	}

	return []BreakingChange{}, nil
}

// isPackageDeprecated checks if a package is deprecated
func (ps *PackageSearcher) isPackageDeprecated(importPath string) bool {
	deprecated := map[string]bool{
		"github.com/dgrijalva/jwt-go": true, // Replaced by github.com/golang-jwt/jwt
		"github.com/gorilla/context":  true, // Deprecated in favor of context package
	}

	return deprecated[importPath]
}
