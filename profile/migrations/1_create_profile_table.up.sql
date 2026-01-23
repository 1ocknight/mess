-- Active: 1766314077599@@127.0.0.1@5430@profile
CREATE TABLE profile (
    subject_id TEXT, 
    alias TEXT NOT NULL,
    version INT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_profile_alias ON profile(alias);

CREATE UNIQUE INDEX uniq_profile_subject_alive
ON profile(subject_id)
WHERE deleted_at IS NULL;