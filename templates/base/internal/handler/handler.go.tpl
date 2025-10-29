package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"{{.ProjectName}}/internal/config"
	"{{.ProjectName}}/internal/service"

{{if eq .HTTP.ID "gin"}}	"github.com/gin-gonic/gin"
{{else if eq .HTTP.ID "echo"}}	"github.com/labstack/echo/v4"
{{else if eq .HTTP.ID "fiber"}}	"github.com/gofiber/fiber/v2"
{{else if eq .HTTP.ID "chi"}}	"github.com/go-chi/chi/v5"
{{end}}{{if .HasFeature "logging"}}	"github.com/rs/zerolog/log"
{{end}})

// Handler contains all HTTP handlers
type Handler struct {
	service *service.Service
	config  *config.Config
}

// New creates a new Handler instance
func New(svc *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		service: svc,
		config:  cfg,
	}
}

{{if eq .HTTP.ID "gin"}}// Gin Handlers

// CreateUser creates a new user
func (h *Handler) CreateUser(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Invalid request body")
{{end}}		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user")
{{end}}		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// GetUsers lists users with pagination
func (h *Handler) GetUsers(c *gin.Context) {
	var req service.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Invalid query parameters")
{{end}}		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := h.service.ListUsers(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list users")
{{end}}		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetProfile gets user profile (requires JWT)
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

{{if eq .DbDriver.ID "mongo-driver"}}	profile, err := h.service.GetUser(ctx, userID.(string))
{{else if eq .DbDriver.ID "gorm"}}	profile, err := h.service.GetUser(ctx, userID.(uint))
{{else}}	profile, err := h.service.GetUser(ctx, userID.(int64))
{{end}}	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to get user profile")
{{end}}		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}

{{else if eq .HTTP.ID "echo"}}// Echo Handlers

// CreateUser creates a new user
func (h *Handler) CreateUser(c echo.Context) error {
	var req service.CreateUserRequest
	if err := c.Bind(&req); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Invalid request body")
{{end}}		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user")
{{end}}		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"data": user})
}

// GetUsers lists users with pagination
func (h *Handler) GetUsers(c echo.Context) error {
	var req service.ListRequest
	if err := c.Bind(&req); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Invalid query parameters")
{{end}}		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid query parameters"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := h.service.ListUsers(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list users")
{{end}}		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": users})
}

// GetProfile gets user profile (requires JWT)
func (h *Handler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User ID not found in token"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

{{if eq .DbDriver.ID "mongo-driver"}}	profile, err := h.service.GetUser(ctx, userID.(string))
{{else if eq .DbDriver.ID "gorm"}}	profile, err := h.service.GetUser(ctx, userID.(uint))
{{else}}	profile, err := h.service.GetUser(ctx, userID.(int64))
{{end}}	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to get user profile")
{{end}}		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"data": profile})
}

{{else if eq .HTTP.ID "fiber"}}// Fiber Handlers

// CreateUser creates a new user
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var req service.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Invalid request body")
{{end}}		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user")
{{end}}		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"data": user})
}

// GetUsers lists users with pagination
func (h *Handler) GetUsers(c *fiber.Ctx) error {
	var req service.ListRequest
	if err := c.QueryParser(&req); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Invalid query parameters")
{{end}}		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid query parameters"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := h.service.ListUsers(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list users")
{{end}}		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": users})
}

// GetProfile gets user profile (requires JWT)
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in token"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

{{if eq .DbDriver.ID "mongo-driver"}}	profile, err := h.service.GetUser(ctx, userID.(string))
{{else if eq .DbDriver.ID "gorm"}}	profile, err := h.service.GetUser(ctx, userID.(uint))
{{else}}	profile, err := h.service.GetUser(ctx, userID.(int64))
{{end}}	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to get user profile")
{{end}}		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": profile})
}

{{else if eq .HTTP.ID "chi"}}// Chi Handlers (net/http style)

// CreateUser creates a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req service.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Invalid request body")
{{end}}		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user")
{{end}}		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": user})
}

// GetUsers lists users with pagination
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req service.ListRequest
	// Parse query parameters
	query := r.URL.Query()
	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}
	if offset := query.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			req.Offset = o
		}
	}
	req.Search = query.Get("search")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := h.service.ListUsers(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list users")
{{end}}		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": users})
}

// GetProfile gets user profile (requires JWT)
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

{{if eq .DbDriver.ID "mongo-driver"}}	profile, err := h.service.GetUser(ctx, userID.(string))
{{else if eq .DbDriver.ID "gorm"}}	profile, err := h.service.GetUser(ctx, userID.(uint))
{{else}}	profile, err := h.service.GetUser(ctx, userID.(int64))
{{end}}	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to get user profile")
{{end}}		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": profile})
}

{{else}}// net/http Handlers

// CreateUser creates a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req service.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Invalid request body")
{{end}}		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.service.CreateUser(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to create user")
{{end}}		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": user})
}

// GetUsers lists users with pagination
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req service.ListRequest
	// Parse query parameters
	query := r.URL.Query()
	if limit := query.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			req.Limit = l
		}
	}
	if offset := query.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			req.Offset = o
		}
	}
	req.Search = query.Get("search")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := h.service.ListUsers(ctx, req)
	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to list users")
{{end}}		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": users})
}

// GetProfile gets user profile (requires JWT)
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		http.Error(w, "User ID not found in token", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

{{if eq .DbDriver.ID "mongo-driver"}}	profile, err := h.service.GetUser(ctx, userID.(string))
{{else if eq .DbDriver.ID "gorm"}}	profile, err := h.service.GetUser(ctx, userID.(uint))
{{else}}	profile, err := h.service.GetUser(ctx, userID.(int64))
{{end}}	if err != nil {
{{if .HasFeature "logging"}}		log.Error().Err(err).Msg("Failed to get user profile")
{{end}}		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": profile})
}

{{end}}
