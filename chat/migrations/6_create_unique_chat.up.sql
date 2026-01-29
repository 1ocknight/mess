CREATE UNIQUE INDEX idx_chat_unique_subjects
ON chat (
    LEAST(first_subject_id, second_subject_id),
    GREATEST(first_subject_id, second_subject_id)
)
WHERE deleted_at IS NULL;