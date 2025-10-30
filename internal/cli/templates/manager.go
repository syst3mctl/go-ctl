package templates

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/syst3mctl/go-ctl/internal/metadata"
)

// CustomTemplate represents a user-defined template
type CustomTemplate struct {
	ID          string                 `yaml:"id" json:"id"`
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description" json:"description"`
	Author      string                 `yaml:"author" json:"author"`
	Version     string                 `yaml:"version" json:"version"`
	Tags        []string               `yaml:"tags" json:"tags"`
	Config      metadata.ProjectConfig `yaml:"config" json:"config"`
	Files       map[string]string      `yaml:"files" json:"files"`
	Templates   map[string]string      `yaml:"templates" json:"templates"`
	Metadata    CustomTemplateMetadata `yaml:"metadata" json:"metadata"`
	CreatedAt   time.Time              `yaml:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `yaml:"updated_at" json:"updated_at"`
}

// CustomTemplateMetadata contains additional template information
type CustomTemplateMetadata struct {
	SourceProject string            `yaml:"source_project" json:"source_project"`
	Repository    string            `yaml:"repository" json:"repository"`
	Homepage      string            `yaml:"homepage" json:"homepage"`
	License       string            `yaml:"license" json:"license"`
	MinGoVersion  string            `yaml:"min_go_version" json:"min_go_version"`
	Keywords      []string          `yaml:"keywords" json:"keywords"`
	Variables     map[string]string `yaml:"variables" json:"variables"`
}

// TemplateManager manages custom templates
type TemplateManager struct {
	templatesDir string
	configDir    string
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() *TemplateManager {
	homeDir, _ := os.UserHomeDir()

	return &TemplateManager{
		templatesDir: filepath.Join(homeDir, ".go-ctl", "templates"),
		configDir:    filepath.Join(homeDir, ".go-ctl"),
	}
}

// NewTemplateManagerWithPath creates a template manager with custom path
func NewTemplateManagerWithPath(templatesDir string) *TemplateManager {
	return &TemplateManager{
		templatesDir: templatesDir,
		configDir:    filepath.Dir(templatesDir),
	}
}

// GetTemplatesDir returns the templates directory path
func (tm *TemplateManager) GetTemplatesDir() string {
	return tm.templatesDir
}

// CreateTemplate creates a new custom template
func (tm *TemplateManager) CreateTemplate(template *CustomTemplate) error {
	if template.ID == "" {
		return fmt.Errorf("template ID is required")
	}

	// Validate template ID
	if err := tm.validateTemplateID(template.ID); err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	// Check if template already exists
	if tm.TemplateExists(template.ID) {
		return fmt.Errorf("template '%s' already exists", template.ID)
	}

	// Set timestamps
	now := time.Now()
	template.CreatedAt = now
	template.UpdatedAt = now

	// Create template directory
	templateDir := tm.getTemplatePath(template.ID)
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}

	// Save template configuration
	if err := tm.saveTemplateConfig(template); err != nil {
		return fmt.Errorf("failed to save template config: %w", err)
	}

	// Save template files
	if err := tm.saveTemplateFiles(template); err != nil {
		return fmt.Errorf("failed to save template files: %w", err)
	}

	return nil
}

// UpdateTemplate updates an existing template
func (tm *TemplateManager) UpdateTemplate(template *CustomTemplate) error {
	if !tm.TemplateExists(template.ID) {
		return fmt.Errorf("template '%s' does not exist", template.ID)
	}

	// Load existing template to preserve created_at
	existing, err := tm.LoadTemplate(template.ID)
	if err != nil {
		return fmt.Errorf("failed to load existing template: %w", err)
	}

	template.CreatedAt = existing.CreatedAt
	template.UpdatedAt = time.Now()

	// Save updated template
	if err := tm.saveTemplateConfig(template); err != nil {
		return fmt.Errorf("failed to save template config: %w", err)
	}

	if err := tm.saveTemplateFiles(template); err != nil {
		return fmt.Errorf("failed to save template files: %w", err)
	}

	return nil
}

// DeleteTemplate deletes a custom template
func (tm *TemplateManager) DeleteTemplate(templateID string) error {
	if !tm.TemplateExists(templateID) {
		return fmt.Errorf("template '%s' does not exist", templateID)
	}

	templateDir := tm.getTemplatePath(templateID)
	if err := os.RemoveAll(templateDir); err != nil {
		return fmt.Errorf("failed to delete template directory: %w", err)
	}

	return nil
}

// LoadTemplate loads a custom template by ID
func (tm *TemplateManager) LoadTemplate(templateID string) (*CustomTemplate, error) {
	configPath := tm.getTemplateConfigPath(templateID)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template config: %w", err)
	}

	var template CustomTemplate
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	// Load template files
	template.Files, template.Templates, err = tm.loadTemplateFiles(templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to load template files: %w", err)
	}

	return &template, nil
}

// ListTemplates returns all custom templates
func (tm *TemplateManager) ListTemplates() ([]*CustomTemplate, error) {
	if _, err := os.Stat(tm.templatesDir); os.IsNotExist(err) {
		return []*CustomTemplate{}, nil
	}

	entries, err := os.ReadDir(tm.templatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []*CustomTemplate
	for _, entry := range entries {
		if entry.IsDir() {
			template, err := tm.LoadTemplate(entry.Name())
			if err != nil {
				continue // Skip invalid templates
			}
			templates = append(templates, template)
		}
	}

	return templates, nil
}

// TemplateExists checks if a template exists
func (tm *TemplateManager) TemplateExists(templateID string) bool {
	configPath := tm.getTemplateConfigPath(templateID)
	_, err := os.Stat(configPath)
	return err == nil
}

// CreateFromProject creates a template from an existing Go project
func (tm *TemplateManager) CreateFromProject(projectPath, templateID, templateName string) error {
	if !filepath.IsAbs(projectPath) {
		var err error
		projectPath, err = filepath.Abs(projectPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}
	}

	// Check if source project exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("source project does not exist: %s", projectPath)
	}

	// Analyze project structure
	config, files, templates, err := tm.analyzeProject(projectPath)
	if err != nil {
		return fmt.Errorf("failed to analyze project: %w", err)
	}

	// Create custom template
	template := &CustomTemplate{
		ID:          templateID,
		Name:        templateName,
		Description: fmt.Sprintf("Template created from %s", filepath.Base(projectPath)),
		Author:      "go-ctl",
		Version:     "1.0.0",
		Tags:        []string{"custom", "generated"},
		Config:      *config,
		Files:       files,
		Templates:   templates,
		Metadata: CustomTemplateMetadata{
			SourceProject: projectPath,
			MinGoVersion:  config.GoVersion,
			Keywords:      []string{"go", "template"},
		},
	}

	return tm.CreateTemplate(template)
}

// ExportTemplate exports a template to a file
func (tm *TemplateManager) ExportTemplate(templateID, outputPath string) error {
	template, err := tm.LoadTemplate(templateID)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// Choose export format based on file extension
	ext := strings.ToLower(filepath.Ext(outputPath))

	var data []byte
	switch ext {
	case ".json":
		data, err = json.MarshalIndent(template, "", "  ")
	case ".yaml", ".yml":
		data, err = yaml.Marshal(template)
	default:
		// Default to YAML
		data, err = yaml.Marshal(template)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write template file: %w", err)
	}

	return nil
}

// ImportTemplate imports a template from a file
func (tm *TemplateManager) ImportTemplate(templatePath string) error {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	var template CustomTemplate

	// Try YAML first, then JSON
	if err := yaml.Unmarshal(data, &template); err != nil {
		if jsonErr := json.Unmarshal(data, &template); jsonErr != nil {
			return fmt.Errorf("failed to parse template (tried YAML and JSON): %w", err)
		}
	}

	// Validate template
	if err := tm.validateTemplate(&template); err != nil {
		return fmt.Errorf("invalid template: %w", err)
	}

	return tm.CreateTemplate(&template)
}

// GetTemplateConfig returns template configuration for project generation
func (tm *TemplateManager) GetTemplateConfig(templateID string) (*metadata.ProjectConfig, error) {
	template, err := tm.LoadTemplate(templateID)
	if err != nil {
		return nil, err
	}

	return &template.Config, nil
}

// ValidateTemplate validates template configuration
func (tm *TemplateManager) ValidateTemplate(templateID string) error {
	template, err := tm.LoadTemplate(templateID)
	if err != nil {
		return err
	}

	return tm.validateTemplate(template)
}

// getTemplatePath returns the full path to a template directory
func (tm *TemplateManager) getTemplatePath(templateID string) string {
	return filepath.Join(tm.templatesDir, templateID)
}

// getTemplateConfigPath returns the path to template config file
func (tm *TemplateManager) getTemplateConfigPath(templateID string) string {
	return filepath.Join(tm.getTemplatePath(templateID), "template.yaml")
}

// saveTemplateConfig saves template configuration to file
func (tm *TemplateManager) saveTemplateConfig(template *CustomTemplate) error {
	configPath := tm.getTemplateConfigPath(template.ID)

	// Create a copy without files and templates for config
	configTemplate := *template
	configTemplate.Files = nil
	configTemplate.Templates = nil

	data, err := yaml.Marshal(&configTemplate)
	if err != nil {
		return fmt.Errorf("failed to marshal template config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write template config: %w", err)
	}

	return nil
}

// saveTemplateFiles saves template files to directory
func (tm *TemplateManager) saveTemplateFiles(template *CustomTemplate) error {
	templateDir := tm.getTemplatePath(template.ID)

	// Save static files
	if template.Files != nil {
		filesDir := filepath.Join(templateDir, "files")
		if err := os.MkdirAll(filesDir, 0755); err != nil {
			return fmt.Errorf("failed to create files directory: %w", err)
		}

		for relativePath, content := range template.Files {
			filePath := filepath.Join(filesDir, relativePath)
			fileDir := filepath.Dir(filePath)

			if err := os.MkdirAll(fileDir, 0755); err != nil {
				return fmt.Errorf("failed to create file directory: %w", err)
			}

			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write file %s: %w", relativePath, err)
			}
		}
	}

	// Save template files
	if template.Templates != nil {
		templatesDir := filepath.Join(templateDir, "templates")
		if err := os.MkdirAll(templatesDir, 0755); err != nil {
			return fmt.Errorf("failed to create templates directory: %w", err)
		}

		for relativePath, content := range template.Templates {
			filePath := filepath.Join(templatesDir, relativePath)
			fileDir := filepath.Dir(filePath)

			if err := os.MkdirAll(fileDir, 0755); err != nil {
				return fmt.Errorf("failed to create template directory: %w", err)
			}

			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write template %s: %w", relativePath, err)
			}
		}
	}

	return nil
}

// loadTemplateFiles loads template files from directory
func (tm *TemplateManager) loadTemplateFiles(templateID string) (map[string]string, map[string]string, error) {
	templateDir := tm.getTemplatePath(templateID)

	files := make(map[string]string)
	templates := make(map[string]string)

	// Load static files
	filesDir := filepath.Join(templateDir, "files")
	if _, err := os.Stat(filesDir); err == nil {
		err := filepath.WalkDir(filesDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			relativePath, err := filepath.Rel(filesDir, path)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			files[relativePath] = string(content)
			return nil
		})

		if err != nil {
			return nil, nil, fmt.Errorf("failed to load files: %w", err)
		}
	}

	// Load template files
	templatesDir := filepath.Join(templateDir, "templates")
	if _, err := os.Stat(templatesDir); err == nil {
		err := filepath.WalkDir(templatesDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			relativePath, err := filepath.Rel(templatesDir, path)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			templates[relativePath] = string(content)
			return nil
		})

		if err != nil {
			return nil, nil, fmt.Errorf("failed to load templates: %w", err)
		}
	}

	return files, templates, nil
}

// analyzeProject analyzes an existing project to extract template data
func (tm *TemplateManager) analyzeProject(projectPath string) (*metadata.ProjectConfig, map[string]string, map[string]string, error) {
	config := &metadata.ProjectConfig{
		ProjectName:    filepath.Base(projectPath),
		GoVersion:      "1.23", // Default
		CustomPackages: []string{},
	}

	files := make(map[string]string)
	templates := make(map[string]string)

	// Analyze go.mod for dependencies and Go version
	goModPath := filepath.Join(projectPath, "go.mod")
	if goModData, err := os.ReadFile(goModPath); err == nil {
		if err := tm.parseGoMod(string(goModData), config); err == nil {
			templates["go.mod.tpl"] = tm.templatizeGoMod(string(goModData))
		}
	}

	// Find and analyze common project files
	commonFiles := []string{
		"README.md", "Dockerfile", "docker-compose.yml", "Makefile",
		".gitignore", ".env.example", ".air.toml",
	}

	for _, fileName := range commonFiles {
		filePath := filepath.Join(projectPath, fileName)
		if content, err := os.ReadFile(filePath); err == nil {
			if tm.isTemplatable(fileName) {
				templates[fileName+".tpl"] = tm.templatizeContent(string(content), config.ProjectName)
			} else {
				files[fileName] = string(content)
			}
		}
	}

	// Analyze project structure for patterns
	err := filepath.WalkDir(projectPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip certain directories
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "vendor" || name == "bin" || name == "tmp" {
				return filepath.SkipDir
			}
		}

		return nil
	})

	return config, files, templates, err
}

// parseGoMod extracts information from go.mod content
func (tm *TemplateManager) parseGoMod(content string, config *metadata.ProjectConfig) error {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "module ") {
			config.ProjectName = strings.TrimPrefix(line, "module ")
		}

		if strings.HasPrefix(line, "go ") {
			config.GoVersion = strings.TrimPrefix(line, "go ")
		}

		if strings.Contains(line, "require") {
			// Parse dependencies - simplified for now
			if strings.Contains(line, "github.com/gin-gonic/gin") {
				config.HttpPackage = metadata.Option{
					ID:   "gin",
					Name: "Gin",
				}
			}
		}
	}

	return nil
}

// templatizeGoMod converts go.mod content to template
func (tm *TemplateManager) templatizeGoMod(content string) string {
	// Replace module name with template variable
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "module ") {
			lines[i] = "module {{.ProjectName}}"
			break
		}
	}
	return strings.Join(lines, "\n")
}

// templatizeContent converts file content to template
func (tm *TemplateManager) templatizeContent(content, projectName string) string {
	// Simple templating - replace project name occurrences
	return strings.ReplaceAll(content, projectName, "{{.ProjectName}}")
}

// isTemplatable checks if a file should be treated as template
func (tm *TemplateManager) isTemplatable(fileName string) bool {
	templateableFiles := []string{
		"README.md", "Dockerfile", "docker-compose.yml",
		"go.mod", "main.go",
	}

	for _, tf := range templateableFiles {
		if fileName == tf {
			return true
		}
	}

	return false
}

// validateTemplateID validates template ID format
func (tm *TemplateManager) validateTemplateID(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("ID cannot be empty")
	}

	if len(id) > 50 {
		return fmt.Errorf("ID too long (max 50 characters)")
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, r := range id {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_') {
			return fmt.Errorf("ID contains invalid character: %c", r)
		}
	}

	// Cannot start with hyphen or underscore
	if id[0] == '-' || id[0] == '_' {
		return fmt.Errorf("ID cannot start with hyphen or underscore")
	}

	return nil
}

// validateTemplate validates template structure and content
func (tm *TemplateManager) validateTemplate(template *CustomTemplate) error {
	if err := tm.validateTemplateID(template.ID); err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	if template.Name == "" {
		return fmt.Errorf("name is required")
	}

	if template.Config.GoVersion == "" {
		return fmt.Errorf("Go version is required in config")
	}

	// Validate configuration using existing metadata validation
	if warnings := metadata.ValidateConfig(template.Config); len(warnings) > 0 {
		return fmt.Errorf("config validation failed: %s", warnings[0])
	}

	return nil
}
