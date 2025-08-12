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

// StreamData is the unit emitted by streaming providers.
// Exactly one of Delta/Done/Err should be set at a time.
type StreamData struct {
	// Delta is the incremental content (may be empty when signalling Done or Err).
	Delta string
	// Done indicates the stream has completed successfully.
	Done bool
	// Err contains an error encountered during streaming.
	Err error
}

// StreamingProvider can stream incremental assistant output.
// Providers that do not implement this interface will be wrapped by Manager
// to provide a compatibility stream that emits a single Delta+Done.
type StreamingProvider interface {
	GenerateStream(ctx context.Context, history []Message) (<-chan StreamData, error)
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

// GenerateStream returns a stream of incremental tokens. If the underlying
// provider does not implement streaming, this method will perform a single
// Generate call and emit one StreamData with the full content followed by Done.
func (m *Manager) GenerateStream(ctx context.Context, history []Message) (<-chan StreamData, error) {
	if sp, ok := m.provider.(StreamingProvider); ok {
		return sp.GenerateStream(ctx, history)
	}

	// Fallback path: non-streaming provider
	ch := make(chan StreamData, 2)
	go func() {
		defer close(ch)
		msg, err := m.provider.Generate(ctx, history)
		if err != nil {
			ch <- StreamData{Err: err}
			return
		}
		if msg.Content != "" {
			ch <- StreamData{Delta: msg.Content}
		}
		ch <- StreamData{Done: true}
	}()
	return ch, nil
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
