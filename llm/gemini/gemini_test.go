package gemini

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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "test-key" {
			t.Errorf("expected key=test-key, got %s", r.URL.Query().Get("key"))
		}

		json.NewDecoder(r.Body).Decode(&receivedReq)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"Hello from Gemini"}]}}]}`))
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

	if result != "Hello from Gemini" {
		t.Errorf("want %q, got %q", "Hello from Gemini", result)
	}
	if receivedReq.Contents[0].Parts[0].Text != "Hello" {
		t.Errorf("want prompt %q, got %q", "Hello", receivedReq.Contents[0].Parts[0].Text)
	}
}

func TestQuery_Error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"code":400,"message":"API key not valid"}}`))
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
	if llmErr.Code != http.StatusBadRequest {
		t.Errorf("want status %d, got %d", http.StatusBadRequest, llmErr.Code)
	}
}

func TestNew_MissingAPIKey(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	_, err := New()
	if err == nil {
		t.Fatal("expected error for missing API key, got nil")
	}
}
