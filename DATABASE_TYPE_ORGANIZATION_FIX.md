# Database Type Organization Fix

This document explains the critical fix applied to organize storage templates by **database type** instead of **driver library**.

## Problem Identified

The user pointed out that storage should be organized by database type (postgres, mysql, sqlite, mongodb, redis) rather than by Go driver library (gorm, sqlx, mongo-driver, redis-client).

## Previous Structure (Incorrect)
```
/storage/
├── db.go (unified initialization)
├── gorm/
│   └── repository.go
├── sqlx/
│   └── repository.go
├── mongo-driver/
│   └── repository.go
└── redis-client/
    └── repository.go
```

## New Structure (Correct)
```
/storage/
├── db.go (unified initialization)
├── postgres/
│   └── repository.go (can use GORM or sqlx)
├── mysql/
│   └── repository.go (can use GORM or sqlx)
├── sqlite/
│   └── repository.go (can use GORM or sqlx)
├── mongodb/
│   └── repository.go (uses mongo-driver)
└── redis/
    └── repository.go (uses redis-client)
```

## Key Changes Made

### 1. **Restructured Template Directories**
- **Deleted**: `templates/storage/gorm/`, `templates/storage/sqlx/`, `templates/storage/mongo-driver/`, `templates/storage/redis-client/`
- **Created**: `templates/storage/postgres/`, `templates/storage/mysql/`, `templates/storage/sqlite/`, `templates/storage/mongodb/`, `templates/storage/redis/`

### 2. **Updated Import Paths**
- **main.go.tpl**: Changed from `{{.DbDriver.ID}}` to `{{.Database.ID}}`
- **Before**: `import "myproject/internal/storage/gorm"`
- **After**: `import "myproject/internal/storage/postgres"`

### 3. **Database-Specific Repository Templates**

#### PostgreSQL Repository (`storage/postgres/repository.go.tpl`)
- Supports both GORM and sqlx drivers
- Uses PostgreSQL-specific SQL features (ILIKE, SERIAL, etc.)
- Contains PostgreSQL-optimized queries

#### MySQL Repository (`storage/mysql/repository.go.tpl`)
- Supports both GORM and sqlx drivers  
- Uses MySQL-specific SQL features (LIKE, AUTO_INCREMENT, etc.)
- Contains MySQL-optimized queries

#### SQLite Repository (`storage/sqlite/repository.go.tpl`)
- Supports both GORM and sqlx drivers
- Uses SQLite-specific SQL features (AUTOINCREMENT, INTEGER PRIMARY KEY, etc.)
- Contains SQLite-optimized queries

#### MongoDB Repository (`storage/mongodb/repository.go.tpl`)
- Uses mongo-driver internally
- Contains MongoDB-specific operations (aggregation, indexes, etc.)
- Handles ObjectIDs and BSON operations

#### Redis Repository (`storage/redis/repository.go.tpl`)
- Uses redis-client internally
- Contains Redis-specific operations (SET, GET, HSET, etc.)
- Handles caching and key-value operations

## Benefits of This Organization

### 1. **Logical Grouping**
- Developers think in terms of "I'm using PostgreSQL" not "I'm using GORM"
- Database choice is the primary concern, driver is implementation detail

### 2. **Driver Flexibility**
- PostgreSQL repository can use either GORM or sqlx
- MySQL repository can use either GORM or sqlx
- SQLite repository can use either GORM or sqlx
- Driver selection becomes an internal implementation choice

### 3. **Database-Specific Optimizations**
- Each repository can use database-specific features
- PostgreSQL can use `ILIKE`, MySQL uses `LIKE`, SQLite uses different data types
- MongoDB uses aggregation pipelines, Redis uses key-value operations

### 4. **Easier Maintenance**
- Adding a new database type (e.g., CockroachDB) means creating `/storage/cockroachdb/`
- Adding a new driver for existing database doesn't require new top-level directory
- Clear separation between "what database" and "how to connect"

## Generated Project Impact

### Main.go Pattern (Consistent)
```go
// Always imports database-specific repository
import "myproject/internal/storage/postgres"  // not "gorm" or "sqlx"

func main() {
    db, err := storage.InitDatabase(cfg)
    repo := postgres.New(db)  // database-specific constructor
    // ... rest is identical
}
```

### Repository Usage (Clear Intent)
```go
// Clear what database we're targeting
repo := postgres.New(db)    // PostgreSQL with GORM or sqlx
repo := mysql.New(db)       // MySQL with GORM or sqlx  
repo := sqlite.New(db)      // SQLite with GORM or sqlx
repo := mongodb.New(client) // MongoDB with mongo-driver
repo := redis.New(client)   // Redis with redis-client
```

## Standard Compliance

This fix ensures the project follows the user's specified standard:

> "when database driver is selected always initialize any driver in /storage/db.go and repository functions is written under /storage/postgres(for example) <- under this package"

✅ **Database initialization**: Always in `/storage/db.go`
✅ **Repository functions**: Under `/storage/postgres`, `/storage/mysql`, etc.
✅ **Database-type organization**: Not driver-library organization

This makes the generated projects more intuitive and maintainable for developers who think in terms of database systems rather than Go driver libraries.