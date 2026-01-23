package storage

const (
	AllLabelsSelect = "*"
	ReturningSuffix = "RETURNING *"
	SkipLocked      = "FOR UPDATE SKIP LOCKED"
	IsNullLabel     = "IS NULL"
	AscSortLabel    = "ASC"
	DescSortLabel   = "DESC"
)

type Table = string

const (
	ProfileTable         Table = "profile"
	AvatarKeyOutboxTable Table = "avatar_outbox"
)

type Label = string

// Profile
const (
	ProfileSubjectIDLabel Label = "subject_id"
	ProfileAliasLabel     Label = "alias"
	ProfileVersionLabel   Label = "version"
	ProfileUpdatedAtLabel Label = "updated_at"
	ProfileCreatedAtLabel Label = "created_at"
	ProfileDeletedAtLabel Label = "deleted_at"
)

// AvatarKeyOutbox
const (
	AvatarKeyOutboxSubjectIDLabel Label = "subject_id"
	AvatarKeyOutboxDeletedAtLabel Label = "deleted_at"
	AvatarKeyOutboxCreatedAtLabel Label = "created_at"
)
