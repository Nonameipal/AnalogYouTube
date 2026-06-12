package db

import (
	"database/sql"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

type Category struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c Category) ToDomain() domain.Category {
	return domain.Category{
		ID: c.ID,
		Name: c.Name,
		Description: c.Description.String,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
