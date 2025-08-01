package llm

import (
	"context"
	"fmt"
	"sync"
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

var (
	registry   = make(map[string]func(Config) (Provider, error))
	registryMu sync.RWMutex
)

func NewManager(providerType string, config Config) (*Manager, error) {
	registryMu.RLock()
	factory, exists := registry[providerType]
	registryMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported provider: %s", providerType)
	}

	provider, err := factory(config)
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

func Register(name string, factory func(Config) (Provider, error)) {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("provider %q already registered", name))
	}
	registry[name] = factory
}
