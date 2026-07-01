package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/repository/postgres/dbmodel"
)

func (r *Repository) CreateComment(comment domain.Comment) (domain.Comment, error) {
	ctx := context.Background()
	var dbComment dbModels.Comment
	err := r.db.QueryRow(ctx,
		`INSERT INTO comments (user_id, video_id, text, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, video_id, text, status, created_at, updated_at`,
		comment.UserID,
		comment.VideoID,
		comment.Text,
		comment.Status,
	).Scan(&dbComment.ID, &dbComment.UserID, &dbComment.VideoID, &dbComment.Text,
		&dbComment.Status, &dbComment.CreatedAt, &dbComment.UpdatedAt)
	if err != nil {
		return domain.Comment{}, r.translateError(err)
	}

	return dbComment.ToDomain(), nil
}

func (r *Repository) GetVideoComments(videoID int) ([]domain.Comment, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, video_id, text, status, created_at, updated_at
		FROM comments
		WHERE video_id = $1 AND status = $2
		ORDER BY created_at DESC`, videoID, domain.CommentStatusActive)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbComments []dbModels.Comment
	for rows.Next() {
		var comment dbModels.Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.VideoID, &comment.Text,
			&comment.Status, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbComments = append(dbComments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	comments := make([]domain.Comment, 0, len(dbComments))
	for _, comment := range dbComments {
		comments = append(comments, comment.ToDomain())
	}

	return comments, nil
}

func (r *Repository) GetCommentByID(commentID int) (domain.Comment, error) {
	ctx := context.Background()
	var dbComment dbModels.Comment
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, video_id, text, status, created_at, updated_at
		FROM comments
		WHERE id = $1 AND status = $2`, commentID, domain.CommentStatusActive).Scan(
		&dbComment.ID, &dbComment.UserID, &dbComment.VideoID, &dbComment.Text,
		&dbComment.Status, &dbComment.CreatedAt, &dbComment.UpdatedAt)
	if err != nil {
		return domain.Comment{}, r.translateError(err)
	}

	return dbComment.ToDomain(), nil
}

func (r *Repository) UpdateComment(comment domain.Comment) (domain.Comment, error) {
	ctx := context.Background()
	var dbComment dbModels.Comment
	err := r.db.QueryRow(ctx,
		`UPDATE comments
		SET text = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND status = $3
		RETURNING id, user_id, video_id, text, status, created_at, updated_at`,
		comment.Text,
		comment.ID,
		domain.CommentStatusActive,
	).Scan(&dbComment.ID, &dbComment.UserID, &dbComment.VideoID, &dbComment.Text,
		&dbComment.Status, &dbComment.CreatedAt, &dbComment.UpdatedAt)
	if err != nil {
		return domain.Comment{}, r.translateError(err)
	}

	return dbComment.ToDomain(), nil
}

func (r *Repository) DeleteComment(commentID int) error {
	ctx := context.Background()
	result, err := r.db.Exec(ctx,
		`UPDATE comments
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND status = $3`, domain.CommentStatusDeleted, commentID, domain.CommentStatusActive)
	if err != nil {
		return r.translateError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
