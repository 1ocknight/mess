CREATE TABLE avatar_outbox (
    key TEXT PRIMARY KEY,
    subject_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
)