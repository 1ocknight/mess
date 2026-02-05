package domain

import (
	"context"
	"fmt"

	"github.com/1ocknight/mess/chat/internal/ctxkey"
	"github.com/1ocknight/mess/chat/internal/model"
)

func (d *Domain) UpdateLastRead(ctx context.Context, chatID int, messageID int) (*model.LastRead, error) {
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %w", err)
	}
	chat, err := d.s.Chat().GetChatByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat by in: %w", err)
	}
	if chat.FirstSubjectID != subj.GetSubjectId() && chat.SecondSubjectID != subj.GetSubjectId() {
		return nil, SubjectNotHaveThisResource
	}

	mess, err := d.s.Message().GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("get message by id: %w", err)
	}

	lastRead, err := d.s.LastRead().UpdateLastRead(ctx, subj.GetSubjectId(), chatID, messageID, mess.Number)
	if err != nil {
		return nil, fmt.Errorf("update last read: %w", err)
	}

	d.lrs.Send(ctx, chat.GetParticipants(), lastRead)

	return lastRead, nil
}
