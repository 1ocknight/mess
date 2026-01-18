package model

import "time"

type Operation int

const (
	UnknownOperation Operation = iota
	AddOperation
	UpdateOperation
)

type MessageOutbox struct {
	ID        int
	ChatID    int
	MessageID int
	Operation Operation
	DeletedAt *time.Time
}
