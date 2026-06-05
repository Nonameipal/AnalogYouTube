package repository

import (
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/models/db"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (r *Repository) CreateComment(comment domain.Comment) (domain.Comment, error) {
	var dbComment dbModels.Comment
	err := r.db.Get(&dbComment, `
		INSERT INTO comments (user_id, video_id, text, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, video_id, text, status, created_at, updated_at`,
		comment.UserID,
		comment.VideoID,
		comment.Text,
		comment.Status,
	)
	if err != nil {
		return domain.Comment{}, r.translateError(err)
	}

	return dbComment.ToDomain(), nil
}

func (r *Repository) GetVideoComments(videoID int) ([]domain.Comment, error) {
	var dbComments []dbModels.Comment
	err := r.db.Select(&dbComments, `
		SELECT id, user_id, video_id, text, status, created_at, updated_at
		FROM comments
		WHERE video_id = $1 AND status = $2
		ORDER BY created_at DESC`, videoID, domain.CommentStatusActive)
	if err != nil {
		return nil, r.translateError(err)
	}

	comments := make([]domain.Comment, 0, len(dbComments))
	for _, comment := range dbComments {
		comments = append(comments, comment.ToDomain())
	}

	return comments, nil
}

func (r *Repository) GetCommentByID(commentID int) (domain.Comment, error) {
	var dbComment dbModels.Comment
	if err := r.db.Get(&dbComment, `
		SELECT id, user_id, video_id, text, status, created_at, updated_at
		FROM comments
		WHERE id = $1 AND status = $2`, commentID, domain.CommentStatusActive); err != nil {
		return domain.Comment{}, r.translateError(err)
	}

	return dbComment.ToDomain(), nil
}

func (r *Repository) UpdateComment(comment domain.Comment) (domain.Comment, error) {
	var dbComment dbModels.Comment
	err := r.db.Get(&dbComment, `
		UPDATE comments
		SET text = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND status = $3
		RETURNING id, user_id, video_id, text, status, created_at, updated_at`,
		comment.Text,
		comment.ID,
		domain.CommentStatusActive,
	)
	if err != nil {
		return domain.Comment{}, r.translateError(err)
	}

	return dbComment.ToDomain(), nil
}

func (r *Repository) DeleteComment(commentID int) error {
	result, err := r.db.Exec(`
		UPDATE comments
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND status = $3`, domain.CommentStatusDeleted, commentID, domain.CommentStatusActive)
	if err != nil {
		return r.translateError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return r.translateError(err)
	}
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
