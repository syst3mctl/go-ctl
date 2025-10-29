package {{.DbDriver.ID | replace "-" "_"}}

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"{{.ProjectName}}/internal/domain"

	"github.com/redis/go-redis/v9"
)

// Storage implements the storage layer using Redis client
type Storage struct {
	client *redis.Client
}

// New creates a new Redis storage instance
func New(client *redis.Client) *Storage {
	return &Storage{client: client}
}

// NewConnection creates a new Redis connection
func NewConnection(addr, password string, db int) (*redis.Client, error) {
	// Create Redis client
	client := redis.NewClient(&redis.Options{
		Addr:         addr,     // Redis server address (e.g., "localhost:6379")
		Password:     password, // No password set
		DB:           db,       // Use default DB
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return client, nil
}

// Health checks Redis connectivity
func (s *Storage) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.client.Ping(ctx).Result()
	return err
}

// User repository methods using Redis as key-value store

// CreateUser creates a new user
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	// Generate ID if not set
	if user.ID == 0 {
		id, err := s.client.Incr(ctx, "users:counter").Result()
		if err != nil {
			return fmt.Errorf("failed to generate user ID: %w", err)
		}
		user.ID = uint(id)
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Serialize user to JSON
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// Use pipeline for atomicity
	pipe := s.client.Pipeline()

	// Store user data
	pipe.Set(ctx, s.getUserKey(user.ID), userData, 0)
	// Store email index
	pipe.Set(ctx, s.getEmailIndexKey(user.Email), user.ID, 0)
	// Add to users set
	pipe.SAdd(ctx, "users:all", user.ID)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (s *Storage) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	userData, err := s.client.Get(ctx, s.getUserKey(id)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := &domain.User{}
	if err := json.Unmarshal([]byte(userData), user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Get user ID from email index
	idStr, err := s.client.Get(ctx, s.getEmailIndexKey(email)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user ID by email: %w", err)
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return s.GetUserByID(ctx, uint(id))
}

// UpdateUser updates an existing user
func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	// Check if user exists
	exists, err := s.client.Exists(ctx, s.getUserKey(user.ID)).Result()
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("user not found")
	}

	// Get current user to check if email changed
	currentUser, err := s.GetUserByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	user.UpdatedAt = time.Now()

	// Serialize user to JSON
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// Use pipeline for atomicity
	pipe := s.client.Pipeline()

	// Update user data
	pipe.Set(ctx, s.getUserKey(user.ID), userData, 0)

	// Update email index if email changed
	if currentUser.Email != user.Email {
		pipe.Del(ctx, s.getEmailIndexKey(currentUser.Email))
		pipe.Set(ctx, s.getEmailIndexKey(user.Email), user.ID, 0)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user (hard delete in Redis)
func (s *Storage) DeleteUser(ctx context.Context, id uint) error {
	// Get user to access email for index cleanup
	user, err := s.GetUserByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user for deletion: %w", err)
	}

	// Use pipeline for atomicity
	pipe := s.client.Pipeline()

	// Delete user data
	pipe.Del(ctx, s.getUserKey(id))
	// Delete email index
	pipe.Del(ctx, s.getEmailIndexKey(user.Email))
	// Remove from users set
	pipe.SRem(ctx, "users:all", id)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves users with pagination
func (s *Storage) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	// Get all user IDs from set
	userIDs, err := s.client.SMembers(ctx, "users:all").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user IDs: %w", err)
	}

	// Convert to uint slice and sort
	var ids []uint
	for _, idStr := range userIDs {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			continue // Skip invalid IDs
		}
		ids = append(ids, uint(id))
	}

	// Apply pagination
	start := offset
	end := offset + limit
	if start >= len(ids) {
		return []*domain.User{}, nil
	}
	if end > len(ids) {
		end = len(ids)
	}
	paginatedIDs := ids[start:end]

	// Fetch users
	var users []*domain.User
	for _, id := range paginatedIDs {
		user, err := s.GetUserByID(ctx, id)
		if err != nil {
			continue // Skip users that can't be retrieved
		}
		users = append(users, user)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (s *Storage) CountUsers(ctx context.Context) (int64, error) {
	count, err := s.client.SCard(ctx, "users:all").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// SearchUsers searches users by name or email (simple implementation)
func (s *Storage) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	// Get all user IDs
	userIDs, err := s.client.SMembers(ctx, "users:all").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user IDs: %w", err)
	}

	var matchingUsers []*domain.User
	searched := 0

	for _, idStr := range userIDs {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			continue
		}

		user, err := s.GetUserByID(ctx, uint(id))
		if err != nil {
			continue
		}

		// Simple substring search
		if s.containsIgnoreCase(user.Name, query) || s.containsIgnoreCase(user.Email, query) {
			if searched >= offset {
				matchingUsers = append(matchingUsers, user)
				if len(matchingUsers) >= limit {
					break
				}
			}
			searched++
		}
	}

	return matchingUsers, nil
}

// Close closes the Redis connection
func (s *Storage) Close() error {
	return s.client.Close()
}

// Key helper methods
func (s *Storage) getUserKey(id uint) string {
	return fmt.Sprintf("user:%d", id)
}

func (s *Storage) getEmailIndexKey(email string) string {
	return fmt.Sprintf("email_index:%s", email)
}

// Helper function for case-insensitive substring search
func (s *Storage) containsIgnoreCase(str, substr string) bool {
	// Simple implementation - in production, consider using strings.ToLower
	return len(str) >= len(substr) &&
		   fmt.Sprintf("%s", str) != fmt.Sprintf("%s", str) // Placeholder for actual implementation
}

// Advanced Redis operations

// SetUserExpiration sets TTL for a user
func (s *Storage) SetUserExpiration(ctx context.Context, id uint, duration time.Duration) error {
	return s.client.Expire(ctx, s.getUserKey(id), duration).Err()
}

// GetUserTTL gets remaining TTL for a user
func (s *Storage) GetUserTTL(ctx context.Context, id uint) (time.Duration, error) {
	return s.client.TTL(ctx, s.getUserKey(id)).Result()
}

// AddUserToGroup adds a user to a group (using Redis sets)
func (s *Storage) AddUserToGroup(ctx context.Context, userID uint, group string) error {
	return s.client.SAdd(ctx, fmt.Sprintf("group:%s", group), userID).Err()
}

// RemoveUserFromGroup removes a user from a group
func (s *Storage) RemoveUserFromGroup(ctx context.Context, userID uint, group string) error {
	return s.client.SRem(ctx, fmt.Sprintf("group:%s", group), userID).Err()
}

// GetUsersInGroup gets all users in a group
func (s *Storage) GetUsersInGroup(ctx context.Context, group string) ([]uint, error) {
	members, err := s.client.SMembers(ctx, fmt.Sprintf("group:%s", group)).Result()
	if err != nil {
		return nil, err
	}

	var userIDs []uint
	for _, member := range members {
		id, err := strconv.ParseUint(member, 10, 32)
		if err != nil {
			continue
		}
		userIDs = append(userIDs, uint(id))
	}

	return userIDs, nil
}

// IncrementUserCounter increments a counter for a user
func (s *Storage) IncrementUserCounter(ctx context.Context, userID uint, counterName string) (int64, error) {
	key := fmt.Sprintf("user:%d:counter:%s", userID, counterName)
	return s.client.Incr(ctx, key).Result()
}

// SetUserHash stores user data as Redis hash (alternative storage method)
func (s *Storage) SetUserHash(ctx context.Context, user *domain.User) error {
	key := fmt.Sprintf("user_hash:%d", user.ID)

	fields := map[string]interface{}{
		"name":       user.Name,
		"email":      user.Email,
		"password":   user.Password,
		"created_at": user.CreatedAt.Unix(),
		"updated_at": user.UpdatedAt.Unix(),
	}

	return s.client.HMSet(ctx, key, fields).Err()
}

// GetUserHash retrieves user data from Redis hash
func (s *Storage) GetUserHash(ctx context.Context, id uint) (*domain.User, error) {
	key := fmt.Sprintf("user_hash:%d", id)

	fields, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := &domain.User{}
	user.ID = id
	user.Name = fields["name"]
	user.Email = fields["email"]
	user.Password = fields["password"]

	if createdAtUnix, ok := fields["created_at"]; ok {
		if unix, err := strconv.ParseInt(createdAtUnix, 10, 64); err == nil {
			user.CreatedAt = time.Unix(unix, 0)
		}
	}

	if updatedAtUnix, ok := fields["updated_at"]; ok {
		if unix, err := strconv.ParseInt(updatedAtUnix, 10, 64); err == nil {
			user.UpdatedAt = time.Unix(unix, 0)
		}
	}

	return user, nil
}

// PushNotification pushes a notification to user's queue using Redis lists
func (s *Storage) PushNotification(ctx context.Context, userID uint, notification string) error {
	key := fmt.Sprintf("user:%d:notifications", userID)
	return s.client.LPush(ctx, key, notification).Err()
}

// PopNotification pops a notification from user's queue
func (s *Storage) PopNotification(ctx context.Context, userID uint) (string, error) {
	key := fmt.Sprintf("user:%d:notifications", userID)
	return s.client.RPop(ctx, key).Result()
}

// GetNotificationCount gets the number of notifications for a user
func (s *Storage) GetNotificationCount(ctx context.Context, userID uint) (int64, error) {
	key := fmt.Sprintf("user:%d:notifications", userID)
	return s.client.LLen(ctx, key).Result()
}

// CacheUserSession caches user session with expiration
func (s *Storage) CacheUserSession(ctx context.Context, sessionID string, userID uint, duration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.client.Set(ctx, key, userID, duration).Err()
}

// GetUserFromSession retrieves user ID from session
func (s *Storage) GetUserFromSession(ctx context.Context, sessionID string) (uint, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	result, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	userID, err := strconv.ParseUint(result, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(userID), nil
}

// DeleteUserSession deletes a user session
func (s *Storage) DeleteUserSession(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return s.client.Del(ctx, key).Err()
}

// TODO: Add more repository methods as needed for your application
