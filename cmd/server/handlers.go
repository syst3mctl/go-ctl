package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// FileItem represents a file or folder in the project structure
type FileItem struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Icon     string      `json:"icon"`
	IsFolder bool        `json:"is_folder"`
	Children []*FileItem `json:"children,omitempty"`
	Level    int         `json:"level"`
}

// TreeNode helper for building file tree
type TreeNode struct {
	Name     string
	Path     string
	Icon     string
	IsFolder bool
	Children map[string]*TreeNode
	Level    int
}

// ProjectStructureData contains file items and project config for templates
type ProjectStructureData struct {
	Files  []FileItem             `json:"files"`
	Config metadata.ProjectConfig `json:"config"`
}

// PkgGoDevResult represents a package search result from pkg.go.dev
type PkgGoDevResult struct {
	Path     string `json:"path"`
	Synopsis string `json:"synopsis"`
}

// Package search cache
type packageSearchCache struct {
	mu      sync.RWMutex
	results map[string]cacheEntry
}

type cacheEntry struct {
	data      []PkgGoDevResult
	timestamp time.Time
}

var searchCache = &packageSearchCache{
	results: make(map[string]cacheEntry),
}

const cacheDuration = 10 * time.Minute

// handleIndex serves the main project generator page
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load and execute the main template
	tmpl := template.Must(template.New("index").Funcs(template.FuncMap{
		"hasFeature": func(features []metadata.Option, featureID string) bool {
			for _, feature := range features {
				if feature.ID == featureID {
					return true
				}
			}
			return false
		},
	}).Parse(indexTemplate))

	data := struct {
		Options *metadata.ProjectOptions
	}{
		Options: appOptions,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleGenerate processes the form submission and generates a project ZIP
func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Build project configuration from form data
	config := metadata.ProjectConfig{
		ProjectName:    r.FormValue("projectName"),
		GoVersion:      r.FormValue("goVersion"),
		HttpPackage:    metadata.FindOption(appOptions.Http, r.FormValue("httpPackage")),
		Database:       metadata.FindOption(appOptions.Databases, r.FormValue("database")),
		DbDriver:       metadata.FindOption(appOptions.DbDrivers, r.FormValue("dbDriver")),
		Features:       metadata.FindOptions(appOptions.Features, r.Form["features"]),
		CustomPackages: r.Form["customPackages"],
	}

	// Validate configuration
	if warnings := metadata.ValidateConfig(config); len(warnings) > 0 {
		// For now, just log warnings - in production you might want to show them to the user
		for _, warning := range warnings {
			fmt.Printf("Warning: %s\n", warning)
		}
	}

	// Set headers for ZIP download
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", config.ProjectName))
	w.Header().Set("Cache-Control", "no-cache")

	// Generate and stream the ZIP file
	if err := gen.GenerateProjectZip(config, w); err != nil {
		// If we haven't written headers yet, we can still return an error
		http.Error(w, "Failed to generate project: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleExplore generates a preview of the project structure
func handleExplore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Build project configuration from form data
	config := metadata.ProjectConfig{
		ProjectName:    r.FormValue("projectName"),
		GoVersion:      r.FormValue("goVersion"),
		HttpPackage:    metadata.FindOption(appOptions.Http, r.FormValue("httpPackage")),
		Database:       metadata.FindOption(appOptions.Databases, r.FormValue("database")),
		DbDriver:       metadata.FindOption(appOptions.DbDrivers, r.FormValue("dbDriver")),
		Features:       metadata.FindOptions(appOptions.Features, r.Form["features"]),
		CustomPackages: r.Form["customPackages"],
	}

	// Generate file items for the file tree
	fileItems := generateFileItems(config)

	// Return HTML snippet for HTMX
	funcMap := template.FuncMap{
		"mul": func(a, b int) int { return a * b },
		"add": func(a, b int) int { return a + b },
	}
	tmpl := template.Must(template.New("explore").Funcs(funcMap).Parse(exploreTemplate))

	data := ProjectStructureData{
		Files:  fileItems,
		Config: config,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

// handleSearchPackages searches pkg.go.dev for packages
func handleSearchPackages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		w.Write([]byte("")) // Return empty response for empty query
		return
	}

	// Search pkg.go.dev (simplified implementation)
	results, err := searchPackages(query)
	if err != nil {
		http.Error(w, "Failed to search packages", http.StatusInternalServerError)
		return
	}

	// Render search results template
	tmpl := template.Must(template.New("search-results").Parse(searchResultsTemplate))

	data := struct {
		Results []PkgGoDevResult
		Query   string
	}{
		Results: results,
		Query:   query,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render search results", http.StatusInternalServerError)
		return
	}
}

// handleAddPackage adds a package to the selected packages list
func handleAddPackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	pkgPath := r.FormValue("pkgPath")
	if pkgPath == "" {
		http.Error(w, "Package path is required", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the package element
	pkgID := strings.ReplaceAll(pkgPath, "/", "-")
	pkgID = strings.ReplaceAll(pkgID, ".", "-")

	// Render selected package item template
	tmpl := template.Must(template.New("selected-package").Parse(selectedPackageTemplate))

	data := struct {
		PkgPath string
		ID      string
	}{
		PkgPath: pkgPath,
		ID:      pkgID,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render package item", http.StatusInternalServerError)
		return
	}
}

// handleFileContent serves individual file content for the modal
func handleFileContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// Parse project configuration from form data or session
	// For now, we'll create a default config and let user override via query params
	config := metadata.ProjectConfig{
		ProjectName: getQueryParam(r, "projectName", "my-go-app"),
		GoVersion:   getQueryParam(r, "goVersion", "1.23"),
		HttpPackage: metadata.Option{
			ID:   getQueryParam(r, "httpPackage", "gin"),
			Name: "Gin",
		},
		Database: metadata.Option{
			ID:   getQueryParam(r, "database", "postgres"),
			Name: "PostgreSQL",
		},
		DbDriver: metadata.Option{
			ID:   getQueryParam(r, "dbDriver", "gorm"),
			Name: "GORM",
		},
	}

	// Generate content based on the file path and configuration
	content := generateFileContentWithConfig(filePath, config)

	// Return raw content for JavaScript to handle
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(content))
}

// getCachedResults retrieves cached search results if they exist and are still valid
func (c *packageSearchCache) get(query string) ([]PkgGoDevResult, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.results[query]
	if !exists {
		return nil, false
	}

	if time.Since(entry.timestamp) > cacheDuration {
		return nil, false
	}

	return entry.data, true
}

// setCachedResults stores search results in cache
func (c *packageSearchCache) set(query string, results []PkgGoDevResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.results[query] = cacheEntry{
		data:      results,
		timestamp: time.Now(),
	}

	// Clean up old entries (simple cleanup)
	if len(c.results) > 100 {
		for k, v := range c.results {
			if time.Since(v.timestamp) > cacheDuration {
				delete(c.results, k)
			}
		}
	}
}

// searchPackages performs package search with caching and fallback
func searchPackages(query string) ([]PkgGoDevResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []PkgGoDevResult{}, nil
	}

	// Check cache first
	if cached, found := searchCache.get(query); found {
		return cached, nil
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Build request URL
	baseURL := "https://pkg.go.dev/search"
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("q", query)
	q.Add("m", "json")
	q.Add("limit", "20")
	req.URL.RawQuery = q.Encode()

	// Set headers
	req.Header.Set("User-Agent", "go-ctl-initializer/1.0")
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		// Fallback to mock results on network error
		log.Printf("Failed to search pkg.go.dev: %v, using fallback", err)
		return searchPackagesFallback(query)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("pkg.go.dev returned status %d, using fallback", resp.StatusCode)
		return searchPackagesFallback(query)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v, using fallback", err)
		return searchPackagesFallback(query)
	}

	// Parse JSON response
	var apiResponse struct {
		Results []struct {
			Path     string `json:"Path"`
			Synopsis string `json:"Synopsis"`
			Version  string `json:"Version"`
		} `json:"Results"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Failed to parse JSON response: %v, using fallback", err)
		return searchPackagesFallback(query)
	}

	// Convert to our result format
	var results []PkgGoDevResult
	for _, result := range apiResponse.Results {
		if len(results) >= 15 { // Limit results
			break
		}
		results = append(results, PkgGoDevResult{
			Path:     result.Path,
			Synopsis: result.Synopsis,
		})
	}

	// If no results from API, try fallback
	if len(results) == 0 {
		results, err := searchPackagesFallback(query)
		if err != nil {
			return nil, err
		}
		// Cache fallback results too
		searchCache.set(query, results)
		return results, nil
	}

	// Cache successful results
	searchCache.set(query, results)
	return results, nil
}

// searchPackagesFallback provides fallback search results when pkg.go.dev is unavailable
func searchPackagesFallback(query string) ([]PkgGoDevResult, error) {
	fallbackResults := []PkgGoDevResult{
		{
			Path:     "github.com/gin-gonic/gin",
			Synopsis: "Gin is a HTTP web framework written in Go (Golang)",
		},
		{
			Path:     "github.com/labstack/echo/v4",
			Synopsis: "High performance, minimalist Go web framework",
		},
		{
			Path:     "github.com/gofiber/fiber/v2",
			Synopsis: "Express inspired web framework written in Go",
		},
		{
			Path:     "github.com/go-chi/chi/v5",
			Synopsis: "Lightweight, idiomatic and composable router for building HTTP services",
		},
		{
			Path:     "gorm.io/gorm",
			Synopsis: "The fantastic ORM library for Golang",
		},
		{
			Path:     "github.com/jmoiron/sqlx",
			Synopsis: "General purpose extensions to golang's database/sql",
		},
		{
			Path:     "go.mongodb.org/mongo-driver",
			Synopsis: "The MongoDB official Go driver",
		},
		{
			Path:     "github.com/redis/go-redis/v9",
			Synopsis: "Redis client for Go",
		},
		{
			Path:     "github.com/golang-jwt/jwt/v5",
			Synopsis: "JWT token authentication library for Go",
		},
		{
			Path:     "github.com/rs/cors",
			Synopsis: "Go net/http configurable handler to handle CORS requests",
		},
		{
			Path:     "github.com/rs/zerolog",
			Synopsis: "Zero Allocation JSON Logger",
		},
		{
			Path:     "github.com/spf13/viper",
			Synopsis: "Go configuration with fangs",
		},
		{
			Path:     "github.com/stretchr/testify",
			Synopsis: "A toolkit with common assertions and mocks for Go tests",
		},
		{
			Path:     "entgo.io/ent",
			Synopsis: "An entity framework for Go",
		},
		{
			Path:     "github.com/lib/pq",
			Synopsis: "Pure Go Postgres driver for database/sql",
		},
		{
			Path:     "github.com/go-sql-driver/mysql",
			Synopsis: "Go MySQL Driver is a MySQL driver for Go's database/sql package",
		},
		{
			Path:     "github.com/mattn/go-sqlite3",
			Synopsis: "sqlite3 driver for go using database/sql",
		},
		{
			Path:     "google.golang.org/grpc",
			Synopsis: "The Go implementation of gRPC: A high-performance RPC framework",
		},
		{
			Path:     "github.com/gorilla/mux",
			Synopsis: "A powerful HTTP router and URL matcher for building Go web servers",
		},
		{
			Path:     "github.com/gorilla/websocket",
			Synopsis: "A fast, well-tested and widely used WebSocket implementation for Go",
		},
	}

	// Filter results based on query
	var results []PkgGoDevResult
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return fallbackResults[:10], nil
	}

	for _, result := range fallbackResults {
		if strings.Contains(strings.ToLower(result.Path), query) ||
			strings.Contains(strings.ToLower(result.Synopsis), query) {
			results = append(results, result)
			if len(results) >= 10 {
				break
			}
		}
	}

	return results, nil
}

// generateFileItems creates file items for the file tree modal
func generateFileItems(config metadata.ProjectConfig) []FileItem {
	// Create list of all file paths
	filePaths := []struct {
		Path string
		Icon string
	}{
		{"go.mod", "fas fa-cube text-green-500"},
		{"README.md", "fab fa-markdown text-blue-600"},
		{fmt.Sprintf("cmd/%s/main.go", config.ProjectName), "fab fa-golang text-blue-500"},
		{"internal/config/config.go", "fas fa-cog text-gray-600"},
		{"internal/domain/model.go", "fab fa-golang text-blue-500"},
		{"internal/service/service.go", "fab fa-golang text-blue-500"},
		{"internal/handler/handler.go", "fab fa-golang text-blue-500"},
	}

	// Add database storage layer if configured
	if config.Database.ID != "" {
		storageFile := fmt.Sprintf("internal/storage/%s/repository.go", config.Database.ID)
		filePaths = append(filePaths, struct {
			Path string
			Icon string
		}{storageFile, "fas fa-database text-purple-500"})

		// Add storage db.go file
		filePaths = append(filePaths, struct {
			Path string
			Icon string
		}{"internal/storage/db.go", "fas fa-database text-purple-500"})
	}

	// Add feature files
	for _, feature := range config.Features {
		switch feature.ID {
		case "gitignore":
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{".gitignore", "fab fa-git-alt text-orange-500"})
		case "makefile":
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{"Makefile", "fas fa-hammer text-gray-600"})
		case "env":
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{".env.example", "fas fa-key text-green-600"})
		case "air":
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{".air.toml", "fas fa-wind text-blue-400"})
		case "docker":
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{"Dockerfile", "fab fa-docker text-blue-500"})
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{"docker-compose.yml", "fab fa-docker text-blue-500"})
		case "logging":
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{"internal/logger/logger.go", "fas fa-file-alt text-yellow-600"})
		case "testing":
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{"internal/testing/testing.go", "fas fa-vial text-green-600"})
			filePaths = append(filePaths, struct {
				Path string
				Icon string
			}{"internal/service/service_test.go", "fas fa-vial text-green-600"})
		}
	}

	// Build tree structure
	return buildFileTree(filePaths)
}

// buildFileTree creates a hierarchical tree structure from file paths
func buildFileTree(filePaths []struct {
	Path string
	Icon string
}) []FileItem {
	root := &TreeNode{
		Name:     "",
		Path:     "",
		IsFolder: true,
		Children: make(map[string]*TreeNode),
		Level:    -1,
	}

	// Build tree by processing each file path
	for _, file := range filePaths {
		parts := strings.Split(file.Path, "/")
		current := root

		// Navigate/create the path
		for i, part := range parts {
			if part == "" {
				continue
			}

			isLastPart := i == len(parts)-1

			if _, exists := current.Children[part]; !exists {
				current.Children[part] = &TreeNode{
					Name:     part,
					Path:     strings.Join(parts[:i+1], "/"),
					IsFolder: !isLastPart,
					Children: make(map[string]*TreeNode),
					Level:    current.Level + 1,
					Icon:     getFolderIcon(part, !isLastPart),
				}

				// Set file-specific icon for files
				if isLastPart {
					current.Children[part].Icon = file.Icon
				}
			}

			current = current.Children[part]
		}
	}

	// Convert tree to flat list for template
	return flattenTree(root, 0)
}

// getFolderIcon returns appropriate icon for folders and files
func getFolderIcon(name string, isFolder bool) string {
	if isFolder {
		switch name {
		case "cmd":
			return "fas fa-terminal text-blue-600"
		case "internal":
			return "fas fa-folder text-blue-500"
		case "config":
			return "fas fa-cog text-gray-600"
		case "domain":
			return "fas fa-cube text-green-600"
		case "service":
			return "fas fa-server text-purple-600"
		case "handler":
			return "fas fa-hand-paper text-orange-600"
		case "storage":
			return "fas fa-database text-red-600"
		default:
			return "fas fa-folder text-yellow-600"
		}
	}
	return "fas fa-file-code text-gray-500"
}

// flattenTree converts tree structure to flat list while maintaining hierarchy
func flattenTree(node *TreeNode, level int) []FileItem {
	var result []FileItem

	// Sort children - folders first, then files
	var childNames []string
	for name := range node.Children {
		childNames = append(childNames, name)
	}

	sort.Slice(childNames, func(i, j int) bool {
		nodeI := node.Children[childNames[i]]
		nodeJ := node.Children[childNames[j]]

		// Folders before files
		if nodeI.IsFolder != nodeJ.IsFolder {
			return nodeI.IsFolder
		}

		// Alphabetical order within same type
		return childNames[i] < childNames[j]
	})

	for _, name := range childNames {
		child := node.Children[name]

		item := FileItem{
			Name:     child.Name,
			Path:     child.Path,
			Icon:     child.Icon,
			IsFolder: child.IsFolder,
			Level:    level,
		}

		result = append(result, item)

		// Add children recursively
		if child.IsFolder {
			children := flattenTree(child, level+1)
			item.Children = make([]*FileItem, len(children))
			for i := range children {
				item.Children[i] = &children[i]
			}
			result = append(result, children...)
		}
	}

	return result
}

// detectLanguage detects the programming language based on file extension
func detectLanguage(filePath string) string {
	ext := strings.ToLower(filePath[strings.LastIndex(filePath, "."):])

	switch ext {
	case ".go":
		return "go"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".md":
		return "markdown"
	case ".sh":
		return "bash"
	case ".dockerfile":
		return "dockerfile"
	case ".env":
		return "bash"
	default:
		if strings.Contains(filePath, "Makefile") {
			return "makefile"
		}
		if strings.Contains(filePath, "Dockerfile") {
			return "dockerfile"
		}
		return "text"
	}
}

// generateFileContent generates content for a specific file path
// getQueryParam gets a query parameter with a default value
func getQueryParam(r *http.Request, key, defaultValue string) string {
	if value := r.URL.Query().Get(key); value != "" {
		return value
	}
	return defaultValue
}

// generateFileContent creates file items for the file tree modal
func generateFileContent(filePath string) string {
	// This is a deprecated function - use generateFileContentWithConfig instead
	return generateFileContentWithConfig(filePath, metadata.ProjectConfig{
		ProjectName: "my-go-app",
		GoVersion:   "1.23",
		HttpPackage: metadata.Option{ID: "gin", Name: "Gin"},
		Database:    metadata.Option{ID: "postgres", Name: "PostgreSQL"},
		DbDriver:    metadata.Option{ID: "gorm", Name: "GORM"},
	})
}

// generateFileContentWithConfig creates file content based on configuration
func generateFileContentWithConfig(filePath string, config metadata.ProjectConfig) string {

	// Handle specific file patterns
	switch {
	case strings.HasSuffix(filePath, "go.mod"):
		return generateGoModContent(config)

	case strings.HasSuffix(filePath, "main.go"):
		return generateMainContent(config)

	case strings.Contains(filePath, "config.go"):
		return `package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server ServerConfig
	App    AppConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Debug       bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host: getEnv("HOST", "localhost"),
			Port: getEnvAsInt("PORT", 8080),
		},
		App: AppConfig{
			Name:        getEnv("APP_NAME", "my-go-app"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("APP_ENV", "development"),
			Debug:       getEnvAsBool("APP_DEBUG", true),
		},
	}

	return config, nil
}

// Address returns the full server address
func (c *Config) Address() string {
	return c.Server.Host + ":" + strconv.Itoa(c.Server.Port)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(name string, fallback int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

func getEnvAsBool(name string, fallback bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return fallback
}`

	case strings.HasSuffix(filePath, "README.md"):
		return generateReadmeContent(filePath)

	case strings.HasSuffix(filePath, ".gitignore"):
		return `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work
go.work.sum

# Build output
bin/
dist/
build/

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
Thumbs.db

# Environment variables
.env
.env.local

# Database files
*.db
*.sqlite

# Logs
*.log
logs/`

	case strings.HasSuffix(filePath, "Makefile"):
		return `# Makefile for my-go-app

BINARY_NAME=my-go-app
MAIN_PATH=cmd/$(BINARY_NAME)/main.go
BUILD_DIR=bin

.PHONY: build
build: ## Build the application
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: run
run: ## Run the application
	@go run $(MAIN_PATH)

.PHONY: dev
dev: ## Run with hot reload (requires Air)
	@air

.PHONY: test
test: ## Run tests
	@go test ./...

.PHONY: clean
clean: ## Clean build artifacts
	@rm -rf $(BUILD_DIR)

.PHONY: help
help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)`

	case strings.Contains(filePath, "/service/"):
		return generateServiceContent(filePath)
	case strings.Contains(filePath, "/handler/"):
		return generateHandlerContent(filePath)
	case strings.Contains(filePath, "/storage/"):
		return generateStorageContentWithConfig(filePath, config)
	case strings.Contains(filePath, "/domain/"):
		return generateDomainContent(filePath)
	case strings.HasSuffix(filePath, ".env.example"):
		return generateEnvContent()
	case strings.HasSuffix(filePath, ".air.toml"):
		return generateAirContent()
	case strings.HasSuffix(filePath, "Dockerfile"):
		return generateDockerfileContent()
	case strings.HasSuffix(filePath, "docker-compose.yml"):
		return generateDockerComposeContent()
	default:
		// Return file-specific content based on extension
		return generateDefaultContent(filePath)
	}
}

// generateGoModContent creates go.mod content based on configuration
func generateGoModContent(config metadata.ProjectConfig) string {
	content := fmt.Sprintf(`module %s

go %s

require (`, config.ProjectName, config.GoVersion)

	// Add HTTP framework dependency
	switch config.HttpPackage.ID {
	case "gin":
		content += "\n\tgithub.com/gin-gonic/gin v1.9.1"
	case "echo":
		content += "\n\tgithub.com/labstack/echo/v4 v4.11.4"
	case "fiber":
		content += "\n\tgithub.com/gofiber/fiber/v2 v2.52.0"
	case "chi":
		content += "\n\tgithub.com/go-chi/chi/v5 v5.0.11"
	}

	// Add database driver dependencies
	switch config.DbDriver.ID {
	case "gorm":
		content += "\n\tgorm.io/gorm v1.25.5"
		switch config.Database.ID {
		case "postgres":
			content += "\n\tgorm.io/driver/postgres v1.5.4"
		case "mysql":
			content += "\n\tgorm.io/driver/mysql v1.5.2"
		case "sqlite":
			content += "\n\tgorm.io/driver/sqlite v1.5.4"
		}
	case "sqlx":
		content += "\n\tgithub.com/jmoiron/sqlx v1.3.5"
		switch config.Database.ID {
		case "postgres":
			content += "\n\tgithub.com/lib/pq v1.10.9"
		case "mysql":
			content += "\n\tgithub.com/go-sql-driver/mysql v1.7.1"
		case "sqlite":
			content += "\n\tgithub.com/mattn/go-sqlite3 v1.14.18"
		}
	case "mongo-driver":
		content += "\n\tgo.mongodb.org/mongo-driver v1.13.1"
	case "redis-client":
		content += "\n\tgithub.com/redis/go-redis/v9 v9.3.1"
	}

	content += "\n)"
	return content
}

// generateMainContent creates main.go content based on configuration
func generateMainContent(config metadata.ProjectConfig) string {
	var imports []string
	var setupCode string
	var serverCode string

	// Add basic imports
	imports = append(imports, "\"log\"")

	// Add HTTP framework specific imports and setup
	switch config.HttpPackage.ID {
	case "gin":
		imports = append(imports, "\"github.com/gin-gonic/gin\"")
		setupCode = `	// Setup Gin router
	r := gin.Default()
	handler.SetupRoutes(r)`
		serverCode = `	// Start server
	log.Printf("Server starting on %s", cfg.Address())
	if err := r.Run(cfg.Address()); err != nil {
		log.Fatal("Server failed to start:", err)
	}`
	case "echo":
		imports = append(imports, "\"github.com/labstack/echo/v4\"")
		setupCode = `	// Setup Echo router
	e := echo.New()
	handler.SetupRoutes(e)`
		serverCode = `	// Start server
	log.Printf("Server starting on %s", cfg.Address())
	if err := e.Start(cfg.Address()); err != nil {
		log.Fatal("Server failed to start:", err)
	}`
	case "fiber":
		imports = append(imports, "\"github.com/gofiber/fiber/v2\"")
		setupCode = `	// Setup Fiber app
	app := fiber.New()
	handler.SetupRoutes(app)`
		serverCode = `	// Start server
	log.Printf("Server starting on %s", cfg.Address())
	if err := app.Listen(cfg.Address()); err != nil {
		log.Fatal("Server failed to start:", err)
	}`
	default:
		imports = append(imports, "\"net/http\"")
		setupCode = `	// Setup HTTP handlers
	mux := http.NewServeMux()
	handler.SetupRoutes(mux)`
		serverCode = `	// Start server
	log.Printf("Server starting on %s", cfg.Address())
	if err := http.ListenAndServe(cfg.Address(), mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}`
	}

	// Add project imports
	imports = append(imports, fmt.Sprintf("\"%s/internal/config\"", config.ProjectName))
	imports = append(imports, fmt.Sprintf("\"%s/internal/handler\"", config.ProjectName))
	imports = append(imports, fmt.Sprintf("\"%s/internal/service\"", config.ProjectName))

	// Build import block
	importBlock := "import (\n"
	for _, imp := range imports {
		importBlock += "\t" + imp + "\n"
	}
	importBlock += ")"

	return fmt.Sprintf(`package main

%s

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize dependencies
	service := service.NewService()
	handler := handler.NewHandler(service)

%s

%s
}`, importBlock, setupCode, serverCode)
}

// generateReadmeContent creates README.md content
func generateReadmeContent(filePath string) string {
	return `# My Go App

This is a Go web application generated using [go-ctl](https://github.com/syst3mctl/go-ctl).

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or later

### Installation

1. Clone this repository
2. Install dependencies:
   ` + "```bash\n   go mod tidy\n   ```" + `

3. Run the application:
   ` + "```bash\n   go run cmd/my-go-app/main.go\n   ```" + `

The server will start on http://localhost:8080

## 📚 API Documentation

### Health Check
- ` + "`GET /health`" + ` - Returns application health status

### API Routes
- ` + "`GET /api/v1/`" + ` - Welcome message

## 🛠️ Technology Stack

- **Language**: Go 1.23
- **Web Framework**: Gin
- **Architecture**: Clean Architecture

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

---

**Generated with ❤️ by go-ctl**`
}

// generateServiceContent creates service layer content
func generateServiceContent(filePath string) string {
	return `package service

import (
	"context"
	"fmt"

	"my-go-app/internal/domain"
)

// Service provides business logic operations
type Service struct {
	// Add your dependencies here (repositories, etc.)
}

// NewService creates a new service instance
func NewService() *Service {
	return &Service{}
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, user *domain.User) error {
	// Add validation logic here
	if user.Name == "" {
		return fmt.Errorf("user name is required")
	}

	if user.Email == "" {
		return fmt.Errorf("user email is required")
	}

	// Add business logic here
	return nil
}

// GetUserByID retrieves a user by ID
func (s *Service) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	// Add business logic here
	return nil, fmt.Errorf("not implemented")
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(ctx context.Context, user *domain.User) error {
	// Add validation and business logic here
	return fmt.Errorf("not implemented")
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, id uint) error {
	// Add business logic here
	return fmt.Errorf("not implemented")
}

// ListUsers retrieves users with pagination
func (s *Service) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	// Add business logic here
	return nil, fmt.Errorf("not implemented")
}`
}

// generateHandlerContent creates handler layer content
func generateHandlerContent(filePath string) string {
	return `package handler

import (
	"net/http"
	"strconv"

	"my-go-app/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler provides HTTP request handlers
type Handler struct {
	service *service.Service
}

// NewHandler creates a new handler instance
func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// SetupRoutes configures all HTTP routes
func (h *Handler) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", h.HealthCheck)
		api.GET("/", h.Welcome)

		// User routes
		users := api.Group("/users")
		{
			users.POST("", h.CreateUser)
			users.GET("/:id", h.GetUser)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
			users.GET("", h.ListUsers)
		}
	}
}

// HealthCheck returns application health status
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"message": "Service is running",
	})
}

// Welcome returns a welcome message
func (h *Handler) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to my-go-app API",
		"version": "1.0.0",
	})
}

// CreateUser handles user creation
func (h *Handler) CreateUser(c *gin.Context) {
	// Implementation here
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// GetUser handles user retrieval
func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	_ = uint(id)
	// Implementation here
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// UpdateUser handles user updates
func (h *Handler) UpdateUser(c *gin.Context) {
	// Implementation here
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// DeleteUser handles user deletion
func (h *Handler) DeleteUser(c *gin.Context) {
	// Implementation here
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

// ListUsers handles user listing
func (h *Handler) ListUsers(c *gin.Context) {
	// Implementation here
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}`
}

// generateStorageContent creates storage layer content
func generateStorageContent(filePath string) string {
	return generateStorageContentWithConfig(filePath, metadata.ProjectConfig{
		ProjectName: "my-go-app",
		Database:    metadata.Option{ID: "postgres", Name: "PostgreSQL"},
		DbDriver:    metadata.Option{ID: "gorm", Name: "GORM"},
	})
}

// generateStorageContentWithConfig creates storage layer content based on configuration
func generateStorageContentWithConfig(filePath string, config metadata.ProjectConfig) string {
	packageName := "storage"
	if strings.Contains(filePath, "/") {
		parts := strings.Split(filePath, "/")
		for i, part := range parts {
			if part == "storage" && i < len(parts)-1 {
				packageName = parts[i+1]
				break
			}
		}
	}

	var imports []string
	var connectionCode string
	var basicMethods string

	imports = append(imports, "\"context\"", "\"fmt\"", fmt.Sprintf("\"%s/internal/domain\"", config.ProjectName))

	switch config.DbDriver.ID {
	case "gorm":
		imports = append(imports, "\"gorm.io/gorm\"")
		switch config.Database.ID {
		case "postgres":
			imports = append(imports, "\"gorm.io/driver/postgres\"")
		case "mysql":
			imports = append(imports, "\"gorm.io/driver/mysql\"")
		case "sqlite":
			imports = append(imports, "\"gorm.io/driver/sqlite\"")
		}
		connectionCode = `// NewConnection creates a new database connection using GORM
func NewConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}`
		basicMethods = `// CreateUser creates a new user in storage
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID from storage
func (s *Storage) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}`
	case "sqlx":
		imports = append(imports, "\"github.com/jmoiron/sqlx\"")
		connectionCode = `// NewConnection creates a new database connection using sqlx
func NewConnection(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}`
		basicMethods = `// CreateUser creates a new user in storage
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"
	err := s.db.GetContext(ctx, &user.ID, query, user.Name, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}`
	case "mongo-driver":
		imports = append(imports, "\"go.mongodb.org/mongo-driver/mongo\"", "\"go.mongodb.org/mongo-driver/bson\"")
		connectionCode = `// NewConnection creates a new MongoDB connection
func NewConnection(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	return client, nil
}`
		basicMethods = `// CreateUser creates a new user in storage
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	collection := s.database.Collection("users")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}`
	default:
		connectionCode = `// NewConnection creates a new database connection
func NewConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}
	return db, nil
}`
		basicMethods = `// CreateUser creates a new user in storage
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	// Add database operations here
	return fmt.Errorf("not implemented")
}`
	}

	// Build import block
	importBlock := "import (\n"
	for _, imp := range imports {
		importBlock += "\t" + imp + "\n"
	}
	importBlock += ")"

	return fmt.Sprintf(`package %s

%s

// Storage provides data persistence operations
type Storage struct {
	// Add your database connection here
}

// NewStorage creates a new storage instance
func NewStorage() *Storage {
	return &Storage{}
}

%s

%s

// Health checks storage connectivity
func (s *Storage) Health(ctx context.Context) error {
	// Add health check logic here
	return fmt.Errorf("not implemented")
}`, packageName, importBlock, connectionCode, basicMethods)
}

// generateDomainContent creates domain layer content
func generateDomainContent(filePath string) string {
	return `package domain

import (
	"context"
	"fmt"
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint      ` + "`json:\"id\"`" + `
	Name      string    ` + "`json:\"name\"`" + `
	Email     string    ` + "`json:\"email\"`" + `
	Password  string    ` + "`json:\"-\"`" + ` // Never include in JSON responses
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}

	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	// Add more validation rules here
	return nil
}`
}

// generateEnvContent creates .env.example content
func generateEnvContent() string {
	return `# Application Configuration
APP_NAME=my-go-app
APP_VERSION=1.0.0
APP_ENV=development
APP_DEBUG=true

# Server Configuration
HOST=localhost
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=myapp
DB_SSLMODE=disable

# Redis Configuration (if using Redis)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration (if using JWT)
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h

# External APIs
API_KEY=your-api-key-here`
}

// generateAirContent creates .air.toml content
func generateAirContent() string {
	return `root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/my-go-app"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true`
}

// generateDockerfileContent creates Dockerfile content
func generateDockerfileContent() string {
	return `# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/my-go-app

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]`
}

// generateDockerComposeContent creates docker-compose.yml content
func generateDockerComposeContent() string {
	return `version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: myuser
      POSTGRES_PASSWORD: mypassword
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:`
}

// generateDefaultContent creates default content based on file extension
func generateDefaultContent(filePath string) string {
	// Extract filename and extension
	parts := strings.Split(filePath, "/")
	filename := parts[len(parts)-1]

	// Check file extension
	if strings.HasSuffix(filename, ".go") {
		// Extract package name from path
		packageName := "main"
		if len(parts) > 1 {
			packageName = parts[len(parts)-2]
		}

		return fmt.Sprintf(`package %s

// %s contains implementation for %s
// This file will be generated based on your project configuration

import (
	"context"
	"fmt"
)

// TODO: Add your implementation here
`, packageName, filename, packageName)
	}

	if strings.HasSuffix(filename, ".json") {
		return `{
  "name": "my-go-app",
  "version": "1.0.0",
  "description": "Generated Go application"
}`
	}

	if strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml") {
		return `# Configuration file for ` + filename + `
name: my-go-app
version: 1.0.0
description: Generated Go application`
	}

	if strings.HasSuffix(filename, ".toml") {
		return `# Configuration file for ` + filename + `
name = "my-go-app"
version = "1.0.0"
description = "Generated Go application"`
	}

	if strings.HasSuffix(filename, ".md") {
		return `# ` + filename + `

This file is part of the my-go-app project.

## Description

This file will be generated based on your project configuration.`
	}

	// Default fallback
	return fmt.Sprintf(`# Content for %s
# This file will be generated based on your project configuration

`, filePath)
}
