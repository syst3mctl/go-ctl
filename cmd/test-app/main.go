package main

import (
	"log"
	"test-app/internal/config"
	"test-app/internal/handler"
	"test-app/internal/service"
	"test-app/internal/storage"
	"test-app/internal/storage/postgres"
)

func main() {
	log.Println("Starting test-app")


	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := storage.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	defer func() {
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}()


	log.Println("Database connection established")


	// Initialize repository
	repo := postgres.New(db)

	// Run migrations
	if err := repo.Migrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrations completed")


	// Initialize services
	svc := service.New(repo)

	// Initialize HTTP handlers
	handlers := handler.New(svc, cfg)

	// Start HTTP server - ALWAYS through handler/http.go
	if err := handler.InitHTTPServer(handlers, cfg); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
