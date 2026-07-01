package dbmodel

import (
	"database/sql"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"time"
)

type Playlist struct {
	ID int `db:"id"`
	Name string `db:"name"`
	UserID int `db:"user_id"`
	Description sql.NullString `db:"description"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (p Playlist) ToDomain() domain.Playlist {
	return domain.Playlist{
		ID: p.ID,
		Name: p.Name,
		UserID: p.UserID,
		Description: p.Description.String,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
