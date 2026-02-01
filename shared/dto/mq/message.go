package mqdto

import "time"

type Operation int

const (
	UnknownOperation Operation = iota
	AddOperation
	UpdateOperation
)

type Message struct {
	ID        int       `json:"id"`
	Number    int       `json:"number"`
	SenderID  string    `json:"sender_id"`
	Version   int       `json:"version"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type SendMessage struct {
	ChatID    int       `json:"chat_id"`
	Message   *Message  `json:"message"`
	Operation Operation `json:"operation"`
}
