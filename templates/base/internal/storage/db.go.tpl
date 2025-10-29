package storage

import (
	"context"
	"fmt"
	"time"

	"{{.ProjectName}}/internal/config"
{{if eq .DbDriver.ID "gorm"}}
	"gorm.io/gorm"
{{if eq .Database.ID "postgres"}}	"gorm.io/driver/postgres"
{{else if eq .Database.ID "mysql"}}	"gorm.io/driver/mysql"
{{else if eq .Database.ID "sqlite"}}	"gorm.io/driver/sqlite"
{{end}}{{else if eq .DbDriver.ID "sqlx"}}
	"github.com/jmoiron/sqlx"
{{if eq .Database.ID "postgres"}}	_ "github.com/lib/pq"
{{else if eq .Database.ID "mysql"}}	_ "github.com/go-sql-driver/mysql"
{{else if eq .Database.ID "sqlite"}}	_ "github.com/mattn/go-sqlite3"
{{end}}{{else if eq .DbDriver.ID "mongo-driver"}}
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
{{else if eq .DbDriver.ID "redis-client"}}
	"github.com/redis/go-redis/v9"
{{end}})

{{if eq .DbDriver.ID "gorm"}}// InitDatabase initializes database connection using GORM
func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
{{if eq .Database.ID "postgres"}}	dsn := cfg.PostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else if eq .Database.ID "mysql"}}	dsn := cfg.MySQLDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else if eq .Database.ID "sqlite"}}	dsn := cfg.SQLiteDSN()
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else}}	dsn := cfg.PostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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

	// Test the connection
	if err := HealthCheck(db); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// HealthCheck checks GORM database connectivity
func HealthCheck(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

{{else if eq .DbDriver.ID "sqlx"}}// InitDatabase initializes database connection using sqlx
func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
{{if eq .Database.ID "postgres"}}	dsn := cfg.PostgresDSN()
	db, err := sqlx.Connect("postgres", dsn)
{{else if eq .Database.ID "mysql"}}	dsn := cfg.MySQLDSN()
	db, err := sqlx.Connect("mysql", dsn)
{{else if eq .Database.ID "sqlite"}}	dsn := cfg.SQLiteDSN()
	db, err := sqlx.Connect("sqlite3", dsn)
{{else}}	dsn := cfg.PostgresDSN()
	db, err := sqlx.Connect("postgres", dsn)
{{end}}	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := HealthCheck(db); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// HealthCheck checks sqlx database connectivity
func HealthCheck(db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

{{else if eq .DbDriver.ID "mongo-driver"}}// InitDatabase initializes database connection using MongoDB driver
func InitDatabase(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := cfg.MongoURI()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	if err := HealthCheck(client); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// HealthCheck checks MongoDB connectivity
func HealthCheck(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Ping(ctx, nil)
}

{{else if eq .DbDriver.ID "redis-client"}}// InitDatabase initializes database connection using Redis client
func InitDatabase(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL(),
		Password: cfg.Database.Password,
		DB:       cfg.Database.DB,
	})

	// Test the connection
	if err := HealthCheck(rdb); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}

// HealthCheck checks Redis connectivity
func HealthCheck(rdb *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	return err
}

{{else}}// InitDatabase initializes database connection
func InitDatabase(cfg *config.Config) (interface{}, error) {
	// Implement your custom database initialization logic here
	return nil, fmt.Errorf("custom database initialization not implemented")
}

// HealthCheck checks database connectivity
func HealthCheck(db interface{}) error {
	// Implement your custom health check logic here
	return fmt.Errorf("custom health check not implemented")
}
{{end}}
