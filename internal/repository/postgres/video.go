package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/repository/postgres/dbmodel"
)

func (r *Repository) CreateVideo(video domain.Video) (domain.Video, error) {
	ctx := context.Background()
	var dbVideo dbModels.Video
	err := r.db.QueryRow(ctx, `
		INSERT INTO videos (author_id, title, description, video_url, thumbnail_url, category_id, status)
		VALUES ($1, $2, NULLIF($3, ''), $4, NULLIF($5, ''), $6, $7)
		RETURNING id, author_id, category_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at`,
		video.AuthorID,
		video.Title,
		video.Description,
		video.VideoURL,
		video.ThumbnailURL,
		video.CategoryID,
		video.Status,
	).Scan(
		&dbVideo.ID, &dbVideo.AuthorID, &dbVideo.CategoryID, &dbVideo.Title, &dbVideo.Description,
		&dbVideo.VideoURL, &dbVideo.ThumbnailURL, &dbVideo.Views, &dbVideo.Status, &dbVideo.CreatedAt, &dbVideo.UpdatedAt)
	if err != nil {
		return domain.Video{}, r.translateError(err)
	}

	return dbVideo.ToDomain(), nil
}

func (r *Repository) GetAllVideos() ([]domain.Video, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, category_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at
		FROM videos
		WHERE status = $1
		ORDER BY created_at DESC`, domain.VideoStatusActive)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbVideos []dbModels.Video
	for rows.Next() {
		var video dbModels.Video
		if err := rows.Scan(
			&video.ID, &video.AuthorID, &video.CategoryID, &video.Title, &video.Description,
			&video.VideoURL, &video.ThumbnailURL, &video.Views, &video.Status, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbVideos = append(dbVideos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	videos := make([]domain.Video, 0, len(dbVideos))
	for _, video := range dbVideos {
		videos = append(videos, video.ToDomain())
	}

	return videos, nil
}

func (r *Repository) GetRecommendedVideos() ([]domain.Video, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT v.id, v.author_id, v.category_id, v.title, v.description, v.video_url, v.thumbnail_url, v.views, v.status, v.created_at, v.updated_at
		FROM videos v
		LEFT JOIN video_likes vl ON vl.video_id = v.id
		WHERE v.status = $1
		GROUP BY v.id
		ORDER BY COUNT(vl.id) DESC, v.views DESC, v.created_at DESC
		LIMIT 20`, domain.VideoStatusActive)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbVideos []dbModels.Video
	for rows.Next() {
		var video dbModels.Video
		if err := rows.Scan(
			&video.ID, &video.AuthorID, &video.CategoryID, &video.Title, &video.Description,
			&video.VideoURL, &video.ThumbnailURL, &video.Views, &video.Status, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbVideos = append(dbVideos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	videos := make([]domain.Video, 0, len(dbVideos))
	for _, video := range dbVideos {
		videos = append(videos, video.ToDomain())
	}

	return videos, nil
}

func (r *Repository) GetVideoByID(id int) (domain.Video, error) {
	ctx := context.Background()
	var dbVideo dbModels.Video
	err := r.db.QueryRow(ctx,
		`SELECT id, author_id, category_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at
		FROM videos
		WHERE id = $1 AND status = $2`, id, domain.VideoStatusActive).Scan(
		&dbVideo.ID, &dbVideo.AuthorID, &dbVideo.CategoryID, &dbVideo.Title, &dbVideo.Description,
		&dbVideo.VideoURL, &dbVideo.ThumbnailURL, &dbVideo.Views, &dbVideo.Status, &dbVideo.CreatedAt, &dbVideo.UpdatedAt)
	if err != nil {
		return domain.Video{}, r.translateError(err)
	}

	return dbVideo.ToDomain(), nil
}

func (r *Repository) IncrementVideoViews(id int) error {
	ctx := context.Background()
	result, err := r.db.Exec(ctx,
		`UPDATE videos
		SET views = views + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND status = $2`, id, domain.VideoStatusActive)
	if err != nil {
		return r.translateError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *Repository) UpdateVideo(video domain.Video) (domain.Video, error) {
	ctx := context.Background()
	var dbVideo dbModels.Video
	err := r.db.QueryRow(ctx,
		`UPDATE videos
		SET title = $1,
		    description = NULLIF($2, ''),
		    thumbnail_url = NULLIF($3, ''),
		    category_id = $4,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND status = $6
		RETURNING id, author_id, category_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at`,
		video.Title,
		video.Description,
		video.ThumbnailURL,
		video.CategoryID,
		video.ID,
		domain.VideoStatusActive,
	).Scan(
		&dbVideo.ID, &dbVideo.AuthorID, &dbVideo.CategoryID, &dbVideo.Title, &dbVideo.Description,
		&dbVideo.VideoURL, &dbVideo.ThumbnailURL, &dbVideo.Views, &dbVideo.Status, &dbVideo.CreatedAt, &dbVideo.UpdatedAt)
	if err != nil {
		return domain.Video{}, r.translateError(err)
	}

	return dbVideo.ToDomain(), nil
}

func (r *Repository) DeleteVideo(id int) error {
	ctx := context.Background()
	result, err := r.db.Exec(ctx,
		`UPDATE videos
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND status = $3`, domain.VideoStatusDeleted, id, domain.VideoStatusActive)
	if err != nil {
		return r.translateError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *Repository) GetVideosByAuthorID(authorID int) ([]domain.Video, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, author_id, category_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at
		FROM videos
		WHERE author_id = $1 AND status = $2
		ORDER BY created_at DESC`, authorID, domain.VideoStatusActive)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbVideos []dbModels.Video
	for rows.Next() {
		var video dbModels.Video
		if err := rows.Scan(
			&video.ID, &video.AuthorID, &video.CategoryID, &video.Title, &video.Description,
			&video.VideoURL, &video.ThumbnailURL, &video.Views, &video.Status, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbVideos = append(dbVideos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	videos := make([]domain.Video, 0, len(dbVideos))
	for _, video := range dbVideos {
		videos = append(videos, video.ToDomain())
	}

	return videos, nil
}

func (r *Repository) SearchVideosByTitle(title string) ([]domain.Video, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, author_id, category_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at
		FROM videos
		WHERE status = $1 AND title ILIKE '%' || $2 || '%'
		ORDER BY created_at DESC`, domain.VideoStatusActive, title)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	var dbVideos []dbModels.Video
	for rows.Next() {
		var video dbModels.Video
		if err := rows.Scan(
			&video.ID, &video.AuthorID, &video.CategoryID, &video.Title, &video.Description,
			&video.VideoURL, &video.ThumbnailURL, &video.Views, &video.Status, &video.CreatedAt, &video.UpdatedAt); err != nil {
			return nil, r.translateError(err)
		}
		dbVideos = append(dbVideos, video)
	}

	if err := rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	videos := make([]domain.Video, 0, len(dbVideos))
	for _, video := range dbVideos {
		videos = append(videos, video.ToDomain())
	}

	return videos, nil
}

func (r *Repository) ArchiveDeletedVideo(id int, archivedVideoURL string) error {
	ctx := context.Background()

	result, err := r.db.Exec(ctx,
		`UPDATE videos
		SET archived_video_url = NULLIF($1, ''),
		    status = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND status = $4`,
		archivedVideoURL,
		domain.VideoStatusDeleted,
		id,
		domain.VideoStatusActive,
	)
	if err != nil {
		return r.translateError(err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}
