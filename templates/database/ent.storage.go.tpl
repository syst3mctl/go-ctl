package {{.DbDriver.ID | replace "-" "_"}}

import (
	"context"
	"fmt"
	"time"

	"{{.ProjectName}}/internal/domain"

	"entgo.io/ent/dialect"
{{if eq .Database.ID "postgres"}}	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
{{else if eq .Database.ID "mysql"}}	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
{{else if eq .Database.ID "sqlite"}}	"entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
{{else}}	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
{{end}})

// NOTE: This template assumes you have generated ent code using:
// go generate ./ent
//
// Make sure to create your schema files in ent/schema/ directory
// Example schema file (ent/schema/user.go):
//
// package schema
//
// import (
// 	"entgo.io/ent"
// 	"entgo.io/ent/schema/field"
// 	"entgo.io/ent/schema/index"
// 	"time"
// )
//
// type User struct {
// 	ent.Schema
// }
//
// func (User) Fields() []ent.Field {
// 	return []ent.Field{
// 		field.String("name").NotEmpty(),
// 		field.String("email").Unique(),
// 		field.String("password").Sensitive(),
// 		field.Time("created_at").Default(time.Now),
// 		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
// 		field.Time("deleted_at").Optional().Nillable(),
// 	}
// }
//
// func (User) Indexes() []ent.Index {
// 	return []ent.Index{
// 		index.Fields("email").Unique(),
// 		index.Fields("deleted_at"),
// 	}
// }

// Storage implements the storage layer using ent
type Storage struct {
	client EntClient
}

// EntClient interface for dependency injection and testing
type EntClient interface {
	User() UserRepositoryInterface
	Close() error
	Schema() SchemaInterface
}

// UserRepositoryInterface defines user repository operations
type UserRepositoryInterface interface {
	Create() UserCreateInterface
	Query() UserQueryInterface
	Update() UserUpdateInterface
	Delete() UserDeleteInterface
	UpdateOne(*domain.User) UserUpdateOneInterface
	Get(ctx context.Context, id int) (*domain.User, error)
	GetX(ctx context.Context, id int) *domain.User
}

// Schema operations interface
type SchemaInterface interface {
	Create(ctx context.Context, opts ...interface{}) error
}

// User operation interfaces (these would be implemented by generated ent code)
type UserCreateInterface interface {
	SetName(string) UserCreateInterface
	SetEmail(string) UserCreateInterface
	SetPassword(string) UserCreateInterface
	SetCreatedAt(time.Time) UserCreateInterface
	SetUpdatedAt(time.Time) UserCreateInterface
	Save(ctx context.Context) (*domain.User, error)
	SaveX(ctx context.Context) *domain.User
}

type UserQueryInterface interface {
	Where(...interface{}) UserQueryInterface
	Select(...string) UserQueryInterface
	Order(...interface{}) UserQueryInterface
	Limit(int) UserQueryInterface
	Offset(int) UserQueryInterface
	All(ctx context.Context) ([]*domain.User, error)
	AllX(ctx context.Context) []*domain.User
	First(ctx context.Context) (*domain.User, error)
	FirstX(ctx context.Context) *domain.User
	Only(ctx context.Context) (*domain.User, error)
	OnlyX(ctx context.Context) *domain.User
	Count(ctx context.Context) (int, error)
	CountX(ctx context.Context) int
}

type UserUpdateInterface interface {
	Where(...interface{}) UserUpdateInterface
	SetName(string) UserUpdateInterface
	SetEmail(string) UserUpdateInterface
	SetPassword(string) UserUpdateInterface
	SetUpdatedAt(time.Time) UserUpdateInterface
	Save(ctx context.Context) (int, error)
	SaveX(ctx context.Context) int
}

type UserUpdateOneInterface interface {
	SetName(string) UserUpdateOneInterface
	SetEmail(string) UserUpdateOneInterface
	SetPassword(string) UserUpdateOneInterface
	SetUpdatedAt(time.Time) UserUpdateOneInterface
	Save(ctx context.Context) (*domain.User, error)
	SaveX(ctx context.Context) *domain.User
}

type UserDeleteInterface interface {
	Where(...interface{}) UserDeleteInterface
	Exec(ctx context.Context) (int, error)
	ExecX(ctx context.Context) int
}

// New creates a new ent storage instance
func New(client EntClient) *Storage {
	return &Storage{client: client}
}

// NewConnection creates a new database connection using ent
func NewConnection(dsn string) (EntClient, error) {
{{if eq .Database.ID "postgres"}}	drv, err := sql.Open(dialect.Postgres, dsn)
{{else if eq .Database.ID "mysql"}}	drv, err := sql.Open(dialect.MySQL, dsn)
{{else if eq .Database.ID "sqlite"}}	drv, err := sql.Open(dialect.SQLite, dsn)
{{else}}	drv, err := sql.Open(dialect.Postgres, dsn)
{{end}}	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db := drv.DB()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Return the generated ent client
	// NOTE: Replace this with your actual generated client
	// Example: return ent.NewClient(ent.Driver(drv)), nil
	return nil, fmt.Errorf("please replace this with your generated ent client")
}

// Migrate runs database migrations using ent schema
func (s *Storage) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Run auto migration
	return s.client.Schema().Create(ctx)
}

// Health checks database connectivity
func (s *Storage) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try a simple query to check connectivity
	_, err := s.client.User().Query().Count(ctx)
	return err
}

// User repository methods

// CreateUser creates a new user
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	createdUser, err := s.client.User().
		Create().
		SetName(user.Name).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetCreatedAt(user.CreatedAt).
		SetUpdatedAt(user.UpdatedAt).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Update the user ID from the created entity
	user.ID = uint(createdUser.ID)

	return nil
}

// GetUserByID retrieves a user by ID
func (s *Storage) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.client.User().Get(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.client.User().
		Query().
		Where(/* user.EmailEQ(email) */).  // Replace with generated predicate
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user
func (s *Storage) UpdateUser(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	_, err := s.client.User().
		UpdateOne(user).
		SetName(user.Name).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetUpdatedAt(user.UpdatedAt).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser soft deletes a user
func (s *Storage) DeleteUser(ctx context.Context, id uint) error {
	now := time.Now()

	affected, err := s.client.User().
		Update().
		Where(/* user.IDEQ(int(id)) */).  // Replace with generated predicate
		SetUpdatedAt(now).
		// SetDeletedAt(&now).  // Uncomment if using soft deletes
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// ListUsers retrieves users with pagination
func (s *Storage) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	users, err := s.client.User().
		Query().
		Where(/* user.DeletedAtIsNil() */).  // Replace with generated predicate for soft deletes
		Order(/* ent.Desc(user.FieldCreatedAt) */).  // Replace with generated field
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// CountUsers returns the total number of users
func (s *Storage) CountUsers(ctx context.Context) (int64, error) {
	count, err := s.client.User().
		Query().
		Where(/* user.DeletedAtIsNil() */).  // Replace with generated predicate
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return int64(count), nil
}

// SearchUsers searches users by name or email
func (s *Storage) SearchUsers(ctx context.Context, query string, limit, offset int) ([]*domain.User, error) {
	searchPattern := "%" + query + "%"

	users, err := s.client.User().
		Query().
		Where(
			/* user.Or(
				user.NameContains(searchPattern),
				user.EmailContains(searchPattern),
			) */
		).  // Replace with generated predicates
		Where(/* user.DeletedAtIsNil() */).  // Replace with generated predicate
		Order(/* ent.Desc(user.FieldCreatedAt) */).  // Replace with generated field
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	return users, nil
}

// Transaction wraps operations in a database transaction
func (s *Storage) Transaction(ctx context.Context, fn func(tx EntClient) error) error {
	// NOTE: Replace with actual ent transaction implementation
	// tx, err := s.client.Tx(ctx)
	// if err != nil {
	//     return err
	// }
	//
	// defer func() {
	//     if r := recover(); r != nil {
	//         tx.Rollback()
	//         panic(r)
	//     }
	// }()
	//
	// if err := fn(tx.Client()); err != nil {
	//     if rerr := tx.Rollback(); rerr != nil {
	//         return fmt.Errorf("rolling back transaction: %v (original error: %w)", rerr, err)
	//     }
	//     return err
	// }
	//
	// if err := tx.Commit(); err != nil {
	//     return fmt.Errorf("committing transaction: %w", err)
	// }

	return fmt.Errorf("transaction implementation depends on generated ent code")
}

// Close closes the database connection
func (s *Storage) Close() error {
	return s.client.Close()
}

// Advanced ent operations (examples)

// GetUsersWithEdges retrieves users with their related entities
func (s *Storage) GetUsersWithEdges(ctx context.Context, limit int) ([]*domain.User, error) {
	users, err := s.client.User().
		Query().
		// WithPosts().  // Uncomment and replace with actual edges
		// WithProfile().
		Limit(limit).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get users with edges: %w", err)
	}

	return users, nil
}

// BulkCreateUsers creates multiple users in a batch
func (s *Storage) BulkCreateUsers(ctx context.Context, users []*domain.User) error {
	if len(users) == 0 {
		return nil
	}

	bulk := make([]UserCreateInterface, len(users))
	now := time.Now()

	for i, user := range users {
		user.CreatedAt = now
		user.UpdatedAt = now

		bulk[i] = s.client.User().
			Create().
			SetName(user.Name).
			SetEmail(user.Email).
			SetPassword(user.Password).
			SetCreatedAt(user.CreatedAt).
			SetUpdatedAt(user.UpdatedAt)
	}

	// NOTE: Replace with actual bulk create implementation
	// _, err := s.client.User().CreateBulk(bulk...).Save(ctx)
	// if err != nil {
	//     return fmt.Errorf("failed to bulk create users: %w", err)
	// }

	return fmt.Errorf("bulk create implementation depends on generated ent code")
}

// GetUsersByIDs retrieves multiple users by their IDs
func (s *Storage) GetUsersByIDs(ctx context.Context, ids []uint) ([]*domain.User, error) {
	if len(ids) == 0 {
		return []*domain.User{}, nil
	}

	intIDs := make([]int, len(ids))
	for i, id := range ids {
		intIDs[i] = int(id)
	}

	users, err := s.client.User().
		Query().
		Where(/* user.IDIn(intIDs...) */).  // Replace with generated predicate
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}

	return users, nil
}

// TODO: Add more repository methods as needed for your application

/*
SETUP INSTRUCTIONS:

1. Install ent:
   go get entgo.io/ent/cmd/ent

2. Initialize ent in your project:
   go run entgo.io/ent/cmd/ent init User

3. Edit the generated schema file (ent/schema/user.go) to match your domain model

4. Generate ent code:
   go generate ./ent

5. Replace the interface implementations with actual generated ent client calls

6. Update import paths to use your generated ent client

Example complete setup:

// ent/generate.go
//go:generate go run entgo.io/ent/cmd/ent generate ./schema

// ent/schema/user.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("email").Unique(),
		field.String("password").Sensitive(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("deleted_at"),
	}
}

Then run:
go generate ./ent

This will generate the actual client code that you can then use in this storage implementation.
*/
