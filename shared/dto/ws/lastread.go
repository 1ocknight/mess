package wsdto

import "encoding/json"

type LastRead struct {
	ChatID        int    `json:"chat_id"`
	SubjectID     string `json:"subject_id"`
	MessageID     int    `json:"message_id"`
	MessageNumber int    `json:"message_number"`
}

func (lr *LastRead) GetData() ([]byte, error) {
	return json.Marshal(lr)
}
