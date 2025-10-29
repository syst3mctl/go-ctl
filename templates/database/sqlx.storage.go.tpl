package {{.DbDriver.ID | replace "-" "_"}}

import (
	"context"
	"fmt"
	"time"

	"{{.ProjectName}}/internal/domain"

{{if eq .Database.ID "postgres"}}	_ "github.com/lib/pq"
{{else if eq .Database.ID "mysql"}}	_ "github.com/go-sql-driver/mysql"
{{else if eq .Database.ID "sqlite"}}	_ "github.com/mattn/go-sqlite3"
{{end}}	"github.com/jmoiron/sqlx"
)

// Storage implements the storage layer using sqlx
type Storage struct {
	db *sqlx.DB
}

// New creates a new sqlx storage instance
func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

// NewConnection creates a new database connection using sqlx
func NewConnection(dsn string) (*sqlx.DB, error) {
{{if eq .Database.ID "postgres"}}	db, err := sqlx.Connect("postgres", dsn)
{{else if eq .Database.ID "mysql"}}	db, err := sqlx.Connect("mysql", dsn)
{{else if eq .Database.ID "sqlite"}}	db, err := sqlx.Connect("sqlite3", dsn)
{{else}}	db, err := sqlx.Connect("postgres", dsn)
{{end}}	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// Migrate runs database migrations (manual schema creation)
func (s *Storage) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create users table - adjust SQL syntax based on database type
{{if eq .Database.ID "postgres"}}	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(320) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
		)`
{{else if eq .Database.ID "mysql"}}	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(320) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL,
			INDEX idx_users_deleted_at (deleted_at)
		)`
{{else if eq .Database.ID "sqlite"}}	query := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME NULL
		)`
{{else}}	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(320) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP NULL
		)`
{{end}}

	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}

// Health checks database connectivity
func (s *Storage) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}

// User repository methods

// CreateUser creates a new user
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
{{if eq .Database.ID "postgres"}}	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (:name, :email, :password, :created_at, :updated_at) RETURNING id`
{{else if eq .Database.ID "mysql"}}	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (:name, :email, :password, :created_at, :updated_at)`
{{else if eq .Database.ID "sqlite"}}	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (:name, :email, :password, :created_at, :updated_at)`
{{else}}	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (:name, :email, :password, :created_at, :updated_at) RETURNING id`
{{end}}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

{{if eq .Database.ID "postgres"}}	rows, err := s.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&user.ID)
		if err != nil {
			return fmt.Errorf("failed to scan user ID: %w", err)
		}
	}
{{else}}	result, err := s.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	user.ID = uint(id)
{{end}}
	return nil
}

// GetUserByID retrieves a user by ID
func (s *Storage) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	err := s.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL`

	user := &domain.User{}
	err := s.db.GetContext(ctx, user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = :name, email = :email, password = :password, updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`

	user.UpdatedAt = time.Now()

	result, err := s.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// DeleteUser soft deletes a user
func (s *Storage) DeleteUser(ctx context.Context, id uint) error {
	query := `UPDATE users SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

	now := time.Now()
	result, err := s.db.ExecContext(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// ListUsers retrieves users with pagination
func (s *Storage) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
{{if eq .Database.ID "postgres"}}	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
{{else}}	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`
{{end}}

	var users []*domain.User
	err := s.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (s *Storage) CountUsers(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`

	var count int64
	err := s.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// SearchUsers searches users by name or email
func (s *Storage) SearchUsers(ctx context.Context, searchQuery string, limit, offset int) ([]*domain.User, error) {
	searchPattern := "%" + searchQuery + "%"

{{if eq .Database.ID "postgres"}}	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE (name ILIKE $1 OR email ILIKE $2) AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`
{{else}}	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE (name LIKE ? OR email LIKE ?) AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`
{{end}}

	var users []*domain.User
	err := s.db.SelectContext(ctx, &users, query, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// Transaction wraps operations in a database transaction
func (s *Storage) Transaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed: %v, rollback failed: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.db.Close()
}

// TODO: Add more repository methods as needed for your application

// Advanced sqlx features examples:

// GetUsersByIDs retrieves multiple users by their IDs using IN clause
func (s *Storage) GetUsersByIDs(ctx context.Context, ids []uint) ([]*domain.User, error) {
	if len(ids) == 0 {
		return []*domain.User{}, nil
	}

	query, args, err := sqlx.In(`
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id IN (?) AND deleted_at IS NULL
		ORDER BY created_at DESC`, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build IN query: %w", err)
	}

	query = s.db.Rebind(query)
	var users []*domain.User
	err = s.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}

	return users, nil
}

// BatchCreateUsers creates multiple users in a single transaction
func (s *Storage) BatchCreateUsers(ctx context.Context, users []*domain.User) error {
	if len(users) == 0 {
		return nil
	}

	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (:name, :email, :password, :created_at, :updated_at)`

	now := time.Now()
	for _, user := range users {
		user.CreatedAt = now
		user.UpdatedAt = now
	}

	_, err := s.db.NamedExecContext(ctx, query, users)
	if err != nil {
		return fmt.Errorf("failed to batch create users: %w", err)
	}

	return nil
}
