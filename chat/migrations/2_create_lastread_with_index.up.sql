CREATE TABLE last_read (
    subject_id TEXT NOT NULL,
    chat_id INT NOT NULL,
    message_id INT NOT NULL DEFAULT 0,
    message_number INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_last_read_subject_chat_not_deleted
ON last_read (subject_id, chat_id)
WHERE deleted_at IS NULL;