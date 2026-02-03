package redisclient

import "fmt"

// channels type - subject:{subj_id}:chat:{chat_id}:type
const (
	allValues         = "*"
	channelKeyChat    = "chat"
	channelKeySubject = "subject"
)

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	ChannelTypeMessage  Type = "message"
	ChannelTypeLastRead Type = "lastread"
)

func GetChannel(subjectID *string, chatID *int, channelType *Type) string {
	subjID := allValues
	if subjectID != nil {
		subjID = *subjectID
	}

	chatIDStr := allValues
	if chatID != nil {
		chatIDStr = fmt.Sprintf("%v", *chatID)
	}

	channelTypeStr := allValues
	if channelType != nil {
		channelTypeStr = channelType.String()
	}

	return fmt.Sprintf("%v:%v:%v:%v:%v",
		channelKeySubject,
		subjID,
		channelKeyChat,
		chatIDStr,
		channelTypeStr,
	)
}
