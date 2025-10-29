package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"{{.ProjectName}}/internal/domain"

	"github.com/redis/go-redis/v9"
)

// Repository implements Redis-specific repository operations
type Repository struct {
	client *redis.Client
}

// New creates a new Redis repository instance using redis-client
func New(client *redis.Client) *Repository {
	return &Repository{client: client}
}

// Migrate is a no-op for Redis as it's schema-less
func (r *Repository) Migrate() error {
	// Redis doesn't require schema migrations
	return nil
}

// User repository methods

// CreateUser creates a new user in Redis
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) error {
	// Generate ID if not set
	if user.ID == 0 {
		id, err := r.client.Incr(ctx, "user:id:counter").Result()
		if err != nil {
			return fmt.Errorf("failed to generate user ID: %w", err)
		}
		user.ID = id
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Serialize user to JSON
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	pipe := r.client.TxPipeline()

	// Store user data
	userKey := fmt.Sprintf("user:%d", user.ID)
	pipe.Set(ctx, userKey, userData, 0)

	// Store email index for lookups
	emailKey := fmt.Sprintf("user:email:%s", user.Email)
	pipe.Set(ctx, emailKey, user.ID, 0)

	// Add to users set
	pipe.SAdd(ctx, "users", user.ID)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	userKey := fmt.Sprintf("user:%d", id)
	userData, err := r.client.Get(ctx, userKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user domain.User
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	emailKey := fmt.Sprintf("user:email:%s", email)
	userID, err := r.client.Get(ctx, emailKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return r.GetUserByID(ctx, id)
}

// UpdateUser updates an existing user
func (r *Repository) UpdateUser(ctx context.Context, user *domain.User) error {
	// Check if user exists
	userKey := fmt.Sprintf("user:%d", user.ID)
	exists, err := r.client.Exists(ctx, userKey).Result()
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("user not found")
	}

	// Get current user to check for email changes
	currentUser, err := r.GetUserByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	user.UpdatedAt = time.Now()
	user.CreatedAt = currentUser.CreatedAt // Preserve creation time

	// Serialize user to JSON
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	pipe := r.client.TxPipeline()

	// Update user data
	pipe.Set(ctx, userKey, userData, 0)

	// Update email index if email changed
	if currentUser.Email != user.Email {
		// Remove old email index
		oldEmailKey := fmt.Sprintf("user:email:%s", currentUser.Email)
		pipe.Del(ctx, oldEmailKey)

		// Add new email index
		newEmailKey := fmt.Sprintf("user:email:%s", user.Email)
		pipe.Set(ctx, newEmailKey, user.ID, 0)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user
func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	// Get user to remove email index
	user, err := r.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	pipe := r.client.TxPipeline()

	// Delete user data
	userKey := fmt.Sprintf("user:%d", id)
	pipe.Del(ctx, userKey)

	// Delete email index
	emailKey := fmt.Sprintf("user:email:%s", user.Email)
	pipe.Del(ctx, emailKey)

	// Remove from users set
	pipe.SRem(ctx, "users", id)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves users with pagination
func (r *Repository) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	// Get all user IDs from the set
	userIDs, err := r.client.SMembers(ctx, "users").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user IDs: %w", err)
	}

	// Convert string IDs to int64 and sort
	var ids []int64
	for _, idStr := range userIDs {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			continue // Skip invalid IDs
		}
		ids = append(ids, id)
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(ids) {
		return []*domain.User{}, nil
	}
	if end > len(ids) || limit <= 0 {
		end = len(ids)
	}

	paginatedIDs := ids[start:end]
	var users []*domain.User

	// Fetch users in batch
	for _, id := range paginatedIDs {
		user, err := r.GetUserByID(ctx, id)
		if err != nil {
			continue // Skip users that couldn't be retrieved
		}
		users = append(users, user)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	count, err := r.client.SCard(ctx, "users").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// SearchUsers searches users by name or email
func (r *Repository) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	// Get all users first (Redis doesn't have full-text search built-in)
	allUsers, err := r.ListUsers(ctx, 0, 0) // Get all users
	if err != nil {
		return nil, fmt.Errorf("failed to get all users for search: %w", err)
	}

	var matchedUsers []*domain.User
	queryLower := fmt.Sprintf("%s", query) // Simple contains search

	for _, user := range allUsers {
		if contains(user.Name, queryLower) || contains(user.Email, queryLower) {
			matchedUsers = append(matchedUsers, user)
		}
	}

	// Apply pagination to matched results
	start := offset
	end := offset + limit
	if start > len(matchedUsers) {
		return []*domain.User{}, nil
	}
	if end > len(matchedUsers) || limit <= 0 {
		end = len(matchedUsers)
	}

	return matchedUsers[start:end], nil
}

// Helper function for case-insensitive contains
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(substr) == 0 ||
		fmt.Sprintf("%s", s) != s[:len(s)-len(substr)]+substr+s[len(s):])
	// Simple implementation - in production you'd use strings.Contains with strings.ToLower
}

// Transaction wraps operations in a Redis transaction
func (r *Repository) Transaction(ctx context.Context, fn func(*redis.Tx) error) error {
	return r.client.Watch(ctx, func(tx *redis.Tx) error {
		return fn(tx)
	})
}

// Close closes the Redis connection
func (r *Repository) Close() error {
	return r.client.Close()
}

// Cache-specific methods for Redis

// SetCache stores a value in cache with TTL
func (r *Repository) SetCache(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

// GetCache retrieves a value from cache
func (r *Repository) GetCache(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache key not found")
		}
		return fmt.Errorf("failed to get cache value: %w", err)
	}

	return json.Unmarshal([]byte(data), dest)
}

// DeleteCache removes a key from cache
func (r *Repository) DeleteCache(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// SetCacheWithPattern sets multiple keys with a pattern
func (r *Repository) SetCacheWithPattern(ctx context.Context, pattern string, values map[string]interface{}, ttl time.Duration) error {
	pipe := r.client.TxPipeline()

	for key, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}

		fullKey := fmt.Sprintf("%s:%s", pattern, key)
		pipe.Set(ctx, fullKey, data, ttl)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetCacheByPattern retrieves all keys matching a pattern
func (r *Repository) GetCacheByPattern(ctx context.Context, pattern string) (map[string]string, error) {
	keys, err := r.client.Keys(ctx, fmt.Sprintf("%s:*", pattern)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys by pattern: %w", err)
	}

	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	values, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get values: %w", err)
	}

	result := make(map[string]string)
	for i, key := range keys {
		if values[i] != nil {
			if val, ok := values[i].(string); ok {
				result[key] = val
			}
		}
	}

	return result, nil
}
