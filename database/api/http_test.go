package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendJSONResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}
	sendJSONResponse(rec, data)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected content-type application/json, got %s", ct)
	}
	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	}
	if resp["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", resp)
	}
}

func TestReadJSONRequest_Success(t *testing.T) {
	rec := httptest.NewRecorder()
	type testStruct struct{ Foo string }
	body := bytes.NewBufferString(`{"Foo":"bar"}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	var dest testStruct
	ok := readJSONRequest(rec, req, &dest)
	if !ok {
		t.Errorf("expected success reading JSON")
	}
	if dest.Foo != "bar" {
		t.Errorf("expected Foo=bar, got %s", dest.Foo)
	}
}

func TestReadJSONRequest_InvalidJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	type testStruct struct{ Foo string }
	body := bytes.NewBufferString(`{"Foo":}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	var dest testStruct
	ok := readJSONRequest(rec, req, &dest)
	if ok {
		t.Errorf("expected failure for invalid JSON")
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}
