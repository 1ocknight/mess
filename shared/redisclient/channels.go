// channels type - subject:{subj_id}:chat:{chat_id}:type
package redisclient

import (
	"fmt"
	"strconv"
	"strings"
)

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

func BuildListenChannel(subjectID *string, chatID *int, channelType *Type) string {
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

func BuildWriteChannel(subjectID string, chatID int, channelType Type) string {
	return fmt.Sprintf("%v:%v:%v:%v:%v",
		channelKeySubject,
		subjectID,
		channelKeyChat,
		chatID,
		channelType.String(),
	)
}

func GetChannelInfo(ch string) (string, int, Type, error) {
	parts := strings.Split(ch, ":")
	if len(parts) != 5 || parts[0] != "subject" || parts[2] != "chat" {
		return "", 0, "", fmt.Errorf("invalid channel format: %s", ch)
	}

	chatID, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", 0, "", fmt.Errorf("invalid chat ID: %s", parts[3])
	}

	return parts[1], chatID, Type(parts[4]), nil
}
