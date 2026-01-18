package model

import "time"

type Message struct {
	ID              int
	ChatID          int
	SenderSubjectID string
	Content         string
	Number          int
	Version         int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
