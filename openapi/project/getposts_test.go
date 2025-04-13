package project

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestProject_GetPosts(t *testing.T) {
	projectId := "1234567890"
	tomemberIds := "1234567890,0987654321"

	var actualProjectId string
	var actualToMemberIds string

	response := `{"id":"1234567890"}`

	mux := http.NewServeMux()
	mux.HandleFunc("/project/v1/projects/{projectId}/posts", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		actualProjectId = r.PathValue("projectId")

		query := r.URL.Query()
		actualToMemberIds = query.Get("toMemberIds")

		response := []byte(response)
		rw.Write(response)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	actualResponse, err := NewProject(server.URL).GetPosts("dooray-api-key", projectId, tomemberIds)
	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	if !reflect.DeepEqual(response, actualResponse) {
		t.Errorf("Response did not match\nwant: %#v\n got: %#v", projectId, actualProjectId)
	}

	if !reflect.DeepEqual(projectId, actualProjectId) {
		t.Errorf("projeceId did not match\nwant: %#v\n got: %#v", projectId, actualProjectId)
	}

	if !reflect.DeepEqual(tomemberIds, actualToMemberIds) {
		t.Errorf("tomemberIds did not match\nwant: %#v\n got: %#v", tomemberIds, actualToMemberIds)
	}

}
