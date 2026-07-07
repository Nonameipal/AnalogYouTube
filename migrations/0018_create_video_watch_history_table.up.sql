CREATE TABLE IF NOT EXISTS video_watch_history (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    video_id INT NOT NULL,
    watched_seconds INT NOT NULL DEFAULT 0,
    duration_seconds INT NOT NULL,
    watched_percent DOUBLE PRECISION NOT NULL DEFAULT 0,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    last_watched_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE,
    UNIQUE (user_id, video_id),
    CHECK (watched_seconds >= 0),
    CHECK (duration_seconds > 0),
    CHECK (watched_percent >= 0 AND watched_percent <= 100)
);

CREATE INDEX IF NOT EXISTS idx_video_watch_history_user_id ON video_watch_history(user_id);
CREATE INDEX IF NOT EXISTS idx_video_watch_history_video_id ON video_watch_history(video_id);
CREATE INDEX IF NOT EXISTS idx_video_watch_history_last_watched_at ON video_watch_history(last_watched_at);
