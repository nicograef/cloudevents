package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicograef/qugo/core"
)

func TestNewEnqueueHandler_Success(t *testing.T) {
	queue := make(chan core.Message, 1)
	handler := NewEnqueueHandler(queue)

	msg := core.Message{Type: "test"}
	payload := EnqueueRequest{Message: msg}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	var resp EnqueueResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if !resp.Ok {
		t.Errorf("expected ok response, got %+v", resp)
	}
	select {
	case got := <-queue:
		if got.Type != "test" {
			t.Errorf("expected message type 'test', got %s", got.Type)
		}
	default:
		t.Errorf("expected message in queue")
	}
}

func TestNewEnqueueHandler_MethodNotAllowed(t *testing.T) {
	queue := make(chan core.Message, 1)
	handler := NewEnqueueHandler(queue)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
