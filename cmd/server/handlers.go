package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// FileItem represents a file in the project structure
type FileItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Icon string `json:"icon"`
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
	tmpl := template.Must(template.New("explore").Parse(exploreTemplate))

	data := ProjectStructureData{
		Files:  fileItems,
		Config: config,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render explore template", http.StatusInternalServerError)
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

	// Parse the session data to rebuild the config (in production, you'd store this in session)
	// For now, we'll generate content based on the file path
	content := generateFileContent(filePath)

	// Detect language for syntax highlighting
	language := detectLanguage(filePath)

	// Return HTML with syntax highlighting
	tmpl := template.Must(template.New("file-content").Parse(fileContentTemplate))

	data := struct {
		Content  string
		Language string
	}{
		Content:  content,
		Language: language,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render file content", http.StatusInternalServerError)
		return
	}
}

// searchPackages performs a simplified package search
// In a production version, this would call the actual pkg.go.dev API
func searchPackages(query string) ([]PkgGoDevResult, error) {
	// Simplified mock implementation
	// In production, you would call: https://pkg.go.dev/search?q=query&m=json

	mockResults := []PkgGoDevResult{
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
			Path:     "gorm.io/gorm",
			Synopsis: "The fantastic ORM library for Golang",
		},
		{
			Path:     "github.com/jmoiron/sqlx",
			Synopsis: "general purpose extensions to golang's database/sql",
		},
	}

	// Filter results based on query (simple contains check)
	var results []PkgGoDevResult
	query = strings.ToLower(query)
	for _, result := range mockResults {
		if strings.Contains(strings.ToLower(result.Path), query) ||
			strings.Contains(strings.ToLower(result.Synopsis), query) {
			results = append(results, result)
		}
	}

	// Limit results to avoid overwhelming the UI
	if len(results) > 10 {
		results = results[:10]
	}

	return results, nil
}

// generateFileItems creates file items for the file tree modal
func generateFileItems(config metadata.ProjectConfig) []FileItem {
	var files []FileItem

	// Base files
	files = append(files, FileItem{
		Name: "go.mod",
		Path: "go.mod",
		Icon: "fas fa-cube text-green-500",
	})

	files = append(files, FileItem{
		Name: "README.md",
		Path: "README.md",
		Icon: "fab fa-markdown text-blue-600",
	})

	// Main application file
	mainFile := fmt.Sprintf("cmd/%s/main.go", config.ProjectName)
	files = append(files, FileItem{
		Name: "main.go",
		Path: mainFile,
		Icon: "fab fa-golang text-blue-500",
	})

	// Internal structure files
	files = append(files, FileItem{
		Name: "config.go",
		Path: "internal/config/config.go",
		Icon: "fas fa-cog text-gray-600",
	})

	files = append(files, FileItem{
		Name: "model.go",
		Path: "internal/domain/model.go",
		Icon: "fab fa-golang text-blue-500",
	})

	files = append(files, FileItem{
		Name: "service.go",
		Path: "internal/service/service.go",
		Icon: "fab fa-golang text-blue-500",
	})

	files = append(files, FileItem{
		Name: "handler.go",
		Path: "internal/handler/handler.go",
		Icon: "fab fa-golang text-blue-500",
	})

	// Database storage layer if configured
	if config.DbDriver.ID != "" {
		storageFile := fmt.Sprintf("internal/storage/%s/%s.go", config.DbDriver.ID, config.DbDriver.ID)
		files = append(files, FileItem{
			Name: fmt.Sprintf("%s.go", config.DbDriver.ID),
			Path: storageFile,
			Icon: "fas fa-database text-purple-500",
		})
	}

	// Feature files
	for _, feature := range config.Features {
		switch feature.ID {
		case "gitignore":
			files = append(files, FileItem{
				Name: ".gitignore",
				Path: ".gitignore",
				Icon: "fab fa-git-alt text-orange-500",
			})
		case "makefile":
			files = append(files, FileItem{
				Name: "Makefile",
				Path: "Makefile",
				Icon: "fas fa-hammer text-gray-600",
			})
		case "env":
			files = append(files, FileItem{
				Name: ".env.example",
				Path: ".env.example",
				Icon: "fas fa-key text-green-600",
			})
		case "air":
			files = append(files, FileItem{
				Name: ".air.toml",
				Path: ".air.toml",
				Icon: "fas fa-wind text-blue-400",
			})
		case "docker":
			files = append(files, FileItem{
				Name: "Dockerfile",
				Path: "Dockerfile",
				Icon: "fab fa-docker text-blue-500",
			})
			files = append(files, FileItem{
				Name: "docker-compose.yml",
				Path: "docker-compose.yml",
				Icon: "fab fa-docker text-blue-500",
			})
		}
	}

	return files
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
func generateFileContent(filePath string) string {
	// This is a simplified version - in production you'd generate actual content
	// based on the project configuration stored in session/context

	switch {
	case strings.HasSuffix(filePath, "go.mod"):
		return `module my-go-app

go 1.23

require (
	github.com/gin-gonic/gin v1.9.1
	// Other dependencies will be added based on your selections
)`

	case strings.HasSuffix(filePath, "main.go"):
		return `package main

import (
	"log"
	"net/http"

	"my-go-app/internal/config"
	"my-go-app/internal/handler"
	"my-go-app/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize services
	svc := service.New()

	// Initialize handlers
	handlers := handler.New(svc, cfg)

	// Setup router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	api.GET("/", handlers.Welcome)

	// Start server
	log.Printf("Starting server on %s", cfg.Address())
	log.Fatal(r.Run(cfg.Address()))
}`

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

	case strings.Contains(filePath, "README.md"):
		return `# My Go App

This is a Go web application generated using [go-ctl](https://github.com/syst3mctl/go-ctl).

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or later

### Installation

1. Clone this repository
2. Install dependencies:
   ` + "```" + `bash
   go mod tidy
   ` + "```" + `

3. Run the application:
   ` + "```" + `bash
   go run cmd/my-go-app/main.go
   ` + "```" + `

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

	default:
		return `// Content for ` + filePath + `
// This file will be generated based on your project configuration

package main

func main() {
	// TODO: Implement
}`
	}
}
