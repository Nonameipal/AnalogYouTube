package postgres

import (
	"errors"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) translateError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.ErrNotFound
	default:
		return err
	}
}
