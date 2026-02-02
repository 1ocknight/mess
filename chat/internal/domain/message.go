package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/1ocknight/mess/chat/internal/ctxkey"
	loglables "github.com/1ocknight/mess/chat/internal/loglables"
	"github.com/1ocknight/mess/chat/internal/model"
	"github.com/1ocknight/mess/chat/internal/storage"
	"github.com/1ocknight/mess/shared/utils"
)

func (d *Domain) GetMessages(ctx context.Context, chatID int, filter *MessagePaginationFilter) ([]*model.Message, error) {
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

	storageFilter := DefaultPaginationMessage
	switch filter.Direction {
	case DirectionBefore:
		storageFilter.Asc = false
	default:
		storageFilter.Asc = true
	}

	storageFilter.LastID = filter.LastMessageID

	messages, err := d.s.Message().GetMessagesByChatID(ctx, chatID, &storageFilter)
	if err != nil {
		return nil, fmt.Errorf("get messages by chat id: %w", err)
	}

	if len(messages) == 0 {
		return messages, nil
	}

	if !storageFilter.Asc {
		utils.ReverseSlice(messages)
	}

	lastMess := messages[len(messages)-1]
	lastRead, err := d.s.LastRead().UpdateLastRead(ctx, subj.GetSubjectId(), chatID, lastMess.ID, lastMess.Number)
	if err != nil && !errors.Is(err, storage.ErrNoRows) {
		return nil, fmt.Errorf("update last read: %w", err)
	}

	err = d.lrs.Send(ctx, chat.GetParticipants(), lastRead)
	if err != nil {
		return nil, fmt.Errorf("send last read: %w", err)
	}

	return messages, nil
}

func (d *Domain) SendMessage(ctx context.Context, chatID int, content string) (*model.Message, error) {
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %w", err)
	}
	lg, err := ctxkey.ExtractLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract logger: %w", err)
	}

	tx, err := d.s.WithTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage with transaction: %w", err)
	}
	defer tx.Rollback()

	chat, err := tx.Chat().IncrementChatMessageNumber(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("increment chat message number: %w", err)
	}
	lg = lg.With(loglables.Chat, *chat)

	message, err := tx.Message().CreateMessage(ctx, chatID, subj.GetSubjectId(), content, chat.MessagesCount)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	lg = lg.With(loglables.Message, *message)

	payload, err := message.ToString()
	if err != nil {
		return nil, fmt.Errorf("marshal message to string: %w", err)
	}
	outbox, err := tx.MessageOutbox().AddMessageOutbox(ctx, chat.ID, chat.GetParticipants(), payload, model.AddOperation)
	if err != nil {
		return nil, fmt.Errorf("add message outbox: %w", err)
	}
	lg = lg.With(loglables.MessageOutbox, *outbox)

	lastRead, err := tx.LastRead().UpdateLastRead(ctx, subj.GetSubjectId(), chatID, message.ID, message.Number)
	if err != nil {
		return nil, fmt.Errorf("update last read: %w", err)
	}
	lg = lg.With(loglables.LastRead, *lastRead)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	lg.Debug("send message")

	err = d.lrs.Send(ctx, chat.GetParticipants(), lastRead)
	if err != nil {
		return nil, fmt.Errorf("send last read: %w", err)
	}

	return message, nil
}

func (d *Domain) UpdateMessage(ctx context.Context, messageID int, content string, version int) (*model.Message, error) {
	lg, err := ctxkey.ExtractLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract logger: %w", err)
	}
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %w", err)
	}

	mess, err := d.s.Message().GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("get chat by in: %w", err)
	}
	if mess.SenderSubjectID != subj.GetSubjectId() {
		return nil, SubjectNotHaveThisResource
	}

	tx, err := d.s.WithTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage with transaction: %w", err)
	}
	defer tx.Rollback()

	message, err := tx.Message().UpdateMessageContent(ctx, messageID, content, version)
	if err != nil {
		return nil, fmt.Errorf("update message content: %w", err)
	}
	lg = lg.With(loglables.Message, *message)

	payload, err := message.ToString()
	if err != nil {
		return nil, fmt.Errorf("marshal message to string: %w", err)
	}
	outbox, err := tx.MessageOutbox().AddMessageOutbox(ctx, message.ChatID, []string{message.SenderSubjectID}, payload, model.UpdateOperation)
	if err != nil {
		return nil, fmt.Errorf("add message outbox: %w", err)
	}
	lg = lg.With(loglables.MessageOutbox, *outbox)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	lg.Debug("update message")

	return message, nil
}
