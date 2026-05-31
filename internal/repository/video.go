package repository

import (
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/models/db"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (r *Repository) CreateVideo(video domain.Video) (domain.Video, error) {
	var dbVideo dbModels.Video
	err := r.db.Get(&dbVideo, `
		INSERT INTO videos (author_id, title, description, video_url, thumbnail_url, status)
		VALUES ($1, $2, NULLIF($3, ''), $4, NULLIF($5, ''), $6)
		RETURNING id, author_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at`,
		video.AuthorID,
		video.Title,
		video.Description,
		video.VideoURL,
		video.ThumbnailURL,
		video.Status,
	)
	if err != nil {
		return domain.Video{}, r.translateError(err)
	}

	return dbVideo.ToDomain(), nil
}

func (r *Repository) GetAllVideos() ([]domain.Video, error) {
	var dbVideos []dbModels.Video
	err := r.db.Select(&dbVideos, `
		SELECT id, author_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at
		FROM videos
		WHERE status = $1
		ORDER BY created_at DESC`, domain.VideoStatusActive)
	if err != nil {
		return nil, r.translateError(err)
	}

	videos := make([]domain.Video, 0, len(dbVideos))
	for _, video := range dbVideos {
		videos = append(videos, video.ToDomain())
	}

	return videos, nil
}

func (r *Repository) GetVideoByID(id int) (domain.Video, error) {
	var dbVideo dbModels.Video
	if err := r.db.Get(&dbVideo, `
		SELECT id, author_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at
		FROM videos
		WHERE id = $1 AND status = $2`, id, domain.VideoStatusActive); err != nil {
		return domain.Video{}, r.translateError(err)
	}

	return dbVideo.ToDomain(), nil
}

func (r *Repository) UpdateVideo(video domain.Video) (domain.Video, error) {
	var dbVideo dbModels.Video
	err := r.db.Get(&dbVideo, `
		UPDATE videos
		SET title = $1,
		    description = NULLIF($2, ''),
		    video_url = $3,
		    thumbnail_url = NULLIF($4, ''),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND status = $6
		RETURNING id, author_id, title, description, video_url, thumbnail_url, views, status, created_at, updated_at`,
		video.Title,
		video.Description,
		video.VideoURL,
		video.ThumbnailURL,
		video.ID,
		domain.VideoStatusActive,
	)
	if err != nil {
		return domain.Video{}, r.translateError(err)
	}

	return dbVideo.ToDomain(), nil
}

func (r *Repository) DeleteVideo(id int) error {
	result, err := r.db.Exec(`
		UPDATE videos
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND status = $3`, domain.VideoStatusDeleted, id, domain.VideoStatusActive)
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
