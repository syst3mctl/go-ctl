package generator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// Generator handles project generation
type Generator struct {
	templates map[string]*template.Template
}

// New creates a new Generator instance
func New() *Generator {
	return &Generator{
		templates: make(map[string]*template.Template),
	}
}

// LoadTemplates loads all templates from the templates directory
func (g *Generator) LoadTemplates() error {
	// Define template files and their paths
	templateFiles := map[string]string{
		// Base templates
		"go.mod":    "templates/base/go.mod.tpl",
		"README.md": "templates/base/README.md.tpl",
		"config.go": "templates/base/config.go.tpl",
		"main.go":   "templates/base/main.go.tpl",

		// Domain and core structure
		"domain.model.go":    "templates/base/internal/domain/model.go.tpl",
		"storage.db.go":      "templates/base/internal/storage/db.go.tpl",
		"handler.handler.go": "templates/base/internal/handler/handler.go.tpl",
		"handler.http.go":    "templates/base/internal/handler/http.go.tpl",
		"service.service.go": "templates/base/internal/service/service.go.tpl",

		// Storage repository templates - organized by database type
		"postgres.repository.go": "templates/storage/postgres/repository.go.tpl",
		"mysql.repository.go":    "templates/storage/mysql/repository.go.tpl",
		"sqlite.repository.go":   "templates/storage/sqlite/repository.go.tpl",
		"mongodb.repository.go":  "templates/storage/mongodb/repository.go.tpl",
		"redis.repository.go":    "templates/storage/redis/repository.go.tpl",

		// Feature templates
		"gitignore":       "templates/features/gitignore.tpl",
		"Makefile":        "templates/features/Makefile.tpl",
		"env.example":     "templates/features/env.example.tpl",
		"air.toml":        "templates/features/air.toml.tpl",
		"zerolog.go":      "templates/features/zerolog.go.tpl",
		"testing.go":      "templates/features/testing.go.tpl",
		"service_test.go": "templates/features/service_test.go.tpl",

		// Legacy database templates (to be removed later)
		"gorm.storage.go":         "templates/database/gorm.storage.go.tpl",
		"sqlx.storage.go":         "templates/database/sqlx.storage.go.tpl",
		"database-sql.storage.go": "templates/database/database-sql.storage.go.tpl",
		"mongo-driver.storage.go": "templates/database/mongo-driver.storage.go.tpl",
		"redis-client.storage.go": "templates/database/redis-client.storage.go.tpl",

		// Net/HTTP + Raw SQL pattern templates
		"main-net-http-raw":                "templates/base/main-net-http-raw.tpl",
		"cmd-config.config.go":             "templates/base/cmd-config/config.go.tpl",
		"internal-db.db.go":                "templates/base/internal-db/db.go.tpl",
		"internal-db.redis.go":             "templates/base/internal-db/redis.go.tpl",
		"internal-store.store.go":          "templates/base/internal-store/store.go.tpl",
		"internal-store.user.go":           "templates/base/internal-store/user.go.tpl",
		"internal-store.redis.go":          "templates/base/internal-store/redis.go.tpl",
		"internal-handlers.handler.go":     "templates/base/internal-handlers/handler.go.tpl",
		"internal-handlers.routes.go":      "templates/base/internal-handlers/routes.go.tpl",
		"internal-handlers.middleware.go":  "templates/base/internal-handlers/middleware.go.tpl",
		"internal-handlers.users.go":       "templates/base/internal-handlers/users.go.tpl",
		"internal-handlers-dto.request.go": "templates/base/internal-handlers/dto/request.go.tpl",
		"internal-validate.validate.go":    "templates/base/internal-validate/validate.go.tpl",
		"internal-validate.response.go":    "templates/base/internal-validate/response.go.tpl",
		"internal-gen.gen.go":              "templates/base/internal-gen/gen.go.tpl",
	}

	// Custom template functions
	funcMap := template.FuncMap{
		"title":   strings.Title,
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"replace": strings.ReplaceAll,
		"hasRedis": func(databases []metadata.DatabaseSelection) bool {
			for _, db := range databases {
				if db.Database.ID == "redis" {
					return true
				}
			}
			return false
		},
	}

	// Load each template
	for name, path := range templateFiles {
		tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(path)
		if err != nil {
			// If template file doesn't exist, create a basic one
			tmpl = template.New(name).Funcs(funcMap)
			g.templates[name] = tmpl
			continue
		}

		// Get the template with the base filename (e.g., "zerolog.go.tpl" from the file)
		baseFileName := filepath.Base(path)
		actualTemplate := tmpl.Lookup(baseFileName)
		if actualTemplate != nil {
			// Store the actual template with our desired name
			g.templates[name] = actualTemplate
		} else {
			// If we can't find the template by filename, store the whole template
			g.templates[name] = tmpl
		}
	}

	return nil
}

// GenerateFileContent generates content for a single file based on the configuration
func (g *Generator) GenerateFileContent(filePath string, config metadata.ProjectConfig) (string, error) {
	// Load templates if not already loaded
	if len(g.templates) == 0 {
		if err := g.LoadTemplates(); err != nil {
			return "", fmt.Errorf("failed to load templates: %w", err)
		}
	}

	// Generate the complete project structure to get file content
	projectStructure := g.generateProjectStructure(config)

	// Look for the file in the generated structure
	if content, exists := projectStructure[filePath]; exists {
		return content, nil
	}

	// If exact path not found, try to match by filename or pattern
	for path, content := range projectStructure {
		if strings.HasSuffix(path, filepath.Base(filePath)) {
			return content, nil
		}
	}

	// If still not found, try to generate specific content based on file path patterns
	if content := g.generateSpecificFileContent(filePath, config); content != "" {
		return content, nil
	}

	return "", fmt.Errorf("file not found in project structure: %s", filePath)
}

// generateSpecificFileContent generates content for specific file patterns
func (g *Generator) generateSpecificFileContent(filePath string, config metadata.ProjectConfig) string {
	switch {
	case strings.Contains(filePath, "storage/") && strings.HasSuffix(filePath, ".go"):
		// Generate database-specific storage content
		for _, db := range config.Databases {
			if strings.Contains(filePath, db.Database.ID) {
				templateName := g.getRepositoryTemplate(db.Database.ID)
				return g.renderTemplate(templateName, config)
			}
		}
		// Default storage content
		return g.renderTemplate("storage.db.go", config)

	case strings.Contains(filePath, "handler/") && strings.HasSuffix(filePath, ".go"):
		if strings.Contains(filePath, "http") {
			return g.renderTemplate("handler.http.go", config)
		}
		return g.renderTemplate("handler.handler.go", config)

	case strings.Contains(filePath, "service/") && strings.HasSuffix(filePath, ".go"):
		if strings.Contains(filePath, "test") {
			return g.renderTemplate("service_test.go", config)
		}
		return g.renderTemplate("service.service.go", config)

	case strings.Contains(filePath, "domain/") && strings.HasSuffix(filePath, ".go"):
		return g.renderTemplate("domain.model.go", config)

	case strings.Contains(filePath, "config/") && strings.HasSuffix(filePath, ".go"):
		return g.renderTemplate("config.go", config)

	case strings.Contains(filePath, "testing/") && strings.HasSuffix(filePath, ".go"):
		return g.renderTemplate("testing.go", config)
	}

	return ""
}

// GenerateProjectZip generates a project ZIP file based on the configuration
func (g *Generator) GenerateProjectZip(config metadata.ProjectConfig, w io.Writer) error {
	// Create a new ZIP archive
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// Load templates if not already loaded
	if len(g.templates) == 0 {
		if err := g.LoadTemplates(); err != nil {
			return fmt.Errorf("failed to load templates: %w", err)
		}
	}

	// Create project structure
	projectStructure := g.generateProjectStructure(config)

	// Generate each file in the project structure
	for filePath, content := range projectStructure {
		if err := g.addFileToZip(zipWriter, filePath, content); err != nil {
			return fmt.Errorf("failed to add file %s to zip: %w", filePath, err)
		}
	}

	return nil
}

// GenerateReactProjectZip generates a React project ZIP file
func (g *Generator) GenerateReactProjectZip(config metadata.ProjectConfig, w io.Writer) error {
	// Create a new ZIP archive
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// Generate React project structure
	projectStructure := g.GenerateReactProject(config)

	// Generate each file in the project structure
	for filePath, content := range projectStructure {
		if err := g.addFileToZip(zipWriter, filePath, content); err != nil {
			return fmt.Errorf("failed to add file %s to zip: %w", filePath, err)
		}
	}

	return nil
}

// isNetHTTPRawSQLPattern checks if the configuration uses net/http with database/sql
func (g *Generator) isNetHTTPRawSQLPattern(config metadata.ProjectConfig) bool {
	if config.HttpPackage.ID != "net-http" {
		return false
	}
	if len(config.Databases) == 0 {
		return false
	}
	// Check if any database uses database-sql driver
	for _, dbSelection := range config.Databases {
		if dbSelection.Driver.ID == "database-sql" {
			return true
		}
	}
	return false
}

// generateProjectStructure creates the complete project file structure
func (g *Generator) generateProjectStructure(config metadata.ProjectConfig) map[string]string {
	// Always use the net/http + raw SQL pattern structure for all projects
	// This provides a consistent structure regardless of HTTP framework or database driver
	return g.generateNetHTTPRawSQLStructure(config)
}

// generateNetHTTPRawSQLStructure creates project structure for net/http + raw SQL pattern
func (g *Generator) generateNetHTTPRawSQLStructure(config metadata.ProjectConfig) map[string]string {
	files := make(map[string]string)

	// Base files
	files["go.mod"] = g.renderTemplate("go.mod", config)
	files["README.md"] = g.renderTemplate("README.md", config)
	files[fmt.Sprintf("cmd/%s/main.go", config.ProjectName)] = g.renderTemplate("main-net-http-raw", config)

	// Configuration in cmd/config
	files["cmd/config/config.go"] = g.renderTemplate("cmd-config.config.go", config)

	// Database connection
	hasRedis := false
	for _, dbSelection := range config.Databases {
		if dbSelection.Database.ID == "redis" {
			hasRedis = true
			break
		}
	}

	if len(config.Databases) > 0 {
		// Only generate SQL database files if there's at least one SQL database
		hasSQLDB := false
		for _, dbSelection := range config.Databases {
			if dbSelection.Database.ID != "redis" && dbSelection.Database.ID != "mongodb" {
				hasSQLDB = true
				break
			}
		}

		if hasSQLDB {
			files["internal/db/db.go"] = g.renderTemplate("internal-db.db.go", config)
			// Store layer
			files["internal/store/store.go"] = g.renderTemplate("internal-store.store.go", config)
			files["internal/store/user.go"] = g.renderTemplate("internal-store.user.go", config)
		}
	}

	// Redis files
	if hasRedis {
		files["internal/db/redis.go"] = g.renderTemplate("internal-db.redis.go", config)
		files["internal/store/redis.go"] = g.renderTemplate("internal-store.redis.go", config)
	}

	// Handlers
	files["internal/handlers/handler.go"] = g.renderTemplate("internal-handlers.handler.go", config)
	files["internal/handlers/routes.go"] = g.renderTemplate("internal-handlers.routes.go", config)
	files["internal/handlers/middleware.go"] = g.renderTemplate("internal-handlers.middleware.go", config)
	files["internal/handlers/users.go"] = g.renderTemplate("internal-handlers.users.go", config)
	files["internal/handlers/dto/request.go"] = g.renderTemplate("internal-handlers-dto.request.go", config)

	// Validation
	files["internal/validate/validate.go"] = g.renderTemplate("internal-validate.validate.go", config)
	files["internal/validate/response.go"] = g.renderTemplate("internal-validate.response.go", config)

	// Utilities
	files["internal/gen/gen.go"] = g.renderTemplate("internal-gen.gen.go", config)

	// Domain models (use existing template)
	files["internal/domain/model.go"] = g.renderTemplate("domain.model.go", config)

	// Feature files
	for _, feature := range config.Features {
		switch feature.ID {
		case "gitignore":
			files[".gitignore"] = g.renderTemplate("gitignore", config)
		case "makefile":
			files["Makefile"] = g.renderTemplate("Makefile", config)
		case "env":
			files[".env.example"] = g.renderTemplate("env.example", config)
		case "air":
			files[".air.toml"] = g.renderTemplate("air.toml", config)
		case "docker":
			files["Dockerfile"] = g.generateDockerfile(config)
			files["docker-compose.yml"] = g.generateDockerCompose(config)
		case "logging":
			files["internal/logger/logger.go"] = g.renderTemplate("zerolog.go", config)
		case "testing":
			files["internal/testing/testing.go"] = g.renderTemplate("testing.go", config)
		}
	}

	return files
}

// getRepositoryTemplate returns the appropriate repository template based on database type
func (g *Generator) getRepositoryTemplate(databaseID string) string {
	switch databaseID {
	case "postgres":
		return "postgres.repository.go"
	case "mysql":
		return "mysql.repository.go"
	case "sqlite":
		return "sqlite.repository.go"
	case "mongodb":
		return "mongodb.repository.go"
	case "redis":
		return "redis.repository.go"
	default:
		return "postgres.repository.go" // Default to PostgreSQL
	}
}

// renderTemplate renders a template with the given configuration
func (g *Generator) renderTemplate(templateName string, config metadata.ProjectConfig) string {
	tmpl, exists := g.templates[templateName]
	if !exists {
		return fmt.Sprintf("// Template %s not found\npackage main\n\nfunc main() {\n\t// TODO: Implement\n}\n", templateName)
	}

	// Create enhanced template with basic function map
	enhancedFuncMap := template.FuncMap{
		"title":   strings.Title,
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"replace": strings.ReplaceAll,
		"hasRedis": func(databases []metadata.DatabaseSelection) bool {
			for _, db := range databases {
				if db.Database.ID == "redis" {
					return true
				}
			}
			return false
		},
	}

	// Clone template with enhanced function map
	enhancedTemplate := template.Must(tmpl.Clone())
	enhancedTemplate.Funcs(enhancedFuncMap)

	// Get primary database and driver for template compatibility
	// For SQL database templates, use the first SQL database (not Redis/MongoDB)
	var database, dbDriver metadata.Option
	if len(config.Databases) > 0 {
		// Check if this is a SQL database template
		isSQLTemplate := strings.Contains(templateName, "internal-db.db.go") ||
			strings.Contains(templateName, "internal-store.store.go") ||
			strings.Contains(templateName, "internal-store.user.go") ||
			strings.Contains(templateName, "cmd-config.config.go")

		if isSQLTemplate {
			// Find first SQL database (not Redis or MongoDB)
			for _, dbSelection := range config.Databases {
				if dbSelection.Database.ID != "redis" && dbSelection.Database.ID != "mongodb" {
					database = dbSelection.Database
					dbDriver = dbSelection.Driver
					break
				}
			}
			// Fallback to first database if no SQL database found
			if database.ID == "" && len(config.Databases) > 0 {
				database = config.Databases[0].Database
				dbDriver = config.Databases[0].Driver
			}
		} else {
			// For non-SQL templates, use first database
			database = config.Databases[0].Database
			dbDriver = config.Databases[0].Driver
		}
	}

	// Create a wrapper with HasFeature method
	dataWithMethods := &TemplateData{
		ProjectConfig: config,
		HTTP:          config.HttpPackage,
		Database:      database,
		DbDriver:      dbDriver,
		generator:     g,
	}

	var buf bytes.Buffer
	if err := enhancedTemplate.Execute(&buf, dataWithMethods); err != nil {
		return fmt.Sprintf("// Error rendering template %s: %v\npackage main\n\nfunc main() {\n\t// TODO: Fix template\n}\n", templateName, err)
	}

	return buf.String()
}

// TemplateData wraps ProjectConfig with methods for template usage
type TemplateData struct {
	metadata.ProjectConfig
	HTTP      metadata.Option
	Database  metadata.Option
	DbDriver  metadata.Option
	generator *Generator
}

// HasFeature checks if a feature is enabled in the configuration
func (td *TemplateData) HasFeature(featureID string) bool {
	for _, feature := range td.Features {
		if feature.ID == featureID {
			return true
		}
	}
	return false
}

// hasFeature checks if a feature is enabled in the configuration (legacy method)
func (g *Generator) hasFeature(config metadata.ProjectConfig, featureID string) bool {
	for _, feature := range config.Features {
		if feature.ID == featureID {
			return true
		}
	}
	return false
}

// generateDockerfile creates a Dockerfile
func (g *Generator) generateDockerfile(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`# Build stage
FROM golang:%s-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o bin/%s main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/%s .

EXPOSE 8080
CMD ["./%s"]
`, config.GoVersion, config.ProjectName, config.ProjectName, config.ProjectName)
}

// generateDockerCompose creates a docker-compose.yml file
func (g *Generator) generateDockerCompose(config metadata.ProjectConfig) string {
	var dbServices []string
	var volumes []string
	var dependsOnServices []string

	for _, dbSelection := range config.Databases {
		switch dbSelection.Database.ID {
		case "postgres":
			dbServices = append(dbServices, `
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: `+config.ProjectName+`
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data`)
			volumes = append(volumes, "postgres_data:")
			dependsOnServices = append(dependsOnServices, "postgres")

		case "mysql":
			dbServices = append(dbServices, `
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: `+config.ProjectName+`
      MYSQL_ROOT_PASSWORD: password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql`)
			volumes = append(volumes, "mysql_data:")
			dependsOnServices = append(dependsOnServices, "mysql")

		case "mongodb":
			dbServices = append(dbServices, `
  mongodb:
    image: mongo:6
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db`)
			volumes = append(volumes, "mongo_data:")
			dependsOnServices = append(dependsOnServices, "mongodb")

		case "redis":
			dbServices = append(dbServices, `
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data`)
			volumes = append(volumes, "redis_data:")
			dependsOnServices = append(dependsOnServices, "redis")
		}
	}

	dbService := strings.Join(dbServices, "")
	volumeService := ""
	if len(volumes) > 0 {
		volumeService = "\n\nvolumes:\n  " + strings.Join(volumes, "\n  ")
	}

	dependsOn := ""
	if len(dependsOnServices) > 0 {
		dependsOn = "    depends_on:\n      - " + strings.Join(dependsOnServices, "\n      - ")
	}

	return fmt.Sprintf(`version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
%s%s%s`, dependsOn, dbService, volumeService)
}

// addFileToZip adds a file with content to the ZIP archive
func (g *Generator) addFileToZip(zipWriter *zip.Writer, filePath, content string) error {
	// Create the file in the ZIP
	fileWriter, err := zipWriter.Create(filePath)
	if err != nil {
		return err
	}

	// Write content to the file
	_, err = fileWriter.Write([]byte(content))
	return err
}
