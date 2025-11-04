package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"{{.ProjectName}}/internal/gen"
{{if eq .DbDriver.ID "gorm"}}	"gorm.io/gorm"
{{else if eq .DbDriver.ID "sqlx"}}	"github.com/jmoiron/sqlx"
{{else}}	"database/sql"
{{end}}
)

type User struct {
	ID        uint64     `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	CreatedAt *time.Time `json:"created_at"`
}

type UserSession struct {
	SessionID string `json:"session_id"`
	UserID    uint64 `json:"user_id"`
}

{{if eq .DbDriver.ID "gorm"}}// {{.Database.ID | title}}UserStore implements user storage using GORM
type {{.Database.ID | title}}UserStore struct {
	db *gorm.DB
}

func New{{.Database.ID | title}}UserStore(db *gorm.DB) *{{.Database.ID | title}}UserStore {
	return &{{.Database.ID | title}}UserStore{db: db}
}

func (p *PostgresUserStore) Create(ctx context.Context, u *User) (*User, error) {
	val, err := gen.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = val

	if err := p.db.WithContext(ctx).Create(u).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, fmt.Errorf("duplicate row: %v", err.Error())
		}
		return nil, err
	}

	return u, nil
}

func (p *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := p.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (p *PostgresUserStore) CreateUserSession(ctx context.Context, us *UserSession) error {
	if err := p.db.WithContext(ctx).Create(us).Error; err != nil {
		return err
	}
	return nil
}

{{else if eq .DbDriver.ID "sqlx"}}// {{.Database.ID | title}}UserStore implements user storage using sqlx
type {{.Database.ID | title}}UserStore struct {
	db *sqlx.DB
}

func New{{.Database.ID | title}}UserStore(db *sqlx.DB) *{{.Database.ID | title}}UserStore {
	return &{{.Database.ID | title}}UserStore{db: db}
}

func (p *PostgresUserStore) Create(ctx context.Context, u *User) (*User, error) {
	val, err := gen.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = val

{{if eq .Database.ID "postgres"}}	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	var id uint64
	err = p.db.QueryRowContext(ctx, query, u.Email, u.Password).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, fmt.Errorf("duplicate row: %v", err.Error())
		}
		return nil, err
	}
	u.ID = id
{{else}}	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	result, err := p.db.ExecContext(ctx, query, u.Email, u.Password)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, fmt.Errorf("duplicate row: %v", err.Error())
		}
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = uint64(id)
{{end}}	return u, nil
}

func (p *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
{{if eq .Database.ID "postgres"}}	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`
{{else}}	query := `SELECT id, email, password, created_at FROM users WHERE email = ?`
{{end}}	err := p.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresUserStore) CreateUserSession(ctx context.Context, us *UserSession) error {
{{if eq .Database.ID "postgres"}}	query := `INSERT INTO user_session (session_id, user_id) VALUES($1, $2)`
{{else}}	query := `INSERT INTO user_session (session_id, user_id) VALUES(?, ?)`
{{end}}	_, err := p.db.ExecContext(ctx, query, us.SessionID, us.UserID)
	if err != nil {
		return err
	}
	return nil
}

{{else}}// {{.Database.ID | title}}UserStore implements user storage using database/sql
type {{.Database.ID | title}}UserStore struct {
	db *sql.DB
}

func New{{.Database.ID | title}}UserStore(db *sql.DB) *{{.Database.ID | title}}UserStore {
	return &{{.Database.ID | title}}UserStore{db: db}
}

func (p *PostgresUserStore) Create(ctx context.Context, u *User) (*User, error) {
	val, err := gen.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = val

{{if eq .Database.ID "postgres"}}	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	var id uint64
	err = p.db.QueryRowContext(ctx, query, u.Email, u.Password).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, fmt.Errorf("duplicate row: %v", err.Error())
		}
		return nil, err
	}
	u.ID = id
{{else if eq .Database.ID "mysql"}}	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	result, err := p.db.ExecContext(ctx, query, u.Email, u.Password)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, fmt.Errorf("duplicate row: %v", err.Error())
		}
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = uint64(id)
{{else if eq .Database.ID "sqlite"}}	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	result, err := p.db.ExecContext(ctx, query, u.Email, u.Password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, fmt.Errorf("duplicate row: %v", err.Error())
		}
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = uint64(id)
{{end}}	return u, nil
}

func (p *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
{{if eq .Database.ID "postgres"}}	query := `SELECT id, email, password, created_at FROM users WHERE email = $1`
{{else}}	query := `SELECT id, email, password, created_at FROM users WHERE email = ?`
{{end}}	err := p.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *PostgresUserStore) CreateUserSession(ctx context.Context, us *UserSession) error {
{{if eq .Database.ID "postgres"}}	query := `INSERT INTO user_session (session_id, user_id) VALUES($1, $2)`
{{else}}	query := `INSERT INTO user_session (session_id, user_id) VALUES(?, ?)`
{{end}}	_, err := p.db.ExecContext(ctx, query, us.SessionID, us.UserID)
	if err != nil {
		return err
	}
	return nil
}
{{end}}
