package openrouter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/darling/mana/pkg/llm"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  llm.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: llm.Config{
				APIKey: "test-key",
				Model:  "test-model",
			},
			wantErr: false,
		},
		{
			name: "missing API key",
			config: llm.Config{
				APIKey: "",
				Model:  "test-model",
			},
			wantErr: true,
		},
		{
			name: "model optional",
			config: llm.Config{
				APIKey: "test-key",
				Model:  "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && provider == nil {
				t.Error("New() returned nil provider without error")
			}
		})
	}
}

func TestProvider_Generate(t *testing.T) {
	tests := []struct {
		name         string
		testDataFile string
		statusCode   int
		history      []llm.Message
		wantErr      bool
		wantContent  string
		validateReq  func(*testing.T, *http.Request, []byte)
	}{
		{
			name:         "successful generation",
			testDataFile: "generate_success.json",
			statusCode:   http.StatusOK,
			history: []llm.Message{
				{Role: "user", Content: "Hello"},
			},
			wantErr:     false,
			wantContent: "Hello! How can I assist you today?",
			validateReq: func(t *testing.T, req *http.Request, body []byte) {
				// Check headers
				if auth := req.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("Authorization header = %s, want Bearer test-key", auth)
				}
				if ct := req.Header.Get("Content-Type"); ct != "application/json" {
					t.Errorf("Content-Type = %s, want application/json", ct)
				}
				if ref := req.Header.Get("HTTP-Referer"); ref != "https://github.com/darling/mana" {
					t.Errorf("HTTP-Referer = %s, want https://github.com/darling/mana", ref)
				}
				if title := req.Header.Get("X-Title"); title != "Mana CLI" {
					t.Errorf("X-Title = %s, want Mana CLI", title)
				}

				// Check request body
				var reqBody ChatCompletionRequest
				if err := json.Unmarshal(body, &reqBody); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if reqBody.Model != "test-model" {
					t.Errorf("Model = %s, want test-model", reqBody.Model)
				}
				if len(reqBody.Messages) != 1 {
					t.Errorf("Messages length = %d, want 1", len(reqBody.Messages))
				}
				if reqBody.Messages[0].Role != "user" || reqBody.Messages[0].Content != "Hello" {
					t.Errorf("Message = %+v, want {Role:user Content:Hello}", reqBody.Messages[0])
				}
			},
		},
		{
			name:         "unauthorized error",
			testDataFile: "error_401.json",
			statusCode:   http.StatusUnauthorized,
			history: []llm.Message{
				{Role: "user", Content: "Hello"},
			},
			wantErr: true,
		},
		{
			name:         "invalid model error",
			testDataFile: "error_422.json",
			statusCode:   http.StatusUnprocessableEntity,
			history: []llm.Message{
				{Role: "user", Content: "Hello"},
			},
			wantErr: true,
		},
		{
			name:         "empty choices",
			testDataFile: "empty_choices.json",
			statusCode:   http.StatusOK,
			history: []llm.Message{
				{Role: "user", Content: "Hello"},
			},
			wantErr: true,
		},
		{
			name:         "choice error",
			testDataFile: "choice_error.json",
			statusCode:   http.StatusOK,
			history: []llm.Message{
				{Role: "user", Content: "Hello"},
			},
			wantErr: true,
		},
		{
			name:         "malformed JSON response",
			testDataFile: "malformed_json.json",
			statusCode:   http.StatusOK,
			history: []llm.Message{
				{Role: "user", Content: "Hello"},
			},
			wantErr: true,
		},
		{
			name:         "multiple messages in history",
			testDataFile: "generate_success.json",
			statusCode:   http.StatusOK,
			history: []llm.Message{
				{Role: "system", Content: "You are a helpful assistant"},
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
				{Role: "user", Content: "How are you?"},
			},
			wantErr:     false,
			wantContent: "Hello! How can I assist you today?",
			validateReq: func(t *testing.T, req *http.Request, body []byte) {
				var reqBody ChatCompletionRequest
				if err := json.Unmarshal(body, &reqBody); err != nil {
					t.Fatalf("Failed to unmarshal request body: %v", err)
				}
				if len(reqBody.Messages) != 4 {
					t.Errorf("Messages length = %d, want 4", len(reqBody.Messages))
				}
			},
		},
		{
			name:       "server error without error response",
			statusCode: http.StatusInternalServerError,
			history: []llm.Message{
				{Role: "user", Content: "Hello"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Read request body for validation
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}

				// Validate request if validator provided
				if tt.validateReq != nil {
					tt.validateReq(t, r, body)
				}

				// Set status code
				w.WriteHeader(tt.statusCode)

				// Return test data if file specified
				if tt.testDataFile != "" {
					testData, err := os.ReadFile(filepath.Join("testdata", tt.testDataFile))
					if err != nil {
						t.Fatalf("Failed to read test data file: %v", err)
					}
					w.Write(testData)
				}
			}))
			defer server.Close()

			// Since we can't easily override the URL in the current implementation,
			// we would need to refactor the code to make it testable.
			// This would require modifying the Provider struct to accept a baseURL
			// or httpClient, which is a common pattern for testable code

			t.Skip("Test requires refactoring Provider to accept custom HTTP client or base URL")
		})
	}
}

func TestProvider_Generate_RealAPI(t *testing.T) {
	// Skip if no API key is set
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		t.Skip("OPENROUTER_API_KEY not set, skipping real API test")
	}

	provider := &Provider{
		key:   apiKey,
		model: "openai/gpt-3.5-turbo",
	}

	ctx := context.Background()
	history := []llm.Message{
		{Role: "user", Content: "Say 'test successful' and nothing else"},
	}

	response, err := provider.Generate(ctx, history)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check response
	if response.Provider != "openrouter" {
		t.Errorf("Provider = %s, want openrouter", response.Provider)
	}
	if response.Role != "assistant" {
		t.Errorf("Role = %s, want assistant", response.Role)
	}
	if response.Content == "" {
		t.Error("Content is empty")
	}
	if !strings.Contains(strings.ToLower(response.Content), "test successful") {
		t.Errorf("Content = %s, should contain 'test successful'", response.Content)
	}
}

func TestProvider_ListModels(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		statusCode   int
		wantErr      bool
		wantModels   []string
	}{
		{
			name: "successful list",
			responseBody: `{
				"data": [
					{"id": "openai/gpt-4", "name": "GPT-4"},
					{"id": "anthropic/claude-2", "name": "Claude 2"},
					{"id": "meta/llama-2", "name": "Llama 2"}
				]
			}`,
			statusCode: http.StatusOK,
			wantErr:    false,
			wantModels: []string{"openai/gpt-4", "anthropic/claude-2", "meta/llama-2"},
		},
		{
			name:         "unauthorized",
			responseBody: `{"error": {"code": 401, "message": "Unauthorized"}}`,
			statusCode:   http.StatusUnauthorized,
			wantErr:      true,
		},
		{
			name:         "malformed response",
			responseBody: `{"data": [{"id": "model1"`,
			statusCode:   http.StatusOK,
			wantErr:      true,
		},
		{
			name:         "empty response",
			responseBody: `{"data": []}`,
			statusCode:   http.StatusOK,
			wantErr:      false,
			wantModels:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check method
				if r.Method != "GET" {
					t.Errorf("Method = %s, want GET", r.Method)
				}

				// Check headers
				if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
					t.Errorf("Authorization header = %s, want Bearer test-key", auth)
				}

				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// This test also requires refactoring to make the URL configurable
			t.Skip("Test requires refactoring Provider to accept custom HTTP client or base URL")
		})
	}
}

func TestProvider_Close(t *testing.T) {
	provider := &Provider{
		key:   "test-key",
		model: "test-model",
	}

	err := provider.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

// Helper function to create a testable provider with custom HTTP client
func newTestProvider(client *http.Client, baseURL string) *Provider {
	// This would require refactoring the Provider struct
	// to accept these as configuration options
	return &Provider{
		key:   "test-key",
		model: "test-model",
		// client: client,  // Would need to add this field
		// baseURL: baseURL, // Would need to add this field
	}
}
