package model

import "time"

type Chat struct {
	ID              int
	FirstSubjectID  string
	SecondSubjectID string
	MessagesCount   int
	UpdatedAt       time.Time
	CreatedAt       time.Time
	DeletedAt       *time.Time
}
