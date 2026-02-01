package storage_test

import (
	"testing"

	"github.com/1ocknight/mess/chat/internal/model"
	"github.com/1ocknight/mess/chat/internal/storage"
)

func TestStorage_AddMessageOutbox(t *testing.T) {
	s, err := storage.New(CFG)
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	initData(t)
	defer cleanupDB(t)

	// Создаем новый MessageOutbox
	chatID := 2
	recipientsID := []string{"subj-1", "subj-3"}
	messagePayload := `{"id":4,"chat_id":2,"sender_subject_id":"subj-1","content":"new test message","version":1}`
	operation := model.AddOperation

	result, err := s.MessageOutbox().AddMessageOutbox(t.Context(), chatID, recipientsID, messagePayload, operation)
	if err != nil {
		t.Fatalf("add message outbox: %v", err)
	}

	// Проверяем результат
	if result.ChatID != chatID {
		t.Errorf("expected chat_id %d, got %d", chatID, result.ChatID)
	}

	if len(result.RecipientsID) != len(recipientsID) {
		t.Errorf("expected %d recipients, got %d", len(recipientsID), len(result.RecipientsID))
	}

	for i, recip := range recipientsID {
		if result.RecipientsID[i] != recip {
			t.Errorf("expected recipient[%d] = %s, got %s", i, recip, result.RecipientsID[i])
		}
	}

	if result.MessagePayload != messagePayload {
		t.Errorf("expected message_payload %s, got %s", messagePayload, result.MessagePayload)
	}

	if result.Operation != operation {
		t.Errorf("expected operation %v, got %v", operation, result.Operation)
	}

	if result.DeletedAt != nil {
		t.Error("expected deleted_at to be nil")
	}

	if result.ID == 0 {
		t.Error("expected ID to be assigned")
	}
}

func TestStorage_GetMessageOutbox(t *testing.T) {
	s, err := storage.New(CFG)
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	initData(t)
	defer cleanupDB(t)

	outbox, err := s.MessageOutbox().GetMessageOutbox(t.Context(), 2)
	if err != nil {
		t.Fatalf("get message outbox: %v", err)
	}

	if len(outbox) != 2 {
		t.Fatalf("wait len 2, have: %v", len(outbox))
	}

	// Проверяем, что данные корректно загружены
	for _, msg := range outbox {
		if msg.ChatID != 1 {
			t.Errorf("expected chat_id 1, got %d", msg.ChatID)
		}
		if len(msg.RecipientsID) == 0 {
			t.Error("expected recipients_id to be populated")
		}
		if msg.MessagePayload == "" {
			t.Error("expected message_payload to be populated")
		}
	}
}

func TestStorage_DeleteMessageOutbox(t *testing.T) {
	s, err := storage.New(CFG)
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	initData(t)
	defer cleanupDB(t)

	del, err := s.MessageOutbox().DeleteMessageOutbox(t.Context(), []int{1, 2})
	if err != nil {
		t.Fatalf("delete message outbox: %v", err)
	}
	if len(del) != 2 {
		t.Fatalf("wait len 2, have: %v", len(del))
	}

	for _, dl := range del {
		if dl.DeletedAt == nil {
			t.Fatalf("not delete: %v", *dl)
		}
	}

	// Проверяем, что повторное удаление не найдет записи
	del2, err := s.MessageOutbox().DeleteMessageOutbox(t.Context(), []int{1, 2})
	if err != nil {
		t.Fatalf("delete message outbox second time: %v", err)
	}
	if len(del2) != 0 {
		t.Fatalf("expected 0 deleted records on second delete, got: %v", len(del2))
	}
}

func TestStorage_GetMessageOutbox_EmptyResult(t *testing.T) {
	s, err := storage.New(CFG)
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	// Не инициализируем данные - чистая БД
	defer cleanupDB(t)

	outbox, err := s.MessageOutbox().GetMessageOutbox(t.Context(), 10)
	if err != nil {
		t.Fatalf("get message outbox: %v", err)
	}

	if len(outbox) != 0 {
		t.Fatalf("expected empty result, got %d records", len(outbox))
	}
}

func TestStorage_DeleteMessageOutbox_EmptyIDs(t *testing.T) {
	s, err := storage.New(CFG)
	if err != nil {
		t.Fatalf("could not construct receiver type: %v", err)
	}

	initData(t)
	defer cleanupDB(t)

	// Передаем пустой массив ID
	del, err := s.MessageOutbox().DeleteMessageOutbox(t.Context(), []int{})
	if err != nil {
		t.Fatalf("delete message outbox with empty ids: %v", err)
	}

	if len(del) != 0 {
		t.Fatalf("expected 0 deleted records, got: %v", len(del))
	}
}
