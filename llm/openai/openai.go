package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dooray-go/dooray-sdk/llm"
)

const (
	defaultBaseURL   = "https://api.openai.com"
	defaultModel     = "gpt-4o"
	defaultMaxTokens = 1024
	envAPIKey        = "OPENAI_API_KEY"
)

// Client is an OpenAI LLM provider.
type Client struct {
	config llm.Config
}

// New creates a new OpenAI provider with the given options.
func New(opts ...llm.Option) (*Client, error) {
	var cfg llm.Config
	cfg.ApplyOptions(opts...)
	cfg.ApplyDefaults(envAPIKey, defaultModel, defaultBaseURL, defaultMaxTokens)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("llm/openai: API key is required (set %s or use llm.WithAPIKey)", envAPIKey)
	}

	return &Client{config: cfg}, nil
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
}

type choice struct {
	Message message `json:"message"`
}

type response struct {
	Choices []choice `json:"choices"`
}

type errorDetail struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

// Query sends a prompt to OpenAI and returns the text response.
func (c *Client) Query(ctx context.Context, prompt string) (string, error) {
	reqBody := request{
		Model:     c.config.Model,
		MaxTokens: c.config.MaxTokens,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	raw, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("llm/openai: marshal failed: %w", err)
	}

	url := c.config.BaseURL + "/v1/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("llm/openai: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm/openai: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("llm/openai: failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error.Message != "" {
			return "", &llm.Error{Provider: "openai", Code: resp.StatusCode, Message: errResp.Error.Message}
		}
		return "", &llm.Error{Provider: "openai", Code: resp.StatusCode, Message: resp.Status}
	}

	var result response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("llm/openai: failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", nil
	}

	return result.Choices[0].Message.Content, nil
}
