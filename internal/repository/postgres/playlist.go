package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/repository/postgres/dbmodel"
)

func (r *Repository) CreatePlaylist(playlist domain.Playlist) (domain.Playlist, error) {
	ctx := context.Background()
	var dbPlaylist dbmodel.Playlist

	err := r.db.QueryRow(ctx,
		`INSERT INTO playlists (user_id, name, description)
		VALUES ($1, $2, NULLIF($3, ''))
		RETURNING id, user_id, name, description, created_at, updated_at`,
		playlist.UserID,
		playlist.Name,
		playlist.Description,
	).Scan(&dbPlaylist.ID, &dbPlaylist.UserID, &dbPlaylist.Name, &dbPlaylist.Description, &dbPlaylist.CreatedAt, &dbPlaylist.UpdatedAt)
	if err != nil {
		return domain.Playlist{}, r.translateError(err)
	}

	return dbPlaylist.ToDomain(), nil
}

func (r *Repository) GetPlaylistByID(id int) (domain.Playlist, error) {
	ctx := context.Background()
	var dbPlaylist dbmodel.Playlist

	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, name, description, created_at, updated_at
		FROM playlists
		WHERE id = $1`, id).Scan(
		&dbPlaylist.ID,
		&dbPlaylist.UserID,
		&dbPlaylist.Name,
		&dbPlaylist.Description,
		&dbPlaylist.CreatedAt,
		&dbPlaylist.UpdatedAt,
	)
	if err != nil {
		return domain.Playlist{}, r.translateError(err)
	}

	return dbPlaylist.ToDomain(), nil
}

func (r *Repository) GetUserPlaylists(userID int) ([]domain.Playlist, error) {
	ctx := context.Background()

	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, name, description, created_at, updated_at
		FROM playlists
		WHERE user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	playlists := make([]domain.Playlist, 0)

	for rows.Next() {
		var dbPlaylist dbmodel.Playlist

		err = rows.Scan(
			&dbPlaylist.ID,
			&dbPlaylist.UserID,
			&dbPlaylist.Name,
			&dbPlaylist.Description,
			&dbPlaylist.CreatedAt,
			&dbPlaylist.UpdatedAt,
		)
		if err != nil {
			return nil, r.translateError(err)
		}

		playlists = append(playlists, dbPlaylist.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return playlists, nil
}

func (r *Repository) UpdatePlaylist(playlist domain.Playlist) (domain.Playlist, error) {
	ctx := context.Background()
	var dbPlaylist dbmodel.Playlist

	err := r.db.QueryRow(ctx,
		`UPDATE playlists
		SET name = $1,
		    description = NULLIF($2, ''),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, user_id, name, description, created_at, updated_at`,
		playlist.Name,
		playlist.Description,
		playlist.ID,
	).Scan(
		&dbPlaylist.ID,
		&dbPlaylist.UserID,
		&dbPlaylist.Name,
		&dbPlaylist.Description,
		&dbPlaylist.CreatedAt,
		&dbPlaylist.UpdatedAt,
	)
	if err != nil {
		return domain.Playlist{}, r.translateError(err)
	}

	return dbPlaylist.ToDomain(), nil
}

func (r *Repository) DeletePlaylist(id int) error {
	ctx := context.Background()

	result, err := r.db.Exec(ctx,
		`DELETE FROM playlists
		WHERE id = $1`, id)
	if err != nil {
		return r.translateError(err)
	}

	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *Repository) AddVideoToPlaylist(playlistID int, videoID int) error {
	ctx := context.Background()

	_, err := r.db.Exec(ctx,
		`INSERT INTO playlist_videos (playlist_id, video_id)
		VALUES ($1, $2)
		ON CONFLICT (playlist_id, video_id) DO NOTHING`,
		playlistID,
		videoID,
	)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}

func (r *Repository) GetPlaylistVideos(playlistID int) ([]domain.Video, error) {
	ctx := context.Background()

	rows, err := r.db.Query(ctx,
		`SELECT v.id, v.author_id, v.category_id, v.title, v.description, v.video_url,
		        v.thumbnail_url, v.views, v.status, v.created_at, v.updated_at
		FROM playlist_videos pv
		JOIN videos v ON v.id = pv.video_id
		WHERE pv.playlist_id = $1 AND v.status = $2
		ORDER BY pv.created_at DESC`,
		playlistID,
		domain.VideoStatusActive,
	)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	videos := make([]domain.Video, 0)

	for rows.Next() {
		var dbVideo dbmodel.Video

		err = rows.Scan(
			&dbVideo.ID,
			&dbVideo.AuthorID,
			&dbVideo.CategoryID,
			&dbVideo.Title,
			&dbVideo.Description,
			&dbVideo.VideoURL,
			&dbVideo.ThumbnailURL,
			&dbVideo.Views,
			&dbVideo.Status,
			&dbVideo.CreatedAt,
			&dbVideo.UpdatedAt,
		)
		if err != nil {
			return nil, r.translateError(err)
		}

		videos = append(videos, dbVideo.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return videos, nil
}

func (r *Repository) RemoveVideoFromPlaylist(playlistID int, videoID int) error {
	ctx := context.Background()
	_, err := r.db.Exec(ctx,
		`DELETE FROM playlist_videos
		WHERE playlist_id = $1 AND video_id = $2`,
		playlistID,
		videoID,
	)
	if err != nil {
		return r.translateError(err)
	}

	return nil
}
