package mongodb

import (
	"context"
	"fmt"
	"time"

	"{{.ProjectName}}/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository implements MongoDB-specific repository operations
type Repository struct {
	client   *mongo.Client
	database *mongo.Database
}

// New creates a new MongoDB repository instance using mongo-driver
func New(client *mongo.Client) *Repository {
	// Default database name - you can make this configurable
	database := client.Database("{{.ProjectName}}")
	return &Repository{
		client:   client,
		database: database,
	}
}

// Migrate creates indexes for collections
func (r *Repository) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create indexes for users collection
	usersCollection := r.database.Collection("users")

	// Email unique index
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := usersCollection.Indexes().CreateOne(ctx, emailIndex); err != nil {
		return fmt.Errorf("failed to create email index: %w", err)
	}

{{if ne .Database.ID "redis"}}	// Create indexes for products collection
	productsCollection := r.database.Collection("products")

	// User ID index
	userIDIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	}

	// Name text index for searching
	nameIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}},
	}

	if _, err := productsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{userIDIndex, nameIndex}); err != nil {
		return fmt.Errorf("failed to create product indexes: %w", err)
	}
{{end}}

	return nil
}

// User repository methods

// CreateUser creates a new user in MongoDB
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	now := time.Now()
	user.ID = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now

	collection := r.database.Collection("users")
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID from MongoDB
func (r *Repository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	collection := r.database.Collection("users")

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email from MongoDB
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	collection := r.database.Collection("users")

	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// UpdateUser updates an existing user in MongoDB
func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()
	collection := r.database.Collection("users")

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"active":     user.Active,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user
func (r *Repository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	collection := r.database.Collection("users")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ListUsers retrieves users with pagination from MongoDB
func (r *Repository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	collection := r.database.Collection("users")

	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	if offset > 0 {
		findOptions.SetSkip(int64(offset))
	}

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users in MongoDB
func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	collection := r.database.Collection("users")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// SearchUsers searches users by name or email
func (r *Repository) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	collection := r.database.Collection("users")

	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": query, "$options": "i"}},
			{"email": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	if offset > 0 {
		findOptions.SetSkip(int64(offset))
	}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}

	return users, nil
}

{{if ne .Database.ID "redis"}}// Product repository methods

// CreateProduct creates a new product in MongoDB
func (r *Repository) CreateProduct(ctx context.Context, product *domain.Product) error {
	now := time.Now()
	product.ID = primitive.NewObjectID()
	product.CreatedAt = now
	product.UpdatedAt = now

	collection := r.database.Collection("products")
	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetProductByID retrieves a product by ID from MongoDB
func (r *Repository) GetProductByID(ctx context.Context, id primitive.ObjectID) (*domain.Product, error) {
	var product domain.Product
	collection := r.database.Collection("products")

	// Use aggregation to populate user data
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id}},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": bson.M{
			"path":                       "$user",
			"preserveNullAndEmptyArrays": true,
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to query product: %w", err)
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, fmt.Errorf("product not found")
	}

	if err := cursor.Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product: %w", err)
	}

	return &product, nil
}

// GetProductsByUserID retrieves products by user ID from MongoDB
func (r *Repository) GetProductsByUserID(ctx context.Context, userID primitive.ObjectID) ([]*domain.Product, error) {
	var products []*domain.Product
	collection := r.database.Collection("products")

	filter := bson.M{"user_id": userID}
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

// UpdateProduct updates an existing product in MongoDB
func (r *Repository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	product.UpdatedAt = time.Now()
	collection := r.database.Collection("products")

	filter := bson.M{"_id": product.ID}
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"available":   product.Available,
			"updated_at":  product.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// DeleteProduct soft deletes a product in MongoDB
func (r *Repository) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	collection := r.database.Collection("products")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// ListProducts retrieves products with pagination from MongoDB
func (r *Repository) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	collection := r.database.Collection("products")

	// Use aggregation to populate user data
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": bson.M{
			"path":                       "$user",
			"preserveNullAndEmptyArrays": true,
		}},
		{"$sort": bson.M{"created_at": -1}},
	}

	if offset > 0 {
		pipeline = append(pipeline, bson.M{"$skip": offset})
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

// SearchProducts searches products by name or description in MongoDB
func (r *Repository) SearchProducts(ctx context.Context, query string, limit, offset int) ([]*domain.Product, error) {
	var products []*domain.Product
	collection := r.database.Collection("products")

	// Use text search if available, otherwise regex
	filter := bson.M{
		"$text": bson.M{"$search": query},
	}

	// Fallback to regex search if text search fails
	pipeline := []bson.M{
		{"$match": filter},
		{"$lookup": bson.M{
			"from":         "users",
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user",
		}},
		{"$unwind": bson.M{
			"path":                       "$user",
			"preserveNullAndEmptyArrays": true,
		}},
		{"$sort": bson.M{"created_at": -1}},
	}

	if offset > 0 {
		pipeline = append(pipeline, bson.M{"$skip": offset})
	}
	if limit > 0 {
		pipeline = append(pipeline, bson.M{"$limit": limit})
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		// Fallback to regex search
		regexFilter := bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": query, "$options": "i"}},
				{"description": bson.M{"$regex": query, "$options": "i"}},
			},
		}

		pipeline[0] = bson.M{"$match": regexFilter}
		cursor, err = collection.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, fmt.Errorf("failed to search products: %w", err)
		}
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

// CountProducts returns the total number of products in MongoDB
func (r *Repository) CountProducts(ctx context.Context) (int64, error) {
	collection := r.database.Collection("products")

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

{{end}}// Transaction wraps operations in a MongoDB session transaction
func (r *Repository) Transaction(ctx context.Context, fn func(mongo.SessionContext) error) error {
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	})

	return err
}

// Close closes the database connection
func (r *Repository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.client.Disconnect(ctx)
}
