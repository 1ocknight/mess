package model

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID              int        `json:"id"`
	ChatID          int        `json:"chat_id"`
	SenderSubjectID string     `json:"sender_subject_id"`
	Content         string     `json:"content"`
	Number          int        `json:"number"`
	Version         int        `json:"version"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

func (m *Message) ToString() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func UnmarshalMessage(data string) (*Message, error) {
	var m Message
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil, err
	}
	return &m, nil
}
