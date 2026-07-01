package domain

import "time"

type Donation struct {
	ID int `json:"id"`
	SenderID int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
	VideoID *int `json:"video_id"`
	Amount float64 `json:"amount"`
	Message string `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
