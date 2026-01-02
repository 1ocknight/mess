DROP INDEX IF EXISTS uniq_profile_subject_alive;

ALTER TABLE profile
ADD CONSTRAINT profile_pkey PRIMARY KEY (subject_id);

ALTER TABLE profile
DROP COLUMN IF EXISTS deleted_at;
