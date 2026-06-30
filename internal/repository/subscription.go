package repository

import "context"

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
