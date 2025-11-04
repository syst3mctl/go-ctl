package handlers

import (
{{if eq .HTTP.ID "gin"}}	"net/http"
	"github.com/gin-gonic/gin"
{{else if eq .HTTP.ID "echo"}}	"net/http"
	"github.com/labstack/echo/v4"
{{else if eq .HTTP.ID "fiber"}}	"github.com/gofiber/fiber/v2"
{{else}}	"log"
	"net/http"
{{end}}

	"{{.ProjectName}}/internal/gen"
{{if or (eq .HTTP.ID "chi") (eq .HTTP.ID "net-http")}}	"{{.ProjectName}}/internal/handlers/dto"
{{end}}	store2 "{{.ProjectName}}/internal/store"
{{if or (eq .HTTP.ID "chi") (eq .HTTP.ID "net-http")}}	"{{.ProjectName}}/internal/validate"
{{else if or (eq .HTTP.ID "gin") (eq .HTTP.ID "echo") (eq .HTTP.ID "fiber")}}	"{{.ProjectName}}/internal/handlers/dto"
	"{{.ProjectName}}/internal/validate"
{{end}}
)

type loginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

{{if eq .HTTP.ID "gin"}}// LoginHandler handles user login using Gin
func (h *Handler) LoginHandler(c *gin.Context) {
	var loginReq loginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dto.Validate.Struct(loginReq); err != nil {
		errs := validate.RangeErrors(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"data": gin.H{
				"message": "unprocessable entity",
				"errors": errs,
			},
		})
		return
	}

	user, err := h.app.Store.Users.GetByEmail(c.Request.Context(), loginReq.Email)
	if err != nil {
		switch err {
		case store2.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	err = gen.CompareHashAndPasswordBcrypt(user.Password, loginReq.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	str, err := gen.GenerateRandomString(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate session"})
		return
	}

	h.app.Store.Users.CreateUserSession(c.Request.Context(), &store2.UserSession{
		SessionID: str,
		UserID:    user.ID,
	})

	c.Header("HX-Redirect", "/v1/dashboard/product")
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "session_id": str})
}

// LogoutHandler handles user logout using Gin
func (h *Handler) LogoutHandler(c *gin.Context) {
	c.Header("HX-Redirect", "/v1/login")
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// HashPasswordHandler handles password hashing using Gin
func (h *Handler) HashPasswordHandler(c *gin.Context) {
	var req struct {
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password field is required"})
		return
	}

	hash, err := gen.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hash": hash})
}

{{else if eq .HTTP.ID "echo"}}// LoginHandler handles user login using Echo
func (h *Handler) LoginHandler(c echo.Context) error {
	var loginReq loginRequest

	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := dto.Validate.Struct(loginReq); err != nil {
		errs := validate.RangeErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"data": map[string]interface{}{
				"message": "unprocessable entity",
				"errors": errs,
			},
		})
	}

	user, err := h.app.Store.Users.GetByEmail(c.Request().Context(), loginReq.Email)
	if err != nil {
		switch err {
		case store2.ErrNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		}
	}

	err = gen.CompareHashAndPasswordBcrypt(user.Password, loginReq.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid credentials"})
	}

	str, err := gen.GenerateRandomString(32)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate session"})
	}

	h.app.Store.Users.CreateUserSession(c.Request().Context(), &store2.UserSession{
		SessionID: str,
		UserID:    user.ID,
	})

	c.Response().Header().Set("HX-Redirect", "/v1/dashboard/product")
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "login successful", "session_id": str})
}

// LogoutHandler handles user logout using Echo
func (h *Handler) LogoutHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", "/v1/login")
	return c.JSON(http.StatusOK, map[string]string{"message": "logout successful"})
}

// HashPasswordHandler handles password hashing using Echo
func (h *Handler) HashPasswordHandler(c echo.Context) error {
	var req struct {
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "password field is required"})
	}

	if err := dto.Validate.Struct(req); err != nil {
		errs := validate.RangeErrors(err)
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"data": map[string]interface{}{
				"message": "unprocessable entity",
				"errors": errs,
			},
		})
	}

	hash, err := gen.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"hash": hash})
}

{{else if eq .HTTP.ID "fiber"}}// LoginHandler handles user login using Fiber
func (h *Handler) LoginHandler(c *fiber.Ctx) error {
	var loginReq loginRequest

	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := dto.Validate.Struct(loginReq); err != nil {
		errs := validate.RangeErrors(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"data": fiber.Map{
				"message": "unprocessable entity",
				"errors": errs,
			},
		})
	}

	user, err := h.app.Store.Users.GetByEmail(c.Context(), loginReq.Email)
	if err != nil {
		switch err {
		case store2.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
	}

	err = gen.CompareHashAndPasswordBcrypt(user.Password, loginReq.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid credentials"})
	}

	str, err := gen.GenerateRandomString(32)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate session"})
	}

	h.app.Store.Users.CreateUserSession(c.Context(), &store2.UserSession{
		SessionID: str,
		UserID:    user.ID,
	})

	c.Set("HX-Redirect", "/v1/dashboard/product")
	return c.JSON(fiber.Map{"message": "login successful", "session_id": str})
}

// LogoutHandler handles user logout using Fiber
func (h *Handler) LogoutHandler(c *fiber.Ctx) error {
	c.Set("HX-Redirect", "/v1/login")
	return c.JSON(fiber.Map{"message": "logout successful"})
}

// HashPasswordHandler handles password hashing using Fiber
func (h *Handler) HashPasswordHandler(c *fiber.Ctx) error {
	var req struct {
		Password string `json:"password" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password field is required"})
	}

	if err := dto.Validate.Struct(req); err != nil {
		errs := validate.RangeErrors(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"data": fiber.Map{
				"message": "unprocessable entity",
				"errors": errs,
			},
		})
	}

	hash, err := gen.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to hash password"})
	}

	return c.JSON(fiber.Map{"hash": hash})
}

{{else if eq .HTTP.ID "chi"}}// LoginHandler handles user login using Chi
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq loginRequest

	err := dto.ParseAndValidate(w, r, &loginReq)
	if err != nil {
		validate.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.app.Store.Users.GetByEmail(r.Context(), loginReq.Email)
	if err != nil {
		switch err {
		case store2.ErrNotFound:
			validate.NotFoundResponse(w, r, err)
		default:
			validate.InternalServerError(w, r, err)
		}
		return
	}

	err = gen.CompareHashAndPasswordBcrypt(user.Password, loginReq.Password)
	if err != nil {
		validate.BadRequestResponse(w, r, err)
		return
	}

	str, err := gen.GenerateRandomString(32)
	if err != nil {
		validate.InternalServerError(w, r, err)
		return
	}

	h.app.Store.Users.CreateUserSession(r.Context(), &store2.UserSession{
		SessionID: str,
		UserID:    user.ID,
	})

	w.Header().Set("HX-Redirect", "/v1/dashboard/product")
	validate.SendResponse(w, r, http.StatusOK, map[string]interface{}{"message": "login successful", "session_id": str})
}

// LogoutHandler handles user logout using Chi
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Redirect", "/v1/login")
	validate.SendResponse(w, r, http.StatusOK, map[string]string{"message": "logout successful"})
}

// HashPasswordHandler handles password hashing using Chi
func (h *Handler) HashPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password" validate:"required"`
	}

	err := dto.ParseAndValidate(w, r, &req)
	if err != nil {
		validate.BadRequestResponse(w, r, err)
		return
	}

	hash, err := gen.HashPassword(req.Password)
	if err != nil {
		log.Printf("ERROR: failed to hash password: %v", err)
		validate.InternalServerError(w, r, err)
		return
	}

	validate.SendResponse(w, r, http.StatusOK, map[string]string{"hash": hash})
}

{{else}}// LoginHandler handles user login using net/http
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq loginRequest

	err := dto.ParseAndValidate(w, r, &loginReq)
	if err != nil {
		validate.BadRequestResponse(w, r, err)
		return
	}

	user, err := h.app.Store.Users.GetByEmail(r.Context(), loginReq.Email)
	if err != nil {
		switch err {
		case store2.ErrNotFound:
			validate.NotFoundResponse(w, r, err)
		default:
			validate.InternalServerError(w, r, err)
		}
		return
	}

	err = gen.CompareHashAndPasswordBcrypt(user.Password, loginReq.Password)
	if err != nil {
		validate.BadRequestResponse(w, r, err)
		return
	}

	str, err := gen.GenerateRandomString(32)
	if err != nil {
		validate.InternalServerError(w, r, err)
		return
	}

	h.app.Store.Users.CreateUserSession(r.Context(), &store2.UserSession{
		SessionID: str,
		UserID:    user.ID,
	})

	w.Header().Set("HX-Redirect", "/v1/dashboard/product")
}

// LogoutHandler handles user logout using net/http
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("HX-Redirect", "/v1/login")
}

// HashPasswordHandler handles password hashing using net/http
func (h *Handler) HashPasswordHandler(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("Password")
	if password == "" {
		http.Error(w, "Password field is required", http.StatusBadRequest)
		return
	}

	hash, err := gen.HashPassword(password)
	if err != nil {
		log.Printf("ERROR: failed to hash password: %v", err)
		validate.InternalServerError(w, r, err)
		return
	}

	log.Printf("Successfully generated hash: %s", hash)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(hash))
}
{{end}}
