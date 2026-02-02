package domain

import (
	"context"

	"github.com/1ocknight/mess/chat/internal/adapter/lastreadsender"
	"github.com/1ocknight/mess/chat/internal/adapter/subjectexist"
	"github.com/1ocknight/mess/chat/internal/model"
	"github.com/1ocknight/mess/chat/internal/storage"
)

type Direction int

const (
	DirectionUnknown Direction = iota
	DirectionAfter
	DirectionBefore
)

type MessagePaginationFilter struct {
	Limit         int
	LastMessageID *int
	Direction     Direction
}

var DefaultPaginationMessage = storage.PaginationFilterIntLastID{
	Limit:     30,
	Asc:       false,
	SortLabel: storage.MessageCreatedAtLabel,
}

type ChatPaginationFilter struct {
	Limit      int
	LastChatID *int
	Direction  Direction
}

var DefaultPaginationChat = storage.PaginationFilterIntLastID{
	Limit:     30,
	Asc:       false,
	SortLabel: storage.ChatUpdatedAtLabel,
}

type Service interface {
	GetChatsMetadata(ctx context.Context, filter *ChatPaginationFilter) ([]*model.ChatMetadata, error)

	AddChat(ctx context.Context, secondSubjectID string) (*model.ChatMetadata, error)
	GetChatMetadataBySubjectID(ctx context.Context, secondSubjectID string) (*model.ChatMetadata, error)
	GetChatMetadataByID(ctx context.Context, chatID int) (*model.ChatMetadata, error)

	UpdateLastRead(ctx context.Context, chatID int, messageID int) (*model.LastRead, error)

	GetMessages(ctx context.Context, chatID int, filter *MessagePaginationFilter) ([]*model.Message, error)
	SendMessage(ctx context.Context, chatID int, content string) (*model.Message, error)
	UpdateMessage(ctx context.Context, messageID int, content string, version int) (*model.Message, error)
}

type Domain struct {
	s   storage.Service
	se  subjectexist.Service
	lrs lastreadsender.Service
}

func New(s storage.Service, se subjectexist.Service, lrs lastreadsender.Service) Service {
	return &Domain{
		s:   s,
		se:  se,
		lrs: lrs,
	}
}
