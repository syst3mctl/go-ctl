package mysql

import (
	"context"
	"fmt"
	"time"

	"{{.ProjectName}}/internal/domain"

{{if eq .DbDriver.ID "gorm"}}	"gorm.io/gorm"
{{else if eq .DbDriver.ID "sqlx"}}	"github.com/jmoiron/sqlx"
{{end}})

// Repository implements MySQL-specific repository operations
type Repository struct {
{{if eq .DbDriver.ID "gorm"}}	db *gorm.DB
{{else if eq .DbDriver.ID "sqlx"}}	db *sqlx.DB
{{end}}
}

{{if eq .DbDriver.ID "gorm"}}// New creates a new MySQL repository instance using GORM
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Migrate runs database migrations for MySQL using GORM
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

// CreateUser creates a new user in MySQL
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID from MySQL
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

// GetUserByEmail retrieves a user by email from MySQL
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

// UpdateUser updates an existing user in MySQL
func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser soft deletes a user in MySQL
func (r *Repository) DeleteUser(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers retrieves users with pagination from MySQL
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

// CountUsers returns the total number of users in MySQL
func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.User{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// SearchUsers searches users by name or email in MySQL
func (r *Repository) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	searchQuery := "%" + query + "%"

	dbQuery := r.db.WithContext(ctx).
		Where("name LIKE ? OR email LIKE ?", searchQuery, searchQuery).
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

// CreateProduct creates a new product in MySQL
func (r *Repository) CreateProduct(ctx context.Context, product *domain.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

// GetProductByID retrieves a product by ID from MySQL
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

// ListProducts retrieves products with pagination from MySQL
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

// Transaction wraps operations in a MySQL transaction using GORM
func (r *Repository) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

{{else if eq .DbDriver.ID "sqlx"}}// New creates a new MySQL repository instance using sqlx
func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Migrate runs database migrations for MySQL using sqlx
func (r *Repository) Migrate() error {
	// Create users table
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(320) UNIQUE NOT NULL,
		active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	if _, err := r.db.Exec(userTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create addresses table
	addressTable := `
	CREATE TABLE IF NOT EXISTS addresses (
		id INT AUTO_INCREMENT PRIMARY KEY,
		street VARCHAR(200) NOT NULL,
		city VARCHAR(100) NOT NULL,
		state VARCHAR(100) NOT NULL,
		zip_code VARCHAR(20) NOT NULL,
		country VARCHAR(100) NOT NULL,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	if _, err := r.db.Exec(addressTable); err != nil {
		return fmt.Errorf("failed to create addresses table: %w", err)
	}

	// Create categories table
	categoryTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		description TEXT,
		active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	if _, err := r.db.Exec(categoryTable); err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}

	// Create products table
	productTable := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(200) NOT NULL,
		description TEXT,
		price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
		available BOOLEAN DEFAULT true,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	if _, err := r.db.Exec(productTable); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	// Create orders table
	orderTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		order_number VARCHAR(50) UNIQUE NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		total DECIMAL(10,2) NOT NULL CHECK (total >= 0),
		order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`

	if _, err := r.db.Exec(orderTable); err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_users_active ON users(active)",
		"CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_user_id ON products(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_products_available ON products(available)",
		"CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status)",
		"CREATE INDEX IF NOT EXISTS idx_categories_active ON categories(active)",
	}

	for _, index := range indexes {
		if _, err := r.db.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// User repository methods using sqlx

// CreateUser creates a new user in MySQL
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Active, now, now)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetUserByID retrieves a user by ID from MySQL
func (r *Repository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}
	query := "SELECT id, name, email, active, created_at, updated_at FROM users WHERE id = ?"

	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email from MySQL
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := "SELECT id, name, email, active, created_at, updated_at FROM users WHERE email = ?"

	err := r.db.GetContext(ctx, user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// UpdateUser updates an existing user in MySQL
func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, active = ?, updated_at = ?
		WHERE id = ?`

	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Active, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user from MySQL (hard delete for sqlx)
func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// ListUsers retrieves users with pagination from MySQL
func (r *Repository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users := []*domain.User{}
	query := `
		SELECT id, name, email, active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// CountUsers returns the total number of users in MySQL
func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	query := "SELECT COUNT(*) FROM users"

	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// SearchUsers searches users by name or email in MySQL
func (r *Repository) SearchUsers(ctx context.Context, searchQuery string, limit, offset int) ([]*domain.User, error) {
	users := []*domain.User{}
	query := `
		SELECT id, name, email, active, created_at, updated_at
		FROM users
		WHERE name LIKE ? OR email LIKE ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	searchPattern := "%" + searchQuery + "%"
	err := r.db.SelectContext(ctx, &users, query, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	return users, nil
}

// Product repository methods using sqlx

// CreateProduct creates a new product in MySQL
func (r *Repository) CreateProduct(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (name, description, price, available, user_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		product.Name, product.Description, product.Price, product.Available, product.UserID, now, now)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get product ID: %w", err)
	}

	product.ID = id
	product.CreatedAt = now
	product.UpdatedAt = now
	return nil
}

// GetProductByID retrieves a product by ID from MySQL
func (r *Repository) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	product := &domain.Product{}
	query := `
		SELECT p.id, p.name, p.description, p.price, p.available, p.user_id, p.created_at, p.updated_at
		FROM products p
		WHERE p.id = ?`

	err := r.db.GetContext(ctx, product, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

// ListProducts retrieves products with pagination from MySQL
func (r *Repository) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	products := []*domain.Product{}
	query := `
		SELECT p.id, p.name, p.description, p.price, p.available, p.user_id, p.created_at, p.updated_at
		FROM products p
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`

	err := r.db.SelectContext(ctx, &products, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

// Transaction wraps operations in a MySQL transaction using sqlx
func (r *Repository) Transaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

{{end}}

// Close closes the database connection
func (r *Repository) Close() error {
{{if eq .DbDriver.ID "gorm"}}	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Close()
{{else if eq .DbDriver.ID "sqlx"}}	return r.db.Close()
{{end}}
}
