package domain

import (
	"time"
{{if eq .DbDriver.ID "gorm"}}
	"gorm.io/gorm"
{{else if eq .DbDriver.ID "mongo-driver"}}
	"go.mongodb.org/mongo-driver/bson/primitive"
{{end}})

// User represents a user entity in the domain
type User struct {
{{if eq .DbDriver.ID "gorm"}}	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
{{else if eq .DbDriver.ID "mongo-driver"}}	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
{{else}}	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
{{end}}
	Name      string `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"name" {{else}}db:"name" {{end}}json:"name"`
	Email     string `{{if eq .DbDriver.ID "gorm"}}gorm:"uniqueIndex;not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"email" {{else}}db:"email" {{end}}json:"email"`
	Active    bool   `{{if eq .DbDriver.ID "gorm"}}gorm:"default:true" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"active" {{else}}db:"active" {{end}}json:"active"`
}

{{if eq .DbDriver.ID "gorm"}}// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}{{end}}

{{if ne .Database.ID "redis"}}// Product represents a product entity in the domain
type Product struct {
{{if eq .DbDriver.ID "gorm"}}	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
{{else if eq .DbDriver.ID "mongo-driver"}}	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
{{else}}	ID          int64     `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
{{end}}
	Name        string  `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"name" {{else}}db:"name" {{end}}json:"name"`
	Description string  `{{if eq .DbDriver.ID "gorm"}}gorm:"type:text" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"description" {{else}}db:"description" {{end}}json:"description"`
	Price       float64 `{{if eq .DbDriver.ID "gorm"}}gorm:"not null;check:price >= 0" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"price" {{else}}db:"price" {{end}}json:"price"`
	Available   bool    `{{if eq .DbDriver.ID "gorm"}}gorm:"default:true" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"available" {{else}}db:"available" {{end}}json:"available"`
{{if eq .DbDriver.ID "gorm"}}	UserID      uint    `gorm:"not null" json:"user_id"`
	User        User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
{{else if eq .DbDriver.ID "mongo-driver"}}	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
{{else}}	UserID      int64   `db:"user_id" json:"user_id"`
{{end}}
}

{{if eq .DbDriver.ID "gorm"}}// TableName returns the table name for GORM
func (Product) TableName() string {
	return "products"
}{{end}}

{{end}}// Address represents an address entity
type Address struct {
{{if eq .DbDriver.ID "gorm"}}	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
{{else if eq .DbDriver.ID "mongo-driver"}}	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
{{else}}	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
{{end}}
	Street   string `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"street" {{else}}db:"street" {{end}}json:"street"`
	City     string `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"city" {{else}}db:"city" {{end}}json:"city"`
	State    string `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"state" {{else}}db:"state" {{end}}json:"state"`
	ZipCode  string `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"zip_code" {{else}}db:"zip_code" {{end}}json:"zip_code"`
	Country  string `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"country" {{else}}db:"country" {{end}}json:"country"`
{{if eq .DbDriver.ID "gorm"}}	UserID   uint   `gorm:"not null" json:"user_id"`
{{else if eq .DbDriver.ID "mongo-driver"}}	UserID   primitive.ObjectID `bson:"user_id" json:"user_id"`
{{else}}	UserID   int64  `db:"user_id" json:"user_id"`
{{end}}
}

{{if eq .DbDriver.ID "gorm"}}// TableName returns the table name for GORM
func (Address) TableName() string {
	return "addresses"
}{{end}}

// Order represents an order entity
type Order struct {
{{if eq .DbDriver.ID "gorm"}}	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
{{else if eq .DbDriver.ID "mongo-driver"}}	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
{{else}}	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
{{end}}
	OrderNumber string    `{{if eq .DbDriver.ID "gorm"}}gorm:"uniqueIndex;not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"order_number" {{else}}db:"order_number" {{end}}json:"order_number"`
	Status      string    `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"status" {{else}}db:"status" {{end}}json:"status"`
	Total       float64   `{{if eq .DbDriver.ID "gorm"}}gorm:"not null;check:total >= 0" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"total" {{else}}db:"total" {{end}}json:"total"`
	OrderDate   time.Time `{{if eq .DbDriver.ID "gorm"}}gorm:"not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"order_date" {{else}}db:"order_date" {{end}}json:"order_date"`
{{if eq .DbDriver.ID "gorm"}}	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
{{else if eq .DbDriver.ID "mongo-driver"}}	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
{{else}}	UserID      int64     `db:"user_id" json:"user_id"`
{{end}}
}

{{if eq .DbDriver.ID "gorm"}}// TableName returns the table name for GORM
func (Order) TableName() string {
	return "orders"
}{{end}}

// Category represents a product category entity
type Category struct {
{{if eq .DbDriver.ID "gorm"}}	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
{{else if eq .DbDriver.ID "mongo-driver"}}	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
{{else}}	ID        int64     `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
{{end}}
	Name        string `{{if eq .DbDriver.ID "gorm"}}gorm:"uniqueIndex;not null" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"name" {{else}}db:"name" {{end}}json:"name"`
	Description string `{{if eq .DbDriver.ID "gorm"}}gorm:"type:text" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"description" {{else}}db:"description" {{end}}json:"description"`
	Active      bool   `{{if eq .DbDriver.ID "gorm"}}gorm:"default:true" {{else if eq .DbDriver.ID "mongo-driver"}}bson:"active" {{else}}db:"active" {{end}}json:"active"`
}

{{if eq .DbDriver.ID "gorm"}}// TableName returns the table name for GORM
func (Category) TableName() string {
	return "categories"
}{{end}}
