package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user entity in the domain
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name      string `gorm:"not null" json:"name"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Active    bool   `gorm:"default:true" json:"active"`
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}

// Product represents a product entity in the domain
type Product struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `gorm:"not null;check:price >= 0" json:"price"`
	Available   bool    `gorm:"default:true" json:"available"`
	UserID      uint    `gorm:"not null" json:"user_id"`
	User        User    `gorm:"foreignKey:UserID" json:"user,omitempty"`

}

// TableName returns the table name for GORM
func (Product) TableName() string {
	return "products"
}

// Address represents an address entity
type Address struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Street   string `gorm:"not null" json:"street"`
	City     string `gorm:"not null" json:"city"`
	State    string `gorm:"not null" json:"state"`
	ZipCode  string `gorm:"not null" json:"zip_code"`
	Country  string `gorm:"not null" json:"country"`
	UserID   uint   `gorm:"not null" json:"user_id"`

}

// TableName returns the table name for GORM
func (Address) TableName() string {
	return "addresses"
}

// Order represents an order entity
type Order struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	OrderNumber string    `gorm:"uniqueIndex;not null" json:"order_number"`
	Status      string    `gorm:"not null" json:"status"`
	Total       float64   `gorm:"not null;check:total >= 0" json:"total"`
	OrderDate   time.Time `gorm:"not null" json:"order_date"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`

}

// TableName returns the table name for GORM
func (Order) TableName() string {
	return "orders"
}

// Category represents a product category entity
type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Active      bool   `gorm:"default:true" json:"active"`
}

// TableName returns the table name for GORM
func (Category) TableName() string {
	return "categories"
}
