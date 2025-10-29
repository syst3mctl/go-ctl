package service

import (
	"context"
	"fmt"
	"time"

	"{{.ProjectName}}/internal/domain"
{{if .HasFeature "logging"}}	"github.com/rs/zerolog/log"
{{end}})

// Service provides business logic operations
type Service struct {
{{if ne .DbDriver.ID ""}}	repo interface{}
{{end}}	config interface{}
}

// New creates a new service instance
{{if ne .DbDriver.ID ""}}func New(repo interface{}) *Service {
	return &Service{
		repo:   repo,
		config: nil,
	}
}
{{else}}func New() *Service {
	return &Service{
		config: nil,
	}
}
{{end}}

// CreateUser creates a new user with validation
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
{{if .HasFeature "logging"}}	log.Info().Str("email", req.Email).Msg("Creating user")
{{end}}

	// Validate request
	if req.Name == "" || len(req.Name) < 2 {
		return nil, fmt.Errorf("name must be at least 2 characters long")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Create domain user
	user := &domain.User{
		Name:   req.Name,
		Email:  req.Email,
		Active: req.Active,
	}

{{if ne .DbDriver.ID ""}}	// Save to repository (implementation depends on selected driver)
	// This would call the appropriate repository method
	// Example: err := s.repo.CreateUser(ctx, user)
	// For now, we'll simulate success
	_ = user // Prevent unused variable error
{{end}}

	// Return response
	return &UserResponse{
{{if eq .DbDriver.ID "mongo-driver"}}		ID:        "generated-id",
{{else if eq .DbDriver.ID "gorm"}}		ID:        1,
{{else}}		ID:        1,
{{end}}		Name:      user.Name,
		Email:     user.Email,
		Active:    user.Active,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// GetUser retrieves a user by ID
{{if eq .DbDriver.ID "mongo-driver"}}func (s *Service) GetUser(ctx context.Context, id string) (*UserResponse, error) {
{{else if eq .DbDriver.ID "gorm"}}func (s *Service) GetUser(ctx context.Context, id uint) (*UserResponse, error) {
{{else}}func (s *Service) GetUser(ctx context.Context, id int64) (*UserResponse, error) {
{{end}}{{if .HasFeature "logging"}}	log.Info().Interface("id", id).Msg("Getting user")
{{end}}

{{if ne .DbDriver.ID ""}}	// Get from repository (implementation depends on selected driver)
	// This would call the appropriate repository method
	// Example: user, err := s.repo.GetUserByID(ctx, id)
{{end}}

	// For now, return a mock response
	return &UserResponse{
		ID:        id,
		Name:      "John Doe",
		Email:     "john@example.com",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// UpdateUser updates an existing user
{{if eq .DbDriver.ID "mongo-driver"}}func (s *Service) UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
{{else if eq .DbDriver.ID "gorm"}}func (s *Service) UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (*UserResponse, error) {
{{else}}func (s *Service) UpdateUser(ctx context.Context, id int64, req UpdateUserRequest) (*UserResponse, error) {
{{end}}{{if .HasFeature "logging"}}	log.Info().Interface("id", id).Msg("Updating user")
{{end}}

	// Validate partial update request
	if req.Name != nil && len(*req.Name) < 2 {
		return nil, fmt.Errorf("name must be at least 2 characters long")
	}

{{if ne .DbDriver.ID ""}}	// Get existing user from repository
	// Update fields that are provided
	// Save back to repository
{{end}}

	// Return updated response
	response := &UserResponse{
		ID:        id,
		CreatedAt: time.Now().Add(-24 * time.Hour), // Mock created time
		UpdatedAt: time.Now(),
	}

	// Apply updates
	if req.Name != nil {
		response.Name = *req.Name
	} else {
		response.Name = "John Doe" // Mock existing name
	}
	if req.Email != nil {
		response.Email = *req.Email
	} else {
		response.Email = "john@example.com" // Mock existing email
	}
	if req.Active != nil {
		response.Active = *req.Active
	} else {
		response.Active = true // Mock existing active state
	}

	return response, nil
}

// DeleteUser soft deletes a user
{{if eq .DbDriver.ID "mongo-driver"}}func (s *Service) DeleteUser(ctx context.Context, id string) error {
{{else if eq .DbDriver.ID "gorm"}}func (s *Service) DeleteUser(ctx context.Context, id uint) error {
{{else}}func (s *Service) DeleteUser(ctx context.Context, id int64) error {
{{end}}{{if .HasFeature "logging"}}	log.Info().Interface("id", id).Msg("Deleting user")
{{end}}

{{if ne .DbDriver.ID ""}}	// Delete from repository (implementation depends on selected driver)
	// This would call the appropriate repository method
	// Example: return s.repo.DeleteUser(ctx, id)
{{end}}

	return nil
}

// ListUsers returns a paginated list of users
func (s *Service) ListUsers(ctx context.Context, req ListRequest) (*PaginatedResult[UserResponse], error) {
{{if .HasFeature "logging"}}	log.Info().Int("limit", req.Limit).Int("offset", req.Offset).Msg("Listing users")
{{end}}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

{{if ne .DbDriver.ID ""}}	// Get from repository with pagination
	// This would call the appropriate repository method
	// Example: users, err := s.repo.ListUsers(ctx, req.Limit, req.Offset)
	// count, err := s.repo.CountUsers(ctx)
{{end}}

	// Mock response
	users := []UserResponse{
		{
{{if eq .DbDriver.ID "mongo-driver"}}			ID:        "1",
{{else}}			ID:        1,
{{end}}			Name:      "John Doe",
			Email:     "john@example.com",
			Active:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
{{if eq .DbDriver.ID "mongo-driver"}}			ID:        "2",
{{else}}			ID:        2,
{{end}}			Name:      "Jane Smith",
			Email:     "jane@example.com",
			Active:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return &PaginatedResult[UserResponse]{
		Items:  users,
		Total:  len(users),
		Offset: req.Offset,
		Limit:  req.Limit,
	}, nil
}

{{if ne .Database.ID "redis"}}// CreateProduct creates a new product with validation
func (s *Service) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
{{if .HasFeature "logging"}}	log.Info().Str("name", req.Name).Msg("Creating product")
{{end}}

	// Validate request
	if req.Name == "" || len(req.Name) < 2 {
		return nil, fmt.Errorf("name must be at least 2 characters long")
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("price must be non-negative")
	}

	// Create domain product
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Available:   req.Available,
{{if eq .DbDriver.ID "mongo-driver"}}		UserID:      req.UserID, // This would be converted from string
{{else}}		UserID:      req.UserID,
{{end}}	}

{{if ne .DbDriver.ID ""}}	// Save to repository
	// Example: err := s.repo.CreateProduct(ctx, product)
	_ = product // Prevent unused variable error
{{end}}

	// Return response
	return &ProductResponse{
{{if eq .DbDriver.ID "mongo-driver"}}		ID:          "generated-id",
{{else if eq .DbDriver.ID "gorm"}}		ID:          1,
{{else}}		ID:          1,
{{end}}		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Available:   product.Available,
		UserID:      req.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// GetProduct retrieves a product by ID
{{if eq .DbDriver.ID "mongo-driver"}}func (s *Service) GetProduct(ctx context.Context, id string) (*ProductResponse, error) {
{{else if eq .DbDriver.ID "gorm"}}func (s *Service) GetProduct(ctx context.Context, id uint) (*ProductResponse, error) {
{{else}}func (s *Service) GetProduct(ctx context.Context, id int64) (*ProductResponse, error) {
{{end}}{{if .HasFeature "logging"}}	log.Info().Interface("id", id).Msg("Getting product")
{{end}}

{{if ne .DbDriver.ID ""}}	// Get from repository
	// Example: product, err := s.repo.GetProductByID(ctx, id)
{{end}}

	// Mock response
	return &ProductResponse{
		ID:          id,
		Name:        "Sample Product",
		Description: "A sample product",
		Price:       99.99,
		Available:   true,
{{if eq .DbDriver.ID "mongo-driver"}}		UserID:      "user-1",
{{else}}		UserID:      1,
{{end}}		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// ListProducts returns a paginated list of products
func (s *Service) ListProducts(ctx context.Context, req ListRequest) (*PaginatedResult[ProductResponse], error) {
{{if .HasFeature "logging"}}	log.Info().Int("limit", req.Limit).Int("offset", req.Offset).Msg("Listing products")
{{end}}

	// Set defaults
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

{{if ne .DbDriver.ID ""}}	// Get from repository with pagination
	// Example: products, err := s.repo.ListProducts(ctx, req.Limit, req.Offset)
{{end}}

	// Mock response
	products := []ProductResponse{
		{
{{if eq .DbDriver.ID "mongo-driver"}}			ID:          "1",
			UserID:      "user-1",
{{else}}			ID:          1,
			UserID:      1,
{{end}}			Name:        "Product 1",
			Description: "First product",
			Price:       99.99,
			Available:   true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
{{if eq .DbDriver.ID "mongo-driver"}}			ID:          "2",
			UserID:      "user-1",
{{else}}			ID:          2,
			UserID:      1,
{{end}}			Name:        "Product 2",
			Description: "Second product",
			Price:       149.99,
			Available:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	return &PaginatedResult[ProductResponse]{
		Items:  products,
		Total:  len(products),
		Offset: req.Offset,
		Limit:  req.Limit,
	}, nil
}

{{end}}// Health check method
func (s *Service) HealthCheck(ctx context.Context) *HealthResponse {
	return &HealthResponse{
		Status:    "ok",
		Service:   "{{.ProjectName}}",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
