package domain

import "time"

const (
	ChatRequestStatusPending  = "pending"
	ChatRequestStatusAccepted = "accepted"
	ChatRequestStatusRejected = "rejected"
)

type Chat struct {
	ID           int       `json:"id"`
	FirstUserID  int       `json:"first_user_id"`
	SecondUserID int       `json:"second_user_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type ChatRequest struct {
	ID         int       `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Status     string    `json:"status"`
	ChatID     *int      `json:"chat_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AcceptedChatRequest struct {
	Request ChatRequest `json:"request"`
	Chat    Chat        `json:"chat"`
}

type ChatMessage struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	SenderID  int       `json:"sender_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
