package chat

import "context"

// Store is the persistence port for chat conversations and messages.
// Implementations may be in-memory, SQLite, or remote.
//
// NOTE: For the short term, we will likely use an in-memory store.
// TODO: Provide a SQLite implementation under pkg/store/sqlite using modernc.org/sqlite.
// TODO: Consider WAL and batching for delta persistence efficiency.
type Store interface {
	CreateConversation(ctx context.Context, meta map[string]string) (ConversationID, error)
	AppendUserMessage(ctx context.Context, conv ConversationID, content string, meta map[string]any) (MessageID, error)
	CreateAssistantMessage(ctx context.Context, conv ConversationID, meta map[string]any) (MessageID, error)
	AppendAssistantDelta(ctx context.Context, conv ConversationID, msg MessageID, delta string) error
	CompleteAssistantMessage(ctx context.Context, conv ConversationID, msg MessageID, usage *Usage) error
	ListConversations(ctx context.Context, limit int) ([]ConversationID, error)
	GetConversation(ctx context.Context, conv ConversationID) ([]Message, error)
}
