package mqdto

import "encoding/json"

type LastRead struct {
	ChatID        int    `json:"chat_id"`
	SubjectID     string `json:"subject_id"`
	MessageID     int    `json:"message_id"`
	MessageNumber int    `json:"message_number"`
}

func (lr *LastRead) Marshal() ([]byte, error) {
	return json.Marshal(lr)
}

func UnmarshalLastRead(data []byte) (*LastRead, error) {
	var lr LastRead
	if err := json.Unmarshal(data, &lr); err != nil {
		return nil, err
	}
	return &lr, nil
}
