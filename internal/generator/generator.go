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
	}

	// Custom template functions
	funcMap := template.FuncMap{
		"title":   strings.Title,
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"replace": strings.ReplaceAll,
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

// generateProjectStructure creates the complete project file structure
func (g *Generator) generateProjectStructure(config metadata.ProjectConfig) map[string]string {
	files := make(map[string]string)

	// Base files
	files["go.mod"] = g.renderTemplate("go.mod", config)
	files["README.md"] = g.renderTemplate("README.md", config)
	files["main.go"] = g.renderTemplate("main.go", config)

	// Core structure
	files["internal/config/config.go"] = g.renderTemplate("config.go", config)
	files["internal/domain/model.go"] = g.renderTemplate("domain.model.go", config)
	files["internal/service/service.go"] = g.renderTemplate("service.service.go", config)
	files["internal/handler/handler.go"] = g.renderTemplate("handler.handler.go", config)
	files["internal/handler/http.go"] = g.renderTemplate("handler.http.go", config)

	// Database layer if specified
	if len(config.Databases) > 0 {
		// Database initialization
		files["internal/storage/db.go"] = g.renderTemplate("storage.db.go", config)

		// Repository implementations - organized by database type
		for _, dbSelection := range config.Databases {
			repositoryFile := g.getRepositoryTemplate(dbSelection.Database.ID)
			files[fmt.Sprintf("internal/storage/%s/repository.go", dbSelection.Database.ID)] = g.renderTemplate(repositoryFile, config)
		}
	}

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
			files["internal/service/service_test.go"] = g.renderTemplate("service_test.go", config)
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

	// Create enhanced template with updated function map
	enhancedFuncMap := template.FuncMap{
		"title":      strings.Title,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"replace":    strings.ReplaceAll,
		"HasFeature": func(featureID string) bool { return g.hasFeature(config, featureID) },
	}

	// Clone template with enhanced function map
	enhancedTemplate := template.Must(tmpl.Clone())
	enhancedTemplate.Funcs(enhancedFuncMap)

	// Create template data with template-expected field names
	data := struct {
		metadata.ProjectConfig
		HTTP metadata.Option
	}{
		ProjectConfig: config,
		HTTP:          config.HttpPackage, // Map HttpPackage to HTTP for template compatibility
	}

	var buf bytes.Buffer
	if err := enhancedTemplate.Execute(&buf, data); err != nil {
		return fmt.Sprintf("// Error rendering template %s: %v\npackage main\n\nfunc main() {\n\t// TODO: Fix template\n}\n", templateName, err)
	}

	return buf.String()
}

// hasFeature checks if a feature is enabled in the configuration
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
