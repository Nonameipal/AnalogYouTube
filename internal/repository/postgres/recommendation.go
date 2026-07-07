package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

func (r *Repository) GetUserLikedVideoIDs(userID int) ([]int, error) {
	return r.getUserVideoIDs(
		`SELECT video_id
		FROM video_likes
		WHERE user_id = $1
		ORDER BY created_at DESC`,
		userID,
	)
}

func (r *Repository) GetUserCommentedVideoIDs(userID int) ([]int, error) {
	return r.getUserVideoIDs(
		`SELECT DISTINCT video_id
		FROM comments
		WHERE user_id = $1 AND status = $2`,
		userID,
		domain.CommentStatusActive,
	)
}

func (r *Repository) GetUserPlaylistVideoIDs(userID int) ([]int, error) {
	return r.getUserVideoIDs(
		`SELECT DISTINCT pv.video_id
		FROM playlist_videos pv
		JOIN playlists p ON p.id = pv.playlist_id
		WHERE p.user_id = $1`,
		userID,
	)
}

func (r *Repository) GetSubscribedAuthorIDs(subscriberID int) ([]int, error) {
	return r.getUserVideoIDs(
		`SELECT author_id
		FROM subscriptions
		WHERE subscriber_id = $1
		ORDER BY created_at DESC`,
		subscriberID,
	)
}

func (r *Repository) GetVideoPopularity(videoIDs []int) (map[int]domain.VideoPopularity, error) {
	popularityByVideoID := make(map[int]domain.VideoPopularity, len(videoIDs))
	if len(videoIDs) == 0 {
		return popularityByVideoID, nil
	}

	for _, id := range videoIDs {
		popularityByVideoID[id] = domain.VideoPopularity{VideoID: id}
	}

	if err := r.fillPopularityCounts(videoIDs, popularityByVideoID, "likes"); err != nil {
		return nil, err
	}
	if err := r.fillPopularityCounts(videoIDs, popularityByVideoID, "comments"); err != nil {
		return nil, err
	}
	if err := r.fillPopularityCounts(videoIDs, popularityByVideoID, "playlists"); err != nil {
		return nil, err
	}

	return popularityByVideoID, nil
}

func (r *Repository) getUserVideoIDs(query string, args ...any) ([]int, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	ids := make([]int, 0)
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return nil, r.translateError(err)
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return ids, nil
}

func (r *Repository) fillPopularityCounts(videoIDs []int, popularityByVideoID map[int]domain.VideoPopularity, source string) error {
	ctx := context.Background()
	query := popularityCountQuery(source)
	rows, err := r.db.Query(ctx, query, videoIDs)
	if err != nil {
		return r.translateError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var videoID int
		var count int
		if err = rows.Scan(&videoID, &count); err != nil {
			return r.translateError(err)
		}

		popularity := popularityByVideoID[videoID]
		switch source {
		case "likes":
			popularity.LikesCount = count
		case "comments":
			popularity.CommentsCount = count
		case "playlists":
			popularity.PlaylistAddsCount = count
		}
		popularityByVideoID[videoID] = popularity
	}

	if err = rows.Err(); err != nil {
		return r.translateError(err)
	}

	return nil
}

func popularityCountQuery(source string) string {
	switch source {
	case "likes":
		return `SELECT video_id, COUNT(*)
			FROM video_likes
			WHERE video_id = ANY($1)
			GROUP BY video_id`
	case "comments":
		return `SELECT video_id, COUNT(*)
			FROM comments
			WHERE video_id = ANY($1) AND status = 'ACTIVE'
			GROUP BY video_id`
	case "playlists":
		return `SELECT video_id, COUNT(*)
			FROM playlist_videos
			WHERE video_id = ANY($1)
			GROUP BY video_id`
	default:
		return `SELECT 0, 0 WHERE FALSE`
	}
}
