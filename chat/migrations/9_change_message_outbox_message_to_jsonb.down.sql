-- Откатываем: возвращаем message_id, удаляем message_payload
ALTER TABLE message_outbox 
DROP COLUMN message_payload,
ADD COLUMN message_id INT NOT NULL;
