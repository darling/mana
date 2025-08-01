package llm

import (
	"context"
	"fmt"
)

type Message struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	Role     string `json:"role"`
	Content  string `json:"content"`
}

// Provider defines the interface for a Language Model (LLM) provider.
type Provider interface {
	// Generate generates a response from the LLM based on the provided messages.
	Generate(ctx context.Context, history []Message) (Message, error)

	// ListModels lists the available models for the LLM provider.
	ListModels(ctx context.Context) ([]string, error)

	// Clean up resources used by the LLM provider.
	Close() error
}

type Config struct {
	APIKey string
	Model  string
}

type Manager struct {
	provider Provider
}

func NewManager(providerType string, config Config) (*Manager, error) {
	var provider Provider
	var err error

	switch providerType {
	case "openrouter":
		provider, err = newOpenRouterProvider(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerType)
	}

	if err != nil {
		return nil, err
	}

	return &Manager{provider: provider}, nil
}

func (m *Manager) Generate(ctx context.Context, history []Message) (Message, error) {
	return m.provider.Generate(ctx, history)
}

func (m *Manager) ListModels(ctx context.Context) ([]string, error) {
	return m.provider.ListModels(ctx)
}

func (m *Manager) Close() error {
	return m.provider.Close()
}

var newOpenRouterProvider func(Config) (Provider, error)

func RegisterOpenRouterProvider(factory func(Config) (Provider, error)) {
	newOpenRouterProvider = factory
}
