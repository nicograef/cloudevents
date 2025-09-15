package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/nicograef/cloudevents/database/database"
	"github.com/nicograef/cloudevents/event"
)

func TestNewAddEventHandler_Success(t *testing.T) {
	db := database.New()
	handler := NewAddEventHandler(*db)

	e := event.Candidate{Type: "com.example.event:v1", Source: "https://example.com", Subject: "/users/123", Data: map[string]any{"k": "v"}}
	body, _ := json.Marshal(e)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	var resp AddEventResponseSuccess
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Ok {
		t.Errorf("expected ok response, got %+v", resp)
	}
	if resp.Event.ID == uuid.Nil {
		t.Errorf("expected valid Event ID, got %v", resp.Event.ID)
	}

}

func TestNewAddEventHandler_MethodNotAllowed(t *testing.T) {
	db := database.New()
	handler := NewAddEventHandler(*db)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
