package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ProjectName}}/internal/config"
{{if ne .DbDriver.ID ""}}	"{{.ProjectName}}/internal/storage/{{.DbDriver.ID}}"
{{end}}	"{{.ProjectName}}/internal/handler"
	"{{.ProjectName}}/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
{{if .HasFeature "cors"}}	"github.com/labstack/echo/v4/middleware"
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
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Failed to load configuration")
{{else}}		panic("Failed to load configuration: " + err.Error())
{{end}}	}

{{if .HasFeature "logging"}}	log.Info().Str("version", cfg.App.Version).Msg("Starting {{.ProjectName}}")
{{else}}	println("Starting", cfg.App.Name, "version", cfg.App.Version)
{{end}}

{{if ne .DbDriver.ID ""}}	// Initialize database
{{if eq .DbDriver.ID "gorm"}}	db, err := initGORMDatabase(cfg)
{{else if eq .DbDriver.ID "sqlx"}}	db, err := initSQLXDatabase(cfg)
{{else if eq .DbDriver.ID "mongo-driver"}}	db, err := initMongoDatabase(cfg)
{{else if eq .DbDriver.ID "redis-client"}}	db, err := initRedisDatabase(cfg)
{{else}}	db, err := initDatabase(cfg)
{{end}}	if err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Failed to initialize database")
{{else}}		panic("Failed to initialize database: " + err.Error())
{{end}}	}
{{if or (eq .DbDriver.ID "gorm") (eq .DbDriver.ID "sqlx")}}	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()
{{else if eq .DbDriver.ID "mongo-driver"}}	defer func() {
		if err := db.Disconnect(context.Background()); err != nil {
{{if .HasFeature "logging"}}			log.Error().Err(err).Msg("Failed to disconnect from MongoDB")
{{else}}			println("Failed to disconnect from MongoDB:", err.Error())
{{end}}		}
	}()
{{else if eq .DbDriver.ID "redis-client"}}	defer db.Close()
{{end}}

{{if .HasFeature "logging"}}	log.Info().Msg("Database connection established")
{{else}}	println("Database connection established")
{{end}}

{{end}}	// Initialize services
{{if ne .DbDriver.ID ""}}	svc := service.New(db)
{{else}}	svc := service.New()
{{end}}

	// Initialize HTTP handlers
	handlers := handler.New(svc, cfg)

	// Create Echo instance
	e := echo.New()

	// Setup middleware and routes
	setupEcho(e, handlers, cfg)

	// Start server in a goroutine
	go func() {
{{if .HasFeature "logging"}}		log.Info().Str("address", cfg.Address()).Msg("Starting HTTP server")
{{else}}		println("Starting HTTP server on", cfg.Address())
{{end}}
		if err := e.Start(":" + string(rune(cfg.Server.Port))); err != nil && err != http.ErrServerClosed {
{{if .HasFeature "logging"}}			log.Fatal().Err(err).Msg("Failed to start server")
{{else}}			panic("Failed to start server: " + err.Error())
{{end}}		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

{{if .HasFeature "logging"}}	log.Info().Msg("Shutting down server...")
{{else}}	println("Shutting down server...")
{{end}}

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Server forced to shutdown")
{{else}}		panic("Server forced to shutdown: " + err.Error())
{{end}}	}

{{if .HasFeature "logging"}}	log.Info().Msg("Server exited")
{{else}}	println("Server exited")
{{end}}
}

// setupEcho configures the Echo server with middleware and routes
func setupEcho(e *echo.Echo, handlers *handler.Handler, cfg *config.Config) {
	// Hide Echo banner in production
	if cfg.IsProduction() {
		e.HideBanner = true
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

{{if .HasFeature "cors"}}	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

{{end}}	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"service": cfg.App.Name,
			"version": cfg.App.Version,
		})
	})

	// API routes
	api := e.Group("/api/v1")

	// Welcome endpoint
	api.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome to " + cfg.App.Name + " API",
			"version": cfg.App.Version,
		})
	})

{{if .HasFeature "jwt"}}	// Protected routes (require JWT)
	protected := api.Group("")
	protected.Use(JWTMiddleware(cfg.JWT.Secret))
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.POST("/data", handlers.CreateData)
	}

{{end}}	// Public routes
	api.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "running"})
	})

	// User CRUD routes (example)
	users := api.Group("/users")
	users.POST("", handlers.CreateUser)
	users.GET("/:id", handlers.GetUser)
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

{{if .HasFeature "jwt"}}// JWTMiddleware validates JWT tokens for Echo
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
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
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				c.Set("user_id", claims["user_id"])
				c.Set("username", claims["username"])
			}

			return next(c)
		}
	}
}
{{end}}
