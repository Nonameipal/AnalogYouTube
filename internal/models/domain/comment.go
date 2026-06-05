package domain

import "time"

const (
	CommentStatusActive  = "ACTIVE"
	CommentStatusDeleted = "DELETED"
)

type Comment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	VideoID   int       `json:"video_id"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
