package domain

import "time"

const (
	VideoStatusActive  = "ACTIVE"
	VideoStatusDeleted = "DELETED"
)

type Video struct {
	ID           int            `json:"id"`
	AuthorID     int            `json:"author_id"`
	CategoryName *string        `json:"category_name"`
	CategoryID   *int           `json:"category_id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	VideoURL     string         `json:"video_url"`
	ThumbnailURL string         `json:"thumbnail_url"`
	Views        int64          `json:"views"`
	Status       string         `json:"status"`
	Qualities    []VideoQuality `json:"qualities,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
