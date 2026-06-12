package repository

func (r *Repository) LikeVideo(userID int, videoID int) error {
	_, err := r.db.Exec(
		`INSERT INTO video_likes (user_id, video_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, video_id) DO NOTHING`, userID, videoID)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) UnlikeVideo(userID int, videoID int) error {
	_, err := r.db.Exec(
		`DELETE FROM video_likes
		WHERE user_id = $1 AND video_id = $2`, userID, videoID)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) GetVideoLikesCount(videoID int) (int, error) {
	var count int
	if err := r.db.Get(&count, 
		`SELECT COUNT(*)
		FROM video_likes
		WHERE video_id = $1`, videoID); err != nil {
		return 0, r.translateError(err)
	}

	return count, nil
}

func (r *Repository) IsVideoLikedByUser(userID int, videoID int) (bool, error) {
	var exists bool
	if err := r.db.Get(&exists, 
		`SELECT EXISTS(
			SELECT 1 FROM video_likes WHERE user_id = $1 AND video_id = $2
		)`, userID, videoID); err != nil {
		return false, r.translateError(err)
	}

	return exists, nil
}
