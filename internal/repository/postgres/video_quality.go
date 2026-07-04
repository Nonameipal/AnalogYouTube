package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	dbModels "github.com/Nonameipal/AnalogYouTube/internal/repository/postgres/dbmodel"
)

func (r *Repository) CreateVideoQuality(quality domain.VideoQuality) (domain.VideoQuality, error) {
	ctx := context.Background()
	var dbQuality dbModels.VideoQuality

	err := r.db.QueryRow(ctx,
		`INSERT INTO video_qualities (video_id, quality, video_url)
		VALUES ($1, $2, $3)
		ON CONFLICT (video_id, quality)
		DO UPDATE SET video_url = EXCLUDED.video_url
		RETURNING id, video_id, quality, video_url, created_at`,
		quality.VideoID,
		quality.Quality,
		quality.VideoURL,
	).Scan(
		&dbQuality.ID,
		&dbQuality.VideoID,
		&dbQuality.Quality,
		&dbQuality.VideoURL,
		&dbQuality.CreatedAt,
	)
	if err != nil {
		return domain.VideoQuality{}, r.translateError(err)
	}

	return dbQuality.ToDomain(), nil
}

func (r *Repository) GetVideoQualities(videoID int) ([]domain.VideoQuality, error) {
	ctx := context.Background()

	rows, err := r.db.Query(ctx,
		`SELECT id, video_id, quality, video_url, created_at
		FROM video_qualities
		WHERE video_id = $1
		ORDER BY
		    CASE quality
		        WHEN '1080p' THEN 1
		        WHEN '720p' THEN 2
		        WHEN '480p' THEN 3
		        WHEN '360p' THEN 4
		        ELSE 5
		    END`, videoID)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	qualities := make([]domain.VideoQuality, 0)

	for rows.Next() {
		var dbQuality dbModels.VideoQuality

		err = rows.Scan(
			&dbQuality.ID,
			&dbQuality.VideoID,
			&dbQuality.Quality,
			&dbQuality.VideoURL,
			&dbQuality.CreatedAt,
		)
		if err != nil {
			return nil, r.translateError(err)
		}

		qualities = append(qualities, dbQuality.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return qualities, nil
}