package utils

import (
	"fmt"

	mqdto "github.com/1ocknight/mess/shared/dto/mq"
	wsdto "github.com/1ocknight/mess/shared/dto/ws"
	"github.com/1ocknight/mess/shared/redisclient"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidChannelType   = fmt.Errorf("unknown channel type")
	ErrInvalidOperationType = fmt.Errorf("unknown operation type")
)

func messageToWSMessage(m string) (*wsdto.WSMessage, error) {
	mqDTO, err := mqdto.UnmarshalSendMessage([]byte(m))
	if err != nil {
		return nil, fmt.Errorf("unmarshal send message: %w", err)
	}

	msg := wsdto.Message{
		ChatID:    mqDTO.ChatID,
		ID:        mqDTO.Message.ID,
		Number:    mqDTO.Message.Number,
		SenderID:  mqDTO.Message.SenderID,
		Version:   mqDTO.Message.Version,
		Content:   mqDTO.Message.Content,
		CreatedAt: mqDTO.Message.CreatedAt,
	}

	payload, err := msg.GetData()
	if err != nil {
		return nil, fmt.Errorf("get data: %w", err)
	}

	wsMsg := wsdto.WSMessage{
		Data: payload,
	}

	switch mqDTO.Operation {
	case mqdto.AddOperation:
		wsMsg.Type = wsdto.SendMessageOperation
	case mqdto.UpdateOperation:
		wsMsg.Type = wsdto.UpdateMessageOperation
	default:
		return nil, ErrInvalidOperationType
	}

	return &wsMsg, nil
}

func lastReadToWSMessage(m string) (*wsdto.WSMessage, error) {
	mqDTO, err := mqdto.UnmarshalLastRead([]byte(m))
	if err != nil {
		return nil, fmt.Errorf("unmarshal last read: %w", err)
	}

	lastRead := wsdto.LastRead{
		ChatID:        mqDTO.ChatID,
		SubjectID:     mqDTO.SubjectID,
		MessageID:     mqDTO.MessageID,
		MessageNumber: mqDTO.MessageNumber,
	}

	payload, err := lastRead.GetData()
	if err != nil {
		return nil, fmt.Errorf("get data: %w", err)
	}

	wsMsg := wsdto.WSMessage{
		Type: wsdto.UpdateLastReadOperation,
		Data: payload,
	}

	return &wsMsg, nil
}

func RedisMessageToWSMessage(rMsg *redis.Message) (string, int, *wsdto.WSMessage, error) {
	subjID, chatID, rType, err := redisclient.GetChannelInfo(rMsg.Channel)
	if err != nil {
		return "", 0, nil, err
	}

	if rType == redisclient.ChannelTypeMessage {
		msg, err := messageToWSMessage(rMsg.Payload)
		if err != nil {
			return "", 0, nil, err
		}
		return subjID, chatID, msg, nil
	}

	if rType == redisclient.ChannelTypeLastRead {
		msg, err := lastReadToWSMessage(rMsg.Payload)
		if err != nil {
			return "", 0, nil, err
		}
		return subjID, chatID, msg, nil
	}

	return "", 0, nil, ErrInvalidChannelType
}
