package {{.DbDriver.ID}}

import (
	"context"
	"fmt"
	"time"

	"{{.ProjectName}}/internal/domain"

	"gorm.io/gorm"
{{if eq .Database.ID "postgres"}}	"gorm.io/driver/postgres"
{{else if eq .Database.ID "mysql"}}	"gorm.io/driver/mysql"
{{else if eq .Database.ID "sqlite"}}	"gorm.io/driver/sqlite"
{{end}})

// Storage implements the storage layer using GORM
type Storage struct {
	db *gorm.DB
}

// New creates a new GORM storage instance
func New(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

// NewConnection creates a new database connection using GORM
func NewConnection(dsn string) (*gorm.DB, error) {
{{if eq .Database.ID "postgres"}}	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else if eq .Database.ID "mysql"}}	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else if eq .Database.ID "sqlite"}}	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else}}	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{end}}	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// Migrate runs database migrations
func (s *Storage) Migrate() error {
	return s.db.AutoMigrate(
		&domain.User{},
		// Add your models here
	)
}

// Health checks database connectivity
func (s *Storage) Health(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// User repository methods

// CreateUser creates a new user
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID
func (s *Storage) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// UpdateUser updates an existing user
func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := s.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser soft deletes a user
func (s *Storage) DeleteUser(ctx context.Context, id uint) error {
	if err := s.db.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers retrieves users with pagination
func (s *Storage) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	query := s.db.WithContext(ctx).Order("created_at DESC")

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

// CountUsers returns the total number of users
func (s *Storage) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&domain.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// SearchUsers searches users by name or email
func (s *Storage) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	searchQuery := "%" + query + "%"

	dbQuery := s.db.WithContext(ctx).
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

// Transaction wraps operations in a database transaction
func (s *Storage) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return s.db.WithContext(ctx).Transaction(fn)
}

// Close closes the database connection
func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
}

// TODO: Add more repository methods as needed for your application
