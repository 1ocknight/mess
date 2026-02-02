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

func (i *IMPL) Send(ctx context.Context, recipients []string, lastRead *model.LastRead) error {
	mess := mqdto.LastRead{
		ChatID:        lastRead.ChatID,
		SubjectID:     lastRead.SubjectID,
		MessageID:     lastRead.MessageID,
		MessageNumber: lastRead.MessageNumber,
	}

	data, err := mess.Marshal()
	if err != nil {
		return fmt.Errorf("marshal last read message: %w", err)
	}

	for _, r := range recipients {
		channel := fmt.Sprintf("%v:%v:%v:%v:%v",
			redisclient.ChannelKeySubject,
			r,
			redisclient.ChannelKeyChat,
			lastRead.ChatID,
			redisclient.ChannelTypeLastRead,
		)

		err := i.client.Publish(ctx, channel, data)
		if err != nil {
			return fmt.Errorf("publish last read message to redis: %v", err)
		}
	}

	return nil
}

func (i *IMPL) Close() error {
	return i.client.Close()
}
