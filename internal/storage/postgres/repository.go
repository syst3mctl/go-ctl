package postgres

import (
	"context"
	"fmt"
	"time"

	"test-app/internal/domain"

	"gorm.io/gorm"
)

// Repository implements PostgreSQL-specific repository operations
type Repository struct {
	db *gorm.DB

}

// New creates a new PostgreSQL repository instance using GORM
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Migrate runs database migrations for PostgreSQL using GORM
func (r *Repository) Migrate() error {
	return r.db.AutoMigrate(
		&domain.User{},
		&domain.Address{},
		&domain.Order{},
		&domain.Category{},
		&domain.Product{},
	)
}

// User repository methods using GORM

// CreateUser creates a new user in PostgreSQL
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID from PostgreSQL
func (r *Repository) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email from PostgreSQL
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// UpdateUser updates an existing user in PostgreSQL
func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser soft deletes a user in PostgreSQL
func (r *Repository) DeleteUser(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers retrieves users with pagination from PostgreSQL
func (r *Repository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	query := r.db.WithContext(ctx).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users in PostgreSQL
func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// SearchUsers searches users by name or email in PostgreSQL
func (r *Repository) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	searchQuery := "%" + query + "%"

	dbQuery := r.db.WithContext(ctx).
		Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).
		Order("created_at DESC")

	if limit > 0 {
		dbQuery = dbQuery.Limit(limit)
	}
	if offset > 0 {
		dbQuery = dbQuery.Offset(offset)
	}

	if err := dbQuery.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// Product repository methods using GORM

// CreateProduct creates a new product in PostgreSQL
func (r *Repository) CreateProduct(ctx context.Context, product *domain.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetProductByID retrieves a product by ID from PostgreSQL
func (r *Repository) GetProductByID(ctx context.Context, id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.WithContext(ctx).Preload("User").First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &product, nil
}

// ListProducts retrieves products with pagination from PostgreSQL
func (r *Repository) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	query := r.db.WithContext(ctx).Preload("User").Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return products, nil
}

// Transaction wraps operations in a PostgreSQL transaction using GORM
func (r *Repository) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}



// Close closes the database connection
func (r *Repository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()

}
