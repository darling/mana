package chat

import (
	"context"
	"sync"
	"time"
)

// memoryStore is a minimal in-memory Store implementation for short-term use.
// It is safe for basic concurrent use via a RWMutex. It is not optimized.
// TODO: Replace or complement with a SQLite-backed implementation.
type memoryStore struct {
	mu            sync.RWMutex
	conversations map[ConversationID][]Message
}

// NewMemoryStore constructs an in-memory Store suitable for initial wiring/tests.
func NewMemoryStore() Store {
	return &memoryStore{conversations: make(map[ConversationID][]Message)}
}

func (s *memoryStore) CreateConversation(ctx context.Context, meta map[string]string) (ConversationID, error) {
	id := ConversationID(time.Now().UTC().Format("20060102T150405.000000000"))
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conversations[id] = nil
	return id, nil
}

func (s *memoryStore) AppendUserMessage(ctx context.Context, conv ConversationID, content string, meta map[string]any) (MessageID, error) {
	msgID := MessageID(time.Now().UTC().Format("20060102T150405.000000000"))
	msg := Message{
		ID:           msgID,
		Conversation: conv,
		Role:         RoleUser,
		Content:      content,
		CreatedAt:    time.Now().UTC(),
		Metadata:     meta,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conversations[conv] = append(s.conversations[conv], msg)
	return msgID, nil
}

func (s *memoryStore) CreateAssistantMessage(ctx context.Context, conv ConversationID, meta map[string]any) (MessageID, error) {
	msgID := MessageID(time.Now().UTC().Format("20060102T150405.000000000"))
	msg := Message{
		ID:           msgID,
		Conversation: conv,
		Role:         RoleAssistant,
		Content:      "",
		CreatedAt:    time.Now().UTC(),
		Metadata:     meta,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conversations[conv] = append(s.conversations[conv], msg)
	return msgID, nil
}

func (s *memoryStore) AppendAssistantDelta(ctx context.Context, conv ConversationID, msg MessageID, delta string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	msgs := s.conversations[conv]
	for i := range msgs {
		if msgs[i].ID == msg {
			msgs[i].Content += delta
			break
		}
	}
	s.conversations[conv] = msgs
	return nil
}

func (s *memoryStore) CompleteAssistantMessage(ctx context.Context, conv ConversationID, msg MessageID, usage *Usage) error {
	completed := time.Now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	msgs := s.conversations[conv]
	for i := range msgs {
		if msgs[i].ID == msg {
			msgs[i].CompletedAt = &completed
			break
		}
	}
	s.conversations[conv] = msgs
	return nil
}

func (s *memoryStore) ListConversations(ctx context.Context, limit int) ([]ConversationID, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]ConversationID, 0, len(s.conversations))
	for id := range s.conversations {
		ids = append(ids, id)
	}
	if limit > 0 && len(ids) > limit {
		ids = ids[:limit]
	}
	return ids, nil
}

func (s *memoryStore) GetConversation(ctx context.Context, conv ConversationID) ([]Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := s.conversations[conv]
	out := make([]Message, len(msgs))
	copy(out, msgs)
	return out, nil
}
