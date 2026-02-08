package messagesenderworker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/1ocknight/mess/chat/internal/adapter/messagesender"
	"github.com/1ocknight/mess/chat/internal/loglables"
	"github.com/1ocknight/mess/chat/internal/model"
	"github.com/1ocknight/mess/chat/internal/storage"
	"github.com/1ocknight/mess/shared/logger"
)

var (
	NoMessagesError = fmt.Errorf("no messages to send")
)

type Config struct {
	Enabled       bool          `yaml:"enabled"`
	GroupsCount   int           `yaml:"groups_count"`
	GroupNumber   int           `yaml:"group_number"`
	MessagesLimit int           `yaml:"messages_limit"`
	Delay         time.Duration `yaml:"delay"`
}

type Worker struct {
	cfg Config
	lg  logger.Logger
	s   storage.Service
	ms  messagesender.Service

	cancel context.CancelFunc
}

func New(cfg Config, s storage.Service, ms messagesender.Service, lg logger.Logger) *Worker {
	return &Worker{
		cfg: cfg,
		s:   s,
		ms:  ms,
		lg:  lg,
	}
}

func (w *Worker) send(ctx context.Context) ([]int, error) {
	messages, err := w.s.MessageOutbox().GetMessageOutbox(ctx, w.cfg.GroupsCount, w.cfg.GroupNumber, w.cfg.MessagesLimit)
	if err != nil {
		return nil, fmt.Errorf("get message outbox: %w", err)
	}

	if len(messages) == 0 {
		return nil, NoMessagesError
	}

	err = w.ms.BatchSend(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("batch send messages: %w", err)
	}

	delMessages, err := w.s.MessageOutbox().DeleteMessageOutbox(ctx, model.GetIDsFromMessageOutboxes(messages))
	if err != nil {
		return nil, fmt.Errorf("delete message outbox: %w", err)
	}

	return model.GetIDsFromMessageOutboxes(delMessages), nil
}

func (w *Worker) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	if !w.cfg.Enabled {
		w.lg.Info("message sender worker is disabled")
		return
	}

	go func() {
		ticker := time.NewTicker(w.cfg.Delay)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				ids, err := w.send(ctx)
				if err == nil {
					w.lg.With(loglables.IDs, ids).Info("sent messages")
					continue
				}
				if errors.Is(err, NoMessagesError) {
					w.lg.Info("no messages to send")
				} else {
					w.lg.Error(err)
				}

				select {
				case <-ticker.C:
					w.lg.Info("waiting delay")
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	w.lg.Info("message sender worker started")
}

func (w *Worker) Stop() {
	if w.cancel != nil {
		return
	}

	w.cancel()
}
