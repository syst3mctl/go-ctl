# Database Driver Configuration Guide

This guide provides comprehensive documentation for configuring database drivers in the go-ctl project generator. Each database driver has specific setup requirements, connection patterns, and usage examples.

## Overview

The go-ctl initializer supports multiple database drivers and databases, allowing you to choose the best combination for your project needs:

**Supported Databases:**
- PostgreSQL
- MySQL
- SQLite
- MongoDB
- Redis
- BigQuery

**Supported Database Drivers/ORMs:**
- `database/sql` - Standard library SQL interface
- `GORM` - Feature-rich ORM with associations and migrations
- `sqlx` - Extensions to database/sql with easier scanning
- `Ent` - Schema-first entity framework with code generation
- `MongoDB Driver` - Official MongoDB Go driver
- `Redis Client` - Redis client with Cluster/Sentinel support

## Driver-Database Compatibility Matrix

| Driver/ORM | PostgreSQL | MySQL | SQLite | MongoDB | Redis | BigQuery |
|------------|------------|-------|--------|---------|-------|----------|
| database/sql | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ |
| GORM | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| sqlx | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ |
| Ent | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| MongoDB Driver | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ |
| Redis Client | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |

## 1. Standard Library (database/sql)

The built-in SQL interface for Go with connection pooling and transaction support.

### Dependencies
```json
{
  "postgres": ["github.com/lib/pq"],
  "mysql": ["github.com/go-sql-driver/mysql"],
  "sqlite": ["github.com/mattn/go-sqlite3"]
}
```

### Connection Examples

#### PostgreSQL
```go
dsn := "postgres://username:password@localhost/dbname?sslmode=disable"
db, err := sql.Open("postgres", dsn)
```

#### MySQL
```go
dsn := "username:password@tcp(localhost:3306)/dbname?parseTime=true"
db, err := sql.Open("mysql", dsn)
```

#### SQLite
```go
dsn := "./database.db"
db, err := sql.Open("sqlite3", dsn)
```

### Features
- Connection pooling with configurable limits
- Manual schema migrations
- Raw SQL queries with parameter binding
- Transaction support
- Context-aware operations
- Prepared statements

### Best For
- Projects requiring fine-grained SQL control
- High-performance applications
- Legacy database schemas
- Custom query optimization

## 2. GORM

A feature-rich ORM library with automatic migrations, associations, and hooks.

### Dependencies
```json
{
  "postgres": ["gorm.io/driver/postgres"],
  "mysql": ["gorm.io/driver/mysql"],
  "sqlite": ["gorm.io/driver/sqlite"]
}
```

### Connection Examples

#### PostgreSQL
```go
dsn := "host=localhost user=username password=password dbname=dbname port=5432 sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

#### MySQL
```go
dsn := "username:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

### Features
- Auto-migrations with `db.AutoMigrate(&Model{})`
- Associations (belongs to, has one, has many, many to many)
- Hooks (before/after create, update, delete)
- Soft deletes with `gorm.DeletedAt`
- Preloading related data
- Scopes for reusable query logic
- Plugin system

### Best For
- Rapid application development
- Projects with complex relationships
- Teams preferring ORM over raw SQL
- Applications requiring auto-migrations

## 3. sqlx

Extensions to the standard database/sql with struct scanning and named parameters.

### Dependencies
```json
{
  "postgres": ["github.com/lib/pq"],
  "mysql": ["github.com/go-sql-driver/mysql"],
  "sqlite": ["github.com/mattn/go-sqlite3"]
}
```

### Connection Examples
```go
db, err := sqlx.Connect("postgres", dsn)
```

### Features
- Struct scanning with `db` tags
- Named parameter binding (`:name`, `:email`)
- `Get()` for single row queries
- `Select()` for multiple row queries
- `NamedExec()` and `NamedQuery()` for named parameters
- `sqlx.In()` for IN clause queries
- Marshal/Unmarshal support

### Best For
- Projects wanting SQL control with convenience
- Teams familiar with database/sql
- Applications requiring complex queries
- Performance-critical applications

## 4. Ent

Schema-first entity framework with type-safe, code-generated APIs.

### Dependencies
```json
{
  "postgres": ["entgo.io/ent/dialect/sql", "github.com/lib/pq"],
  "mysql": ["entgo.io/ent/dialect/sql", "github.com/go-sql-driver/mysql"],
  "sqlite": ["entgo.io/ent/dialect/sql", "github.com/mattn/go-sqlite3"]
}
```

### Setup Process
1. **Install Ent CLI:**
   ```bash
   go get entgo.io/ent/cmd/ent
   ```

2. **Initialize Schema:**
   ```bash
   go run entgo.io/ent/cmd/ent init User
   ```

3. **Define Schema (ent/schema/user.go):**
   ```go
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
   ```

4. **Generate Code:**
   ```bash
   go generate ./ent
   ```

### Features
- Schema-first development with type safety
- Code generation for CRUD operations
- Graph-based queries and mutations
- Automatic migrations
- Edge (relationship) support
- Hooks and interceptors
- Privacy policies
- Custom predicates and modifiers

### Best For
- Large-scale applications
- Teams requiring type safety
- Complex data relationships
- GraphQL backends
- Projects with evolving schemas

## 5. MongoDB Driver

Official MongoDB driver for document-based operations.

### Dependencies
```json
{
  "mongodb": [
    "go.mongodb.org/mongo-driver/mongo",
    "go.mongodb.org/mongo-driver/bson",
    "go.mongodb.org/mongo-driver/mongo/options"
  ]
}
```

### Connection Examples
```go
uri := "mongodb://username:password@localhost:27017"
client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
database := client.Database("myapp")
collection := database.Collection("users")
```

### Features
- BSON document operations
- Aggregation pipelines
- GridFS for file storage
- Change streams
- Transactions (replica sets/sharded clusters)
- Connection pooling
- Automatic failover

### Best For
- Document-oriented applications
- Flexible, evolving schemas
- Real-time applications
- Content management systems
- Rapid prototyping

## 6. Redis Client

High-performance in-memory data store client.

### Dependencies
```json
{
  "redis": ["github.com/redis/go-redis/v9"]
}
```

### Connection Examples
```go
rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password
    DB:       0,  // default DB
})
```

### Features
- Key-value operations
- Data structures (strings, hashes, lists, sets, sorted sets)
- Pub/Sub messaging
- Lua scripting
- Pipelines for batch operations
- Cluster support
- Sentinel for high availability

### Best For
- Caching layers
- Session storage
- Real-time applications
- Message queues
- Rate limiting
- Analytics and counters

## Configuration Examples

### Environment Variables
```env
# PostgreSQL
DATABASE_URL=postgres://user:pass@localhost/dbname?sslmode=disable

# MySQL
DATABASE_URL=user:pass@tcp(localhost:3306)/dbname?parseTime=true

# SQLite
DATABASE_URL=./app.db

# MongoDB
MONGODB_URI=mongodb://localhost:27017

# Redis
REDIS_URL=redis://localhost:6379/0
```

### Connection Pool Settings
```go
// For SQL databases
db.SetMaxOpenConns(25)        // Maximum open connections
db.SetMaxIdleConns(5)         // Maximum idle connections
db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime

// For MongoDB
clientOptions.SetMaxPoolSize(100)
clientOptions.SetMinPoolSize(5)
clientOptions.SetMaxConnIdleTime(30 * time.Second)

// For Redis
&redis.Options{
    PoolSize:     10,
    PoolTimeout:  30 * time.Second,
    IdleTimeout:  time.Minute,
}
```

## Migration Strategies

### GORM Auto-Migration
```go
db.AutoMigrate(&User{}, &Post{}, &Comment{})
```

### Manual SQL Migrations (database/sql, sqlx)
```go
func (s *Storage) Migrate() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS users (...)`,
        `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
    }
    
    for _, query := range queries {
        if _, err := s.db.ExecContext(ctx, query); err != nil {
            return err
        }
    }
    return nil
}
```

### Ent Schema Migration
```go
if err := client.Schema.Create(ctx); err != nil {
    return fmt.Errorf("failed creating schema resources: %v", err)
}
```

## Performance Considerations

### Query Optimization
- Use prepared statements for repeated queries
- Implement proper indexing strategies
- Use connection pooling effectively
- Consider read replicas for heavy read workloads

### Monitoring and Observability
- Track connection pool metrics
- Monitor slow queries
- Implement proper logging
- Use database-specific monitoring tools

### Caching Strategies
- Use Redis for frequently accessed data
- Implement cache-aside pattern
- Consider write-through/write-behind caching
- Set appropriate TTL values

## Security Best Practices

### Connection Security
- Use SSL/TLS for database connections
- Store credentials in environment variables
- Rotate database passwords regularly
- Use connection string encryption

### Query Security
- Always use parameterized queries
- Validate and sanitize input data
- Implement proper access controls
- Use least privilege principle

### Data Protection
- Encrypt sensitive data at rest
- Hash passwords properly (bcrypt)
- Implement audit logging
- Use database-level encryption features

## Troubleshooting Common Issues

### Connection Problems
- Check network connectivity
- Verify credentials and permissions
- Ensure database server is running
- Review firewall and security group settings

### Performance Issues
- Analyze query execution plans
- Check for missing indexes
- Monitor connection pool usage
- Review application-level caching

### Migration Failures
- Backup database before migrations
- Test migrations in development first
- Use transactional migrations when possible
- Implement rollback strategies

## Choosing the Right Driver

### Decision Matrix

**Choose `database/sql` when:**
- Maximum performance is critical
- Complex queries with fine-grained control needed
- Working with legacy schemas
- Team has strong SQL expertise

**Choose `GORM` when:**
- Rapid development is priority
- Complex relationships between entities
- Team prefers ORM over raw SQL
- Auto-migrations are beneficial

**Choose `sqlx` when:**
- Want SQL control with convenience features
- Need better struct scanning than database/sql
- Performance is important but some convenience is desired
- Team is familiar with database/sql patterns

**Choose `Ent` when:**
- Type safety is paramount
- Large-scale application with evolving schema
- GraphQL integration needed
- Team values code generation benefits

**Choose `MongoDB Driver` when:**
- Document-oriented data model fits use case
- Schema flexibility is important
- Horizontal scaling requirements
- Working with unstructured data

**Choose `Redis Client` when:**
- Need high-performance caching
- Real-time features required
- Session management needed
- Working with simple key-value data

This guide should help you choose the appropriate database driver combination for your Go project using the go-ctl initializer.