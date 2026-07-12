ALTER TABLE comments ADD COLUMN parent_id INT NULL REFERENCES comments(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
