package dbmodel

import (
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

type VideoQuality struct {
	ID int `db:"id"`
	VideoID int `db:"video_id"`
	Quality string `db:"quality"`
	VideoURL string `db:"video_url"`
	CreatedAt time.Time `db:"created_at"`
}

func (q VideoQuality) ToDomain() domain.VideoQuality {
	return domain.VideoQuality{
		ID: q.ID,
		VideoID: q.VideoID,
		Quality: q.Quality,
		VideoURL: q.VideoURL,
		CreatedAt: q.CreatedAt,
	}
}