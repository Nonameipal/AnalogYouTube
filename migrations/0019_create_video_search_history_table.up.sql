CREATE TABLE IF NOT EXISTS video_search_history (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    query TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_video_search_history_user_id ON video_search_history(user_id);
CREATE INDEX IF NOT EXISTS idx_video_search_history_created_at ON video_search_history(created_at);
