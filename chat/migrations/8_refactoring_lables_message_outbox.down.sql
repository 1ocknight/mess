ALTER TABLE message_outbox
ADD COLUMN recipient_id TEXT,
DROP COLUMN chat_id,
DROP COLUMN recipients_id,
DROP COLUMN created_at;