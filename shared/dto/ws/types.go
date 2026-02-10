package wsdto

type Operation string

const (
	UnknownOperation        Operation = "unknown"
	SendMessageOperation    Operation = "send_message"
	UpdateMessageOperation  Operation = "update_message"
	UpdateLastReadOperation Operation = "update_last_read"
)
