package generator

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
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

		// Feature templates
		"gitignore":   "templates/features/gitignore.tpl",
		"Makefile":    "templates/features/Makefile.tpl",
		"env.example": "templates/features/env.example.tpl",
		"air.toml":    "templates/features/air.toml.tpl",

		// HTTP framework templates
		"gin.main.go":      "templates/http/gin.main.go.tpl",
		"echo.main.go":     "templates/http/echo.main.go.tpl",
		"fiber.main.go":    "templates/http/fiber.main.go.tpl",
		"chi.main.go":      "templates/http/chi.main.go.tpl",
		"net-http.main.go": "templates/http/net-http.main.go.tpl",

		// Database templates
		"gorm.storage.go":         "templates/database/gorm.storage.go.tpl",
		"sqlx.storage.go":         "templates/database/sqlx.storage.go.tpl",
		"database-sql.storage.go": "templates/database/database-sql.storage.go.tpl",
		"mongo-driver.storage.go": "templates/database/mongo-driver.storage.go.tpl",
		"redis-client.storage.go": "templates/database/redis-client.storage.go.tpl",
	}

	// Custom template functions
	funcMap := template.FuncMap{
		"hasFeature": g.hasFeature,
		"title":      strings.Title,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"replace":    strings.ReplaceAll,
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
		g.templates[name] = tmpl
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
	files["internal/config/config.go"] = g.renderTemplate("config.go", config)

	// Main application file based on HTTP framework
	mainFile := g.getMainTemplate(config.HttpPackage.ID)
	files[fmt.Sprintf("cmd/%s/main.go", config.ProjectName)] = g.renderTemplate(mainFile, config)

	// Database layer if specified
	if config.DbDriver.ID != "" {
		storageFile := g.getStorageTemplate(config.DbDriver.ID)
		files[fmt.Sprintf("internal/storage/%s/%s.go", config.DbDriver.ID, config.DbDriver.ID)] = g.renderTemplate(storageFile, config)
	}

	// Domain layer (basic structure)
	files["internal/domain/model.go"] = g.generateDomainModel(config)
	files["internal/service/service.go"] = g.generateService(config)
	files["internal/handler/handler.go"] = g.generateHandler(config)

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
		}
	}

	return files
}

// getMainTemplate returns the appropriate main.go template based on HTTP framework
func (g *Generator) getMainTemplate(httpID string) string {
	switch httpID {
	case "gin":
		return "gin.main.go"
	case "echo":
		return "echo.main.go"
	case "fiber":
		return "fiber.main.go"
	case "chi":
		return "chi.main.go"
	case "net-http":
		return "net-http.main.go"
	default:
		return "gin.main.go" // Default to Gin
	}
}

// getStorageTemplate returns the appropriate storage template based on database driver
func (g *Generator) getStorageTemplate(driverID string) string {
	switch driverID {
	case "gorm":
		return "gorm.storage.go"
	case "sqlx":
		return "sqlx.storage.go"
	case "database-sql":
		return "database-sql.storage.go"
	case "mongo-driver":
		return "mongo-driver.storage.go"
	case "redis-client":
		return "redis-client.storage.go"
	default:
		return "gorm.storage.go" // Default to GORM
	}
}

// renderTemplate renders a template with the given configuration
func (g *Generator) renderTemplate(templateName string, config metadata.ProjectConfig) string {
	tmpl, exists := g.templates[templateName]
	if !exists {
		return fmt.Sprintf("// Template %s not found\npackage main\n\nfunc main() {\n\t// TODO: Implement\n}\n", templateName)
	}

	// Create template data with helper methods
	data := struct {
		metadata.ProjectConfig
		HasFeature func(string) bool
	}{
		ProjectConfig: config,
		HasFeature:    func(featureID string) bool { return g.hasFeature(config, featureID) },
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
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

// generateDomainModel creates a basic domain model file
func (g *Generator) generateDomainModel(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`package domain

import (
	"time"
%s)

// User represents a user in the system
type User struct {
	ID        uint      `+"`json:\"id\"`"+`
	CreatedAt time.Time `+"`json:\"created_at\"`"+`
	UpdatedAt time.Time `+"`json:\"updated_at\"`"+`
	Name      string    `+"`json:\"name\"`"+`
	Email     string    `+"`json:\"email\"`"+`
	Active    bool      `+"`json:\"active\"`"+`
}

// TODO: Add your domain models here
`, g.getImportsForDomain(config))
}

// generateService creates a basic service layer file
func (g *Generator) generateService(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`package service

import (
	"context"

	"%s/internal/domain"
%s)

// Service defines the business logic interface
type Service interface {
	// User operations
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id uint) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context) ([]*domain.User, error)

	// TODO: Add your service methods here
}

// service implements the Service interface
type service struct {
%s}

// New creates a new service instance
func New(%s) Service {
	return &service{
%s	}
}

// CreateUser creates a new user
func (s *service) CreateUser(ctx context.Context, user *domain.User) error {
	// TODO: Implement user creation logic
	return nil
}

// GetUser retrieves a user by ID
func (s *service) GetUser(ctx context.Context, id uint) (*domain.User, error) {
	// TODO: Implement user retrieval logic
	return nil, nil
}

// UpdateUser updates an existing user
func (s *service) UpdateUser(ctx context.Context, user *domain.User) error {
	// TODO: Implement user update logic
	return nil
}

// DeleteUser deletes a user by ID
func (s *service) DeleteUser(ctx context.Context, id uint) error {
	// TODO: Implement user deletion logic
	return nil
}

// ListUsers retrieves all users
func (s *service) ListUsers(ctx context.Context) ([]*domain.User, error) {
	// TODO: Implement user listing logic
	return nil, nil
}
`, config.ProjectName, g.getImportsForService(config), g.getServiceFields(config), g.getServiceConstructorParams(config), g.getServiceConstructorFields(config))
}

// generateHandler creates a basic handler file
func (g *Generator) generateHandler(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`package handler

import (
	"net/http"
	"strconv"

	"%s/internal/config"
	"%s/internal/service"
	"%s/internal/domain"
%s)

// Handler contains all HTTP handlers
type Handler struct {
	service service.Service
	config  *config.Config
}

// New creates a new Handler instance
func New(svc service.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: svc,
		config:  cfg,
	}
}

%s

// Example handlers for User operations
%s
`, config.ProjectName, config.ProjectName, config.ProjectName, g.getImportsForHandler(config), g.getHandlerMethods(config), g.getUserHandlers(config))
}

// generateDockerfile creates a Dockerfile
func (g *Generator) generateDockerfile(config metadata.ProjectConfig) string {
	return fmt.Sprintf(`# Build stage
FROM golang:%s-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-w -s" -o bin/%s cmd/%s/main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/%s .
%s

EXPOSE 8080
CMD ["./%s"]
`, config.GoVersion, config.ProjectName, config.ProjectName, config.ProjectName, g.getDockerEnvFiles(config), config.ProjectName)
}

// generateDockerCompose creates a docker-compose.yml file
func (g *Generator) generateDockerCompose(config metadata.ProjectConfig) string {
	services := fmt.Sprintf(`version: '3.8'

services:
  %s:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
%s    depends_on:
%s
%s`, config.ProjectName, g.getComposeEnvVars(config), g.getComposeDependencies(config), g.getComposeServices(config))

	return services
}

// Helper methods for generating imports and fields based on configuration

func (g *Generator) getImportsForDomain(config metadata.ProjectConfig) string {
	var imports []string

	if config.DbDriver.ID == "gorm" {
		imports = append(imports, "\t\"gorm.io/gorm\"")
	}

	if len(imports) > 0 {
		return "\n" + strings.Join(imports, "\n")
	}
	return ""
}

func (g *Generator) getImportsForService(config metadata.ProjectConfig) string {
	var imports []string

	if config.DbDriver.ID != "" {
		imports = append(imports, fmt.Sprintf("\t\"%s/internal/storage/%s\"", config.ProjectName, config.DbDriver.ID))
	}

	if len(imports) > 0 {
		return "\n" + strings.Join(imports, "\n")
	}
	return ""
}

func (g *Generator) getImportsForHandler(config metadata.ProjectConfig) string {
	var imports []string

	switch config.HttpPackage.ID {
	case "gin":
		imports = append(imports, "\t\"github.com/gin-gonic/gin\"")
	case "echo":
		imports = append(imports, "\t\"github.com/labstack/echo/v4\"")
	case "fiber":
		imports = append(imports, "\t\"github.com/gofiber/fiber/v2\"")
	case "chi":
		imports = append(imports, "\t\"github.com/go-chi/chi/v5\"")
	default:
		imports = append(imports, "\t\"encoding/json\"")
	}

	if len(imports) > 0 {
		return "\n" + strings.Join(imports, "\n")
	}
	return ""
}

func (g *Generator) getServiceFields(config metadata.ProjectConfig) string {
	if config.DbDriver.ID != "" {
		switch config.DbDriver.ID {
		case "gorm":
			return "\tdb *gorm.DB"
		case "sqlx":
			return "\tdb *sqlx.DB"
		case "mongo-driver":
			return "\tdb *mongo.Client"
		case "redis-client":
			return "\tdb *redis.Client"
		default:
			return "\tdb interface{}"
		}
	}
	return "\t// Add your dependencies here"
}

func (g *Generator) getServiceConstructorParams(config metadata.ProjectConfig) string {
	if config.DbDriver.ID != "" {
		switch config.DbDriver.ID {
		case "gorm":
			return "db *gorm.DB"
		case "sqlx":
			return "db *sqlx.DB"
		case "mongo-driver":
			return "db *mongo.Client"
		case "redis-client":
			return "db *redis.Client"
		default:
			return "db interface{}"
		}
	}
	return ""
}

func (g *Generator) getServiceConstructorFields(config metadata.ProjectConfig) string {
	if config.DbDriver.ID != "" {
		return "\t\tdb: db,\n"
	}
	return ""
}

func (g *Generator) getHandlerMethods(config metadata.ProjectConfig) string {
	switch config.HttpPackage.ID {
	case "gin":
		return `// GetProfile returns user profile (example protected endpoint)
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "Profile endpoint",
	})
}

// CreateData creates new data (example protected endpoint)
func (h *Handler) CreateData(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Data created successfully"})
}`
	default:
		return `// GetProfile returns user profile (example protected endpoint)
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile endpoint"})
}

// CreateData creates new data (example protected endpoint)
func (h *Handler) CreateData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Data created successfully"})
}`
	}
}

func (g *Generator) getUserHandlers(config metadata.ProjectConfig) string {
	switch config.HttpPackage.ID {
	case "gin":
		return `
func (h *Handler) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}`
	default:
		return `
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path - implementation depends on router
	idStr := r.URL.Path[len("/api/v1/users/"):]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(r.Context(), uint(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}`
	}
}

func (g *Generator) getDockerEnvFiles(config metadata.ProjectConfig) string {
	if g.hasFeature(config, "env") {
		return "COPY .env.example .env"
	}
	return ""
}

func (g *Generator) getComposeEnvVars(config metadata.ProjectConfig) string {
	var envVars []string

	if config.Database.ID != "" {
		switch config.Database.ID {
		case "postgres":
			envVars = append(envVars, "      - DATABASE_URL=postgres://postgres:password@postgres:5432/"+config.ProjectName+"_db?sslmode=disable")
		case "mysql":
			envVars = append(envVars, "      - DATABASE_URL=root:password@tcp(mysql:3306)/"+config.ProjectName+"_db?parseTime=true")
		case "mongodb":
			envVars = append(envVars, "      - MONGO_URI=mongodb://mongo:27017/"+config.ProjectName+"_db")
		case "redis":
			envVars = append(envVars, "      - REDIS_URL=redis://redis:6379/0")
		}
	}

	if len(envVars) > 0 {
		return strings.Join(envVars, "\n") + "\n"
	}
	return ""
}

func (g *Generator) getComposeDependencies(config metadata.ProjectConfig) string {
	var deps []string

	if config.Database.ID != "" {
		switch config.Database.ID {
		case "postgres":
			deps = append(deps, "      - postgres")
		case "mysql":
			deps = append(deps, "      - mysql")
		case "mongodb":
			deps = append(deps, "      - mongo")
		case "redis":
			deps = append(deps, "      - redis")
		}
	}

	if len(deps) > 0 {
		return strings.Join(deps, "\n")
	}
	return "      # Add database dependencies here"
}

func (g *Generator) getComposeServices(config metadata.ProjectConfig) string {
	var services []string

	if config.Database.ID != "" {
		switch config.Database.ID {
		case "postgres":
			services = append(services, `
  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=`+config.ProjectName+`_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data`)
		case "mysql":
			services = append(services, `
  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_DATABASE=`+config.ProjectName+`_db
      - MYSQL_ROOT_PASSWORD=password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql`)
		case "mongodb":
			services = append(services, `
  mongo:
    image: mongo:6.0
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db`)
		case "redis":
			services = append(services, `
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data`)
		}
	}

	if len(services) > 0 {
		result := strings.Join(services, "")
		result += "\n\nvolumes:"

		if config.Database.ID == "postgres" {
			result += "\n  postgres_data:"
		} else if config.Database.ID == "mysql" {
			result += "\n  mysql_data:"
		} else if config.Database.ID == "mongodb" {
			result += "\n  mongo_data:"
		} else if config.Database.ID == "redis" {
			result += "\n  redis_data:"
		}

		return result
	}

	return ""
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
