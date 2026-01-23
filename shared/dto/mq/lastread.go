package mqdto

type LastRead struct {
	ChatID      int    `json:"chat_id"`
	RecipientID string `json:"recipient_id"`
	SubjectID   string `json:"subject_id"`
	MessageID   int    `json:"message_id"`
}
