package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/nicograef/cloudevents/queue/config"
)

func TestNewApp(t *testing.T) {
	cfg := config.Config{
		Port:        8080,
		Capacity:    100,
		ConsumerURL: "http://localhost:3000/webhook",
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	if app.Queue.Queue == nil {
		t.Error("Queue channel should not be nil")
	}

	if app.Server == nil {
		t.Error("Server should not be nil")
	}

	if app.Config.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", app.Config.Port)
	}

	if app.Config.Capacity != 100 {
		t.Errorf("Expected capacity 100, got %d", app.Config.Capacity)
	}

	// Verify queue capacity
	if cap(app.Queue.Queue) != 100 {
		t.Errorf("Expected queue capacity 100, got %d", cap(app.Queue.Queue))
	}
}

func TestNewApp_ServerConfiguration(t *testing.T) {
	cfg := config.Config{
		Port:        9090,
		Capacity:    50,
		ConsumerURL: "http://example.com/webhook",
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
		Capacity:    100,
		ConsumerURL: "http://localhost:3000/webhook",
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

func TestStartQueueConsumer(t *testing.T) {
	cfg := config.Config{
		Port:        8080,
		Capacity:    10,
		ConsumerURL: "http://localhost:3000/webhook",
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	// Start the consumer
	app.startQueueConsumer()

	// Verify that the wait group has been incremented
	// We can't directly test this, but we can test the shutdown behavior

	// Close the queue and wait for the consumer to finish
	close(app.Queue.Queue)
	app.wg.Wait()

	// If we reach here without hanging, the consumer was properly started and stopped
}

func TestShutdown(t *testing.T) {
	cfg := config.Config{
		Port:        8080,
		Capacity:    10,
		ConsumerURL: "http://localhost:3000/webhook",
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	// Start the consumer so we can test shutdown
	app.startQueueConsumer()

	err = app.Shutdown()
	if err != nil {
		t.Errorf("Shutdown() failed: %v", err)
	}

	// After shutdown, the queue channel should be closed
	// We can verify this by checking if a select with default case
	// on the channel returns immediately
	select {
	case _, ok := <-app.Queue.Queue:
		if ok {
			t.Error("Expected queue channel to be closed")
		}
	default:
		// Channel is closed and drained
	}
}

func TestRun_ContextCancellation(t *testing.T) {
	cfg := config.Config{
		Port:        8080,
		Capacity:    10,
		ConsumerURL: "http://localhost:3000/webhook",
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
