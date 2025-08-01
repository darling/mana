package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
