package domain

import "time"

type Chat struct {
	ID int `json:"id"`
	FirstUserID int `json:"first_user_id"`
	SecondUserID int `json:"second_user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMessage struct {
	ID int `json:"id"`
	ChatID int `json:"chat_id"`
	SenderID int `json:"sender_id"`
	Text string `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
