package {{.DbDriver.ID | replace "-" "_"}}

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"{{.ProjectName}}/internal/domain"

{{if eq .Database.ID "postgres"}}	_ "github.com/lib/pq"
{{else if eq .Database.ID "mysql"}}	_ "github.com/go-sql-driver/mysql"
{{else if eq .Database.ID "sqlite"}}	_ "github.com/mattn/go-sqlite3"
{{end}})

// Storage implements the storage layer using database/sql
type Storage struct {
	db *sql.DB
}

// New creates a new database/sql storage instance
func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// NewConnection creates a new database connection using database/sql
func NewConnection(dsn string) (*sql.DB, error) {
{{if eq .Database.ID "postgres"}}	db, err := sql.Open("postgres", dsn)
{{else if eq .Database.ID "mysql"}}	db, err := sql.Open("mysql", dsn)
{{else if eq .Database.ID "sqlite"}}	db, err := sql.Open("sqlite3", dsn)
{{else}}	db, err := sql.Open("postgres", dsn)
{{end}}	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
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
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
{{else if eq .Database.ID "mysql"}}	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`
{{else if eq .Database.ID "sqlite"}}	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`
{{else}}	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
{{end}}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

{{if eq .Database.ID "postgres"}}	err := s.db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
{{else}}	result, err := s.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err == nil {
		id, err := result.LastInsertId()
		if err == nil {
			user.ID = uint(id)
		}
	}
{{end}}	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (s *Storage) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
{{if eq .Database.ID "postgres"}}	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL`
{{else}}	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL`
{{end}}

	user := &domain.User{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
{{if eq .Database.ID "postgres"}}	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1 AND deleted_at IS NULL`
{{else}}	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = ? AND deleted_at IS NULL`
{{end}}

	user := &domain.User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
{{if eq .Database.ID "postgres"}}	query := `
		UPDATE users
		SET name = $2, email = $3, password = $4, updated_at = $5
		WHERE id = $1 AND deleted_at IS NULL`
{{else}}	query := `
		UPDATE users
		SET name = ?, email = ?, password = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL`
{{end}}

	user.UpdatedAt = time.Now()

{{if eq .Database.ID "postgres"}}	result, err := s.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.Password, user.UpdatedAt)
{{else}}	result, err := s.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
{{end}}	if err != nil {
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
{{if eq .Database.ID "postgres"}}	query := `UPDATE users SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`
{{else}}	query := `UPDATE users SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`
{{end}}

	now := time.Now()
{{if eq .Database.ID "postgres"}}	result, err := s.db.ExecContext(ctx, query, id, now)
{{else}}	result, err := s.db.ExecContext(ctx, query, now, id)
{{end}}	if err != nil {
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
{{else if eq .Database.ID "mysql"}}	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`
{{else}}	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`
{{end}}

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Password,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (s *Storage) CountUsers(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`

	var count int64
	err := s.db.QueryRowContext(ctx, query).Scan(&count)
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

	rows, err := s.db.QueryContext(ctx, query, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Password,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// Transaction wraps operations in a database transaction
func (s *Storage) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
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
