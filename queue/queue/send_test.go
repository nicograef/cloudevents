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

func TestSendToWebhook_InvalidMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	// Create an event with an invalid field (e.g., a channel)
	msg := event.Event{
		Type: "test",
		Data: make(chan int), // channels cannot be JSON serialized
	}

	_, err := SendToWebhook(ts.URL, msg)
	if err == nil {
		t.Errorf("expected error for invalid message data")
	}
}

func TestSendToWebhook_InvalidResponseBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a server error by closing the connection immediately
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, err := hj.Hijack()
			if err == nil {
				conn.Close()
			}
		}
	}))
	defer ts.Close()

	msg := event.Event{Type: "test"}
	resp, err := SendToWebhook(ts.URL, msg)
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if resp != "" {
		t.Errorf("expected empty response body on server error, got %s", resp)
	}
}
