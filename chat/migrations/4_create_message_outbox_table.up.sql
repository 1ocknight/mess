CREATE TABLE message_outbox (
    id SERIAL PRIMARY KEY,
    recipient_id TEXT NOT NULL,
    message_id INT NOT NULL,
    operation INT NOT NULL,
    deleted_at TIMESTAMPTZ
);