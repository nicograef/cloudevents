package queue

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nicograef/cloudevents/event"
)

func TestSendToWebhook_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	msg := event.Event{Type: "test"}
	resp, err := SendToWebhook(ts.URL, msg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != "ok" {
		t.Errorf("expected response 'ok', got %s", resp)
	}
}

func TestSendToWebhook_BadURL(t *testing.T) {
	msg := event.Event{Type: "test"}
	_, err := SendToWebhook("http://bad url", msg)
	if err == nil {
		t.Errorf("expected error for bad url")
	}
}
