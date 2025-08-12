package chat

import (
	"time"
)

// ConversationID uniquely identifies a conversation/thread.
type ConversationID string

// MessageID uniquely identifies a message within a conversation.
type MessageID string

// Role represents the author of a message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Usage represents token accounting reported by a provider.
type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// Message is the domain representation of a chat message.
type Message struct {
	ID           MessageID
	Conversation ConversationID
	Role         Role
	Content      string
	Provider     string
	CreatedAt    time.Time
	CompletedAt  *time.Time
	// Metadata can store provider-specific fields (e.g., tool calls, annotations).
	Metadata map[string]any
}

// StreamEvent conveys incremental generation state from the assistant.
// Exactly one of Delta/Done/Err should be set per event.
type StreamEvent struct {
	Conversation ConversationID
	MessageID    MessageID
	Delta        string
	Usage        *Usage
	Done         bool
	Err          error
}

// GenerateOptions are optional controls for a generation.
type GenerateOptions struct {
	// Model override for this request only.
	Model string
	// SystemPrompt allows injecting/overriding a system message.
	SystemPrompt string
	// Temperature, MaxTokens etc. can be added as needed.
	// TODO: add common sampling/limits here and map per provider in the service.
}
