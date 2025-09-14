package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicograef/cloudevents/event"
	"github.com/nicograef/cloudevents/queue/queue"
)

func TestNewEnqueueHandler_Success(t *testing.T) {
	q := &queue.Queue{Queue: make(chan queue.QueueMessage, 1)}
	handler := NewEnqueueHandler(*q)

	e, err := event.New("com.example.event:v1", "https://example.com", "/users/123", map[string]any{"k": "v"})
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
	var resp EnqueueResponseSuccess
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Ok {
		t.Errorf("expected ok response, got %+v", resp)
	}
	if resp.QueueSize != 1 {
		t.Errorf("expected QueueSize 1, got %d", resp.QueueSize)
	}
	select {
	case got := <-q.Queue:
		if got.Message.Type != "com.example.event:v1" {
			t.Errorf("expected message type 'com.example.event:v1', got %s", got.Message.Type)
		}
		if got.Attempts != 0 {
			t.Errorf("expected Attempts 0, got %d", got.Attempts)
		}
	default:
		t.Errorf("expected message in queue")
	}
}

func TestNewEnqueueHandler_MethodNotAllowed(t *testing.T) {
	q := &queue.Queue{Queue: make(chan queue.QueueMessage, 1)}
	handler := NewEnqueueHandler(*q)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
