package postgres

import (
	"context"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetAllTags() ([]domain.Tag, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT id, name, created_at
		FROM tags
		ORDER BY name`)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	tags := make([]domain.Tag, 0)
	for rows.Next() {
		var tag domain.Tag
		if err = rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, r.translateError(err)
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return tags, nil
}

func (r *Repository) GetVideoTags(videoID int) ([]domain.Tag, error) {
	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT t.id, t.name, t.created_at
		FROM video_tags vt
		JOIN tags t ON t.id = vt.tag_id
		WHERE vt.video_id = $1
		ORDER BY t.name`,
		videoID,
	)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	tags := make([]domain.Tag, 0)
	for rows.Next() {
		var tag domain.Tag
		if err = rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, r.translateError(err)
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return tags, nil
}

func (r *Repository) GetTagsByVideoIDs(videoIDs []int) (map[int][]domain.Tag, error) {
	tagsByVideoID := make(map[int][]domain.Tag, len(videoIDs))
	if len(videoIDs) == 0 {
		return tagsByVideoID, nil
	}

	ctx := context.Background()
	rows, err := r.db.Query(ctx,
		`SELECT vt.video_id, t.id, t.name, t.created_at
		FROM video_tags vt
		JOIN tags t ON t.id = vt.tag_id
		WHERE vt.video_id = ANY($1)
		ORDER BY vt.video_id, t.name`,
		videoIDs,
	)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var videoID int
		var tag domain.Tag
		if err = rows.Scan(&videoID, &tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, r.translateError(err)
		}
		tagsByVideoID[videoID] = append(tagsByVideoID[videoID], tag)
	}

	if err = rows.Err(); err != nil {
		return nil, r.translateError(err)
	}

	return tagsByVideoID, nil
}

func (r *Repository) ReplaceVideoTags(videoID int, tagNames []string) ([]domain.Tag, error) {
	ctx := context.Background()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, r.translateError(err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM video_tags WHERE video_id = $1`, videoID); err != nil {
		return nil, r.translateError(err)
	}

	tags := make([]domain.Tag, 0, len(tagNames))
	for _, name := range tagNames {
		tag, err := r.upsertTag(ctx, tx, name)
		if err != nil {
			return nil, err
		}

		if _, err = tx.Exec(ctx,
			`INSERT INTO video_tags (video_id, tag_id)
			VALUES ($1, $2)
			ON CONFLICT (video_id, tag_id) DO NOTHING`,
			videoID,
			tag.ID,
		); err != nil {
			return nil, r.translateError(err)
		}

		tags = append(tags, tag)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, r.translateError(err)
	}

	return tags, nil
}

func (r *Repository) upsertTag(ctx context.Context, tx pgx.Tx, name string) (domain.Tag, error) {
	var tag domain.Tag
	err := tx.QueryRow(ctx,
		`WITH inserted AS (
			INSERT INTO tags (name)
			VALUES ($1)
			ON CONFLICT DO NOTHING
			RETURNING id, name, created_at
		)
		SELECT id, name, created_at FROM inserted
		UNION ALL
		SELECT id, name, created_at FROM tags WHERE LOWER(name) = LOWER($1)
		LIMIT 1`,
		name,
	).Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
	if err != nil {
		return domain.Tag{}, r.translateError(err)
	}

	return tag, nil
}
