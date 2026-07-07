package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/repository/postgres/dbmodel"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetChatByID(chatID int) (domain.Chat, error) {
	ctx := context.Background()
	var dbChat dbmodel.Chat
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

	var dbChats []dbmodel.Chat
	for rows.Next() {
		var chat dbmodel.Chat
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

func (r *Repository) GetChatBetweenUsers(firstUserID int, secondUserID int) (domain.Chat, error) {
	ctx := context.Background()
	var dbChat dbmodel.Chat
	err := r.db.QueryRow(ctx,
		`SELECT id, first_user_id, second_user_id, created_at
		FROM chats
		WHERE LEAST(first_user_id, second_user_id) = LEAST($1, $2)
		AND GREATEST(first_user_id, second_user_id) = GREATEST($1, $2)`,
		firstUserID,
		secondUserID,
	).Scan(&dbChat.ID, &dbChat.FirstUserID, &dbChat.SecondUserID, &dbChat.CreatedAt)
	if err != nil {
		return domain.Chat{}, r.translateError(err)
	}

	return dbChat.ToDomain(), nil
}

func (r *Repository) CreateChatRequest(request domain.ChatRequest) (domain.ChatRequest, error) {
	ctx := context.Background()
	var dbRequest dbmodel.ChatRequest
	err := r.db.QueryRow(ctx,
		`INSERT INTO chat_requests (sender_id, receiver_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, sender_id, receiver_id, status, chat_id, created_at, updated_at`,
		request.SenderID,
		request.ReceiverID,
		domain.ChatRequestStatusPending,
	).Scan(
		&dbRequest.ID,
		&dbRequest.SenderID,
		&dbRequest.ReceiverID,
		&dbRequest.Status,
		&dbRequest.ChatID,
		&dbRequest.CreatedAt,
		&dbRequest.UpdatedAt,
	)
	if err != nil {
		return domain.ChatRequest{}, r.translateError(err)
	}

	return dbRequest.ToDomain(), nil
}

func (r *Repository) GetChatRequestByID(requestID int) (domain.ChatRequest, error) {
	ctx := context.Background()
	var dbRequest dbmodel.ChatRequest
	err := r.db.QueryRow(ctx,
		`SELECT id, sender_id, receiver_id, status, chat_id, created_at, updated_at
		FROM chat_requests
		WHERE id = $1`, requestID).Scan(
		&dbRequest.ID,
		&dbRequest.SenderID,
		&dbRequest.ReceiverID,
		&dbRequest.Status,
		&dbRequest.ChatID,
		&dbRequest.CreatedAt,
		&dbRequest.UpdatedAt,
	)
	if err != nil {
		return domain.ChatRequest{}, r.translateError(err)
	}

	return dbRequest.ToDomain(), nil
}

func (r *Repository) GetIncomingChatRequests(userID int) ([]domain.ChatRequest, error) {
	return r.getChatRequests(
		`SELECT id, sender_id, receiver_id, status, chat_id, created_at, updated_at
		FROM chat_requests
		WHERE receiver_id = $1 AND status = $2
		ORDER BY created_at DESC`,
		userID,
		domain.ChatRequestStatusPending,
	)
}

func (r *Repository) GetOutgoingChatRequests(userID int) ([]domain.ChatRequest, error) {
	return r.getChatRequests(
		`SELECT id, sender_id, receiver_id, status, chat_id, created_at, updated_at
		FROM chat_requests
		WHERE sender_id = $1 AND status = $2
		ORDER BY created_at DESC`,
		userID,
		domain.ChatRequestStatusPending,
	)
}

func (r *Repository) AcceptChatRequest(requestID int, firstUserID int, secondUserID int) (domain.ChatRequest, domain.Chat, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.ChatRequest{}, domain.Chat{}, err
	}
	defer tx.Rollback(ctx)

	chat, err := r.createOrGetChatTx(ctx, tx, firstUserID, secondUserID)
	if err != nil {
		return domain.ChatRequest{}, domain.Chat{}, r.translateError(err)
	}

	var dbRequest dbmodel.ChatRequest
	err = tx.QueryRow(ctx,
		`UPDATE chat_requests
		SET status = $1, chat_id = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND status = $4
		RETURNING id, sender_id, receiver_id, status, chat_id, created_at, updated_at`,
		domain.ChatRequestStatusAccepted,
		chat.ID,
		requestID,
		domain.ChatRequestStatusPending,
	).Scan(
		&dbRequest.ID,
		&dbRequest.SenderID,
		&dbRequest.ReceiverID,
		&dbRequest.Status,
		&dbRequest.ChatID,
		&dbRequest.CreatedAt,
		&dbRequest.UpdatedAt,
	)
	if err != nil {
		return domain.ChatRequest{}, domain.Chat{}, r.translateError(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return domain.ChatRequest{}, domain.Chat{}, err
	}

	return dbRequest.ToDomain(), chat, nil
}

func (r *Repository) RejectChatRequest(requestID int) (domain.ChatRequest, error) {
	ctx := context.Background()
	var dbRequest dbmodel.ChatRequest
	err := r.db.QueryRow(ctx,
		`UPDATE chat_requests
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND status = $3
		RETURNING id, sender_id, receiver_id, status, chat_id, created_at, updated_at`,
		domain.ChatRequestStatusRejected,
		requestID,
		domain.ChatRequestStatusPending,
	).Scan(
		&dbRequest.ID,
		&dbRequest.SenderID,
		&dbRequest.ReceiverID,
		&dbRequest.Status,
		&dbRequest.ChatID,
		&dbRequest.CreatedAt,
		&dbRequest.UpdatedAt,
	)
	if err != nil {
		return domain.ChatRequest{}, r.translateError(err)
	}

	return dbRequest.ToDomain(), nil
}

func (r *Repository) CreateChatMessage(message domain.ChatMessage) (domain.ChatMessage, error) {
	ctx := context.Background()
	var dbMessage dbmodel.ChatMessage
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

	var dbMessages []dbmodel.ChatMessage
	for rows.Next() {
		var message dbmodel.ChatMessage
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

func (r *Repository) createOrGetChatTx(ctx context.Context, tx pgx.Tx, firstUserID int, secondUserID int) (domain.Chat, error) {
	var dbChat dbmodel.Chat
	err := tx.QueryRow(ctx,
		`INSERT INTO chats (first_user_id, second_user_id)
		VALUES ($1, $2)
		ON CONFLICT (LEAST(first_user_id, second_user_id), GREATEST(first_user_id, second_user_id))
		DO UPDATE SET first_user_id = chats.first_user_id
		RETURNING id, first_user_id, second_user_id, created_at`,
		firstUserID,
		secondUserID,
	).Scan(&dbChat.ID, &dbChat.FirstUserID, &dbChat.SecondUserID, &dbChat.CreatedAt)
	if err != nil {
		return domain.Chat{}, err
	}

	return dbChat.ToDomain(), nil
}

func (r *Repository) getChatRequests(query string, args ...any) ([]domain.ChatRequest, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbRequests []dbmodel.ChatRequest
	for rows.Next() {
		var request dbmodel.ChatRequest
		if err := rows.Scan(
			&request.ID,
			&request.SenderID,
			&request.ReceiverID,
			&request.Status,
			&request.ChatID,
			&request.CreatedAt,
			&request.UpdatedAt,
		); err != nil {
			return nil, r.translateError(err)
		}
		dbRequests = append(dbRequests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	requests := make([]domain.ChatRequest, 0, len(dbRequests))
	for _, request := range dbRequests {
		requests = append(requests, request.ToDomain())
	}

	return requests, nil
}
