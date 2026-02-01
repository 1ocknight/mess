package mqdto

import (
	"encoding/json"
	"time"
)

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

func (m *SendMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func UnmarshalSendMessage(data []byte) (*SendMessage, error) {
	var sm SendMessage
	if err := json.Unmarshal(data, &sm); err != nil {
		return nil, err
	}
	return &sm, nil
}
