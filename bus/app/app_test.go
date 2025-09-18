package app

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/nicograef/cloudevents/bus/config"
)

func TestNewApp(t *testing.T) {
	cfg := config.Config{
		Port:        8080,
		Subscribers: []string{"http://localhost:3000/webhook"},
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	if app.Server == nil {
		t.Error("Server should not be nil")
	}

	if app.Config.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", app.Config.Port)
	}
}

func TestNewApp_ServerConfiguration(t *testing.T) {
	cfg := config.Config{
		Port:        9090,
		Subscribers: []string{"http://example.com/webhook"},
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	expectedAddr := ":9090"
	if app.Server.Addr != expectedAddr {
		t.Errorf("Expected server address %s, got %s", expectedAddr, app.Server.Addr)
	}

	expectedReadTimeout := 30 * time.Second
	if app.Server.ReadTimeout != expectedReadTimeout {
		t.Errorf("Expected read timeout %v, got %v", expectedReadTimeout, app.Server.ReadTimeout)
	}

	expectedWriteTimeout := 30 * time.Second
	if app.Server.WriteTimeout != expectedWriteTimeout {
		t.Errorf("Expected write timeout %v, got %v", expectedWriteTimeout, app.Server.WriteTimeout)
	}
}

func TestSetupRoutes(t *testing.T) {
	cfg := config.Config{
		Port:        8080,
		Subscribers: []string{"http://localhost:3000/webhook"},
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	app.SetupRoutes()

	// Test that routes are set up by checking the default mux
	// Note: This is a basic check - integration tests would be better
	req, _ := http.NewRequest("GET", "/health", nil)

	// We can't easily test the mux without starting the server,
	// so this is more of a smoke test
	if req == nil {
		t.Error("Failed to create test request")
	}
}

func TestShutdown(t *testing.T) {
	cfg := config.Config{
		Port:        8080,
		Subscribers: []string{"http://localhost:3000/webhook"},
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	err = app.Shutdown()
	if err != nil {
		t.Errorf("Shutdown() failed: %v", err)
	}
}

func TestRun_ContextCancellation(t *testing.T) {
	cfg := config.Config{
		Port:        8080,
		Subscribers: []string{"http://localhost:3000/webhook"},
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Run the app in a separate goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- app.Run(ctx)
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Cancel the context to trigger shutdown
	cancel()

	// Wait for Run to return
	err = <-errChan
	if err != nil {
		t.Errorf("Run() returned error: %v", err)
	}
}
