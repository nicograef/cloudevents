package event

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNew_Success(t *testing.T) {
	e, err := New("com.example.event:v1", "https://example.com", "/users/123", map[string]any{"k": "v"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if e.ID == uuid.Nil {
		t.Errorf("expected non-nil UUID")
	}
	if e.Type != "com.example.event:v1" {
		t.Errorf("unexpected type: %s", e.Type)
	}
	if e.Source != "https://example.com" {
		t.Errorf("unexpected source: %s", e.Source)
	}
	if e.Subject != "/users/123" {
		t.Errorf("unexpected subject: %s", e.Subject)
	}
	if e.Data == nil {
		t.Errorf("expected data to be set")
	}
	if time.Since(e.Time) > time.Minute {
		t.Errorf("unexpected event time: %v", e.Time)
	}
}

func TestValidate_Errors(t *testing.T) {
	cases := []struct {
		name     string
		mutate   func(*Event)
		expected string
	}{
		{"nil uuid", func(e *Event) { e.ID = uuid.Nil }, "event ID cannot be nil"},
		{"short type", func(e *Event) { e.Type = "aaa" }, "event type must be at least 5 characters long"},
		{"zero time", func(e *Event) { e.Time = time.Time{} }, "event time cannot be zero"},
		{"short source", func(e *Event) { e.Source = "abc" }, "event source must be at least 5 characters long"},
		{"bad source scheme", func(e *Event) { e.Source = "ftp://example.com" }, "event source must be a valid URI starting with http:// or https://"},
		{"short subject", func(e *Event) { e.Subject = "abc" }, "event subject must be at least 5 characters long"},
		{"nil data", func(e *Event) { e.Data = nil }, "event data cannot be nil"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := &Event{
				ID:      uuid.New(),
				Type:    "com.example.event:v1",
				Time:    time.Now().UTC(),
				Source:  "https://example.com",
				Subject: "/users/123",
				Data:    map[string]any{"k": "v"},
			}
			// mutate to make invalid
			tc.mutate(e)
			if err := e.Validate(); err == nil || err.Error() != tc.expected {
				t.Fatalf("expected error %q, got %v", tc.expected, err)
			}
		})
	}
}

func TestFromJSON_Success(t *testing.T) {
	validJSON := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"type": "com.example.event:v1",
		"time": "2023-01-01T12:00:00Z",
		"source": "https://example.com",
		"subject": "/users/123",
		"data": {"key": "value"}
	}`

	e, err := FromJSON(validJSON)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if e.Type != "com.example.event:v1" {
		t.Errorf("unexpected type: %s", e.Type)
	}
	if e.Source != "https://example.com" {
		t.Errorf("unexpected source: %s", e.Source)
	}
	if e.Subject != "/users/123" {
		t.Errorf("unexpected subject: %s", e.Subject)
	}
}

func TestFromJSON_InvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid": json}`
	_, err := FromJSON(invalidJSON)
	if err == nil {
		t.Errorf("expected error for invalid JSON")
	}
}

func TestFromJSON_InvalidEvent(t *testing.T) {
	invalidEventJSON := `{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"type": "abc",
		"time": "2023-01-01T12:00:00Z",
		"source": "https://example.com",
		"subject": "/users/123",
		"data": {"key": "value"}
	}`

	_, err := FromJSON(invalidEventJSON)
	if err == nil {
		t.Errorf("expected validation error for short type")
	}
	if err != nil && err.Error() != "event type must be at least 5 characters long" {
		t.Errorf("unexpected error: %v", err)
	}
}
