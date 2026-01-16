package workers_test

import (
	"context"
	"fmt"
	"testing"

	avatarmocks "github.com/TATAROmangol/mess/profile/internal/adapter/avatar/mocks"
	"github.com/TATAROmangol/mess/profile/internal/ctxkey"
	"github.com/TATAROmangol/mess/profile/internal/loglables"
	"github.com/TATAROmangol/mess/profile/internal/model"
	storagemocks "github.com/TATAROmangol/mess/profile/internal/storage/mocks"
	workers "github.com/TATAROmangol/mess/profile/internal/wokers"
	loggermocks "github.com/TATAROmangol/mess/shared/logger/mocks"
	"github.com/golang/mock/gomock"
)

func TestAvatarDeleter_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := storagemocks.NewMockService(ctrl)
	profile := storagemocks.NewMockProfile(ctrl)
	outbox := storagemocks.NewMockAvatarOutbox(ctrl)
	tx := storagemocks.NewMockServiceTransaction(ctrl)

	storage.EXPECT().Profile().Return(profile).AnyTimes()
	storage.EXPECT().AvatarOutbox().Return(outbox).AnyTimes()
	storage.EXPECT().WithTransaction(gomock.Any()).Return(tx, nil).AnyTimes()
	tx.EXPECT().Profile().Return(profile).AnyTimes()
	tx.EXPECT().AvatarOutbox().Return(outbox).AnyTimes()
	tx.EXPECT().Commit().Return(nil).AnyTimes()
	tx.EXPECT().Rollback().Return(fmt.Errorf("test")).AnyTimes()

	avatar := avatarmocks.NewMockService(ctrl)

	lg := loggermocks.NewMockLogger(ctrl)
	ctx := ctxkey.WithLogger(context.Background(), lg)

	keys := []*model.AvatarOutbox{
		{Key: "avatar1"},
		{Key: "avatar2"},
	}

	outboxKeys := model.GetOutboxKeys(keys)

	outbox.EXPECT().GetKeys(ctx, workers.DeleteAvatarsLimit).Return(keys, nil)
	outbox.EXPECT().GetKeys(ctx, workers.DeleteAvatarsLimit).Return([]*model.AvatarOutbox{}, nil)
	avatar.EXPECT().DeleteObjects(ctx, outboxKeys).Return(nil)
	outbox.EXPECT().DeleteKeys(ctx, outboxKeys).Return(keys, nil)
	lg.EXPECT().With(loglables.DeletedAvatarKeys, outboxKeys).Return(lg)
	lg.EXPECT().Info("success delete")
	lg.EXPECT().Info(gomock.Any())

	svc := workers.AvatarDeleter{
		Avatar:  avatar,
		Storage: storage,
	}

	err := svc.Delete(ctx)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
}
