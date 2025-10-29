package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ProjectName}}/internal/config"
{{if ne .DbDriver.ID ""}}	"{{.ProjectName}}/internal/storage/{{.DbDriver.ID}}"
{{end}}	"{{.ProjectName}}/internal/handler"
	"{{.ProjectName}}/internal/service"

	"github.com/gin-gonic/gin"
{{if .HasFeature "cors"}}	"github.com/rs/cors"
{{end}}{{if .HasFeature "logging"}}	"github.com/rs/zerolog/log"
{{end}}{{if .HasFeature "jwt"}}	"github.com/golang-jwt/jwt/v5"
{{end}}{{if eq .DbDriver.ID "gorm"}}	"gorm.io/gorm"
{{if eq .Database.ID "postgres"}}	"gorm.io/driver/postgres"
{{else if eq .Database.ID "mysql"}}	"gorm.io/driver/mysql"
{{else if eq .Database.ID "sqlite"}}	"gorm.io/driver/sqlite"
{{end}}{{else if eq .DbDriver.ID "sqlx"}}	"github.com/jmoiron/sqlx"
{{if eq .Database.ID "postgres"}}	_ "github.com/lib/pq"
{{else if eq .Database.ID "mysql"}}	_ "github.com/go-sql-driver/mysql"
{{else if eq .Database.ID "sqlite"}}	_ "github.com/mattn/go-sqlite3"
{{end}}{{else if eq .DbDriver.ID "mongo-driver"}}	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
{{else if eq .DbDriver.ID "redis-client"}}	"github.com/redis/go-redis/v9"
{{end}})

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

{{if .HasFeature "logging"}}	log.Info().Str("version", cfg.App.Version).Msg("Starting {{.ProjectName}}")
{{else}}	log.Printf("Starting %s version %s", cfg.App.Name, cfg.App.Version)
{{end}}

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

{{if ne .DbDriver.ID ""}}	// Initialize database
{{if eq .DbDriver.ID "gorm"}}	db, err := initGORMDatabase(cfg)
{{else if eq .DbDriver.ID "sqlx"}}	db, err := initSQLXDatabase(cfg)
{{else if eq .DbDriver.ID "mongo-driver"}}	db, err := initMongoDatabase(cfg)
{{else if eq .DbDriver.ID "redis-client"}}	db, err := initRedisDatabase(cfg)
{{else}}	db, err := initDatabase(cfg)
{{end}}	if err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Failed to initialize database")
{{else}}		log.Fatal("Failed to initialize database:", err)
{{end}}	}
{{if or (eq .DbDriver.ID "gorm") (eq .DbDriver.ID "sqlx")}}	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()
{{else if eq .DbDriver.ID "mongo-driver"}}	defer func() {
		if err := db.Disconnect(context.Background()); err != nil {
{{if .HasFeature "logging"}}			log.Error().Err(err).Msg("Failed to disconnect from MongoDB")
{{else}}			log.Printf("Failed to disconnect from MongoDB: %v", err)
{{end}}		}
	}()
{{else if eq .DbDriver.ID "redis-client"}}	defer db.Close()
{{end}}

{{if .HasFeature "logging"}}	log.Info().Msg("Database connection established")
{{else}}	log.Println("Database connection established")
{{end}}

{{end}}	// Initialize services
{{if ne .DbDriver.ID ""}}	svc := service.New(db)
{{else}}	svc := service.New()
{{end}}

	// Initialize HTTP handlers
	handlers := handler.New(svc, cfg)

	// Setup router
	router := setupRouter(handlers, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:    cfg.Address(),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
{{if .HasFeature "logging"}}		log.Info().Str("address", cfg.Address()).Msg("Starting HTTP server")
{{else}}		log.Printf("Starting HTTP server on %s", cfg.Address())
{{end}}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
{{if .HasFeature "logging"}}			log.Fatal().Err(err).Msg("Failed to start server")
{{else}}			log.Fatal("Failed to start server:", err)
{{end}}		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

{{if .HasFeature "logging"}}	log.Info().Msg("Shutting down server...")
{{else}}	log.Println("Shutting down server...")
{{end}}

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Server forced to shutdown")
{{else}}		log.Fatal("Server forced to shutdown:", err)
{{end}}	}

{{if .HasFeature "logging"}}	log.Info().Msg("Server exited")
{{else}}	log.Println("Server exited")
{{end}}
}

// setupRouter configures the Gin router with middleware and routes
func setupRouter(handlers *handler.Handler, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

{{if .HasFeature "cors"}}	// CORS middleware
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(func() gin.HandlerFunc {
		return gin.WrapH(corsConfig.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
	}())

{{end}}	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Example endpoints
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Welcome to " + cfg.App.Name + " API",
				"version": cfg.App.Version,
			})
		})

{{if .HasFeature "jwt"}}		// Protected routes (require JWT)
		protected := api.Group("/")
		protected.Use(JWTMiddleware(cfg.JWT.Secret))
		{
			protected.GET("/profile", handlers.GetProfile)
			protected.POST("/data", handlers.CreateData)
		}

{{end}}		// Public routes
		api.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "running"})
		})
	}

	return r
}

{{if ne .DbDriver.ID ""}}{{if eq .DbDriver.ID "gorm"}}// initGORMDatabase initializes GORM database connection
func initGORMDatabase(cfg *config.Config) (*gorm.DB, error) {
{{if eq .Database.ID "postgres"}}	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{})
{{else if eq .Database.ID "mysql"}}	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN()), &gorm.Config{})
{{else if eq .Database.ID "sqlite"}}	db, err := gorm.Open(sqlite.Open(cfg.SQLiteDSN()), &gorm.Config{})
{{else}}	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{})
{{end}}	if err != nil {
		return nil, err
	}

	// Auto-migrate your models here
	// err = db.AutoMigrate(&model.User{})
	// if err != nil {
	//     return nil, err
	// }

	return db, nil
}
{{else if eq .DbDriver.ID "sqlx"}}// initSQLXDatabase initializes sqlx database connection
func initSQLXDatabase(cfg *config.Config) (*sqlx.DB, error) {
{{if eq .Database.ID "postgres"}}	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
{{else if eq .Database.ID "mysql"}}	db, err := sqlx.Connect("mysql", cfg.MySQLDSN())
{{else if eq .Database.ID "sqlite"}}	db, err := sqlx.Connect("sqlite3", cfg.SQLiteDSN())
{{else}}	db, err := sqlx.Connect("postgres", cfg.PostgresDSN())
{{end}}	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
{{else if eq .DbDriver.ID "mongo-driver"}}// initMongoDatabase initializes MongoDB connection
func initMongoDatabase(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI()))
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}
{{else if eq .DbDriver.ID "redis-client"}}// initRedisDatabase initializes Redis connection
func initRedisDatabase(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL(),
		Password: cfg.Database.Password,
		DB:       cfg.Database.DB,
	})

	// Test the connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
{{else}}// initDatabase initializes database connection
func initDatabase(cfg *config.Config) (interface{}, error) {
	// Implement your database initialization logic here
	return nil, nil
}
{{end}}{{end}}

{{if .HasFeature "jwt"}}// JWTMiddleware validates JWT tokens
func JWTMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("username", claims["username"])
		}

		c.Next()
	}
}
{{end}}
