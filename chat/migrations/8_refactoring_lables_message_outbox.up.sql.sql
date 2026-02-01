-- Active: 1766314077599@@127.0.0.1@5430@chat
ALTER TABLE message_outbox
ADD COLUMN chat_id INT,
ADD COLUMN created_at TIMESTAMPTZ DEFAULT NOW(),
ADD COLUMN recipients_id TEXT[],
DROP COLUMN recipient_id;