package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ProjectName}}/internal/config"
{{if eq .HTTP.ID "gin"}}	"github.com/gin-gonic/gin"
{{else if eq .HTTP.ID "echo"}}	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
{{else if eq .HTTP.ID "fiber"}}	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
{{else if eq .HTTP.ID "chi"}}	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
{{end}}{{if .HasFeature "cors"}}{{if eq .HTTP.ID "gin"}}	"github.com/rs/cors"
{{else if eq .HTTP.ID "chi"}}	"github.com/go-chi/cors"
{{end}}{{end}}{{if .HasFeature "logging"}}	"github.com/rs/zerolog/log"
{{end}}{{if .HasFeature "jwt"}}	"github.com/golang-jwt/jwt/v5"
{{end}})

// InitHTTPServer initializes and starts the HTTP server for any framework
func InitHTTPServer(handlers *Handler, cfg *config.Config) error {
{{if eq .HTTP.ID "gin"}}	return initGinServer(handlers, cfg)
{{else if eq .HTTP.ID "echo"}}	return initEchoServer(handlers, cfg)
{{else if eq .HTTP.ID "fiber"}}	return initFiberServer(handlers, cfg)
{{else if eq .HTTP.ID "chi"}}	return initChiServer(handlers, cfg)
{{else}}	return initNetHTTPServer(handlers, cfg)
{{end}}
}

{{if eq .HTTP.ID "gin"}}// initGinServer initializes and starts the Gin HTTP server
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
{{if .HasFeature "logging"}}		log.Info().Str("address", cfg.Address()).Msg("Starting HTTP server")
{{else}}		fmt.Printf("Starting HTTP server on %s\n", cfg.Address())
{{end}}		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
{{if .HasFeature "logging"}}			log.Fatal().Err(err).Msg("Failed to start server")
{{else}}			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
{{end}}		}
	}()

	return gracefulShutdown(srv)
}

// setupGinRouter configures the Gin router with middleware and routes
func setupGinRouter(handlers *Handler, cfg *config.Config) *gin.Engine {
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

{{if .HasFeature "jwt"}}		// Protected routes (require JWT)
		protected := api.Group("/")
		protected.Use(ginJWTMiddleware(cfg.JWT.Secret))
		{
			protected.GET("/profile", handlers.GetProfile)
			protected.POST("/users", handlers.CreateUser)
			protected.GET("/users", handlers.GetUsers)
		}

{{end}}		// Public routes
		api.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "running"})
		})
	}

	return r
}

{{if .HasFeature "jwt"}}// ginJWTMiddleware validates JWT tokens for Gin
func ginJWTMiddleware(secret string) gin.HandlerFunc {
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

{{else if eq .HTTP.ID "echo"}}// initEchoServer initializes and starts the Echo HTTP server
func initEchoServer(handlers *Handler, cfg *config.Config) error {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

{{if .HasFeature "cors"}}	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{"*"},
	}))

{{end}}	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "ok",
			"service":   cfg.App.Name,
			"version":   cfg.App.Version,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API routes
	api := e.Group("/api/v1")
	{
		// Example endpoints
		api.GET("/", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "Welcome to " + cfg.App.Name + " API",
				"version": cfg.App.Version,
			})
		})

{{if .HasFeature "jwt"}}		// Protected routes (require JWT)
		protected := api.Group("")
		protected.Use(echoJWTMiddleware(cfg.JWT.Secret))
		{
			protected.GET("/profile", func(c echo.Context) error { return handlers.GetProfile(c) })
			protected.POST("/users", func(c echo.Context) error { return handlers.CreateUser(c) })
			protected.GET("/users", func(c echo.Context) error { return handlers.GetUsers(c) })
		}

{{end}}		// Public routes
		api.GET("/status", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "running"})
		})
	}

	// Start server in a goroutine
	go func() {
{{if .HasFeature "logging"}}		log.Info().Str("address", cfg.Address()).Msg("Starting HTTP server")
{{else}}		fmt.Printf("Starting HTTP server on %s\n", cfg.Address())
{{end}}		if err := e.Start(cfg.Address()); err != nil && err != http.ErrServerClosed {
{{if .HasFeature "logging"}}			log.Fatal().Err(err).Msg("Failed to start server")
{{else}}			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
{{end}}		}
	}()

	return gracefulShutdownEcho(e)
}

{{if .HasFeature "jwt"}}// echoJWTMiddleware validates JWT tokens for Echo
func echoJWTMiddleware(secret string) echo.MiddlewareFunc {
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

// gracefulShutdownEcho handles graceful shutdown for Echo
func gracefulShutdownEcho(e *echo.Echo) error {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

{{if .HasFeature "logging"}}	log.Info().Msg("Shutting down server...")
{{else}}	fmt.Println("Shutting down server...")
{{end}}

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Server forced to shutdown")
{{else}}		fmt.Printf("Server forced to shutdown: %v\n", err)
{{end}}		return err
	}

{{if .HasFeature "logging"}}	log.Info().Msg("Server exited")
{{else}}	fmt.Println("Server exited")
{{end}}
	return nil
}

{{else if eq .HTTP.ID "fiber"}}// initFiberServer initializes and starts the Fiber HTTP server
func initFiberServer(handlers *Handler, cfg *config.Config) error {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

{{if .HasFeature "cors"}}	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:8080",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

{{end}}	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"service":   cfg.App.Name,
			"version":   cfg.App.Version,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API routes
	api := app.Group("/api/v1")
	{
		// Example endpoints
		api.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "Welcome to " + cfg.App.Name + " API",
				"version": cfg.App.Version,
			})
		})

{{if .HasFeature "jwt"}}		// Protected routes (require JWT)
		protected := api.Group("", fiberJWTMiddleware(cfg.JWT.Secret))
		{
			protected.Get("/profile", func(c *fiber.Ctx) error { return handlers.GetProfile(c) })
			protected.Post("/users", func(c *fiber.Ctx) error { return handlers.CreateUser(c) })
			protected.Get("/users", func(c *fiber.Ctx) error { return handlers.GetUsers(c) })
		}

{{end}}		// Public routes
		api.Get("/status", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"status": "running"})
		})
	}

	// Start server in a goroutine
	go func() {
{{if .HasFeature "logging"}}		log.Info().Str("address", cfg.Address()).Msg("Starting HTTP server")
{{else}}		fmt.Printf("Starting HTTP server on %s\n", cfg.Address())
{{end}}		if err := app.Listen(cfg.Address()); err != nil {
{{if .HasFeature "logging"}}			log.Fatal().Err(err).Msg("Failed to start server")
{{else}}			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
{{end}}		}
	}()

	return gracefulShutdownFiber(app)
}

{{if .HasFeature "jwt"}}// fiberJWTMiddleware validates JWT tokens for Fiber
func fiberJWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header required"})
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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_id", claims["user_id"])
			c.Locals("username", claims["username"])
		}

		return c.Next()
	}
}
{{end}}

// gracefulShutdownFiber handles graceful shutdown for Fiber
func gracefulShutdownFiber(app *fiber.App) error {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

{{if .HasFeature "logging"}}	log.Info().Msg("Shutting down server...")
{{else}}	fmt.Println("Shutting down server...")
{{end}}

	if err := app.Shutdown(); err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Server forced to shutdown")
{{else}}		fmt.Printf("Server forced to shutdown: %v\n", err)
{{end}}		return err
	}

{{if .HasFeature "logging"}}	log.Info().Msg("Server exited")
{{else}}	fmt.Println("Server exited")
{{end}}
	return nil
}

{{else if eq .HTTP.ID "chi"}}// initChiServer initializes and starts the Chi HTTP server
func initChiServer(handlers *Handler, cfg *config.Config) error {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

{{if .HasFeature "cors"}}	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

{{end}}	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"%s","version":"%s","timestamp":"%s"}`,
			cfg.App.Name, cfg.App.Version, time.Now().UTC().Format(time.RFC3339))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Example endpoints
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message":"Welcome to %s API","version":"%s"}`,
				cfg.App.Name, cfg.App.Version)
		})

{{if .HasFeature "jwt"}}		// Protected routes (require JWT)
		r.Group(func(r chi.Router) {
			r.Use(chiJWTMiddleware(cfg.JWT.Secret))
			r.Get("/profile", func(w http.ResponseWriter, r *http.Request) { handlers.GetProfile(w, r) })
			r.Post("/users", func(w http.ResponseWriter, r *http.Request) { handlers.CreateUser(w, r) })
			r.Get("/users", func(w http.ResponseWriter, r *http.Request) { handlers.GetUsers(w, r) })
		})

{{end}}		// Public routes
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"status":"running"}`)
		})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    cfg.Address(),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
{{if .HasFeature "logging"}}		log.Info().Str("address", cfg.Address()).Msg("Starting HTTP server")
{{else}}		fmt.Printf("Starting HTTP server on %s\n", cfg.Address())
{{end}}		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
{{if .HasFeature "logging"}}			log.Fatal().Err(err).Msg("Failed to start server")
{{else}}			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
{{end}}		}
	}()

	return gracefulShutdown(srv)
}

{{if .HasFeature "jwt"}}// chiJWTMiddleware validates JWT tokens for Chi
func chiJWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, `{"error":"Authorization header required"}`)
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
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, `{"error":"Invalid token"}`)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
				ctx = context.WithValue(ctx, "username", claims["username"])
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
{{end}}

{{else}}// initNetHTTPServer initializes and starts the net/http server
func initNetHTTPServer(handlers *Handler, cfg *config.Config) error {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"%s","version":"%s","timestamp":"%s"}`,
			cfg.App.Name, cfg.App.Version, time.Now().UTC().Format(time.RFC3339))
	})

	// API routes
	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"message":"Welcome to %s API","version":"%s"}`,
				cfg.App.Name, cfg.App.Version)
			return
		}
		http.NotFound(w, r)
	})

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"running"}`)
	})

{{if .HasFeature "jwt"}}	// Protected routes
	mux.HandleFunc("/api/v1/profile", netHTTPJWTMiddleware(cfg.JWT.Secret, func(w http.ResponseWriter, r *http.Request) {
		handlers.GetProfile(w, r)
	}))
	mux.HandleFunc("/api/v1/users", netHTTPJWTMiddleware(cfg.JWT.Secret, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateUser(w, r)
		case http.MethodGet:
			handlers.GetUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
{{end}}

	// Create HTTP server
	srv := &http.Server{
		Addr:    cfg.Address(),
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
{{if .HasFeature "logging"}}		log.Info().Str("address", cfg.Address()).Msg("Starting HTTP server")
{{else}}		fmt.Printf("Starting HTTP server on %s\n", cfg.Address())
{{end}}		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
{{if .HasFeature "logging"}}			log.Fatal().Err(err).Msg("Failed to start server")
{{else}}			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
{{end}}		}
	}()

	return gracefulShutdown(srv)
}

{{if .HasFeature "jwt"}}// netHTTPJWTMiddleware validates JWT tokens for net/http
func netHTTPJWTMiddleware(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error":"Authorization header required"}`)
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"error":"Invalid token"}`)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			ctx = context.WithValue(ctx, "username", claims["username"])
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	}
}
{{end}}

{{end}}

// gracefulShutdown handles graceful shutdown for standard HTTP server
func gracefulShutdown(srv *http.Server) error {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

{{if .HasFeature "logging"}}	log.Info().Msg("Shutting down server...")
{{else}}	fmt.Println("Shutting down server...")
{{end}}

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
{{if .HasFeature "logging"}}		log.Fatal().Err(err).Msg("Server forced to shutdown")
{{else}}		fmt.Printf("Server forced to shutdown: %v\n", err)
{{end}}		return err
	}

{{if .HasFeature "logging"}}	log.Info().Msg("Server exited")
{{else}}	fmt.Println("Server exited")
{{end}}
	return nil
}
