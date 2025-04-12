package messenger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPostWebhook_OK(t *testing.T) {
	mux := http.NewServeMux()
	var receivedRequest DirectSendRequest
	mux.HandleFunc("/messenger/v1/channels/direct-send", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&receivedRequest)
		if err != nil {
			t.Errorf("Request contained invalid JSON, %s", err)
		}

		response := []byte(`{}`)
		rw.Write(response)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	payload := &DirectSendRequest{
		Text:                 "Test Text",
		OrganizationMemberId: "12321321321321",
	}

	err := NewMessenger(server.URL).DirectSend("dooray-api-key", payload)
	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	if !reflect.DeepEqual(payload, &receivedRequest) {
		t.Errorf("Payload did not match\nwant: %#v\n got: %#v", payload, receivedRequest)
	}
}
