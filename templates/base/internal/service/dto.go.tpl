package service

import (
	"time"
)

// Common domain errors
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Predefined domain errors
var (
	ErrUserNotFound     = Error{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrUserExists       = Error{Code: "USER_EXISTS", Message: "User already exists"}
	ErrInvalidEmail     = Error{Code: "INVALID_EMAIL", Message: "Invalid email format"}
	ErrInvalidInput     = Error{Code: "INVALID_INPUT", Message: "Invalid input provided"}
{{if ne .Database.ID "redis"}}	ErrProductNotFound  = Error{Code: "PRODUCT_NOT_FOUND", Message: "Product not found"}
	ErrInvalidPrice     = Error{Code: "INVALID_PRICE", Message: "Price must be non-negative"}
{{end}}	ErrUnauthorized     = Error{Code: "UNAUTHORIZED", Message: "Unauthorized access"}
	ErrInternalError    = Error{Code: "INTERNAL_ERROR", Message: "Internal server error"}
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"min=1,max=100"`
}

// PaginatedResult represents a paginated result
type PaginatedResult[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=100"`
	Email  string `json:"email" validate:"required,email"`
	Active bool   `json:"active"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`
	Active *bool   `json:"active,omitempty"`
}

// UserResponse represents the response for user data
type UserResponse struct {
	ID        {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

{{if ne .Database.ID "redis"}}// CreateProductRequest represents the request to create a product
type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Available   bool    `json:"available"`
	UserID      {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"user_id" validate:"required"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=200"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	Available   *bool    `json:"available,omitempty"`
}

// ProductResponse represents the response for product data
type ProductResponse struct {
	ID          {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Available   bool      `json:"available"`
	UserID      {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"user_id"`
	User        *UserResponse `json:"user,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

{{end}}// ListRequest represents a generic list request with pagination
type ListRequest struct {
	Offset int    `json:"offset" form:"offset" validate:"min=0"`
	Limit  int    `json:"limit" form:"limit" validate:"min=1,max=100"`
	Search string `json:"search" form:"search"`
}

// APIResponse represents a standard API response
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// CreateAddressRequest represents the request to create an address
type CreateAddressRequest struct {
	Street  string `json:"street" validate:"required,min=5,max=200"`
	City    string `json:"city" validate:"required,min=2,max=100"`
	State   string `json:"state" validate:"required,min=2,max=100"`
	ZipCode string `json:"zip_code" validate:"required,min=5,max=20"`
	Country string `json:"country" validate:"required,min=2,max=100"`
	UserID  {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"user_id" validate:"required"`
}

// UpdateAddressRequest represents the request to update an address
type UpdateAddressRequest struct {
	Street  *string `json:"street,omitempty" validate:"omitempty,min=5,max=200"`
	City    *string `json:"city,omitempty" validate:"omitempty,min=2,max=100"`
	State   *string `json:"state,omitempty" validate:"omitempty,min=2,max=100"`
	ZipCode *string `json:"zip_code,omitempty" validate:"omitempty,min=5,max=20"`
	Country *string `json:"country,omitempty" validate:"omitempty,min=2,max=100"`
}

// AddressResponse represents the response for address data
type AddressResponse struct {
	ID        {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"id"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	ZipCode   string    `json:"zip_code"`
	Country   string    `json:"country"`
	UserID    {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateOrderRequest represents the request to create an order
type CreateOrderRequest struct {
	OrderNumber string  `json:"order_number" validate:"required,min=5,max=50"`
	Status      string  `json:"status" validate:"required,oneof=pending processing shipped delivered cancelled"`
	Total       float64 `json:"total" validate:"required,min=0"`
	UserID      {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"user_id" validate:"required"`
}

// UpdateOrderRequest represents the request to update an order
type UpdateOrderRequest struct {
	Status *string  `json:"status,omitempty" validate:"omitempty,oneof=pending processing shipped delivered cancelled"`
	Total  *float64 `json:"total,omitempty" validate:"omitempty,min=0"`
}

// OrderResponse represents the response for order data
type OrderResponse struct {
	ID          {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"id"`
	OrderNumber string    `json:"order_number"`
	Status      string    `json:"status"`
	Total       float64   `json:"total"`
	OrderDate   time.Time `json:"order_date"`
	UserID      {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"user_id"`
	User        *UserResponse `json:"user,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

// CategoryResponse represents the response for category data
type CategoryResponse struct {
	ID          {{if eq .DbDriver.ID "mongo-driver"}}string{{else if eq .DbDriver.ID "gorm"}}uint{{else}}int64{{end}} `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
