// Package handlers provide way to work with HTTP frameworks
package handlers

import (
{{if eq .HTTP.ID "gin"}}	"github.com/gin-gonic/gin"
{{else if eq .HTTP.ID "echo"}}	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
{{else if eq .HTTP.ID "fiber"}}	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
{{else if eq .HTTP.ID "chi"}}	"github.com/go-chi/chi/v5"
{{else}}	"net/http"
{{end}}
)

{{if eq .HTTP.ID "gin"}}// Routes returns the Gin router with all routes configured
func Routes() *gin.Engine {
	router := gin.New()
	
	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API routes
	api := router.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/login", H.LoginHandler)
			authRoutes.POST("/logout", H.LogoutHandler)
			authRoutes.POST("/hash", H.HashPasswordHandler)
		}
	}

	return router
}

{{else if eq .HTTP.ID "echo"}}// Routes returns the Echo router with all routes configured
func Routes() *echo.Echo {
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API routes
	api := e.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/login", H.LoginHandler)
			authRoutes.POST("/logout", H.LogoutHandler)
			authRoutes.POST("/hash", H.HashPasswordHandler)
		}
	}

	return e
}

{{else if eq .HTTP.ID "fiber"}}// Routes returns the Fiber app with all routes configured
func Routes() *fiber.App {
	app := fiber.New()

	// Add middleware
	app.Use(fiberLogger.New())
	app.Use(fiberRecover.New())

	// API routes
	api := app.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.Post("/login", H.LoginHandler)
			authRoutes.Post("/logout", H.LogoutHandler)
			authRoutes.Post("/hash", H.HashPasswordHandler)
		}
	}

	return app
}

{{else if eq .HTTP.ID "chi"}}// Routes returns the Chi router with all routes configured
func Routes() *chi.Mux {
	r := chi.NewRouter()

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", H.LoginHandler)
			r.Post("/logout", H.LogoutHandler)
			r.Post("/hash", H.HashPasswordHandler)
		})
	})

	return r
}

{{else}}// Routes returns the net/http ServeMux with all routes configured
func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	authRoutes := authRoutes()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", authRoutes))
	return mux
}

func authRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/login", H.LoginHandler)
	mux.HandleFunc("POST /auth/logout", H.LogoutHandler)
	mux.HandleFunc("POST /auth/hash", H.HashPasswordHandler)

	return mux
}
{{end}}
