package storage

import (
	"context"
	"fmt"
	"time"

	"test-app/internal/config"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

// InitDatabase initializes database connection using GORM
func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.PostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
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


