package main

import (
	"testing"
{{if eq .HTTP.ID "gin"}}	"net/http"
	"net/http/httptest"
	"strings"

	"{{.ProjectName}}/internal/config"
	"{{.ProjectName}}/internal/handler"
	"{{.ProjectName}}/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
{{else if eq .HTTP.ID "echo"}}	"net/http"
	"net/http/httptest"
	"strings"

	"{{.ProjectName}}/internal/config"
	"{{.ProjectName}}/internal/handler"
	"{{.ProjectName}}/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
{{else if eq .HTTP.ID "fiber"}}	"strings"

	"{{.ProjectName}}/internal/config"
	"{{.ProjectName}}/internal/handler"
	"{{.ProjectName}}/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
{{else}}	"net/http"
	"net/http/httptest"
	"strings"

	"{{.ProjectName}}/internal/config"
	"{{.ProjectName}}/internal/handler"
	"{{.ProjectName}}/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
{{end}})

// IntegrationTestSuite provides a test suite for integration tests
type IntegrationTestSuite struct {
	suite.Suite
	Config   *config.Config
	Services *service.Services
{{if eq .HTTP.ID "gin"}}	Router   *gin.Engine
{{else if eq .HTTP.ID "echo"}}	Server   *echo.Echo
{{else if eq .HTTP.ID "fiber"}}	App      *fiber.App
{{else}}	Handler  http.Handler
{{end}}
}

// SetupSuite initializes the test suite
func (suite *IntegrationTestSuite) SetupSuite() {
	// Load test configuration
	suite.Config = &config.Config{
		App: config.AppConfig{
			Name:    "{{.ProjectName}}-test",
			Version: "test",
			Env:     "test",
		},
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
{{if ne .DbDriver.ID ""}}		Database: config.DatabaseConfig{
{{if eq .Database.ID "postgres"}}			Host:     "localhost",
			Port:     5432,
			Name:     "{{.ProjectName}}_test",
			User:     "postgres",
			Password: "password",
			SSLMode:  "disable",
{{else if eq .Database.ID "mysql"}}			Host:     "localhost",
			Port:     3306,
			Name:     "{{.ProjectName}}_test",
			User:     "root",
			Password: "password",
{{else if eq .Database.ID "sqlite"}}			Name: ":memory:",
{{else if eq .Database.ID "mongodb"}}			Host: "localhost",
			Port: 27017,
			Name: "{{.ProjectName}}_test",
{{else if eq .Database.ID "redis"}}			Host: "localhost",
			Port: 6379,
			DB:   1,
{{end}}		},
{{end}}{{if .HasFeature "jwt"}}		JWT: config.JWTConfig{
			Secret:     "test-secret-key-for-jwt-tokens",
			Expiration: "24h",
		},
{{end}}	}

{{if ne .DbDriver.ID ""}}	// Initialize database for tests
	// Note: In real tests, you'd set up a test database connection here
{{end}}
	// Initialize services
{{if ne .DbDriver.ID ""}}	suite.Services = service.New(nil) // Would pass test DB connection
{{else}}	suite.Services = service.New()
{{end}}

{{if eq .HTTP.ID "gin"}}	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize handlers
	handlers := handler.New(suite.Services, suite.Config)

	// Setup router
	suite.Router = gin.New()
	suite.setupRoutes(handlers)

{{else if eq .HTTP.ID "echo"}}	// Initialize handlers
	handlers := handler.New(suite.Services, suite.Config)

	// Setup Echo server
	suite.Server = echo.New()
	suite.Server.HideBanner = true
	suite.setupRoutes(handlers)

{{else if eq .HTTP.ID "fiber"}}	// Initialize handlers
	handlers := handler.New(suite.Services, suite.Config)

	// Setup Fiber app
	suite.App = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	suite.setupRoutes(handlers)

{{else}}	// Initialize handlers
	handlers := handler.New(suite.Services, suite.Config)

	// Setup HTTP handler
	suite.Handler = suite.setupRoutes(handlers)
{{end}}
}

{{if eq .HTTP.ID "gin"}}// setupRoutes configures the Gin routes for testing
func (suite *IntegrationTestSuite) setupRoutes(handlers *handler.Handler) {
	// Health check
	suite.Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": suite.Config.App.Name,
			"version": suite.Config.App.Version,
		})
	})

	// API routes
	api := suite.Router.Group("/api/v1")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to " + suite.Config.App.Name + " API",
				"version": suite.Config.App.Version,
			})
		})

		api.POST("/users", handlers.CreateUser)
		api.GET("/users/:id", handlers.GetUser)

{{if .HasFeature "jwt"}}		// Protected routes
		protected := api.Group("/")
		protected.Use(func(c *gin.Context) {
			// Mock JWT middleware for testing
			c.Set("user_id", uint(1))
			c.Set("username", "testuser")
			c.Next()
		})
		protected.GET("/profile", handlers.GetProfile)
{{end}}	}
}

{{else if eq .HTTP.ID "echo"}}// setupRoutes configures the Echo routes for testing
func (suite *IntegrationTestSuite) setupRoutes(handlers *handler.Handler) {
	// Health check
	suite.Server.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"service": suite.Config.App.Name,
			"version": suite.Config.App.Version,
		})
	})

	// API routes
	api := suite.Server.Group("/api/v1")
	api.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to " + suite.Config.App.Name + " API",
			"version": suite.Config.App.Version,
		})
	})

	api.POST("/users", handlers.CreateUser)
	api.GET("/users/:id", handlers.GetUser)

{{if .HasFeature "jwt"}}	// Protected routes
	protected := api.Group("")
	protected.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Mock JWT middleware for testing
			c.Set("user_id", uint(1))
			c.Set("username", "testuser")
			return next(c)
		}
	})
	protected.GET("/profile", handlers.GetProfile)
{{end}}}

{{else if eq .HTTP.ID "fiber"}}// setupRoutes configures the Fiber routes for testing
func (suite *IntegrationTestSuite) setupRoutes(handlers *handler.Handler) {
	// Health check
	suite.App.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": suite.Config.App.Name,
			"version": suite.Config.App.Version,
		})
	})

	// API routes
	api := suite.App.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to " + suite.Config.App.Name + " API",
			"version": suite.Config.App.Version,
		})
	})

	api.Post("/users", handlers.CreateUser)
	api.Get("/users/:id", handlers.GetUser)

{{if .HasFeature "jwt"}}	// Protected routes with mock JWT
	api.Get("/profile", func(c *fiber.Ctx) error {
		// Mock JWT middleware for testing
		c.Locals("user_id", uint(1))
		c.Locals("username", "testuser")
		return handlers.GetProfile(c)
	})
{{end}}}

{{else}}// setupRoutes configures the HTTP routes for testing
func (suite *IntegrationTestSuite) setupRoutes(handlers *handler.Handler) http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"` + suite.Config.App.Name + `","version":"` + suite.Config.App.Version + `"}`))
	})

	// API routes
	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/" && r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Welcome to ` + suite.Config.App.Name + ` API","version":"` + suite.Config.App.Version + `"}`))
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/api/v1/users", handlers.CreateUser)
	mux.HandleFunc("/api/v1/users/", handlers.GetUser)

{{if .HasFeature "jwt"}}	// Protected routes would need JWT middleware wrapper
{{end}}
	return mux
}
{{end}}

// TearDownSuite cleans up after all tests
func (suite *IntegrationTestSuite) TearDownSuite() {
	// Clean up any resources
{{if ne .DbDriver.ID ""}}	// Close database connections, etc.
{{end}}
}

// TestHealthEndpoint tests the health check endpoint
func (suite *IntegrationTestSuite) TestHealthEndpoint() {
{{if eq .HTTP.ID "gin"}}	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "ok")
	assert.Contains(suite.T(), w.Body.String(), suite.Config.App.Name)

{{else if eq .HTTP.ID "echo"}}	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := suite.Server.NewContext(req, rec)

	err := suite.Server.Router().Find(http.MethodGet, "/health", c)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Contains(suite.T(), rec.Body.String(), "ok")

{{else if eq .HTTP.ID "fiber"}}	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := suite.App.Test(req)

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

{{else}}	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	suite.Handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "ok")
{{end}}}

// TestAPIWelcomeEndpoint tests the API welcome endpoint
func (suite *IntegrationTestSuite) TestAPIWelcomeEndpoint() {
{{if eq .HTTP.ID "gin"}}	req, _ := http.NewRequest("GET", "/api/v1/", nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Welcome to")
	assert.Contains(suite.T(), w.Body.String(), suite.Config.App.Name)

{{else if eq .HTTP.ID "echo"}}	req := httptest.NewRequest(http.MethodGet, "/api/v1/", nil)
	rec := httptest.NewRecorder()
	c := suite.Server.NewContext(req, rec)

	err := suite.Server.Router().Find(http.MethodGet, "/api/v1/", c)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Contains(suite.T(), rec.Body.String(), "Welcome to")

{{else if eq .HTTP.ID "fiber"}}	req := httptest.NewRequest("GET", "/api/v1/", nil)
	resp, err := suite.App.Test(req)

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

{{else}}	req, _ := http.NewRequest("GET", "/api/v1/", nil)
	w := httptest.NewRecorder()
	suite.Handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Welcome to")
{{end}}}

// TestCreateUser tests user creation endpoint
func (suite *IntegrationTestSuite) TestCreateUser() {
{{if eq .HTTP.ID "gin"}}	userJSON := `{"name":"Test User","email":"test@example.com","active":true}`
	req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	// Note: This will fail until you implement the actual handlers
	// assert.Equal(suite.T(), http.StatusCreated, w.Code)

{{else if eq .HTTP.ID "echo"}}	userJSON := `{"name":"Test User","email":"test@example.com","active":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := suite.Server.NewContext(req, rec)

	err := suite.Server.Router().Find(http.MethodPost, "/api/v1/users", c)
	require.NoError(suite.T(), err)

	// Note: This will fail until you implement the actual handlers
	// assert.Equal(suite.T(), http.StatusCreated, rec.Code)

{{else if eq .HTTP.ID "fiber"}}	userJSON := `{"name":"Test User","email":"test@example.com","active":true}`
	req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := suite.App.Test(req)

	require.NoError(suite.T(), err)
	// Note: This will fail until you implement the actual handlers
	// assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

{{else}}	userJSON := `{"name":"Test User","email":"test@example.com","active":true}`
	req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.Handler.ServeHTTP(w, req)

	// Note: This will fail until you implement the actual handlers
	// assert.Equal(suite.T(), http.StatusCreated, w.Code)
{{end}}}

{{if .HasFeature "jwt"}}// TestProtectedEndpoint tests JWT-protected endpoints
func (suite *IntegrationTestSuite) TestProtectedEndpoint() {
{{if eq .HTTP.ID "gin"}}	req, _ := http.NewRequest("GET", "/api/v1/profile", nil)
	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)

	// Since we're using mock JWT middleware, this should succeed
	// assert.Equal(suite.T(), http.StatusOK, w.Code)

{{else if eq .HTTP.ID "echo"}}	req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
	rec := httptest.NewRecorder()
	c := suite.Server.NewContext(req, rec)

	err := suite.Server.Router().Find(http.MethodGet, "/api/v1/profile", c)
	require.NoError(suite.T(), err)

	// Since we're using mock JWT middleware, this should succeed
	// assert.Equal(suite.T(), http.StatusOK, rec.Code)

{{else if eq .HTTP.ID "fiber"}}	req := httptest.NewRequest("GET", "/api/v1/profile", nil)
	resp, err := suite.App.Test(req)

	require.NoError(suite.T(), err)
	// Since we're using mock JWT middleware, this should succeed
	// assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

{{else}}	// Protected routes need JWT middleware implementation for net/http
{{end}}}
{{end}}

// TestIntegration runs the integration test suite
func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Example of table-driven tests
func (suite *IntegrationTestSuite) TestValidation() {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "valid input",
			input:    `{"name":"John Doe","email":"john@example.com","active":true}`,
			expected: http.StatusCreated,
		},
		{
			name:     "invalid email",
			input:    `{"name":"John Doe","email":"invalid-email","active":true}`,
			expected: http.StatusBadRequest,
		},
		{
			name:     "missing name",
			input:    `{"email":"john@example.com","active":true}`,
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
{{if eq .HTTP.ID "gin"}}			req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			suite.Router.ServeHTTP(w, req)

			// Note: Uncomment when handlers are implemented
			// assert.Equal(suite.T(), tt.expected, w.Code, "Test case: %s", tt.name)

{{else if eq .HTTP.ID "echo"}}			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := suite.Server.NewContext(req, rec)

			err := suite.Server.Router().Find(http.MethodPost, "/api/v1/users", c)
			require.NoError(suite.T(), err)

			// Note: Uncomment when handlers are implemented
			// assert.Equal(suite.T(), tt.expected, rec.Code, "Test case: %s", tt.name)

{{else if eq .HTTP.ID "fiber"}}			req := httptest.NewRequest("POST", "/api/v1/users", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")
			resp, err := suite.App.Test(req)

			require.NoError(suite.T(), err)
			// Note: Uncomment when handlers are implemented
			// assert.Equal(suite.T(), tt.expected, resp.StatusCode, "Test case: %s", tt.name)

{{else}}			req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(tt.input))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			suite.Handler.ServeHTTP(w, req)

			// Note: Uncomment when handlers are implemented
			// assert.Equal(suite.T(), tt.expected, w.Code, "Test case: %s", tt.name)
{{end}}		})
	}
}

// Mock functions and helpers for testing
func (suite *IntegrationTestSuite) createTestUser() map[string]interface{} {
	return map[string]interface{}{
		"id":     1,
		"name":   "Test User",
		"email":  "test@example.com",
		"active": true,
	}
}

// Helper function to make authenticated requests (when JWT is enabled)
{{if .HasFeature "jwt"}}func (suite *IntegrationTestSuite) makeAuthenticatedRequest(method, url, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer mock-jwt-token-for-testing")
	req.Header.Set("Content-Type", "application/json")

{{if eq .HTTP.ID "gin"}}	w := httptest.NewRecorder()
	suite.Router.ServeHTTP(w, req)
	return w
{{else}}	// Implement for other frameworks
	return nil
{{end}}}
{{end}}
