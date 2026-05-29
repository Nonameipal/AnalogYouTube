package repository

import (
	"database/sql"
	"errors"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) translateError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errs.ErrNotFound
	default:
		return err
	}
}
