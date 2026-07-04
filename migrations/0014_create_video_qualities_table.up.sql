CREATE TABLE IF NOT EXISTS video_qualities (
    id SERIAL PRIMARY KEY,
    video_id INT NOT NULL,
    quality VARCHAR(20) NOT NULL,
    video_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE,
    UNIQUE(video_id, quality)
);

CREATE INDEX IF NOT EXISTS idx_video_qualities_video_id ON video_qualities(video_id);