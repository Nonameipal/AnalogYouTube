ALTER TABLE videos
ADD COLUMN IF NOT EXISTS category_id INT NULL REFERENCES categories(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_videos_category_id ON videos(category_id);

