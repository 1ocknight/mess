package model

import "time"

type Operation int

const (
	UnknownOperation Operation = iota
	AddOperation
	UpdateOperation
)

type MessageOutbox struct {
	ID          int
	RecipientID string
	MessageID   int
	Operation   Operation
	DeletedAt   *time.Time
}

func GetIDsFromMessageOutboxes(outboxes []*MessageOutbox) []int {
	res := make([]int, 0, len(outboxes))
	for _, out := range outboxes {
		res = append(res, out.ID)
	}
	return res
}

func GetMessageIDsFromMessageOutboxes(outboxes []*MessageOutbox) []int {
	res := make([]int, 0, len(outboxes))
	for _, out := range outboxes {
		res = append(res, out.MessageID)
	}
	return res
}
