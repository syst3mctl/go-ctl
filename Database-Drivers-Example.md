# Go (Golang) Database & Cache Examples

This file provides standard "getting started" examples for common database drivers and clients in Go. Each example shows the necessary imports, the connection/setup code, and a basic operation (like a ping, create, or read).

---

## 1. Standard Library (`database/sql`)

This is Go's built-in interface for SQL databases. It requires a specific **driver** for the database you want to connect to (e.g., MySQL, PostgreSQL, SQLite).

**Install (using MySQL driver):**
```bash
go get [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)


1.  build int database/sql package

    ```go
        package main

        import (
       	"database/sql"
       	"fmt"
       	"log"
       	"time"

       	// The blank import is required to register the driver
       	_ "[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)"
        )

        func main() {
       	// DSN: Data Source Name
       	// Format: "username:password@protocol(address)/dbname?param=value"
       	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?parseTime=true"

       	// sql.Open just validates its arguments, it doesn't create a connection.
       	db, err := sql.Open("mysql", dsn)
       	if err != nil {
      		log.Fatalf("Error opening database: %v", err)
       	}
       	// It's important to close the database when the application exits.
       	defer db.Close()

       	// Set connection pool settings (optional but recommended)
       	db.SetConnMaxLifetime(time.Minute * 3)
       	db.SetMaxOpenConns(10)
       	db.SetMaxIdleConns(10)

       	// db.Ping verifies the connection is alive and valid.
       	err = db.Ping()
       	if err != nil {
      		log.Fatalf("Error connecting to database: %v", err)
       	}

       	fmt.Println("Successfully connected to MySQL!")

       	// Example: Perform a simple query
       	var username string
       	// QueryRow is used for queries that are expected to return at most one row.
       	err = db.QueryRow("SELECT username FROM users WHERE id = ?", 1).Scan(&username)
       	if err != nil {
      		if err == sql.ErrNoRows {
     			fmt.Println("No user found with that ID.")
      		} else {
     			log.Printf("Query error: %v", err)
      		}
       	} else {
      		fmt.Printf("Username for ID 1 is: %s\n", username)
       	}
        }
    ```


2. GORM
GORM is a full-featured ORM (Object Relational Mapper) for Go. It simplifies database interactions by mapping Go structs to database tables.

Install (using MySQL driver):

 ```go
	 package main

	 import (
		"fmt"
		"log"

		"gorm.io/driver/mysql"
		"gorm.io/gorm"
	 )

	 // Product model maps to the 'products' table
	 type Product struct {
		gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt
		Code       string
		Price      uint
	 }

	 func main() {
		dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

		// Open connection
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to database")
		}

		fmt.Println("Successfully connected with GORM!")

		// AutoMigrate will create/update the 'products' table based on the Product struct.
		db.AutoMigrate(&Product{})

		// --- Basic Operations ---

		// Create
		fmt.Println("Creating product...")
		product := Product{Code: "D42", Price: 100}
		result := db.Create(&product) // product.ID will be set after creation
		if result.Error != nil {
			log.Printf("Failed to create product: %v", result.Error)
		}
		fmt.Printf("Created product with ID: %d\n", product.ID)

		// Read
		var readProduct Product
		// Find product with ID 1
		db.First(&readProduct, 1)
		fmt.Printf("Read product (ID 1): %s\n", readProduct.Code)

		// Find product with code "D42"
		db.First(&readProduct, "code = ?", "D42")
		fmt.Printf("Read product (Code D42): %s - Price: %d\n", readProduct.Code, readProduct.Price)

		// Update
		db.Model(&readProduct).Update("Price", 200)
		fmt.Printf("Updated product price to: %d\n", readProduct.Price)
	 }
 ```

3. sqlx
sqlx is an extension to the standard database/sql library. Its main features are powerful, clean scanning of query results into structs and maps, and named query support.

Install (using MySQL driver):
```go
package main

import (
	"fmt"
	"log"

	_ "[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)"
	"[github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx)"
)

// User struct with `db` tags to map database columns to fields.
type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
}

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/testdb?parseTime=true"

	// sqlx.Connect combines sql.Open and db.Ping
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected with sqlx!")

	// --- Basic Operations ---

	// Example: Get a single user (db.Get)
	var user User
	err = db.Get(&user, "SELECT id, username, email FROM users WHERE id = ?", 1)
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		fmt.Printf("Got user (Get): %+v\n", user)
	}

	// Example: Get multiple users (db.Select)
	var users []User
	err = db.Select(&users, "SELECT id, username, email FROM users WHERE id > ?", 0)
	if err != nil {
		log.Printf("Error selecting users: %v", err)
	} else {
		fmt.Println("Got all users (Select):")
		for _, u := range users {
			fmt.Printf("- %+v\n", u)
		}
	}
}
```

4. ent
ent is an entity framework for Go. It uses code generation to create a type-safe, graph-based API for your database. The setup is more involved than other libraries.

```go
// ent/schema/user.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("unknown"),
		field.Int("age").
			Positive(),
		field.String("email").
			Unique(),
	}
}
```

5. MongoDB Official Driver
This is the official driver from MongoDB for working with MongoDB collections and documents.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Use context.Background() with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Defer disconnect
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// Ping the primary
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	// --- Basic Operations ---

	// Get a handle for your collection
	collection := client.Database("testdb").Collection("users")

	// Insert a document
	// bson.M is an unordered map. For ordered, use bson.D
	userDoc := bson.M{"name": "Charlie", "age": 40, "email": "charlie@example.com"}
	insertResult, err := collection.InsertOne(ctx, userDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// Find a document
	var result bson.M
	err = collection.FindOne(ctx, bson.M{"name": "Charlie"}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("No document found")
		} else {
			log.Fatal(err)
		}
	} else {
		fmt.Printf("Found document: %+v\n", result)
	}
}
```

6. Redis Client (go-redis)
This is the most popular client for working with the Redis in-memory data store.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"[github.com/redis/go-redis/v9](https://github.com/redis/go-redis/v9)"
)

func main() {
	// Use a background context
	ctx := context.Background()

	// Create a new client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Ping the server to check the connection
	status, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully connected to Redis! Ping status: %s\n", status)

	// --- Basic Operations ---

	// Set a key-value pair
	// Set "mykey" to "myvalue", with an expiration of 10 seconds
	err = rdb.Set(ctx, "mykey", "myvalue", 10*time.Second).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("SET 'mykey' to 'myvalue'")

	// Get the value for the key
	val, err := rdb.Get(ctx, "mykey").Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("'mykey' does not exist")
		} else {
			panic(err)
		}
	} else {
		fmt.Printf("GET 'mykey': %s\n", val)
	}

	// Get a key that doesn't exist
	val2, err := rdb.Get(ctx, "nonexistent_key").Result()
	if err == redis.Nil {
		fmt.Println("GET 'nonexistent_key': (key does not exist, as expected)")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("nonexistent_key:", val2)
	}
}
```
