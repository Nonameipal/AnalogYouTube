package db

import (
	"database/sql"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

type Video struct {
	ID           int            `db:"id"`
	AuthorID     int            `db:"author_id"`
	Title        string         `db:"title"`
	Description  sql.NullString `db:"description"`
	VideoURL     string         `db:"video_url"`
	ThumbnailURL sql.NullString `db:"thumbnail_url"`
	Views        int64          `db:"views"`
	Status       string         `db:"status"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

func (v Video) ToDomain() domain.Video {
	return domain.Video{
		ID:           v.ID,
		AuthorID:     v.AuthorID,
		Title:        v.Title,
		Description:  v.Description.String,
		VideoURL:     v.VideoURL,
		ThumbnailURL: v.ThumbnailURL.String,
		Views:        v.Views,
		Status:       v.Status,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}
