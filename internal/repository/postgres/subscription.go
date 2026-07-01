package postgres

import (
	"context"
	"database/sql"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

func (r *Repository) GetSubscribers(authorID int) ([]domain.User, error) {
	ctx := context.Background()

	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.username, u.email, u.role, u.avatar_url, u.description, u.created_at, u.updated_at
		FROM subscriptions s
		JOIN users u ON u.id = s.subscriber_id
		WHERE s.author_id = $1
		ORDER BY s.created_at DESC`, authorID)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		var email sql.NullString
		var avatarURL sql.NullString
		var description sql.NullString

		err = rows.Scan(
			&user.ID,
			&user.Username,
			&email,
			&user.Role,
			&avatarURL,
			&description,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, r.translateError(err)
		}

		user.Email = email.String
		user.AvatarURL = avatarURL.String
		user.Description = description.String

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return users, nil
}

func (r *Repository) GetSubscriptions(subscriberID int) ([]domain.User, error) {
	ctx := context.Background()

	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.username, u.email, u.role, u.avatar_url, u.description, u.created_at, u.updated_at
		FROM subscriptions s
		JOIN users u ON u.id = s.author_id
		WHERE s.subscriber_id = $1
		ORDER BY s.created_at DESC`, subscriberID)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		var email sql.NullString
		var avatarURL sql.NullString
		var description sql.NullString

		err = rows.Scan(
			&user.ID,
			&user.Username,
			&email,
			&user.Role,
			&avatarURL,
			&description,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, r.translateError(err)
		}

		user.Email = email.String
		user.AvatarURL = avatarURL.String
		user.Description = description.String

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return users, nil
}

func (r *Repository) SubscribeToUser(subscriberID int, authorID int) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`INSERT INTO subscriptions (subscriber_id, author_id)
		VALUES ($1, $2)
		ON CONFLICT (subscriber_id, author_id) DO NOTHING`, subscriberID, authorID)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) UnsubscribeFromUser(subscriberID int, authorID int) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`DELETE FROM subscriptions
		WHERE subscriber_id = $1 AND author_id = $2`, subscriberID, authorID)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) GetSubscribersCount(authorID int) (int, error) {
	ctx := context.Background()
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		FROM subscriptions
		WHERE author_id = $1`, authorID).Scan(&count)
	if err != nil {
		return 0, r.translateError(err)
	}

	return count, nil
}

func (r *Repository) GetSubscriptionsCount(subscriberID int) (int, error) {
	ctx := context.Background()
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		FROM subscriptions
		WHERE subscriber_id = $1`, subscriberID).Scan(&count)
	if err != nil {
		return 0, r.translateError(err)
	}

	return count, nil
}

func (r *Repository) IsSubscribed(subscriberID int, authorID int) (bool, error) {
	ctx := context.Background()
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM subscriptions WHERE subscriber_id = $1 AND author_id = $2
		)`, subscriberID, authorID).Scan(&exists)
	if err != nil {
		return false, r.translateError(err)
	}

	return exists, nil
}
