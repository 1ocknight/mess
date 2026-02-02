package httpdto

import "time"

type MessageResponse struct {
	ID        int       `json:"id"`
	Number    int       `json:"number"`
	Version   int       `json:"version"`
	Content   string    `json:"content"`
	SenderID  string    `json:"sender_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LastReadResponse struct {
	MessageID     int `json:"message_id"`
	MessageNumber int `json:"message_number"`
}

type ChatResponse struct {
	ChatID          int    `json:"chat_id"`
	SecondSubjectID string `json:"second_subject_id"`

	LastReads     map[string]LastReadResponse `json:"last_reads"`
	MessagesCount int                         `json:"messages_count"`

	LastMessage MessageResponse `json:"last_message"`

	UpdatedAt time.Time `json:"updated_at"`
}

type AddMessageRequest struct {
	Content string `json:"content"`
}

type UpdateMessageRequest struct {
	Content string `json:"content"`
	Version int    `json:"version"`
}

type UpdateLastReadRequest struct {
	MessageID int `json:"message_id"`
}
