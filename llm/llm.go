package llm

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	DefaultConnectTimeout = 5 * time.Second
	DefaultReadTimeout    = 60 * time.Second
)

// Provider is the core LLM abstraction. Any LLM backend implements this.
type Provider interface {
	Query(ctx context.Context, prompt string) (string, error)
}

// Config holds common configuration for all LLM providers.
type Config struct {
	APIKey     string
	Model      string
	BaseURL    string
	MaxTokens  int
	HTTPClient *http.Client
}

// Option configures a provider.
type Option func(*Config)

func WithAPIKey(key string) Option {
	return func(c *Config) { c.APIKey = key }
}

func WithModel(model string) Option {
	return func(c *Config) { c.Model = model }
}

func WithBaseURL(url string) Option {
	return func(c *Config) { c.BaseURL = url }
}

func WithMaxTokens(n int) Option {
	return func(c *Config) { c.MaxTokens = n }
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) { c.HTTPClient = client }
}

// ApplyOptions applies functional options to the config.
func (c *Config) ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

// ApplyDefaults fills in zero-value fields with provider-specific defaults.
func (c *Config) ApplyDefaults(envKey, defaultModel, defaultBaseURL string, defaultMaxTokens int) {
	if c.APIKey == "" {
		c.APIKey = os.Getenv(envKey)
	}
	if c.Model == "" {
		c.Model = defaultModel
	}
	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.MaxTokens == 0 {
		c.MaxTokens = defaultMaxTokens
	}
	if c.HTTPClient == nil {
		c.HTTPClient = newDefaultHTTPClient()
	}
}

func newDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: DefaultReadTimeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: DefaultConnectTimeout,
			}).DialContext,
			TLSHandshakeTimeout: DefaultConnectTimeout,
		},
	}
}

// Error represents an LLM provider error.
type Error struct {
	Provider string
	Code     int
	Message  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("llm/%s: %d %s", e.Provider, e.Code, e.Message)
}

func (e *Error) Retryable() bool {
	return e.Code >= 500 || e.Code == http.StatusTooManyRequests
}
