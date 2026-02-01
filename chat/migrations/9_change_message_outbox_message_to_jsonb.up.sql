-- Заменяем message_id на message_payload (JSONB для хранения полного сообщения)
ALTER TABLE message_outbox 
DROP COLUMN message_id,
ADD COLUMN message_payload JSONB NOT NULL DEFAULT '{}'::jsonb;
