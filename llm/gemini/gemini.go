package gemini

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
	defaultBaseURL   = "https://generativelanguage.googleapis.com"
	defaultModel     = "gemini-2.0-flash"
	defaultMaxTokens = 1024
	envAPIKey        = "GEMINI_API_KEY"
)

// Client is a Google Gemini LLM provider.
type Client struct {
	config llm.Config
}

// New creates a new Gemini provider with the given options.
func New(opts ...llm.Option) (*Client, error) {
	var cfg llm.Config
	cfg.ApplyOptions(opts...)
	cfg.ApplyDefaults(envAPIKey, defaultModel, defaultBaseURL, defaultMaxTokens)

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("llm/gemini: API key is required (set %s or use llm.WithAPIKey)", envAPIKey)
	}

	return &Client{config: cfg}, nil
}

type part struct {
	Text string `json:"text"`
}

type content struct {
	Parts []part `json:"parts"`
}

type generationConfig struct {
	MaxOutputTokens int `json:"maxOutputTokens"`
}

type request struct {
	Contents         []content        `json:"contents"`
	GenerationConfig generationConfig `json:"generationConfig"`
}

type candidate struct {
	Content content `json:"content"`
}

type response struct {
	Candidates []candidate `json:"candidates"`
}

type errorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorDetail `json:"error"`
}

// Query sends a prompt to Gemini and returns the text response.
func (c *Client) Query(ctx context.Context, prompt string) (string, error) {
	reqBody := request{
		Contents: []content{
			{Parts: []part{{Text: prompt}}},
		},
		GenerationConfig: generationConfig{
			MaxOutputTokens: c.config.MaxTokens,
		},
	}

	raw, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("llm/gemini: marshal failed: %w", err)
	}

	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s",
		c.config.BaseURL, c.config.Model, c.config.APIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("llm/gemini: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm/gemini: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("llm/gemini: failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error.Message != "" {
			return "", &llm.Error{Provider: "gemini", Code: resp.StatusCode, Message: errResp.Error.Message}
		}
		return "", &llm.Error{Provider: "gemini", Code: resp.StatusCode, Message: resp.Status}
	}

	var result response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("llm/gemini: failed to decode response: %w", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", nil
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
