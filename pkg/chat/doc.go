// Package chat defines domain types and interfaces for conversations
// and streaming generation, decoupling UI from provider specifics.
//
// This layer sits above pkg/llm and below UIs. It is responsible for
// mapping provider deltas into domain events, and coordinating persistence
// via a Store.
//
// Short term: we will add an in-memory Store and a default Service
// implementation that adapts pkg/llm.Manager, while keeping TODOs for
// SQLite and richer features.
package chat
