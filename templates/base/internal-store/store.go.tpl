package store

import (
	"context"
	"errors"
	"time"
{{if eq .DbDriver.ID "gorm"}}	"gorm.io/gorm"
{{else if eq .DbDriver.ID "sqlx"}}	"database/sql"
	"github.com/jmoiron/sqlx"
{{else}}	"database/sql"
{{end}}
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record conflict")

	QueryTimeoutDuration = 5 * time.Second
)

type Store struct {
	Users interface {
		Create(ctx context.Context, u *User) (*User, error)
		GetByEmail(ctx context.Context, email string) (*User, error)
		CreateUserSession(ctx context.Context, us *UserSession) error
	}
}

{{if eq .DbDriver.ID "gorm"}}func NewStorage(db *gorm.DB) Store {
	return Store{
		Users: New{{.Database.ID | title}}UserStore(db),
	}
}
{{else if eq .DbDriver.ID "sqlx"}}func NewStorage(db *sqlx.DB) Store {
	return Store{
		Users: New{{.Database.ID | title}}UserStore(db),
	}
}
{{else}}func NewStorage(db *sql.DB) Store {
	return Store{
		Users: New{{.Database.ID | title}}UserStore(db),
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		err = tx.Rollback()
		if err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}
{{end}}

