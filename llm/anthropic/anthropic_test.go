package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dooray-go/dooray-sdk/llm"
)

func TestQuery_OK(t *testing.T) {
	mux := http.NewServeMux()
	var receivedReq request
	mux.HandleFunc("/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("expected x-api-key test-key, got %s", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("anthropic-version") != apiVersion {
			t.Errorf("expected anthropic-version %s, got %s", apiVersion, r.Header.Get("anthropic-version"))
		}

		json.NewDecoder(r.Body).Decode(&receivedReq)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"content":[{"type":"text","text":"Hello from Claude"}]}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(llm.WithAPIKey("test-key"), llm.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := client.Query(context.Background(), "Hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Hello from Claude" {
		t.Errorf("want %q, got %q", "Hello from Claude", result)
	}
	if receivedReq.Messages[0].Content != "Hello" {
		t.Errorf("want prompt %q, got %q", "Hello", receivedReq.Messages[0].Content)
	}
}

func TestQuery_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":{"type":"authentication_error","message":"invalid api key"}}`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := New(llm.WithAPIKey("bad-key"), llm.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.Query(context.Background(), "Hello")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	llmErr, ok := err.(*llm.Error)
	if !ok {
		t.Fatalf("expected *llm.Error, got %T", err)
	}
	if llmErr.Code != http.StatusUnauthorized {
		t.Errorf("want status %d, got %d", http.StatusUnauthorized, llmErr.Code)
	}
}

func TestNew_MissingAPIKey(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	_, err := New()
	if err == nil {
		t.Fatal("expected error for missing API key, got nil")
	}
}
