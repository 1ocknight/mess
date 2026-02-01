package lastreadsender

import (
	"context"

	"github.com/1ocknight/mess/chat/internal/model"
)

type Service interface {
	BatchSend(ctx context.Context, recipients []string, messages []model.MessageOutbox) error
	Close() error
}
