CREATE TABLE avatar_outbox (
    subject_id TEXT NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);