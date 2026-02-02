package model

import "time"

type ChatMetadata struct {
	ChatID          int
	SecondSubjectID string

	LastReads     map[string]*LastRead
	MessagesCount int

	LastMessage *Message

	UpdatedAt time.Time
}
