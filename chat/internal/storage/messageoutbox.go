package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/1ocknight/mess/chat/internal/model"
	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	deletedATIsNullMessageOutboxFilter = fmt.Sprintf("%v %v", MessageOutboxDeletedAtLabel, IsNullLabel)
)

func (s *Storage) doAndReturnMessageOutbox(ctx context.Context, query string, args []interface{}) (*model.MessageOutbox, error) {
	var entity MessageOutboxEntity
	err := sqlx.GetContext(ctx, s.exec, &entity, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return entity.ToModel(), nil
}

func (s *Storage) doAndReturnMessageOutboxes(ctx context.Context, query string, args []interface{}) ([]*model.MessageOutbox, error) {
	var entities []*MessageOutboxEntity
	err := sqlx.SelectContext(ctx, s.exec, &entities, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db get: %w", err)
	}

	return MessageOutboxEntitiesToModels(entities), nil
}

func (s *Storage) AddMessageOutbox(ctx context.Context, chatID int, recipientsID []string, messagePayload string, operation model.Operation) (*model.MessageOutbox, error) {
	query, args, err := sq.
		Insert(MessageOutboxTable).
		Columns(
			MessageOutboxChatIDLabel,
			MessageOutboxRecipientsIDLabel,
			MessageOutboxMessagePayloadLabel,
			MessageOutboxOperationLabel,
		).
		Values(chatID, pq.StringArray(recipientsID), messagePayload, operation).
		Suffix(ReturningSuffix).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build sql: %w", err)
	}

	return s.doAndReturnMessageOutbox(ctx, query, args)
}

func (s *Storage) GetMessageOutbox(ctx context.Context, groupsCnt int, groupNumber int, limit int) ([]*model.MessageOutbox, error) {
	query, args, err := sq.
		Select(AllLabelsSelect).
		From(MessageOutboxTable).
		Where(sq.Expr(fmt.Sprintf("%s %% ? = ?", MessageOutboxChatIDLabel), groupsCnt, groupNumber)).
		Where(sq.Expr(deletedATIsNullMessageOutboxFilter)).
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build main query sql: %w", err)
	}

	return s.doAndReturnMessageOutboxes(ctx, query, args)
}

func (s *Storage) DeleteMessageOutbox(ctx context.Context, ids []int) ([]*model.MessageOutbox, error) {
	if len(ids) == 0 {
		return []*model.MessageOutbox{}, nil
	}

	query, args, err := sq.
		Update(MessageOutboxTable).
		Set(MessageOutboxDeletedAtLabel, time.Now().UTC()).
		Where(sq.Eq{MessageOutboxIDLabel: ids}).
		Where(sq.Expr(deletedATIsNullMessageOutboxFilter)).
		Suffix(ReturningSuffix).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build sql: %w", err)
	}

	return s.doAndReturnMessageOutboxes(ctx, query, args)
}
