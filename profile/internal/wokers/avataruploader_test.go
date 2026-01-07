package workers_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/TATAROmangol/mess/profile/internal/ctxkey"
	"github.com/TATAROmangol/mess/profile/internal/loglables"
	"github.com/TATAROmangol/mess/profile/internal/model"
	storagemocks "github.com/TATAROmangol/mess/profile/internal/storage/mocks"
	workers "github.com/TATAROmangol/mess/profile/internal/wokers"
	loggermocks "github.com/TATAROmangol/mess/shared/logger/mocks"
	"github.com/golang/mock/gomock"
)

func TestAvatarUploader_Upload_SuccessUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumer := mqmocks.NewMockConsumer(ctrl)
	storage := storagemocks.NewMockService(ctrl)
	profileRepo := storagemocks.NewMockProfile(ctrl)
	outboxRepo := storagemocks.NewMockAvatarOutbox(ctrl)
	tx := storagemocks.NewMockServiceTransaction(ctrl)

	lg := loggermocks.NewMockLogger(ctrl)
	ctx := ctxkey.WithLogger(context.Background(), lg)

	au := workers.NewAvatarUploader(workers.AvatarUploaderConfig{}, storage)

	msg := workers.AvatarUploaderMessage{
		Key: "user/123/avatar/new.png",
	}
	msgBytes, _ := json.Marshal(msg)

	mqMsg := mqmocks.NewMockMessage(ctrl)
	subjectID := "123"
	prevKey := "user/123/avatar/old.png"

	profileBefore := &model.Profile{
		SubjectID: subjectID,
		AvatarKey: &prevKey,
	}

	profileAfter := &model.Profile{
		SubjectID: subjectID,
		AvatarKey: &msg.Key,
	}

	outbox := &model.AvatarOutbox{
		Key: prevKey,
	}

	consumer.EXPECT().ReadMessage(ctx).Return(mqMsg, nil)

	mqMsg.EXPECT().Value().Return(msgBytes)

	storage.EXPECT().Profile().Return(profileRepo)

	profileRepo.EXPECT().GetProfileFromSubjectID(ctx, subjectID).Return(profileBefore, nil)

	storage.EXPECT().WithTransaction(ctx).Return(tx, nil)

	tx.EXPECT().Profile().Return(profileRepo)

	profileRepo.EXPECT().UpdateAvatarKey(ctx, subjectID, msg.Key).Return(profileAfter, nil)

	lg.EXPECT().With(loglables.Profile, *profileAfter).Return(lg)

	storage.EXPECT().AvatarOutbox().Return(outboxRepo)

	outboxRepo.EXPECT().AddKey(ctx, subjectID, prevKey).Return(outbox, nil)

	lg.EXPECT().With(loglables.AvatarOutbox, *outbox).Return(lg)

	tx.EXPECT().Commit().Return(nil)

	tx.EXPECT().Rollback()

	consumer.EXPECT().Commit(ctx, mqMsg).Return(nil)
	lg.EXPECT().Info("success update")

	err := au.Upload(ctx)

	if err != nil {
		t.Fatalf("upload: %v", err)
	}
}
