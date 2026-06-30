CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    first_user_id INT NOT NULL,
    second_user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (first_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (second_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CHECK (first_user_id != second_user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_chats_unique_pair
ON chats (LEAST(first_user_id, second_user_id), GREATEST(first_user_id, second_user_id));

CREATE INDEX IF NOT EXISTS idx_chats_first_user_id ON chats(first_user_id);
CREATE INDEX IF NOT EXISTS idx_chats_second_user_id ON chats(second_user_id);
