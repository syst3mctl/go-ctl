package testing

import (
	"testing"
{{if eq .DbDriver.ID "gorm"}}
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
{{else if eq .DbDriver.ID "sqlx"}}	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
{{else if eq .DbDriver.ID "mongo-driver"}}	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
{{else if eq .DbDriver.ID "redis-client"}}	"context"
	"github.com/redis/go-redis/v9"
{{end}}	"{{.ProjectName}}/internal/config"
	"{{.ProjectName}}/internal/domain"
{{if ne .DbDriver.ID ""}}	"{{.ProjectName}}/internal/repository"
{{end}}	"{{.ProjectName}}/internal/service"

{{if .HasFeature "testify"}}	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
{{end}})

// TestConfig returns a test configuration
func TestConfig() *config.Config {
	return &config.Config{
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
			DB:   1, // Use different DB for tests
{{end}}		},
{{end}}{{if .HasFeature "jwt"}}		JWT: config.JWTConfig{
			Secret:     "test-secret-key",
			Expiration: "24h",
		},
{{end}}	}
}

{{if eq .DbDriver.ID "gorm"}}// SetupTestDB creates a test database connection using SQLite in-memory
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
{{if .HasFeature "testify"}}	require.NoError(t, err)
{{else}}	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
{{end}}
	// Auto-migrate test models
	err = db.AutoMigrate(&domain.User{})
{{if ne .Database.ID "redis"}}	if err == nil {
		err = db.AutoMigrate(&domain.Product{})
	}
{{end}}{{if .HasFeature "testify"}}	require.NoError(t, err)
{{else}}	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
{{end}}
	return db
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
{{if .HasFeature "testify"}}	require.NoError(t, err)
{{else}}	if err != nil {
		t.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
{{end}}	sqlDB.Close()
}

{{else if eq .DbDriver.ID "sqlx"}}// SetupTestDB creates a test database connection using SQLite in-memory
func SetupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", ":memory:")
{{if .HasFeature "testify"}}	require.NoError(t, err)
{{else}}	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
{{end}}
	// Create test tables
	createTablesSQL := `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			active BOOLEAN DEFAULT TRUE,
			created_at DATETIME,
			updated_at DATETIME
		);
{{if ne .Database.ID "redis"}}
		CREATE TABLE products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			price REAL NOT NULL,
			available BOOLEAN DEFAULT TRUE,
			user_id INTEGER NOT NULL,
			created_at DATETIME,
			updated_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
{{end}}	`
	_, err = db.Exec(createTablesSQL)
{{if .HasFeature "testify"}}	require.NoError(t, err)
{{else}}	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}
{{end}}
	return db
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T, db *sqlx.DB) {
	db.Close()
}

{{else if eq .DbDriver.ID "mongo-driver"}}// SetupTestDB creates a test MongoDB connection
func SetupTestDB(t *testing.T) *mongo.Client {
	// For testing, you might want to use a test MongoDB instance or testcontainers
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
{{if .HasFeature "testify"}}	require.NoError(t, err)
{{else}}	if err != nil {
		t.Fatalf("Failed to connect to test MongoDB: %v", err)
	}
{{end}}
	// Test the connection
	err = client.Ping(context.Background(), nil)
{{if .HasFeature "testify"}}	require.NoError(t, err)
{{else}}	if err != nil {
		t.Fatalf("Failed to ping test MongoDB: %v", err)
	}
{{end}}
	return client
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T, client *mongo.Client) {
	// Clean up test data
	db := client.Database("{{.ProjectName}}_test")
	db.Drop(context.Background())
	client.Disconnect(context.Background())
}

{{else if eq .DbDriver.ID "redis-client"}}// SetupTestDB creates a test Redis connection
func SetupTestDB(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // Use DB 1 for tests
	})

	// Test the connection
	_, err := client.Ping(context.Background()).Result()
{{if .HasFeature "testify"}}	require.NoError(t, err)
{{else}}	if err != nil {
		t.Fatalf("Failed to connect to test Redis: %v", err)
	}
{{end}}
	return client
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(t *testing.T, client *redis.Client) {
	// Clear test data
	client.FlushDB(context.Background())
	client.Close()
}
{{end}}

{{if .HasFeature "testify"}}// TestSuite is a base test suite for integration tests
type TestSuite struct {
	suite.Suite
	Config   *config.Config
{{if ne .DbDriver.ID ""}}	DB       {{if eq .DbDriver.ID "gorm"}}*gorm.DB{{else if eq .DbDriver.ID "sqlx"}}*sqlx.DB{{else if eq .DbDriver.ID "mongo-driver"}}*mongo.Client{{else if eq .DbDriver.ID "redis-client"}}*redis.Client{{else}}interface{}{{end}}
	Repos    *TestRepositories
{{end}}	Services *service.Services
}

{{if ne .DbDriver.ID ""}}// TestRepositories holds test repository instances
type TestRepositories struct {
	User    domain.UserRepository
{{if ne .Database.ID "redis"}}	Product domain.ProductRepository
{{end}}
}
{{end}}

// SetupSuite runs before all tests in the suite
func (s *TestSuite) SetupSuite() {
	s.Config = TestConfig()
{{if ne .DbDriver.ID ""}}	s.DB = SetupTestDB(s.T())

	// Initialize test repositories
	s.Repos = &TestRepositories{
		User:    repository.NewUserRepository(s.DB),
{{if ne .Database.ID "redis"}}		Product: repository.NewProductRepository(s.DB),
{{end}}	}

	// Initialize services with test repositories
	s.Services = &service.Services{
		User:    service.NewUserService(s.Repos.User),
{{if ne .Database.ID "redis"}}		Product: service.NewProductService(s.Repos.Product, s.Repos.User),
{{end}}	}
{{else}}	s.Services = service.New()
{{end}}
}

// TearDownSuite runs after all tests in the suite
func (s *TestSuite) TearDownSuite() {
{{if ne .DbDriver.ID ""}}	TeardownTestDB(s.T(), s.DB)
{{end}}
}

// SetupTest runs before each test
func (s *TestSuite) SetupTest() {
	// Add any per-test setup here
}

// TearDownTest runs after each test
func (s *TestSuite) TearDownTest() {
{{if ne .DbDriver.ID ""}}	// Clean up test data between tests
{{if eq .DbDriver.ID "gorm"}}	s.DB.Exec("DELETE FROM users")
{{if ne .Database.ID "redis"}}	s.DB.Exec("DELETE FROM products")
{{end}}{{else if eq .DbDriver.ID "sqlx"}}	s.DB.Exec("DELETE FROM users")
{{if ne .Database.ID "redis"}}	s.DB.Exec("DELETE FROM products")
{{end}}{{else if eq .DbDriver.ID "mongo-driver"}}	ctx := context.Background()
	db := s.DB.Database("{{.ProjectName}}_test")
	db.Collection("users").Drop(ctx)
{{if ne .Database.ID "redis"}}	db.Collection("products").Drop(ctx)
{{end}}{{else if eq .DbDriver.ID "redis-client"}}	s.DB.FlushDB(context.Background())
{{end}}{{end}}
}

// CreateTestUser creates a test user
func (s *TestSuite) CreateTestUser() *domain.User {
	user := &domain.User{
		Name:   "Test User",
		Email:  "test@example.com",
		Active: true,
	}
	user.BeforeCreate()

	err := s.Repos.User.Create(user)
	require.NoError(s.T(), err)

	return user
}

{{if ne .Database.ID "redis"}}// CreateTestProduct creates a test product
func (s *TestSuite) CreateTestProduct(userID {{if eq .DbDriver.ID "mongo-driver"}}primitive.ObjectID{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}}) *domain.Product {
	product := &domain.Product{
		Name:        "Test Product",
		Description: "A test product",
		Price:       19.99,
		Available:   true,
		UserID:      userID,
	}
	product.BeforeCreate()

	err := s.Repos.Product.Create(product)
	require.NoError(s.T(), err)

	return product
}
{{end}}

// AssertUser asserts that two users are equal
func (s *TestSuite) AssertUser(expected, actual *domain.User) {
	assert.Equal(s.T(), expected.Name, actual.Name)
	assert.Equal(s.T(), expected.Email, actual.Email)
	assert.Equal(s.T(), expected.Active, actual.Active)
}

{{if ne .Database.ID "redis"}}// AssertProduct asserts that two products are equal
func (s *TestSuite) AssertProduct(expected, actual *domain.Product) {
	assert.Equal(s.T(), expected.Name, actual.Name)
	assert.Equal(s.T(), expected.Description, actual.Description)
	assert.Equal(s.T(), expected.Price, actual.Price)
	assert.Equal(s.T(), expected.Available, actual.Available)
	assert.Equal(s.T(), expected.UserID, actual.UserID)
}
{{end}}
{{else}}// Helper functions for tests without testify

// CreateTestUser creates a test user for testing
func CreateTestUser(t *testing.T, {{if ne .DbDriver.ID ""}}repo domain.UserRepository{{end}}) *domain.User {
	user := &domain.User{
		Name:   "Test User",
		Email:  "test@example.com",
		Active: true,
	}
	user.BeforeCreate()

{{if ne .DbDriver.ID ""}}	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
{{end}}
	return user
}

{{if ne .Database.ID "redis"}}// CreateTestProduct creates a test product for testing
func CreateTestProduct(t *testing.T, {{if ne .DbDriver.ID ""}}repo domain.ProductRepository, {{end}}userID {{if eq .DbDriver.ID "mongo-driver"}}primitive.ObjectID{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}}) *domain.Product {
	product := &domain.Product{
		Name:        "Test Product",
		Description: "A test product",
		Price:       19.99,
		Available:   true,
		UserID:      userID,
	}
	product.BeforeCreate()

{{if ne .DbDriver.ID ""}}	err := repo.Create(product)
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}
{{end}}
	return product
}
{{end}}

// AssertUser asserts that two users are equal
func AssertUser(t *testing.T, expected, actual *domain.User) {
	if expected.Name != actual.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, actual.Name)
	}
	if expected.Email != actual.Email {
		t.Errorf("Expected email %s, got %s", expected.Email, actual.Email)
	}
	if expected.Active != actual.Active {
		t.Errorf("Expected active %t, got %t", expected.Active, actual.Active)
	}
}

{{if ne .Database.ID "redis"}}// AssertProduct asserts that two products are equal
func AssertProduct(t *testing.T, expected, actual *domain.Product) {
	if expected.Name != actual.Name {
		t.Errorf("Expected name %s, got %s", expected.Name, actual.Name)
	}
	if expected.Description != actual.Description {
		t.Errorf("Expected description %s, got %s", expected.Description, actual.Description)
	}
	if expected.Price != actual.Price {
		t.Errorf("Expected price %f, got %f", expected.Price, actual.Price)
	}
	if expected.Available != actual.Available {
		t.Errorf("Expected available %t, got %t", expected.Available, actual.Available)
	}
	if expected.UserID != actual.UserID {
		t.Errorf("Expected user_id %v, got %v", expected.UserID, actual.UserID)
	}
}
{{end}}
{{end}}

// MockHTTPClient can be used for testing HTTP clients
type MockHTTPClient struct {
	responses map[string]string
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		responses: make(map[string]string),
	}
}

// AddResponse adds a mock response for a given URL
func (m *MockHTTPClient) AddResponse(url, response string) {
	m.responses[url] = response
}

// Get returns the mock response for the given URL
func (m *MockHTTPClient) Get(url string) (string, error) {
	if response, exists := m.responses[url]; exists {
		return response, nil
	}
	return "", domain.ErrInternalError
}
```

Now let's create example test files:
