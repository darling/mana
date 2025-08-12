package chat

import (
	"context"
)

// Service is the domain-level API used by UIs to manage chats.
// It coordinates with the provider layer and the Store.
//
// The Service is designed to be streaming-first. It exposes a StartAssistantStream
// method that returns a channel of StreamEvent values.
//
// In the short term, this is an interface; an implementation will be added and wired later.
// TODO: Implement a default Service that uses pkg/llm.Manager for providers and a Store.
// TODO: Provide an in-memory Store implementation for initial integration tests.
type Service interface {
	// NewConversation creates a conversation with optional metadata and returns its ID.
	NewConversation(ctx context.Context, meta map[string]string) (ConversationID, error)

	// AddUserMessage appends a user message to the conversation.
	AddUserMessage(ctx context.Context, conv ConversationID, content string, meta map[string]any) (MessageID, error)

	// StartAssistantStream starts assistant generation for the given conversation.
	// It returns a read-only channel of StreamEvent and a cancel function to terminate the stream.
	StartAssistantStream(ctx context.Context, conv ConversationID, opts GenerateOptions) (<-chan StreamEvent, context.CancelFunc, error)

	// GetConversation returns the materialized conversation messages.
	GetConversation(ctx context.Context, conv ConversationID) ([]Message, error)

	// List returns a subset of conversation IDs for navigation.
	List(ctx context.Context, limit int) ([]ConversationID, error)
}
