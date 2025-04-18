package project

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestProject_GetProjects(t *testing.T) {
	projectType := "1234567890"
	scope := "1234567890,0987654321"
	state := "1234567890,0987654321"

	var actualType string
	var actualScope string
	var actualState string

	response := `{"id":"1234567890"}`

	mux := http.NewServeMux()
	mux.HandleFunc("/project/v1/projects", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		actualType = query.Get("type")
		actualScope = query.Get("scope")
		actualState = query.Get("state")

		response := []byte(response)
		rw.Write(response)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	actualResponse, err := NewProject(server.URL).GetProjects("dooray-api-key", projectType, scope, state)
	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	if !reflect.DeepEqual(response, actualResponse) {
		t.Errorf("Response did not match\nwant: %#v\n got: %#v", response, actualResponse)
	}

	if !reflect.DeepEqual(projectType, actualType) {
		t.Errorf("projectType did not match\nwant: %#v\n got: %#v", projectType, actualType)
	}

	if !reflect.DeepEqual(scope, actualScope) {
		t.Errorf("scope did not match\nwant: %#v\n got: %#v", scope, actualScope)
	}

	if !reflect.DeepEqual(state, actualState) {
		t.Errorf("state did not match\nwant: %#v\n got: %#v", state, actualState)
	}

}
