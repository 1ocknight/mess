package redisclient

// channels type - subject:{subj_id}:chat:{chat_id}:type
const (
	AllValues = "*"

	ChannelKeyChat    = "chat"
	ChannelKeySubject = "subject"

	ChannelTypeMessage  = "message"
	ChannelTypeLastRead = "lastread"
)
