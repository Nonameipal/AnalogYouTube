package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/repository/postgres/dbmodel"
)

func (r *Repository) CreateOrGetChat(firstUserID int, secondUserID int) (domain.Chat, error) {
	ctx := context.Background()
	var dbChat dbModels.Chat
	err := r.db.QueryRow(ctx,
		`INSERT INTO chats (first_user_id, second_user_id)
		VALUES ($1, $2)
		ON CONFLICT (LEAST(first_user_id, second_user_id), GREATEST(first_user_id, second_user_id))
		DO UPDATE SET first_user_id = chats.first_user_id
		RETURNING id, first_user_id, second_user_id, created_at`, firstUserID, secondUserID).Scan(
		&dbChat.ID, &dbChat.FirstUserID, &dbChat.SecondUserID, &dbChat.CreatedAt)
	if err != nil {
		return domain.Chat{}, r.translateError(err)
	}

	return dbChat.ToDomain(), nil
}

func (r *Repository) GetChatByID(chatID int) (domain.Chat, error) {
	ctx := context.Background()
	var dbChat dbModels.Chat
	err := r.db.QueryRow(ctx,
		`SELECT id, first_user_id, second_user_id, created_at
		FROM chats
		WHERE id = $1`, chatID).Scan(&dbChat.ID, &dbChat.FirstUserID, &dbChat.SecondUserID, &dbChat.CreatedAt)
	if err != nil {
		return domain.Chat{}, r.translateError(err)
	}

	return dbChat.ToDomain(), nil
}

func (r *Repository) GetUserChats(userID int) ([]domain.Chat, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, first_user_id, second_user_id, created_at
		FROM chats
		WHERE first_user_id = $1 OR second_user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbChats []dbModels.Chat
	for rows.Next() {
		var chat dbModels.Chat
		if err := rows.Scan(&chat.ID, &chat.FirstUserID, &chat.SecondUserID, &chat.CreatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbChats = append(dbChats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	chats := make([]domain.Chat, 0, len(dbChats))
	for _, chat := range dbChats {
		chats = append(chats, chat.ToDomain())
	}

	return chats, nil
}

func (r *Repository) CreateChatMessage(message domain.ChatMessage) (domain.ChatMessage, error) {
	ctx := context.Background()
	var dbMessage dbModels.ChatMessage
	err := r.db.QueryRow(ctx,
		`INSERT INTO chat_messages (chat_id, sender_id, text)
		VALUES ($1, $2, $3)
		RETURNING id, chat_id, sender_id, text, created_at`,
		message.ChatID,
		message.SenderID,
		message.Text,
	).Scan(&dbMessage.ID, &dbMessage.ChatID, &dbMessage.SenderID, &dbMessage.Text, &dbMessage.CreatedAt)
	if err != nil {
		return domain.ChatMessage{}, r.translateError(err)
	}

	return dbMessage.ToDomain(), nil
}

func (r *Repository) GetChatMessages(chatID int) ([]domain.ChatMessage, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, chat_id, sender_id, text, created_at
		FROM chat_messages
		WHERE chat_id = $1
		ORDER BY created_at ASC`, chatID)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbMessages []dbModels.ChatMessage
	for rows.Next() {
		var message dbModels.ChatMessage
		if err := rows.Scan(&message.ID, &message.ChatID, &message.SenderID, &message.Text, &message.CreatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbMessages = append(dbMessages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	messages := make([]domain.ChatMessage, 0, len(dbMessages))
	for _, message := range dbMessages {
		messages = append(messages, message.ToDomain())
	}

	return messages, nil
}
