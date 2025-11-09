package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/syst3mctl/go-ctl/internal/generator"
	"github.com/syst3mctl/go-ctl/internal/metadata"
	"github.com/syst3mctl/go-ctl/internal/storage"
)

var (
	appOptions *metadata.ProjectOptions
	gen        *generator.Generator
)

func main() {
	// Load options from options.json
	var err error
	appOptions, err = metadata.LoadOptions()
	if err != nil {
		log.Fatal("Failed to load options:", err)
	}

	// Initialize generator
	gen = generator.New()
	if err := gen.LoadTemplates(); err != nil {
		log.Fatal("Failed to load templates:", err)
	}

	// Initialize analytics database
	if err := storage.InitAnalyticsDB(); err != nil {
		log.Printf("Warning: Failed to initialize analytics database: %v", err)
		log.Println("Analytics tracking will be disabled")
	}

	// Setup HTTP routes
	http.HandleFunc("/", handleLanding)           // Landing page
	http.HandleFunc("/generator", handleIndex)    // Project generator interface
	http.HandleFunc("/generate", handleGenerate)
	http.HandleFunc("/explore", handleExplore)
	http.HandleFunc("/search-packages", handleSearchPackages) // Legacy endpoint
	http.HandleFunc("/fetch-packages", handleFetchPackages)   // New dynamic endpoint
	http.HandleFunc("/add-package", handleAddPackage)
	http.HandleFunc("/search-npm-packages", handleSearchNpmPackages) // npm package search
	http.HandleFunc("/add-npm-package", handleAddNpmPackage)         // Add npm package
	http.HandleFunc("/file-content", handleFileContent)

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	fmt.Printf("🚀 go-ctl server starting on http://localhost:%s\n", port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("   GET  /              - Landing page")
	fmt.Println("   GET  /generator      - Project generator interface")
	fmt.Println("   POST /generate       - Generate and download project ZIP")
	fmt.Println("   POST /explore         - Preview project structure")
	fmt.Println("   GET  /search-packages - Search pkg.go.dev for packages (legacy)")
	fmt.Println("   GET  /fetch-packages  - Dynamic package search API (supports JSON & HTML)")
	fmt.Println("   POST /add-package     - Add package to selection")
	fmt.Println("   GET  /search-npm-packages - Search npm registry for packages")
	fmt.Println("   POST /add-npm-package - Add npm package to selection")
	fmt.Println("   GET  /file-content    - Get individual file content for preview")

	// Create server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	// Handle graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("\n🛑 Shutting down server...")
		srv.Close()
		os.Exit(0)
	}()

	// Start server
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
