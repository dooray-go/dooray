package calendar

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dooray-go/dooray/utils"
)

func TestCreateEvent(t *testing.T) {
	// Mock server to simulate the API response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the HTTP method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check the headers
		if r.Header.Get("Authorization") != "dooray-api test-api-key" {
			t.Errorf("Expected Authorization header 'dooray-api test-api-key', got %s", r.Header.Get("Authorization"))
		}

		// Check the Content-Type header
		if r.Header.Get("Content-Type") != "application/json;charset=utf-8" {
			t.Errorf("Expected Content-Type 'application/json;charset=utf-8', got %s", r.Header.Get("Content-Type"))
		}

		fmt.Println("Request received at mock server")
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		defer r.Body.Close()
		// Check the request body

		fmt.Print("Request body: ", string(body))

		// Respond with a mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "12345", "status": "created"}`))
	}))
	defer mockServer.Close()

	// Create a Calendar instance with the mock server endpoint
	calendar := &Calendar{endPoint: mockServer.URL}

	// Create a sample EventRequest
	event := EventRequest{
		Users:   Users{},
		Subject: "Test Event",
		Body: Body{
			MimeType: "text/html",
			Content:  "This is a test event."},
		StartedAt:    utils.JsonTime(time.Now()),
		EndedAt:      utils.JsonTime(time.Now().Add(1 * time.Hour)),
		WholeDayFlag: false,
		Location:     "Test Location",
	}

	// Call the CreateEvent method
	response, err := calendar.CreateEvent("test-api-key", "test-calendar-id", event)
	if err != nil {
		t.Fatalf("CreateEvent returned an error: %v", err)
	}

	// Verify the response
	expectedResponse := `{"id": "12345", "status": "created"}`
	if response != expectedResponse {
		t.Errorf("Expected response %s, got %s", expectedResponse, response)
	}
}
