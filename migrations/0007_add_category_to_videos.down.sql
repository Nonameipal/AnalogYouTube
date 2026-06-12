DROP INDEX IF EXISTS videos_category_id;

ALTER TABLE videos
DROP COLUMN IF EXISTS category_id;

