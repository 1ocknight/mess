CREATE TABLE message_outbox (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL,
    message_id INT NOT NULL,
    operation INT NOT NULL,
    deleted_at TIMESTAMPTZ
);