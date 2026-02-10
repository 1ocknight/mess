package messagesender

import (
	"context"
	"fmt"

	"github.com/1ocknight/mess/chat/internal/model"
	mqdto "github.com/1ocknight/mess/shared/dto/mq"
	"github.com/1ocknight/mess/shared/redisclient"
	"github.com/redis/go-redis/v9"
)

type IMPL struct {
	client *redis.Client
}

func New(cfg redisclient.Config) Service {
	c := redisclient.NewClient(cfg)
	return &IMPL{
		client: c,
	}
}

func (i *IMPL) BatchSend(ctx context.Context, messages []*model.MessageOutbox) error {
	if len(messages) == 0 {
		return nil
	}

	pipe := i.client.Pipeline()

	tp := redisclient.ChannelTypeMessage

	for _, m := range messages {
		messageModel, err := model.UnmarshalMessage(m.MessagePayload)
		if err != nil {
			return fmt.Errorf("unmarshal message payload: %w", err)
		}

		mess := mqdto.Message{
			ID:        messageModel.ID,
			Number:    messageModel.Number,
			SenderID:  messageModel.SenderSubjectID,
			Content:   messageModel.Content,
			Version:   messageModel.Version,
			CreatedAt: messageModel.CreatedAt,
		}

		res := &mqdto.SendMessage{
			ChatID:    messageModel.ChatID,
			Message:   &mess,
			Operation: mqdto.Operation(m.Operation),
		}

		data, err := res.Marshal()
		if err != nil {
			return fmt.Errorf("marshal message dto: %w", err)
		}

		for _, r := range m.RecipientsID {
			channel := redisclient.BuildWriteChannel(r, messageModel.ChatID, tp)
			pipe.Publish(ctx, channel, data)
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("execute redis pipeline: %w", err)
	}

	return nil
}

func (i *IMPL) Close() error {
	return i.client.Close()
}
