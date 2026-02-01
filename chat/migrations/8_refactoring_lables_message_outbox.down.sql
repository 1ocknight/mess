ALTER TABLE message_outbox
ADD COLUMN recipients_id INT,
DROP COLUMN chat_id,
DROP COLUMN recipients_id,
DROP COLUMN created_at;