package workers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/1ocknight/mess/profile/internal/ctxkey"
	"github.com/1ocknight/mess/profile/internal/loglables"
	"github.com/1ocknight/mess/profile/internal/storage"
	"github.com/1ocknight/mess/shared/adapter/subjectdeletereader"
)

type ProfileDeleterConfig struct {
	Kafka subjectdeletereader.Config `yaml:"kafka"`
	Delay time.Duration              `yaml:"delay"`
}

type ProfileDeleter struct {
	CFG      ProfileDeleterConfig
	Consumer subjectdeletereader.Service
	Storage  storage.Service
}

func NewProfileDeleter(cfg ProfileDeleterConfig, s storage.Service) (*ProfileDeleter, error) {
	consumer, err := subjectdeletereader.New(cfg.Kafka)
	if err != nil {
		return nil, fmt.Errorf("new subject delete reader: %w", err)
	}

	return &ProfileDeleter{
		CFG:      cfg,
		Consumer: consumer,
		Storage:  s,
	}, nil
}

func (pd *ProfileDeleter) Delete(ctx context.Context) error {
	lg, err := ctxkey.ExtractLogger(ctx)
	if err != nil {
		return fmt.Errorf("extract logger: %w", err)
	}

	msg, err := pd.Consumer.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("fetch message: %w", err)
	}

	subjectID := msg.GetSubjectID()

	tx, err := pd.Storage.WithTransaction(ctx)
	if err != nil {
		return fmt.Errorf("with transaction: %w", err)
	}
	defer tx.Rollback()

	prof, err := tx.Profile().DeleteProfile(ctx, subjectID)
	if err != nil && errors.Is(err, storage.ErrNoRows) {
		if err := pd.Consumer.Commit(msg); err != nil {
			return fmt.Errorf("commit message: %w", err)
		}
		lg.Info("profile not found, nothing to delete")
		return nil
	}
	if err != nil {
		return fmt.Errorf("delete profile: %w", err)
	}

	out, err := tx.AvatarOutbox().AddKey(ctx, prof.SubjectID)
	if err != nil {
		return fmt.Errorf("add key: %w", err)
	}

	lg = lg.With(loglables.AvatarOutbox, *out)

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	if err := pd.Consumer.Commit(msg); err != nil {
		return fmt.Errorf("commit message: %w", err)
	}

	lg = lg.With(loglables.Profile, *prof)
	lg.Info("success deleted")

	return nil
}

func (pd *ProfileDeleter) Start(ctx context.Context) error {
	lg, err := ctxkey.ExtractLogger(ctx)
	if err != nil {
		return fmt.Errorf("extract logger: %w", err)
	}

	go func() {
		for {
			err := pd.Delete(ctx)
			if err == nil {
				continue
			}

			lg.Error(fmt.Errorf("profile delete: %w", err))

			select {
			case <-time.After(pd.CFG.Delay):
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
