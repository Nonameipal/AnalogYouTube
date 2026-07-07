package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

func (r *Repository) SaveVideoWatchProgress(progress domain.VideoWatchHistory) (domain.VideoWatchHistory, error) {
	ctx := context.Background()
	var saved domain.VideoWatchHistory

	err := r.db.QueryRow(ctx,
		`INSERT INTO video_watch_history (
			user_id, video_id, watched_seconds, duration_seconds, watched_percent, is_completed, last_watched_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, video_id) DO UPDATE
		SET watched_seconds = EXCLUDED.watched_seconds,
		    duration_seconds = EXCLUDED.duration_seconds,
		    watched_percent = GREATEST(video_watch_history.watched_percent, EXCLUDED.watched_percent),
		    is_completed = video_watch_history.is_completed OR EXCLUDED.is_completed,
		    last_watched_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, video_id, watched_seconds, duration_seconds, watched_percent, is_completed,
		          last_watched_at, created_at, updated_at`,
		progress.UserID,
		progress.VideoID,
		progress.WatchedSeconds,
		progress.DurationSeconds,
		progress.WatchedPercent,
		progress.IsCompleted,
	).Scan(
		&saved.ID,
		&saved.UserID,
		&saved.VideoID,
		&saved.WatchedSeconds,
		&saved.DurationSeconds,
		&saved.WatchedPercent,
		&saved.IsCompleted,
		&saved.LastWatchedAt,
		&saved.CreatedAt,
		&saved.UpdatedAt,
	)
	if err != nil {
		return domain.VideoWatchHistory{}, r.translateError(err)
	}

	return saved, nil
}

func (r *Repository) GetUserWatchHistory(userID int, limit int) ([]domain.VideoWatchHistory, error) {
	if limit <= 0 {
		limit = 200
	}

	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, video_id, watched_seconds, duration_seconds, watched_percent, is_completed,
		        last_watched_at, created_at, updated_at
		FROM video_watch_history
		WHERE user_id = $1
		ORDER BY last_watched_at DESC
		LIMIT $2`,
		userID,
		limit,
	)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	history := make([]domain.VideoWatchHistory, 0)
	for rows.Next() {
		var item domain.VideoWatchHistory
		if err = rows.Scan(
			&item.ID,
			&item.UserID,
			&item.VideoID,
			&item.WatchedSeconds,
			&item.DurationSeconds,
			&item.WatchedPercent,
			&item.IsCompleted,
			&item.LastWatchedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, r.translateError(err)
		}
		history = append(history, item)
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return history, nil
}

func (r *Repository) SaveVideoSearchHistory(history domain.VideoSearchHistory) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`INSERT INTO video_search_history (user_id, query)
		VALUES ($1, $2)`,
		history.UserID,
		history.Query,
	)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) GetUserSearchHistory(userID int, limit int) ([]domain.VideoSearchHistory, error) {
	if limit <= 0 {
		limit = 100
	}

	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, query, created_at
		FROM video_search_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`,
		userID,
		limit,
	)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	history := make([]domain.VideoSearchHistory, 0)
	for rows.Next() {
		var item domain.VideoSearchHistory
		if err = rows.Scan(&item.ID, &item.UserID, &item.Query, &item.CreatedAt); err != nil {
			return nil, r.translateError(err)
		}
		history = append(history, item)
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return history, nil
}
