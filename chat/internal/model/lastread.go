package model

import "time"

type LastRead struct {
	SubjectID     string
	ChatID        int
	MessageID     int
	MessageNumber int
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
