package repository

import "context"

func (r *Repository) LikeVideo(userID int, videoID int) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`INSERT INTO video_likes (user_id, video_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, video_id) DO NOTHING`, userID, videoID)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) UnlikeVideo(userID int, videoID int) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`DELETE FROM video_likes
		WHERE user_id = $1 AND video_id = $2`, userID, videoID)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) GetVideoLikesCount(videoID int) (int, error) {
	ctx := context.Background()
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*)
		FROM video_likes
		WHERE video_id = $1`, videoID).Scan(&count)
	if err != nil {
		return 0, r.translateError(err)
	}

	return count, nil
}

func (r *Repository) IsVideoLikedByUser(userID int, videoID int) (bool, error) {
	ctx := context.Background()
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM video_likes WHERE user_id = $1 AND video_id = $2
		)`, userID, videoID).Scan(&exists)
	if err != nil {
		return false, r.translateError(err)
	}

	return exists, nil
}
