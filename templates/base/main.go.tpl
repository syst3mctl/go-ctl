package main

import (
{{if eq .DbDriver.ID "mongo-driver"}}	"context"
{{end}}	"log"
	"{{.ProjectName}}/internal/config"
	"{{.ProjectName}}/internal/handler"
	"{{.ProjectName}}/internal/service"
	"{{.ProjectName}}/internal/storage"
{{if ne .DbDriver.ID ""}}	"{{.ProjectName}}/internal/storage/{{.Database.ID}}"
{{end}}{{if .HasFeature "logging"}}	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
{{end}})

func main() {
{{if .HasFeature "logging"}}	// Setup logging
	if os.Getenv("ENVIRONMENT") != "production" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Starting {{.ProjectName}}")
{{else}}	log.Println("Starting {{.ProjectName}}")
{{end}}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Failed to load configuration")
{{else}}		log.Fatal("Failed to load configuration:", err)
{{end}}	}

{{if ne .DbDriver.ID ""}}	// Initialize database
	db, err := storage.InitDatabase(cfg)
	if err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Failed to initialize database")
{{else}}		log.Fatal("Failed to initialize database:", err)
{{end}}	}

{{if or (eq .DbDriver.ID "gorm") (eq .DbDriver.ID "sqlx")}}	defer func() {
{{if eq .DbDriver.ID "gorm"}}		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
{{else}}		db.Close()
{{end}}	}()
{{else if eq .DbDriver.ID "mongo-driver"}}	defer func() {
		ctx := context.Background()
		if err := db.Disconnect(ctx); err != nil {
{{if .HasFeature "logging"}}			log.Error().Err(err).Msg("Failed to disconnect from MongoDB")
{{else}}			log.Printf("Failed to disconnect from MongoDB: %v", err)
{{end}}		}
	}()
{{else if eq .DbDriver.ID "redis-client"}}	defer db.Close()
{{end}}

{{if .HasFeature "logging"}}	log.Info().Msg("Database connection established")
{{else}}	log.Println("Database connection established")
{{end}}

	// Initialize repository
	repo := {{.Database.ID}}.New(db)

	// Run migrations
	if err := repo.Migrate(); err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Failed to run migrations")
{{else}}		log.Fatal("Failed to run migrations:", err)
{{end}}	}

{{if .HasFeature "logging"}}	log.Info().Msg("Database migrations completed")
{{else}}	log.Println("Database migrations completed")
{{end}}

	// Initialize services
	svc := service.New(repo)

{{else}}	// Initialize services
	svc := service.New()

{{end}}	// Initialize HTTP handlers
	handlers := handler.New(svc, cfg)

	// Start HTTP server - ALWAYS through handler/http.go
	if err := handler.InitHTTPServer(handlers, cfg); err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Failed to start HTTP server")
{{else}}		log.Fatal("Failed to start HTTP server:", err)
{{end}}	}
}
