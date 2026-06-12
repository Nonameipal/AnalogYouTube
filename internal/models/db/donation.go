package db

import (
	"database/sql"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

type Donation struct {
	ID int `db:"id"`
	SenderID int `db:"sender_id"`
	ReceiverID int `db:"receiver_id"`
	VideoID sql.NullInt64 `db:"video_id"`
	Amount float64 `db:"amount"`
	Message sql.NullString `db:"message"`
	CreatedAt  time.Time  `db:"created_at"`
}

func (d Donation) ToDomain() domain.Donation {
	var videoID *int
	if d.VideoID.Valid {
		value := int(d.VideoID.Int64)
		videoID = &value
	}

	return domain.Donation{
		ID:         d.ID,
		SenderID:   d.SenderID,
		ReceiverID: d.ReceiverID,
		VideoID:    videoID,
		Amount:     d.Amount,
		Message:    d.Message.String,
		CreatedAt:  d.CreatedAt,
	}
}
