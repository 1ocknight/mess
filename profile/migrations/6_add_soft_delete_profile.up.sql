ALTER TABLE profile
ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE profile DROP CONSTRAINT profile_pkey;

CREATE UNIQUE INDEX uniq_profile_subject_alive
ON profile(subject_id)
WHERE deleted_at IS NULL;
