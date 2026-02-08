CREATE TABLE message_outbox (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL,
    recipients_id TEXT[],
    message_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    operation INT NOT NULL,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);