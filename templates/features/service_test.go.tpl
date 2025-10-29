package service

import (
	"context"
	"testing"
	"time"

	"{{.ProjectName}}/internal/domain"
{{if .HasFeature "testify"}}	"{{.ProjectName}}/internal/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
{{end}}{{if eq .DbDriver.ID "mongo-driver"}}	"go.mongodb.org/mongo-driver/bson/primitive"
{{end}})

{{if .HasFeature "testify"}}// UserServiceTestSuite is the test suite for UserService
type UserServiceTestSuite struct {
	testing.TestSuite
}

// TestUserService runs the user service test suite
func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// TestCreateUser tests user creation
func (s *UserServiceTestSuite) TestCreateUser() {
	tests := []struct {
		name    string
		user    *domain.User
		wantErr error
	}{
		{
			name: "valid user",
			user: &domain.User{
				Name:   "John Doe",
				Email:  "john@example.com",
				Active: true,
			},
			wantErr: nil,
		},
		{
			name: "invalid email",
			user: &domain.User{
				Name:   "Jane Doe",
				Email:  "", // empty email
				Active: true,
			},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "invalid name",
			user: &domain.User{
				Name:   "", // empty name
				Email:  "jane@example.com",
				Active: true,
			},
			wantErr: domain.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			err := s.Services.User.CreateUser(ctx, tt.user)

			if tt.wantErr != nil {
				assert.Equal(s.T(), tt.wantErr, err)
			} else {
				require.NoError(s.T(), err)
				{{if eq .DbDriver.ID "gorm"}}assert.NotZero(s.T(), tt.user.ID)
				{{else if eq .DbDriver.ID "mongo-driver"}}assert.False(s.T(), tt.user.ID.IsZero())
				{{else}}assert.NotZero(s.T(), tt.user.ID)
				{{end}}assert.False(s.T(), tt.user.CreatedAt.IsZero())
				assert.False(s.T(), tt.user.UpdatedAt.IsZero())
			}
		})
	}
}

// TestCreateUserDuplicate tests creating a user with duplicate email
func (s *UserServiceTestSuite) TestCreateUserDuplicate() {
	ctx := context.Background()

	// Create first user
	user1 := s.CreateTestUser()

	// Try to create user with same email
	user2 := &domain.User{
		Name:   "Different Name",
		Email:  user1.Email, // same email
		Active: true,
	}

	err := s.Services.User.CreateUser(ctx, user2)
	assert.Equal(s.T(), domain.ErrUserExists, err)
}

// TestGetUser tests user retrieval
func (s *UserServiceTestSuite) TestGetUser() {
	ctx := context.Background()

	// Create test user
	testUser := s.CreateTestUser()

	// Get the user
	{{if eq .DbDriver.ID "mongo-driver"}}retrievedUser, err := s.Services.User.GetUser(ctx, testUser.ID)
	{{else if eq .DbDriver.ID "gorm"}}retrievedUser, err := s.Services.User.GetUser(ctx, testUser.ID)
	{{else}}retrievedUser, err := s.Services.User.GetUser(ctx, testUser.ID)
	{{end}}require.NoError(s.T(), err)

	s.AssertUser(testUser, retrievedUser)
}

// TestGetUserNotFound tests getting non-existent user
func (s *UserServiceTestSuite) TestGetUserNotFound() {
	ctx := context.Background()

	{{if eq .DbDriver.ID "mongo-driver"}}nonExistentID := primitive.NewObjectID()
	{{else if eq .DbDriver.ID "gorm"}}nonExistentID := uint(999999)
	{{else}}nonExistentID := int64(999999)
	{{end}}_, err := s.Services.User.GetUser(ctx, nonExistentID)
	assert.Equal(s.T(), domain.ErrUserNotFound, err)
}

// TestGetUserByEmail tests user retrieval by email
func (s *UserServiceTestSuite) TestGetUserByEmail() {
	ctx := context.Background()

	// Create test user
	testUser := s.CreateTestUser()

	// Get the user by email
	retrievedUser, err := s.Services.User.GetUserByEmail(ctx, testUser.Email)
	require.NoError(s.T(), err)

	s.AssertUser(testUser, retrievedUser)
}

// TestGetUserByEmailNotFound tests getting user by non-existent email
func (s *UserServiceTestSuite) TestGetUserByEmailNotFound() {
	ctx := context.Background()

	_, err := s.Services.User.GetUserByEmail(ctx, "nonexistent@example.com")
	assert.Equal(s.T(), domain.ErrUserNotFound, err)
}

// TestUpdateUser tests user update
func (s *UserServiceTestSuite) TestUpdateUser() {
	ctx := context.Background()

	// Create test user
	testUser := s.CreateTestUser()
	originalUpdatedAt := testUser.UpdatedAt

	// Wait a moment to ensure updated_at changes
	time.Sleep(time.Millisecond * 10)

	// Update the user
	testUser.Name = "Updated Name"
	testUser.Email = "updated@example.com"

	err := s.Services.User.UpdateUser(ctx, testUser)
	require.NoError(s.T(), err)

	// Verify the update
	{{if eq .DbDriver.ID "mongo-driver"}}updatedUser, err := s.Services.User.GetUser(ctx, testUser.ID)
	{{else if eq .DbDriver.ID "gorm"}}updatedUser, err := s.Services.User.GetUser(ctx, testUser.ID)
	{{else}}updatedUser, err := s.Services.User.GetUser(ctx, testUser.ID)
	{{end}}require.NoError(s.T(), err)

	assert.Equal(s.T(), "Updated Name", updatedUser.Name)
	assert.Equal(s.T(), "updated@example.com", updatedUser.Email)
	assert.True(s.T(), updatedUser.UpdatedAt.After(originalUpdatedAt))
}

// TestUpdateUserNotFound tests updating non-existent user
func (s *UserServiceTestSuite) TestUpdateUserNotFound() {
	ctx := context.Background()

	user := &domain.User{
		{{if eq .DbDriver.ID "mongo-driver"}}ID:     primitive.NewObjectID(),
		{{else if eq .DbDriver.ID "gorm"}}ID:     uint(999999),
		{{else}}ID:     int64(999999),
		{{end}}Name:   "Non-existent User",
		Email:  "nonexistent@example.com",
		Active: true,
	}

	err := s.Services.User.UpdateUser(ctx, user)
	assert.Equal(s.T(), domain.ErrUserNotFound, err)
}

// TestDeleteUser tests user deletion
func (s *UserServiceTestSuite) TestDeleteUser() {
	ctx := context.Background()

	// Create test user
	testUser := s.CreateTestUser()

	// Delete the user
	{{if eq .DbDriver.ID "mongo-driver"}}err := s.Services.User.DeleteUser(ctx, testUser.ID)
	{{else if eq .DbDriver.ID "gorm"}}err := s.Services.User.DeleteUser(ctx, testUser.ID)
	{{else}}err := s.Services.User.DeleteUser(ctx, testUser.ID)
	{{end}}require.NoError(s.T(), err)

	// Verify the user is deleted
	{{if eq .DbDriver.ID "mongo-driver"}}_, err = s.Services.User.GetUser(ctx, testUser.ID)
	{{else if eq .DbDriver.ID "gorm"}}_, err = s.Services.User.GetUser(ctx, testUser.ID)
	{{else}}_, err = s.Services.User.GetUser(ctx, testUser.ID)
	{{end}}assert.Equal(s.T(), domain.ErrUserNotFound, err)
}

// TestDeleteUserNotFound tests deleting non-existent user
func (s *UserServiceTestSuite) TestDeleteUserNotFound() {
	ctx := context.Background()

	{{if eq .DbDriver.ID "mongo-driver"}}nonExistentID := primitive.NewObjectID()
	{{else if eq .DbDriver.ID "gorm"}}nonExistentID := uint(999999)
	{{else}}nonExistentID := int64(999999)
	{{end}}err := s.Services.User.DeleteUser(ctx, nonExistentID)
	assert.Equal(s.T(), domain.ErrUserNotFound, err)
}

// TestListUsers tests user listing with pagination
func (s *UserServiceTestSuite) TestListUsers() {
	ctx := context.Background()

	// Create multiple test users
	users := make([]*domain.User, 5)
	for i := 0; i < 5; i++ {
		user := &domain.User{
			Name:   fmt.Sprintf("User %d", i+1),
			Email:  fmt.Sprintf("user%d@example.com", i+1),
			Active: true,
		}
		user.BeforeCreate()
		err := s.Repos.User.Create(user)
		require.NoError(s.T(), err)
		users[i] = user
	}

	// Test pagination
	params := domain.PaginationParams{
		Offset: 0,
		Limit:  3,
	}

	result, err := s.Services.User.ListUsers(ctx, params)
	require.NoError(s.T(), err)

	assert.Len(s.T(), result.Items, 3)
	assert.Equal(s.T(), 0, result.Offset)
	assert.Equal(s.T(), 3, result.Limit)
}

{{if ne .Database.ID "redis"}}// ProductServiceTestSuite is the test suite for ProductService
type ProductServiceTestSuite struct {
	testing.TestSuite
}

// TestProductService runs the product service test suite
func TestProductService(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

// TestCreateProduct tests product creation
func (s *ProductServiceTestSuite) TestCreateProduct() {
	ctx := context.Background()

	// Create test user first
	testUser := s.CreateTestUser()

	product := &domain.Product{
		Name:        "Test Product",
		Description: "A test product",
		Price:       29.99,
		Available:   true,
		{{if eq .DbDriver.ID "mongo-driver"}}UserID:      testUser.ID,
		{{else if eq .DbDriver.ID "gorm"}}UserID:      testUser.ID,
		{{else}}UserID:      testUser.ID,
		{{end}}
	}

	err := s.Services.Product.CreateProduct(ctx, product)
	require.NoError(s.T(), err)

	{{if eq .DbDriver.ID "gorm"}}assert.NotZero(s.T(), product.ID)
	{{else if eq .DbDriver.ID "mongo-driver"}}assert.False(s.T(), product.ID.IsZero())
	{{else}}assert.NotZero(s.T(), product.ID)
	{{end}}assert.False(s.T(), product.CreatedAt.IsZero())
	assert.False(s.T(), product.UpdatedAt.IsZero())
}

// TestCreateProductInvalidPrice tests creating product with invalid price
func (s *ProductServiceTestSuite) TestCreateProductInvalidPrice() {
	ctx := context.Background()

	// Create test user first
	testUser := s.CreateTestUser()

	product := &domain.Product{
		Name:        "Invalid Product",
		Description: "A product with invalid price",
		Price:       -10.0, // negative price
		Available:   true,
		{{if eq .DbDriver.ID "mongo-driver"}}UserID:      testUser.ID,
		{{else if eq .DbDriver.ID "gorm"}}UserID:      testUser.ID,
		{{else}}UserID:      testUser.ID,
		{{end}}
	}

	err := s.Services.Product.CreateProduct(ctx, product)
	assert.Equal(s.T(), domain.ErrInvalidInput, err)
}

// TestCreateProductUserNotFound tests creating product with non-existent user
func (s *ProductServiceTestSuite) TestCreateProductUserNotFound() {
	ctx := context.Background()

	product := &domain.Product{
		Name:        "Orphan Product",
		Description: "A product without a valid user",
		Price:       19.99,
		Available:   true,
		{{if eq .DbDriver.ID "mongo-driver"}}UserID:      primitive.NewObjectID(),
		{{else if eq .DbDriver.ID "gorm"}}UserID:      uint(999999),
		{{else}}UserID:      int64(999999),
		{{end}}
	}

	err := s.Services.Product.CreateProduct(ctx, product)
	assert.Equal(s.T(), domain.ErrUserNotFound, err)
}
{{end}}

{{else}}// Test functions without testify

// TestUserServiceCreate tests user creation without testify
func TestUserServiceCreate(t *testing.T) {
{{if ne .DbDriver.ID ""}}	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
{{else}}	userService := service.NewUserService(nil)
{{end}}
	ctx := context.Background()

	user := &domain.User{
		Name:   "Test User",
		Email:  "test@example.com",
		Active: true,
	}

	err := userService.CreateUser(ctx, user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	{{if eq .DbDriver.ID "gorm"}}if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
	{{else if eq .DbDriver.ID "mongo-driver"}}if user.ID.IsZero() {
		t.Error("Expected user ID to be set")
	}
	{{else}}if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
	{{end}}
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

// TestUserServiceCreateInvalid tests user creation with invalid data
func TestUserServiceCreateInvalid(t *testing.T) {
{{if ne .DbDriver.ID ""}}	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
{{else}}	userService := service.NewUserService(nil)
{{end}}
	ctx := context.Background()

	user := &domain.User{
		Name:   "", // empty name should be invalid
		Email:  "test@example.com",
		Active: true,
	}

	err := userService.CreateUser(ctx, user)
	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
}

// TestUserServiceGetNotFound tests getting non-existent user
func TestUserServiceGetNotFound(t *testing.T) {
{{if ne .DbDriver.ID ""}}	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
{{else}}	userService := service.NewUserService(nil)
{{end}}
	ctx := context.Background()

	{{if eq .DbDriver.ID "mongo-driver"}}nonExistentID := primitive.NewObjectID()
	{{else if eq .DbDriver.ID "gorm"}}nonExistentID := uint(999999)
	{{else}}nonExistentID := int64(999999)
	{{end}}_, err := userService.GetUser(ctx, nonExistentID)
	if err != domain.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

{{if ne .Database.ID "redis"}}// TestProductServiceCreate tests product creation
func TestProductServiceCreate(t *testing.T) {
{{if ne .DbDriver.ID ""}}	db := SetupTestDB(t)
	defer TeardownTestDB(t, db)

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo, userRepo)

	// Create test user first
	testUser := CreateTestUser(t, userRepo)
{{else}}	productService := service.NewProductService(nil, nil)
	testUser := CreateTestUser(t)
{{end}}
	ctx := context.Background()

	product := &domain.Product{
		Name:        "Test Product",
		Description: "A test product",
		Price:       29.99,
		Available:   true,
		{{if eq .DbDriver.ID "mongo-driver"}}UserID:      testUser.ID,
		{{else if eq .DbDriver.ID "gorm"}}UserID:      testUser.ID,
		{{else}}UserID:      testUser.ID,
		{{end}}
	}

	err := productService.CreateProduct(ctx, product)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	{{if eq .DbDriver.ID "gorm"}}if product.ID == 0 {
		t.Error("Expected product ID to be set")
	}
	{{else if eq .DbDriver.ID "mongo-driver"}}if product.ID.IsZero() {
		t.Error("Expected product ID to be set")
	}
	{{else}}if product.ID == 0 {
		t.Error("Expected product ID to be set")
	}
	{{end}}
}
{{end}}
{{end}}

// BenchmarkUserServiceCreate benchmarks user creation
func BenchmarkUserServiceCreate(b *testing.B) {
{{if ne .DbDriver.ID ""}}	db := SetupTestDB(&testing.T{}) // This is a hack for benchmarking
	defer TeardownTestDB(&testing.T{}, db)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
{{else}}	userService := service.NewUserService(nil)
{{end}}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &domain.User{
			Name:   fmt.Sprintf("Benchmark User %d", i),
			Email:  fmt.Sprintf("benchmark%d@example.com", i),
			Active: true,
		}
		userService.CreateUser(ctx, user)
	}
}
