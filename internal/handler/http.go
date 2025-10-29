package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test-app/internal/config"
	"github.com/gin-gonic/gin"
)

// InitHTTPServer initializes and starts the HTTP server for any framework
func InitHTTPServer(handlers *Handler, cfg *config.Config) error {
	return initGinServer(handlers, cfg)

}

// initGinServer initializes and starts the Gin HTTP server
func initGinServer(handlers *Handler, cfg *config.Config) error {
	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup router
	router := setupGinRouter(handlers, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:    cfg.Address(),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting HTTP server on %s\n", cfg.Address())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	return gracefulShutdown(srv)
}

// setupGinRouter configures the Gin router with middleware and routes
func setupGinRouter(handlers *Handler, cfg *config.Config) *gin.Engine {
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   cfg.App.Name,
			"version":   cfg.App.Version,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
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

		// Public routes
		api.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "running"})
		})
	}

	return r
}





// gracefulShutdown handles graceful shutdown for standard HTTP server
func gracefulShutdown(srv *http.Server) error {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")


	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
		return err
	}

	fmt.Println("Server exited")

	return nil
}
