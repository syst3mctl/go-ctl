package {{.DbDriver.ID | replace "-" "_"}}

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

// Storage implements the storage layer using MongoDB driver
type Storage struct {
	client   *mongo.Client
	database *mongo.Database
}

// New creates a new MongoDB storage instance
func New(client *mongo.Client, databaseName string) *Storage {
	return &Storage{
		client:   client,
		database: client.Database(databaseName),
	}
}

// NewConnection creates a new MongoDB connection
func NewConnection(uri string) (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Set connection pool settings
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMinPoolSize(5)
	clientOptions.SetMaxConnIdleTime(30 * time.Second)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the primary
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// Migrate runs database migrations (create indexes)
func (s *Storage) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.database.Collection("users")

	// Create unique index on email
	emailIndex := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}

	// Create index on created_at for sorting
	createdAtIndex := mongo.IndexModel{
		Keys: bson.D{{"created_at", -1}},
	}

	// Create indexes
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		emailIndex,
		createdAtIndex,
	})
	if err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}

// Health checks database connectivity
func (s *Storage) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.client.Ping(ctx, nil)
}

// User repository methods

// CreateUser creates a new user
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	collection := s.database.Collection("users")

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Convert to BSON document
	doc := bson.M{
		"name":       user.Name,
		"email":      user.Email,
		"password":   user.Password,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = uint(time.Now().Unix()) // Convert ObjectID to uint for compatibility
		user.ObjectID = oid.Hex()
	}

	return nil
}

// GetUserByID retrieves a user by ID (using ObjectID)
func (s *Storage) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	// For MongoDB, we would typically use ObjectID, but for compatibility
	// we'll search by the created_at timestamp converted to uint
	collection := s.database.Collection("users")

	filter := bson.M{
		"$or": []bson.M{
			{"_id": id},
			{"created_at": bson.M{"$gte": time.Unix(int64(id), 0), "$lt": time.Unix(int64(id+1), 0)}},
		},
		"deleted_at": bson.M{"$exists": false},
	}

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := &domain.User{}
	if err := s.decodeBSONToUser(result, user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	collection := s.database.Collection("users")

	filter := bson.M{
		"email":      email,
		"deleted_at": bson.M{"$exists": false},
	}

	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	user := &domain.User{}
	if err := s.decodeBSONToUser(result, user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return user, nil
}

// GetUserByObjectID retrieves a user by MongoDB ObjectID
func (s *Storage) GetUserByObjectID(ctx context.Context, objectID string) (*domain.User, error) {
	collection := s.database.Collection("users")

	oid, err := primitive.ObjectIDFromHex(objectID)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %w", err)
	}

	filter := bson.M{
		"_id":        oid,
		"deleted_at": bson.M{"$exists": false},
	}

	var result bson.M
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by ObjectID: %w", err)
	}

	user := &domain.User{}
	if err := s.decodeBSONToUser(result, user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	collection := s.database.Collection("users")

	user.UpdatedAt = time.Now()

	var filter bson.M
	if user.ObjectID != "" {
		oid, err := primitive.ObjectIDFromHex(user.ObjectID)
		if err != nil {
			return fmt.Errorf("invalid ObjectID: %w", err)
		}
		filter = bson.M{
			"_id":        oid,
			"deleted_at": bson.M{"$exists": false},
		}
	} else {
		filter = bson.M{
			"email":      user.Email,
			"deleted_at": bson.M{"$exists": false},
		}
	}

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"password":   user.Password,
			"updated_at": user.UpdatedAt,
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// DeleteUser soft deletes a user
func (s *Storage) DeleteUser(ctx context.Context, id uint) error {
	collection := s.database.Collection("users")

	filter := bson.M{
		"$or": []bson.M{
			{"_id": id},
			{"created_at": bson.M{"$gte": time.Unix(int64(id), 0), "$lt": time.Unix(int64(id+1), 0)}},
		},
		"deleted_at": bson.M{"$exists": false},
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// ListUsers retrieves users with pagination
func (s *Storage) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	collection := s.database.Collection("users")

	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		user := &domain.User{}
		if err := s.decodeBSONToUser(result, user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (s *Storage) CountUsers(ctx context.Context) (int64, error) {
	collection := s.database.Collection("users")

	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// SearchUsers searches users by name or email using regex
func (s *Storage) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	collection := s.database.Collection("users")

	filter := bson.M{
		"$and": []bson.M{
			{
				"$or": []bson.M{
					{"name": bson.M{"$regex": query, "$options": "i"}},
					{"email": bson.M{"$regex": query, "$options": "i"}},
				},
			},
			{"deleted_at": bson.M{"$exists": false}},
		},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", -1}})
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		user := &domain.User{}
		if err := s.decodeBSONToUser(result, user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}

// Transaction wraps operations in a MongoDB transaction
func (s *Storage) Transaction(ctx context.Context, fn func(mongo.SessionContext) error) error {
	session, err := s.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		return nil, fn(sc)
	})

	return err
}

// Close closes the database connection
func (s *Storage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.client.Disconnect(ctx)
}

// Helper function to decode BSON to User struct
func (s *Storage) decodeBSONToUser(result bson.M, user *domain.User) error {
	// Handle ObjectID
	if oid, ok := result["_id"].(primitive.ObjectID); ok {
		user.ObjectID = oid.Hex()
		user.ID = uint(time.Now().Unix()) // Convert for compatibility
	}

	// Handle other fields
	if name, ok := result["name"].(string); ok {
		user.Name = name
	}
	if email, ok := result["email"].(string); ok {
		user.Email = email
	}
	if password, ok := result["password"].(string); ok {
		user.Password = password
	}
	if createdAt, ok := result["created_at"].(primitive.DateTime); ok {
		user.CreatedAt = createdAt.Time()
	}
	if updatedAt, ok := result["updated_at"].(primitive.DateTime); ok {
		user.UpdatedAt = updatedAt.Time()
	}

	return nil
}

// Advanced MongoDB operations

// GetUsersByTags retrieves users by tags using array operations
func (s *Storage) GetUsersByTags(ctx context.Context, tags []string) ([]*domain.User, error) {
	collection := s.database.Collection("users")

	filter := bson.M{
		"tags": bson.M{"$in": tags},
		"deleted_at": bson.M{"$exists": false},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by tags: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		user := &domain.User{}
		if err := s.decodeBSONToUser(result, user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// AggregateUsersByCreationDate aggregates users by creation date
func (s *Storage) AggregateUsersByCreationDate(ctx context.Context) ([]bson.M, error) {
	collection := s.database.Collection("users")

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"deleted_at": bson.M{"$exists": false},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"year":  bson.M{"$year": "$created_at"},
					"month": bson.M{"$month": "$created_at"},
					"day":   bson.M{"$dayOfMonth": "$created_at"},
				},
				"count": bson.M{"$sum": 1},
			},
		},
		{
			"$sort": bson.M{"_id": -1},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate users: %w", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode aggregation results: %w", err)
	}

	return results, nil
}

// BulkCreateUsers creates multiple users using bulk operations
func (s *Storage) BulkCreateUsers(ctx context.Context, users []*domain.User) error {
	if len(users) == 0 {
		return nil
	}

	collection := s.database.Collection("users")

	var operations []mongo.WriteModel
	now := time.Now()

	for _, user := range users {
		user.CreatedAt = now
		user.UpdatedAt = now

		doc := bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"password":   user.Password,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		}

		operation := mongo.NewInsertOneModel().SetDocument(doc)
		operations = append(operations, operation)
	}

	_, err := collection.BulkWrite(ctx, operations)
	if err != nil {
		return fmt.Errorf("failed to bulk create users: %w", err)
	}

	return nil
}

// TODO: Add more repository methods as needed for your application
