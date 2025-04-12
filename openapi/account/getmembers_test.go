package account

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAccount_GetMembers_OK(t *testing.T) {
	name := "Manty"
	userCode := "test"

	expectResponse := `{"id":"1234567890"}`

	var actualName string
	var actualUserCode string
	mux := http.NewServeMux()
	mux.HandleFunc("/common/v1/members", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		actualName = query.Get("name")
		actualUserCode = query.Get("userCode")

		response := []byte(expectResponse)
		rw.Write(response)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	response, err := NewAccount(server.URL).GetMembers("dooray-api-key", name, userCode)
	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	if !reflect.DeepEqual(name, actualName) {
		t.Errorf("Response did not match\nwant: %#v\n got: %#v", name, actualName)
	}

	if !reflect.DeepEqual(expectResponse, response) {
		t.Errorf("Response did not match\nwant: %#v\n got: %#v", userCode, actualUserCode)
	}
}
