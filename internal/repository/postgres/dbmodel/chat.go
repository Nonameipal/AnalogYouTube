package dbmodel

import (
	"database/sql"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

type Chat struct {
	ID           int       `db:"id"`
	FirstUserID  int       `db:"first_user_id"`
	SecondUserID int       `db:"second_user_id"`
	CreatedAt    time.Time `db:"created_at"`
}

func (c Chat) ToDomain() domain.Chat {
	return domain.Chat{
		ID:           c.ID,
		FirstUserID:  c.FirstUserID,
		SecondUserID: c.SecondUserID,
		CreatedAt:    c.CreatedAt,
	}
}

type ChatRequest struct {
	ID         int           `db:"id"`
	SenderID   int           `db:"sender_id"`
	ReceiverID int           `db:"receiver_id"`
	Status     string        `db:"status"`
	ChatID     sql.NullInt64 `db:"chat_id"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
}

func (r ChatRequest) ToDomain() domain.ChatRequest {
	var chatID *int
	if r.ChatID.Valid {
		value := int(r.ChatID.Int64)
		chatID = &value
	}

	return domain.ChatRequest{
		ID:         r.ID,
		SenderID:   r.SenderID,
		ReceiverID: r.ReceiverID,
		Status:     r.Status,
		ChatID:     chatID,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

type ChatMessage struct {
	ID        int       `db:"id"`
	ChatID    int       `db:"chat_id"`
	SenderID  int       `db:"sender_id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}

func (m ChatMessage) ToDomain() domain.ChatMessage {
	return domain.ChatMessage{
		ID:        m.ID,
		ChatID:    m.ChatID,
		SenderID:  m.SenderID,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
	}
}
