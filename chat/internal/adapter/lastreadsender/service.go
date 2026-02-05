package lastreadsender

import (
	"context"

	"github.com/1ocknight/mess/chat/internal/model"
)

type Service interface {
	Send(ctx context.Context, recipients []string, lastRead *model.LastRead)
	Close() error
}
