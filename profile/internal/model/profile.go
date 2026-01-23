package model

import "time"

type Profile struct {
	SubjectID string
	Alias     string
	Version   int
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}
