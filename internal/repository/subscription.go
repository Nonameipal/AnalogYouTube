package repository

func (r *Repository) SubscribeToUser(subscriberID int, authorID int) error {
	_, err := r.db.Exec(`
		INSERT INTO subscriptions (subscriber_id, author_id)
		VALUES ($1, $2)
		ON CONFLICT (subscriber_id, author_id) DO NOTHING`, subscriberID, authorID)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) UnsubscribeFromUser(subscriberID int, authorID int) error {
	_, err := r.db.Exec(`
		DELETE FROM subscriptions
		WHERE subscriber_id = $1 AND author_id = $2`, subscriberID, authorID)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) GetSubscribersCount(authorID int) (int, error) {
	var count int
	if err := r.db.Get(&count, `
		SELECT COUNT(*)
		FROM subscriptions
		WHERE author_id = $1`, authorID); err != nil {
		return 0, r.translateError(err)
	}

	return count, nil
}

func (r *Repository) GetSubscriptionsCount(subscriberID int) (int, error) {
	var count int
	if err := r.db.Get(&count, `
		SELECT COUNT(*)
		FROM subscriptions
		WHERE subscriber_id = $1`, subscriberID); err != nil {
		return 0, r.translateError(err)
	}

	return count, nil
}

func (r *Repository) IsSubscribed(subscriberID int, authorID int) (bool, error) {
	var exists bool
	if err := r.db.Get(&exists, `
		SELECT EXISTS(
			SELECT 1 FROM subscriptions WHERE subscriber_id = $1 AND author_id = $2
		)`, subscriberID, authorID); err != nil {
		return false, r.translateError(err)
	}

	return exists, nil
}
