package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicograef/cloudevents/event"
)

func TestNewPublishHandler_Success(t *testing.T) {
	publish := func(e event.Event) error { return nil }
	handler := NewPublishHandler(publish)

	e, err := event.New(event.Candidate{Type: "com.example.event:v1", Source: "https://example.com", Subject: "/users/123", Data: map[string]any{"k": "v"}})
	if err != nil {
		t.Fatalf("failed to create event: %v", err)
	}
	body, _ := json.Marshal(e)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	var resp PublishResponseSuccess
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Ok {
		t.Errorf("expected ok response, got %+v", resp)
	}
}

func TestNewPublishHandler_MethodNotAllowed(t *testing.T) {
	publish := func(e event.Event) error { return nil }
	handler := NewPublishHandler(publish)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestNewPublishHandler_InvalidJSON(t *testing.T) {
	publish := func(e event.Event) error { return nil }
	handler := NewPublishHandler(publish)
	body := bytes.NewBufferString(`{"invalid_json":}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestNewPublishHandler_InvalidEvent(t *testing.T) {
	publish := func(e event.Event) error { return nil }
	handler := NewPublishHandler(publish)
	body := bytes.NewBufferString(`{"type":"", "source":""}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp PublishResponseError
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Ok {
		t.Errorf("expected error response, got %+v", resp)
	}
	if resp.Error == "" {
		t.Errorf("expected error message, got empty")
	}
}
