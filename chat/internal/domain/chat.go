package domain

import (
	"context"
	"fmt"

	"github.com/1ocknight/mess/chat/internal/ctxkey"
	loglables "github.com/1ocknight/mess/chat/internal/loglables"
	"github.com/1ocknight/mess/chat/internal/model"
)

func (d *Domain) GetChatsMetadata(ctx context.Context, filter *ChatPaginationFilter) ([]*model.ChatMetadata, error) {
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %w", err)
	}

	storageFilter := DefaultPaginationChat
	switch filter.Direction {
	case DirectionBefore:
		storageFilter.Asc = false
	default:
		storageFilter.Asc = true
	}

	storageFilter.LastID = filter.LastChatID

	chats, err := d.s.Chat().GetChatsBySubjectID(ctx, subj.GetSubjectId(), &storageFilter)
	if err != nil {
		return nil, fmt.Errorf("get chats bu subject id: %w", err)
	}
	if len(chats) == 0 {
		return []*model.ChatMetadata{}, nil
	}

	lastReads, err := d.s.LastRead().GetLastReadsByChatIDs(ctx, model.GetChatsID(chats))
	if err != nil {
		return nil, fmt.Errorf("get last read by chat ids: %w", err)
	}
	lastReadsMap := map[int]map[string]*model.LastRead{}
	for _, read := range lastReads {
		reads, ok := lastReadsMap[read.ChatID]
		if !ok {
			lastReadsMap[read.ChatID] = map[string]*model.LastRead{}
		}

		reads[read.SubjectID] = read
	}

	lastMessages, err := d.s.Message().GetLastMessagesByChatsID(ctx, model.GetChatsID(chats))
	if err != nil {
		return nil, fmt.Errorf("get last messages by chats id: %w", err)
	}

	lastMessagesMap := map[int]*model.Message{}
	for _, mes := range lastMessages {
		lastMessagesMap[mes.ChatID] = mes
	}

	res := make([]*model.ChatMetadata, 0, len(lastMessages))
	for _, chat := range chats {
		lastMessage, ok := lastMessagesMap[chat.ID]
		if !ok {
			continue
		}

		res = append(res, &model.ChatMetadata{
			ChatID:          chat.ID,
			SecondSubjectID: chat.GetSecondSubject(subj.GetSubjectId()),

			LastReads:     lastReadsMap[chat.ID],
			MessagesCount: chat.MessagesCount,

			LastMessage: lastMessage,

			UpdatedAt: chat.UpdatedAt,
		})
	}

	return res, nil
}

func (d *Domain) AddChat(ctx context.Context, secondSubjectID string) (*model.ChatMetadata, error) {
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %w", err)
	}
	lg, err := ctxkey.ExtractLogger(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract logger: %w", err)
	}

	exists, err := d.se.Exists(ctx, secondSubjectID)
	if err != nil {
		return nil, fmt.Errorf("check subject exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("second subject does not exist")
	}

	tx, err := d.s.WithTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage with transaction: %w", err)
	}
	defer tx.Rollback()

	chat, err := tx.Chat().CreateChat(ctx, subj.GetSubjectId(), secondSubjectID)
	if err != nil {
		return nil, fmt.Errorf("create chat: %w", err)
	}
	lg = lg.With(loglables.Chat, *chat)

	lastReadSubj, err := tx.LastRead().CreateLastRead(ctx, subj.GetSubjectId(), chat.ID)
	if err != nil {
		return nil, fmt.Errorf("create last read subj: %w", err)
	}
	lg = lg.With(loglables.LastReadSubject, *lastReadSubj)

	lastReadSecond, err := tx.LastRead().CreateLastRead(ctx, secondSubjectID, chat.ID)
	if err != nil {
		return nil, fmt.Errorf("create last read second subj: %w", err)
	}
	lg = lg.With(loglables.LastReadSecond, *lastReadSecond)

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	lg.Debug("add last reads with chat")

	return &model.ChatMetadata{
		ChatID:          chat.ID,
		SecondSubjectID: secondSubjectID,
		LastReads: map[string]*model.LastRead{
			subj.GetSubjectId(): lastReadSubj,
			secondSubjectID:     lastReadSecond,
		},
		MessagesCount: 0,
		UpdatedAt:     chat.UpdatedAt,
	}, nil
}

func (d *Domain) GetChatMetadataBySubjectID(ctx context.Context, secondSubjectID string) (*model.ChatMetadata, error) {
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %w", err)
	}

	chat, err := d.s.Chat().GetChatIDBySubjects(ctx, subj.GetSubjectId(), secondSubjectID)
	if err != nil {
		return nil, fmt.Errorf("get chat by subjects: %w", err)
	}

	return d.returnChatMetadata(ctx, subj.GetSubjectId(), chat)
}

func (d *Domain) GetChatMetadataByID(ctx context.Context, chatID int) (*model.ChatMetadata, error) {
	subj, err := ctxkey.ExtractSubject(ctx)
	if err != nil {
		return nil, fmt.Errorf("extract subject: %w", err)
	}

	chat, err := d.s.Chat().GetChatByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat by subjects: %w", err)
	}

	return d.returnChatMetadata(ctx, subj.GetSubjectId(), chat)
}

func (d *Domain) returnChatMetadata(ctx context.Context, subjID string, chat *model.Chat) (*model.ChatMetadata, error) {
	lastReads, err := d.s.LastRead().GetLastReadsByChatID(ctx, chat.ID)
	if err != nil {
		return nil, fmt.Errorf("get last reads by chat id: %w", err)
	}
	lastReadsMap := map[string]*model.LastRead{}
	for _, lr := range lastReads {
		lastReadsMap[lr.SubjectID] = lr
	}

	lastMessage, err := d.s.Message().GetMessageByNumber(ctx, chat.ID, chat.MessagesCount)
	if err != nil {
		return nil, fmt.Errorf("get last message by chat id: %w", err)
	}

	return &model.ChatMetadata{
		ChatID:          chat.ID,
		SecondSubjectID: chat.GetSecondSubject(subjID),
		LastReads:       lastReadsMap,
		MessagesCount:   chat.MessagesCount,
		LastMessage:     lastMessage,
		UpdatedAt:       chat.UpdatedAt,
	}, nil
}
