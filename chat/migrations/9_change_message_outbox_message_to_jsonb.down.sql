-- Active: 1766314077599@@127.0.0.1@5430@chat
-- Откатываем: возвращаем message_id, удаляем message_payload
ALTER TABLE message_outbox 
DROP COLUMN message_payload,
ADD COLUMN message_id INT NOT NULL DEFAULT 0;
