package domain

import "time"

type VideoWatchHistory struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	VideoID         int       `json:"video_id"`
	WatchedSeconds  int       `json:"watched_seconds"`
	DurationSeconds int       `json:"duration_seconds"`
	WatchedPercent  float64   `json:"watched_percent"`
	IsCompleted     bool      `json:"is_completed"`
	LastWatchedAt   time.Time `json:"last_watched_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type VideoSearchHistory struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Query     string    `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}
