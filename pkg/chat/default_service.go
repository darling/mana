package chat

import (
	"context"

	"github.com/darling/mana/pkg/llm"
)

// DefaultService is a simple implementation that adapts llm.Manager and a Store.
// It materializes user messages and assistant deltas into the Store.
//
// Streaming ergonomics:
// - If the underlying provider supports streaming (future in llm.Manager), this will forward deltas.
// - Otherwise, it will call Generate once and emit a single Done event with the full text.
//
// NOTE: In this initial commit, llm.Manager does not yet expose streaming. We wire the
// non-streaming path and leave TODOs for streaming integration.

type DefaultService struct {
	store   Store
	manager *llm.Manager
}

func NewDefaultService(store Store, manager *llm.Manager) *DefaultService {
	return &DefaultService{store: store, manager: manager}
}

func (s *DefaultService) NewConversation(ctx context.Context, meta map[string]string) (ConversationID, error) {
	return s.store.CreateConversation(ctx, meta)
}

func (s *DefaultService) AddUserMessage(ctx context.Context, conv ConversationID, content string, meta map[string]any) (MessageID, error) {
	return s.store.AppendUserMessage(ctx, conv, content, meta)
}

func (s *DefaultService) StartAssistantStream(ctx context.Context, conv ConversationID, opts GenerateOptions) (<-chan StreamEvent, context.CancelFunc, error) {
	// Create assistant placeholder message
	assistantID, err := s.store.CreateAssistantMessage(ctx, conv, map[string]any{"model": opts.Model})
	if err != nil {
		return nil, func() {}, err
	}

	// Channel for events
	events := make(chan StreamEvent, 16)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(events)
		// Load conversation to construct history for llm layer
		msgs, err := s.store.GetConversation(ctx, conv)
		if err != nil {
			select {
			case events <- StreamEvent{Conversation: conv, MessageID: assistantID, Err: err}:
			case <-ctx.Done():
			}
			return
		}

		// Map chat.Message -> llm.Message history
		llmHistory := make([]llm.Message, 0, len(msgs))
		for _, m := range msgs {
			llmHistory = append(llmHistory, llm.Message{Role: string(m.Role), Content: m.Content})
		}

		stream, err := s.manager.GenerateStream(ctx, llmHistory)
		if err != nil {
			select {
			case events <- StreamEvent{Conversation: conv, MessageID: assistantID, Err: err}:
			case <-ctx.Done():
			}
			return
		}

		for data := range stream {
			if data.Err != nil {
				// Surface error and stop
				select {
				case events <- StreamEvent{Conversation: conv, MessageID: assistantID, Err: data.Err}:
				case <-ctx.Done():
				}
				return
			}
			if data.Delta != "" {
				_ = s.store.AppendAssistantDelta(ctx, conv, assistantID, data.Delta)
				select {
				case events <- StreamEvent{Conversation: conv, MessageID: assistantID, Delta: data.Delta}:
				case <-ctx.Done():
					return
				}
			}
			if data.Done {
				_ = s.store.CompleteAssistantMessage(ctx, conv, assistantID, nil)
				select {
				case events <- StreamEvent{Conversation: conv, MessageID: assistantID, Done: true}:
				case <-ctx.Done():
				}
				return
			}
		}
		// If stream channel closes without Done, complete anyway
		_ = s.store.CompleteAssistantMessage(ctx, conv, assistantID, nil)
		select {
		case events <- StreamEvent{Conversation: conv, MessageID: assistantID, Done: true}:
		case <-ctx.Done():
		}
	}()

	return events, cancel, nil
}

func (s *DefaultService) GetConversation(ctx context.Context, conv ConversationID) ([]Message, error) {
	return s.store.GetConversation(ctx, conv)
}

func (s *DefaultService) List(ctx context.Context, limit int) ([]ConversationID, error) {
	return s.store.ListConversations(ctx, limit)
}
