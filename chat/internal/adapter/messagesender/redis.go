package lastreadsender

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

func (i *IMPL) BatchSend(ctx context.Context, recipients []string, messages []model.MessageOutbox) error {
	if len(messages) == 0 {
		return nil
	}

	chatID := messages[0].ChatID

	sendDTO := make([][]byte, 0, len(messages))
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
			ChatID:    chatID,
			Message:   &mess,
			Operation: mqdto.Operation(m.Operation),
		}

		data, err := res.Marshal()
		if err != nil {
			return fmt.Errorf("marshal message dto: %w", err)
		}

		sendDTO = append(sendDTO, data)
	}

	pipe := i.client.Pipeline()
	for _, r := range recipients {
		channel := fmt.Sprintf("%v:%v:%v:%v:%v",
			redisclient.ChannelKeySubject,
			r,
			redisclient.ChannelKeyChat,
			chatID,
			redisclient.ChannelTypeMessage,
		)

		for _, data := range sendDTO {
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
