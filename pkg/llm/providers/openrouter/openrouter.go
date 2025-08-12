package openrouter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/darling/mana/pkg/llm"
)

func init() {
	llm.Register("openrouter", New)
}

type Provider struct {
	key   string
	model string
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream,omitempty"`
}

type ChatCompletionChoice struct {
	FinishReason string `json:"finish_reason"`
	Message      struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Error *ErrorResponse `json:"error,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   *ResponseUsage         `json:"usage,omitempty"`
}

type ResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ErrorResponse struct {
	Code     int                    `json:"code"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func New(cfg llm.Config) (llm.Provider, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("API key is required")
	}
	return &Provider{
		key:   cfg.APIKey,
		model: cfg.Model,
	}, nil
}

func (p *Provider) Generate(ctx context.Context, history []llm.Message) (llm.Message, error) {
	// Convert llm.Message to ChatMessage format
	messages := make([]ChatMessage, len(history))
	for i, msg := range history {
		messages[i] = ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// Create request payload
	request := ChatCompletionRequest{
		Model:    p.model,
		Messages: messages,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return llm.Message{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return llm.Message{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+p.key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/darling/mana")
	req.Header.Set("X-Title", "Mana CLI")

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return llm.Message{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return llm.Message{}, fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}
		return llm.Message{}, fmt.Errorf("API error (code %d): %s", errorResp.Code, errorResp.Message)
	}

	// Parse response
	var chatResp ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return llm.Message{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for errors in choices
	if len(chatResp.Choices) == 0 {
		return llm.Message{}, errors.New("no choices returned from API")
	}

	choice := chatResp.Choices[0]
	if choice.Error != nil {
		return llm.Message{}, fmt.Errorf("choice error (code %d): %s", choice.Error.Code, choice.Error.Message)
	}

	// Convert response to llm.Message
	return llm.Message{
		ID:       chatResp.ID,
		Provider: "openrouter",
		Role:     choice.Message.Role,
		Content:  choice.Message.Content,
	}, nil
}

// streamChunk mirrors the structure of OpenAI-style streaming chunks used by OpenRouter.
// Only the relevant subset is modeled here.
type streamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// Ensure Provider implements llm.StreamingProvider
var _ llm.StreamingProvider = (*Provider)(nil)

func (p *Provider) GenerateStream(ctx context.Context, history []llm.Message) (<-chan llm.StreamData, error) {
	// Convert history
	messages := make([]ChatMessage, len(history))
	for i, msg := range history {
		messages[i] = ChatMessage{Role: msg.Role, Content: msg.Content}
	}

	request := ChatCompletionRequest{Model: p.model, Messages: messages, Stream: true}
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/darling/mana")
	req.Header.Set("X-Title", "Mana CLI")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// We'll stream from resp.Body and close it when the stream ends.
	ch := make(chan llm.StreamData, 32)

	go func() {
		defer close(ch)
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			var errorResp ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
				ch <- llm.StreamData{Err: fmt.Errorf("API request failed with status %d", resp.StatusCode)}
				return
			}
			ch <- llm.StreamData{Err: fmt.Errorf("API error (code %d): %s", errorResp.Code, errorResp.Message)}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		// Increase buffer for large SSE lines
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			// Expect lines in form: "data: {...}" or "data: [DONE]"
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if payload == "[DONE]" {
				ch <- llm.StreamData{Done: true}
				return
			}
			var chunk streamChunk
			if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
				ch <- llm.StreamData{Err: fmt.Errorf("failed to decode stream chunk: %w", err)}
				return
			}
			if len(chunk.Choices) == 0 {
				continue
			}
			delta := chunk.Choices[0].Delta.Content
			if delta != "" {
				select {
				case ch <- llm.StreamData{Delta: delta}:
				case <-ctx.Done():
					return
				}
			}
			if chunk.Choices[0].FinishReason != "" {
				// Some providers signal finish via reason on last chunk.
				ch <- llm.StreamData{Done: true}
				return
			}
		}
		if err := scanner.Err(); err != nil {
			ch <- llm.StreamData{Err: err}
			return
		}
		// If the loop exits without [DONE], consider it done.
		ch <- llm.StreamData{Done: true}
	}()

	return ch, nil
}

type modelsResponse struct {
	Data []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

func (p *Provider) ListModels(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var modelsResp modelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	models := make([]string, len(modelsResp.Data))
	for i, model := range modelsResp.Data {
		models[i] = model.ID
	}

	return models, nil
}

func (provider *Provider) Close() error {
	return nil
}
