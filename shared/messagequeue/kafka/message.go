package kafka

import "time"

type Message struct {
	Key       []byte
	Val       []byte
	Topic     string
	Partition int
	Offset    int64
	Time      time.Time
}

func (m *Message) Value() []byte {
	return m.Val
}
