CREATE UNIQUE INDEX IF NOT EXISTS idx_chat_requests_unique_pending_pair
ON chat_requests (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id))
WHERE status = 'pending';
