package repository

import (
	"context"
	"{{.ProjectName}}/internal/domain"
{{if eq .DbDriver.ID "gorm"}}
	"gorm.io/gorm"
{{else if eq .DbDriver.ID "sqlx"}}	"database/sql"
	"github.com/jmoiron/sqlx"
{{else if eq .DbDriver.ID "mongo-driver"}}	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
{{else if eq .DbDriver.ID "redis-client"}}	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
{{else if eq .DbDriver.ID "database-sql"}}	"database/sql"
{{end}}{{if .HasFeature "logging"}}	"github.com/rs/zerolog/log"
{{end}})

{{if eq .DbDriver.ID "gorm"}}// userRepository implements domain.UserRepository using GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM-based user repository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user")
{{end}}		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Uint("id", id).Msg("Failed to get user by ID")
{{end}}		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("email", email).Msg("Failed to get user by email")
{{end}}		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(user *domain.User) error {
	if err := r.db.Save(user).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Uint("id", user.ID).Msg("Failed to update user")
{{end}}		return err
	}
	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Uint("id", id).Msg("Failed to delete user")
{{end}}		return err
	}
	return nil
}

// List retrieves a paginated list of users
func (r *userRepository) List(offset, limit int) ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list users")
{{end}}		return nil, err
	}
	return users, nil
}

{{if ne .Database.ID "redis"}}// productRepository implements domain.ProductRepository using GORM
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new GORM-based product repository
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(product *domain.Product) error {
	if err := r.db.Create(product).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create product")
{{end}}		return err
	}
	return nil
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(id uint) (*domain.Product, error) {
	var product domain.Product
	if err := r.db.Preload("User").First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrProductNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Uint("id", id).Msg("Failed to get product by ID")
{{end}}		return nil, err
	}
	return &product, nil
}

// GetByUserID retrieves products by user ID
func (r *productRepository) GetByUserID(userID uint) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.Where("user_id = ?", userID).Find(&products).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Uint("user_id", userID).Msg("Failed to get products by user ID")
{{end}}		return nil, err
	}
	return products, nil
}

// Update updates an existing product
func (r *productRepository) Update(product *domain.Product) error {
	if err := r.db.Save(product).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Uint("id", product.ID).Msg("Failed to update product")
{{end}}		return err
	}
	return nil
}

// Delete deletes a product by ID
func (r *productRepository) Delete(id uint) error {
	if err := r.db.Delete(&domain.Product{}, id).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Uint("id", id).Msg("Failed to delete product")
{{end}}		return err
	}
	return nil
}

// List retrieves a paginated list of products
func (r *productRepository) List(offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.Preload("User").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list products")
{{end}}		return nil, err
	}
	return products, nil
}

// Search searches products by query string
func (r *productRepository) Search(query string, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	if err := r.db.Preload("User").Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("query", query).Msg("Failed to search products")
{{end}}		return nil, err
	}
	return products, nil
}
{{end}}

{{else if eq .DbDriver.ID "sqlx"}}// userRepository implements domain.UserRepository using sqlx
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new sqlx-based user repository
func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (name, email, active, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowx(query, user.Name, user.Email, user.Active, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user")
{{end}}		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, name, email, active, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", id).Msg("Failed to get user by ID")
{{end}}		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `SELECT id, name, email, active, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("email", email).Msg("Failed to get user by email")
{{end}}		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(user *domain.User) error {
	query := `UPDATE users SET name = $1, email = $2, active = $3, updated_at = $4 WHERE id = $5`
	_, err := r.db.Exec(query, user.Name, user.Email, user.Active, user.UpdatedAt, user.ID)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", user.ID).Msg("Failed to update user")
{{end}}		return err
	}
	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", id).Msg("Failed to delete user")
{{end}}		return err
	}
	return nil
}

// List retrieves a paginated list of users
func (r *userRepository) List(offset, limit int) ([]*domain.User, error) {
	var users []*domain.User
	query := `SELECT id, name, email, active, created_at, updated_at FROM users ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.Select(&users, query, limit, offset)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list users")
{{end}}		return nil, err
	}
	return users, nil
}

{{if ne .Database.ID "redis"}}// productRepository implements domain.ProductRepository using sqlx
type productRepository struct {
	db *sqlx.DB
}

// NewProductRepository creates a new sqlx-based product repository
func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

// Create creates a new product
func (r *productRepository) Create(product *domain.Product) error {
	query := `INSERT INTO products (name, description, price, available, user_id, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRowx(query, product.Name, product.Description, product.Price, product.Available, product.UserID, product.CreatedAt, product.UpdatedAt).Scan(&product.ID)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create product")
{{end}}		return err
	}
	return nil
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(id int64) (*domain.Product, error) {
	var product domain.Product
	query := `SELECT id, name, description, price, available, user_id, created_at, updated_at FROM products WHERE id = $1`
	err := r.db.Get(&product, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrProductNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", id).Msg("Failed to get product by ID")
{{end}}		return nil, err
	}
	return &product, nil
}

// GetByUserID retrieves products by user ID
func (r *productRepository) GetByUserID(userID int64) ([]*domain.Product, error) {
	var products []*domain.Product
	query := `SELECT id, name, description, price, available, user_id, created_at, updated_at FROM products WHERE user_id = $1`
	err := r.db.Select(&products, query, userID)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("user_id", userID).Msg("Failed to get products by user ID")
{{end}}		return nil, err
	}
	return products, nil
}

// Update updates an existing product
func (r *productRepository) Update(product *domain.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, available = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Available, product.UpdatedAt, product.ID)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", product.ID).Msg("Failed to update product")
{{end}}		return err
	}
	return nil
}

// Delete deletes a product by ID
func (r *productRepository) Delete(id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", id).Msg("Failed to delete product")
{{end}}		return err
	}
	return nil
}

// List retrieves a paginated list of products
func (r *productRepository) List(offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	query := `SELECT id, name, description, price, available, user_id, created_at, updated_at FROM products ORDER BY id LIMIT $1 OFFSET $2`
	err := r.db.Select(&products, query, limit, offset)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list products")
{{end}}		return nil, err
	}
	return products, nil
}

// Search searches products by query string
func (r *productRepository) Search(query string, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	sqlQuery := `SELECT id, name, description, price, available, user_id, created_at, updated_at
				 FROM products
				 WHERE name ILIKE $1 OR description ILIKE $1
				 ORDER BY id LIMIT $2 OFFSET $3`
	searchTerm := "%" + query + "%"
	err := r.db.Select(&products, sqlQuery, searchTerm, limit, offset)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("query", query).Msg("Failed to search products")
{{end}}		return nil, err
	}
	return products, nil
}
{{end}}

{{else if eq .DbDriver.ID "mongo-driver"}}// userRepository implements domain.UserRepository using MongoDB driver
type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new MongoDB-based user repository
func NewUserRepository(db *mongo.Client) domain.UserRepository {
	collection := db.Database("{{.ProjectName}}").Collection("users")
	return &userRepository{collection: collection}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	result, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user")
{{end}}		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("id", id.Hex()).Msg("Failed to get user by ID")
{{end}}		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("email", email).Msg("Failed to get user by email")
{{end}}		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(user *domain.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("id", user.ID.Hex()).Msg("Failed to update user")
{{end}}		return err
	}
	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("id", id.Hex()).Msg("Failed to delete user")
{{end}}		return err
	}
	return nil
}

// List retrieves a paginated list of users
func (r *userRepository) List(offset, limit int) ([]*domain.User, error) {
	var users []*domain.User
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit))
	cursor, err := r.collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list users")
{{end}}		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}
	return users, nil
}

{{if ne .Database.ID "redis"}}// productRepository implements domain.ProductRepository using MongoDB driver
type productRepository struct {
	collection *mongo.Collection
}

// NewProductRepository creates a new MongoDB-based product repository
func NewProductRepository(db *mongo.Client) domain.ProductRepository {
	collection := db.Database("{{.ProjectName}}").Collection("products")
	return &productRepository{collection: collection}
}

// Create creates a new product
func (r *productRepository) Create(product *domain.Product) error {
	result, err := r.collection.InsertOne(context.Background(), product)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create product")
{{end}}		return err
	}
	product.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(id primitive.ObjectID) (*domain.Product, error) {
	var product domain.Product
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrProductNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("id", id.Hex()).Msg("Failed to get product by ID")
{{end}}		return nil, err
	}
	return &product, nil
}

// GetByUserID retrieves products by user ID
func (r *productRepository) GetByUserID(userID primitive.ObjectID) ([]*domain.Product, error) {
	var products []*domain.Product
	cursor, err := r.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("user_id", userID.Hex()).Msg("Failed to get products by user ID")
{{end}}		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			continue
		}
		products = append(products, &product)
	}
	return products, nil
}

// Update updates an existing product
func (r *productRepository) Update(product *domain.Product) error {
	filter := bson.M{"_id": product.ID}
	update := bson.M{"$set": product}
	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("id", product.ID.Hex()).Msg("Failed to update product")
{{end}}		return err
	}
	return nil
}

// Delete deletes a product by ID
func (r *productRepository) Delete(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("id", id.Hex()).Msg("Failed to delete product")
{{end}}		return err
	}
	return nil
}

// List retrieves a paginated list of products
func (r *productRepository) List(offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit))
	cursor, err := r.collection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list products")
{{end}}		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			continue
		}
		products = append(products, &product)
	}
	return products, nil
}

// Search searches products by query string
func (r *productRepository) Search(query string, offset, limit int) ([]*domain.Product, error) {
	var products []*domain.Product
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"description": bson.M{"$regex": query, "$options": "i"}},
		},
	}
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit))
	cursor, err := r.collection.Find(context.Background(), filter, opts)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("query", query).Msg("Failed to search products")
{{end}}		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			continue
		}
		products = append(products, &product)
	}
	return products, nil
}
{{end}}

{{else if eq .DbDriver.ID "redis-client"}}// userRepository implements domain.UserRepository using Redis
type userRepository struct {
	client *redis.Client
}

// NewUserRepository creates a new Redis-based user repository
func NewUserRepository(client *redis.Client) domain.UserRepository {
	return &userRepository{client: client}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	// Generate ID if not set
	if user.ID == 0 {
		id, err := r.client.Incr(context.Background(), "user:counter").Result()
		if err != nil {
			return err
		}
		user.ID = id
	}

	data, err := json.Marshal(user)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to marshal user")
{{end}}		return err
	}

	key := fmt.Sprintf("user:%d", user.ID)
	err = r.client.Set(context.Background(), key, data, 0).Err()
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user in Redis")
{{end}}		return err
	}

	// Index by email
	emailKey := fmt.Sprintf("user:email:%s", user.Email)
	err = r.client.Set(context.Background(), emailKey, user.ID, 0).Err()
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to index user by email")
{{end}}		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	key := fmt.Sprintf("user:%d", id)
	data, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, domain.ErrUserNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", id).Msg("Failed to get user by ID from Redis")
{{end}}		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to unmarshal user")
{{end}}		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	emailKey := fmt.Sprintf("user:email:%s", email)
	idStr, err := r.client.Get(context.Background(), emailKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, domain.ErrUserNotFound
		}
{{if .HasFeature "logging"}}		log.Error().Err(err).Str("email", email).Msg("Failed to get user ID by email from Redis")
{{end}}		return nil, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

// Update updates an existing user
func (r *userRepository) Update(user *domain.User) error {
	// Check if user exists
	_, err := r.GetByID(user.ID)
	if err != nil {
		return err
	}

	data, err := json.Marshal(user)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to marshal user for update")
{{end}}		return err
	}

	key := fmt.Sprintf("user:%d", user.ID)
	err = r.client.Set(context.Background(), key, data, 0).Err()
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", user.ID).Msg("Failed to update user in Redis")
{{end}}		return err
	}

	return nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(id int64) error {
	// Get user first to remove email index
	user, err := r.GetByID(id)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("user:%d", id)
	err = r.client.Del(context.Background(), key).Err()
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Int64("id", id).Msg("Failed to delete user from Redis")
{{end}}		return err
	}

	// Remove email index
	emailKey := fmt.Sprintf("user:email:%s", user.Email)
	r.client.Del(context.Background(), emailKey)

	return nil
}

// List retrieves a paginated list of users
func (r *userRepository) List(offset, limit int) ([]*domain.User, error) {
	// Get all user keys
	keys, err := r.client.Keys(context.Background(), "user:[0-9]*").Result()
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to get user keys from Redis")
{{end}}		return nil, err
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if end > len(keys) {
		end = len(keys)
	}
	if start >= len(keys) {
		return []*domain.User{}, nil
	}

	var users []*domain.User
	for i := start; i < end; i++ {
		data, err := r.client.Get(context.Background(), keys[i]).Result()
		if err != nil {
			continue
		}

		var user domain.User
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}

{{else}}// userRepository is a placeholder implementation for other database drivers
type userRepository struct {
	// Add your database connection here
}

// NewUserRepository creates a new user repository
func NewUserRepository(db interface{}) domain.UserRepository {
	return &userRepository{}
}

// Create creates a new user
func (r *userRepository) Create(user *domain.User) error {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("User repository Create method not implemented")
{{end}}	return domain.ErrInternalError
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id {{if eq .DbDriver.ID "mongo-driver"}}primitive.ObjectID{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}}) (*domain.User, error) {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("User repository GetByID method not implemented")
{{end}}	return nil, domain.ErrUserNotFound
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("User repository GetByEmail method not implemented")
{{end}}	return nil, domain.ErrUserNotFound
}

// Update updates an existing user
func (r *userRepository) Update(user *domain.User) error {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("User repository Update method not implemented")
{{end}}	return domain.ErrInternalError
}

// Delete deletes a user by ID
func (r *userRepository) Delete(id {{if eq .DbDriver.ID "mongo-driver"}}primitive.ObjectID{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}}) error {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("User repository Delete method not implemented")
{{end}}	return domain.ErrInternalError
}

// List retrieves a paginated list of users
func (r *userRepository) List(offset, limit int) ([]*domain.User, error) {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("User repository List method not implemented")
{{end}}	return nil, domain.ErrInternalError
}

{{if ne .Database.ID "redis"}}// productRepository is a placeholder implementation for other database drivers
type productRepository struct {
	// Add your database connection here
}

// NewProductRepository creates a new product repository
func NewProductRepository(db interface{}) domain.ProductRepository {
	return &productRepository{}
}

// Create creates a new product
func (r *productRepository) Create(product *domain.Product) error {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("Product repository Create method not implemented")
{{end}}	return domain.ErrInternalError
}

// GetByID retrieves a product by ID
func (r *productRepository) GetByID(id {{if eq .DbDriver.ID "mongo-driver"}}primitive.ObjectID{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}}) (*domain.Product, error) {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("Product repository GetByID method not implemented")
{{end}}	return nil, domain.ErrProductNotFound
}

// GetByUserID retrieves products by user ID
func (r *productRepository) GetByUserID(userID {{if eq .DbDriver.ID "mongo-driver"}}primitive.ObjectID{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}}) ([]*domain.Product, error) {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("Product repository GetByUserID method not implemented")
{{end}}	return nil, domain.ErrInternalError
}

// Update updates an existing product
func (r *productRepository) Update(product *domain.Product) error {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("Product repository Update method not implemented")
{{end}}	return domain.ErrInternalError
}

// Delete deletes a product by ID
func (r *productRepository) Delete(id {{if eq .DbDriver.ID "mongo-driver"}}primitive.ObjectID{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}}) error {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("Product repository Delete method not implemented")
{{end}}	return domain.ErrInternalError
}

// List retrieves a paginated list of products
func (r *productRepository) List(offset, limit int) ([]*domain.Product, error) {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("Product repository List method not implemented")
{{end}}	return nil, domain.ErrInternalError
}

// Search searches products by query string
func (r *productRepository) Search(query string, offset, limit int) ([]*domain.Product, error) {
	// Implement based on your database driver
{{if .HasFeature "logging"}}	log.Warn().Msg("Product repository Search method not implemented")
{{end}}	return nil, domain.ErrInternalError
}
{{end}}
{{end}}
