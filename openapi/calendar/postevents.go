package calendar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dooray-go/dooray/utils"
)

// Member represents a member in the event.
type Member struct {
	OrganizationMemberId string `json:"organizationMemberId,omitempty"`
}

// EmailUser represents an email user in the event.
type EmailUser struct {
	EmailAddress string `json:"emailAddress,omitempty"`
	Name         string `json:"name,omitempty"`
}

// Recipient represents a recipient (to or cc) in the event.
type Recipient struct {
	Type      string     `json:"type"`
	Member    *Member    `json:"member,omitempty"`
	EmailUser *EmailUser `json:"emailUser,omitempty"`
}

// Users represents the users involved in the event.
type Users struct {
	To []Recipient `json:"to,omitempty"`
	Cc []Recipient `json:"cc,omitempty"`
}

// Body represents the body of the event.
type Body struct {
	MimeType string `json:"mimeType"`
	Content  string `json:"content"`
}

// RecurrenceRule represents the recurrence rule for the event.
type RecurrenceRule struct {
	Frequency    string `json:"frequency"`
	Interval     int    `json:"interval"`
	Until        string `json:"until"`
	Byday        string `json:"byday"`
	Bymonth      string `json:"bymonth"`
	Bymonthday   string `json:"bymonthday"`
	TimezoneName string `json:"timezoneName"`
}

// Alarm represents an alarm setting for the event.
type Alarm struct {
	Action  string `json:"action"`
	Trigger string `json:"trigger"`
}

// PersonalSettings represents personal settings for the event.
type PersonalSettings struct {
	Alarms []Alarm `json:"alarms"`
	Busy   bool    `json:"busy"`
	Class  string  `json:"class"`
}

// EventRequest represents the payload for creating a calendar event.
type EventRequest struct {
	Users            Users             `json:"users"`
	Subject          string            `json:"subject"`
	Body             Body              `json:"body"`
	StartedAt        utils.JsonTime    `json:"startedAt"`
	EndedAt          utils.JsonTime    `json:"endedAt"`
	WholeDayFlag     bool              `json:"wholeDayFlag"`
	Location         string            `json:"location"`
	RecurrenceRule   *RecurrenceRule   `json:"recurrenceRule,omitempty"`
	PersonalSettings *PersonalSettings `json:"personalSettings,omitempty"`
}

func (c *Calendar) CreateEvent(apikey string, calendarID string, event EventRequest) (string, error) {
	return c.CreateEventCustomHTTPContext(context.Background(), apikey, http.DefaultClient, calendarID, event)
}
func (c *Calendar) CreateEventContext(ctx context.Context, apikey string, calendarID string, event EventRequest) (string, error) {
	return c.CreateEventCustomHTTPContext(ctx, apikey, http.DefaultClient, calendarID, event)
}
func (c *Calendar) CreateEventCustomHTTP(apikey string, httpClient *http.Client, calendarID string, event EventRequest) (string, error) {
	return c.CreateEventCustomHTTPContext(context.Background(), apikey, httpClient, calendarID, event)
}

// CreateEvent sends a POST request to create a calendar event.
func (c *Calendar) CreateEventCustomHTTPContext(ctx context.Context, apikey string, httpClient *http.Client, calendarID string, event EventRequest) (string, error) {
	url := fmt.Sprintf("%s/calendar/v1/calendars/%s/events", c.endPoint, calendarID)

	// Serialize the event request to JSON
	payload, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("failed to marshal event request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Authorization", "dooray-api "+apikey)

	// Send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to create event, status: %s, body: %s", resp.Status, string(body))
	}

	return string(body), nil
}
