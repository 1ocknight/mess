package messagesender

import (
	"context"

	"github.com/1ocknight/mess/chat/internal/model"
)

type Service interface {
	BatchSend(ctx context.Context, messages []model.MessageOutbox) error
	Close() error
}
