package lastreadsender

import (
	"context"
	"fmt"

	"github.com/1ocknight/mess/chat/internal/model"
	mqdto "github.com/1ocknight/mess/shared/dto/mq"
	"github.com/1ocknight/mess/shared/logger"
	"github.com/1ocknight/mess/shared/redisclient"
	"github.com/redis/go-redis/v9"
)

type IMPL struct {
	lg     logger.Logger
	client *redis.Client
}

func New(cfg redisclient.Config, lg logger.Logger) Service {
	c := redisclient.NewClient(cfg)
	return &IMPL{
		lg:     lg,
		client: c,
	}
}

func (i *IMPL) Send(ctx context.Context, recipients []string, lastRead *model.LastRead) {
	lastReadDTO := mqdto.LastRead{
		ChatID:        lastRead.ChatID,
		SubjectID:     lastRead.SubjectID,
		MessageID:     lastRead.MessageID,
		MessageNumber: lastRead.MessageNumber,
	}

	data, err := lastReadDTO.Marshal()
	if err != nil {
		i.lg.Error(fmt.Errorf("marshal last read message: %w", err))
	}

	op := redisclient.ChannelTypeLastRead

	for _, r := range recipients {
		channel := redisclient.GetChannel(&r, &lastRead.ChatID, &op)

		resp := i.client.Publish(ctx, channel, data)
		if resp.Err() != nil {
			i.lg.Error(fmt.Errorf("publish last read message to redis: %w", resp.Err()))
		}
	}
}

func (i *IMPL) Close() error {
	return i.client.Close()
}
