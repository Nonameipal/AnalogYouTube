CREATE TABLE IF NOT EXISTS donations (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    video_id INT NULL,
    amount DECIMAL(10,2) NOT NULL,
    message TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE SET NULL,
    CHECK (amount > 0),
    CHECK (sender_id != receiver_id)
);

CREATE INDEX IF NOT EXISTS idx_donations_sender_id ON donations(sender_id);
CREATE INDEX IF NOT EXISTS idx_donations_receiver_id ON donations(receiver_id);
CREATE INDEX IF NOT EXISTS idx_donations_video_id ON donations(video_id);
CREATE INDEX IF NOT EXISTS idx_donations_created_at ON donations(created_at);

