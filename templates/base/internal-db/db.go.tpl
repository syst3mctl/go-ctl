package db

import (
	"context"
	"time"

{{if eq .DbDriver.ID "gorm"}}	"fmt"
	"gorm.io/gorm"
{{if eq .Database.ID "postgres"}}	"gorm.io/driver/postgres"
{{else if eq .Database.ID "mysql"}}	"gorm.io/driver/mysql"
{{else if eq .Database.ID "sqlite"}}	"gorm.io/driver/sqlite"
{{end}}
{{else if eq .DbDriver.ID "sqlx"}}	"database/sql"
	"fmt"
{{if eq .Database.ID "postgres"}}	_ "github.com/lib/pq"
{{else if eq .Database.ID "mysql"}}	_ "github.com/go-sql-driver/mysql"
{{else if eq .Database.ID "sqlite"}}	_ "github.com/mattn/go-sqlite3"
{{end}}	"github.com/jmoiron/sqlx"
{{else if eq .DbDriver.ID "ent"}}	"fmt"
{{if eq .Database.ID "postgres"}}	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
{{else if eq .Database.ID "mysql"}}	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
{{else if eq .Database.ID "sqlite"}}	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
{{end}}	"entgo.io/ent/dialect"
{{else}}	"database/sql"
{{if eq .Database.ID "postgres"}}	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
{{else if eq .Database.ID "mysql"}}	_ "github.com/go-sql-driver/mysql"
{{else if eq .Database.ID "sqlite"}}	_ "github.com/mattn/go-sqlite3"
{{end}}
{{end}})

{{if eq .DbDriver.ID "gorm"}}// NewConn creates a new GORM database connection
func NewConn(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

{{if eq .Database.ID "postgres"}}	db, err = gorm.Open(postgres.Open(addr), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else if eq .Database.ID "mysql"}}	db, err = gorm.Open(mysql.Open(addr), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else if eq .Database.ID "sqlite"}}	db, err = gorm.Open(sqlite.Open(addr), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
{{else}}	db, err = gorm.Open(postgres.Open(addr), &gorm.Config{
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

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), maxIdleTime)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

{{else if eq .DbDriver.ID "sqlx"}}// NewConn creates a new sqlx database connection
func NewConn(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

{{if eq .Database.ID "postgres"}}	db, err = sqlx.Connect("postgres", addr)
{{else if eq .Database.ID "mysql"}}	db, err = sqlx.Connect("mysql", addr)
{{else if eq .Database.ID "sqlite"}}	db, err = sqlx.Connect("sqlite3", addr)
{{else}}	db, err = sqlx.Connect("postgres", addr)
{{end}}	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), maxIdleTime)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

{{else if eq .DbDriver.ID "ent"}}// NewConn creates a new ent database connection
func NewConn(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*entsql.DB, error) {
	var drv *entsql.Driver
	var err error

{{if eq .Database.ID "postgres"}}	drv, err = entsql.Open(dialect.Postgres, addr)
{{else if eq .Database.ID "mysql"}}	drv, err = entsql.Open(dialect.MySQL, addr)
{{else if eq .Database.ID "sqlite"}}	drv, err = entsql.Open(dialect.SQLite, addr)
{{else}}	drv, err = entsql.Open(dialect.Postgres, addr)
{{end}}	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := drv.DB()
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), maxIdleTime)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

{{else}}// NewConn creates a new database/sql connection
func NewConn(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {
{{if eq .Database.ID "postgres"}}	db, err := sql.Open("postgres", addr)
{{else if eq .Database.ID "mysql"}}	db, err := sql.Open("mysql", addr)
{{else if eq .Database.ID "sqlite"}}	db, err := sql.Open("sqlite3", addr)
{{else}}	db, err := sql.Open("postgres", addr)
{{end}}	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), maxIdleTime)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
{{end}}

