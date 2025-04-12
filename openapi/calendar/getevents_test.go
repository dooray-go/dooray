package calendar

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestCalendar_GetEvents_OK(t *testing.T) {
	calendars := "1234567890"
	timeMin := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	timeMax := time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)

	var actualCalendars string
	var actualTimeMin string
	var actualTimeMax string

	expectResponse := `{"id":"1234567890"}`

	mux := http.NewServeMux()
	mux.HandleFunc("/calendar/v1/calendars/*/events", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		actualCalendars = query.Get("calendars")
		actualTimeMin = query.Get("timeMin")
		actualTimeMax = query.Get("timeMax")

		response := []byte(expectResponse)
		rw.Write(response)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	response, err := NewCalendar(server.URL).GetEvents("dooray-api-key", calendars, timeMin, timeMax)
	if err != nil {
		t.Errorf("Expected not to receive error: %s", err)
	}

	if !reflect.DeepEqual(calendars, actualCalendars) {
		t.Errorf("Response did not match\nwant: %#v\n got: %#v", expectResponse, response)
	}

	if !reflect.DeepEqual(timeMin.Format(iso8601), actualTimeMin) {
		t.Errorf("Response did not match\nwant: %#v\n got: %#v", timeMin.Format(iso8601), actualTimeMin)
	}

	if !reflect.DeepEqual(timeMax.Format(iso8601), actualTimeMax) {
		t.Errorf("Response did not match\nwant: %#v\n got: %#v", timeMax.Format(iso8601), actualTimeMax)
	}
}
