package anthropic

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
	defaultBaseURL   = "https://api.anthropic.com"
	defaultModel     = "claude-sonnet-4-20250514"
	defaultMaxTokens = 1024
	apiVersion       = "2023-06-01"
	envAPIKey        = "ANTHROPIC_API_KEY"
)

// Client is an Anthropic Claude LLM provider.
type Client struct {
	config llm.Config
}

// New creates a new Anthropic provider with the given options.
func New(opts ...llm.Option) (*Client, error) {
	var cfg llm.Config
	cfg.ApplyOptions(opts...)
	cfg.ApplyDefaults(envAPIKey, defaultModel, defaultBaseURL, defaultMaxTokens)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("llm/anthropic: API key is required (set %s or use llm.WithAPIKey)", envAPIKey)
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

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type response struct {
	Content []contentBlock `json:"content"`
}

type errorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

// Query sends a prompt to Claude and returns the text response.
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
		return "", fmt.Errorf("llm/anthropic: marshal failed: %w", err)
	}

	url := c.config.BaseURL + "/v1/messages"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("llm/anthropic: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", apiVersion)

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm/anthropic: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("llm/anthropic: failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error.Message != "" {
			return "", &llm.Error{Provider: "anthropic", Code: resp.StatusCode, Message: errResp.Error.Message}
		}
		return "", &llm.Error{Provider: "anthropic", Code: resp.StatusCode, Message: resp.Status}
	}

	var result response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("llm/anthropic: failed to decode response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", nil
	}

	return result.Content[0].Text, nil
}
