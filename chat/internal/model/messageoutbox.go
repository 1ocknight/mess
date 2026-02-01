package model

import "time"

type Operation int

const (
	UnknownOperation Operation = iota
	AddOperation
	UpdateOperation
)

type MessageOutbox struct {
	ID             int
	ChatID         int
	RecipientsID   []string
	MessagePayload string
	Operation      Operation
	DeletedAt      *time.Time
	CreatedAt      time.Time
}

func GetIDsFromMessageOutboxes(outboxes []*MessageOutbox) []int {
	res := make([]int, 0, len(outboxes))
	for _, out := range outboxes {
		res = append(res, out.ID)
	}
	return res
}
