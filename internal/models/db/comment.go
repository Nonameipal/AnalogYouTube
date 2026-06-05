package db

import (
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

type Comment struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	VideoID   int       `db:"video_id"`
	Text      string    `db:"text"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c Comment) ToDomain() domain.Comment {
	return domain.Comment{
		ID:        c.ID,
		UserID:    c.UserID,
		VideoID:   c.VideoID,
		Text:      c.Text,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
