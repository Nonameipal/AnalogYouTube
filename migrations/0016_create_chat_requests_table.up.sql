CREATE TABLE IF NOT EXISTS chat_requests (
    id SERIAL PRIMARY KEY,
    sender_id INT NOT NULL,
    receiver_id INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    chat_id INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE SET NULL,
    CHECK (sender_id != receiver_id),
    CHECK (status IN ('pending', 'accepted', 'rejected'))
);

CREATE INDEX IF NOT EXISTS idx_chat_requests_sender_id ON chat_requests(sender_id);
CREATE INDEX IF NOT EXISTS idx_chat_requests_receiver_id ON chat_requests(receiver_id);
CREATE INDEX IF NOT EXISTS idx_chat_requests_status ON chat_requests(status);
