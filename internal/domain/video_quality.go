package domain

import "time"

type VideoQuality struct {
	ID int `json:"id"`
	VideoID int `json:"video_id"`
	Quality string `json:"quality"`
	VideoURL string `json:"video_url"`
	CreatedAt time.Time `json:"created_at"`
}