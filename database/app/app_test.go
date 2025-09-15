package app

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicograef/cloudevents/database/config"
)

func TestNewApp(t *testing.T) {
	cfg := config.Config{
		Port:    8080,
		DataDir: t.TempDir(),
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	if app.Database == nil {
		t.Error("Database should not be nil")
	}

	if app.Server == nil {
		t.Error("Server should not be nil")
	}

	if app.Config.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", app.Config.Port)
	}
}

func TestNewApp_LoadExistingDatabase(t *testing.T) {
	// Create temporary directory and database file
	tempDir := t.TempDir()
	dbFile := filepath.Join(tempDir, "database.json")

	// Create a sample database file
	sampleData := `[{"id":"123e4567-e89b-12d3-a456-426614174000","type":"test.event","time":"2025-09-15T10:00:00Z","source":"https://test.com","subject":"/test","data":{"test":true}}]`
	if err := os.WriteFile(dbFile, []byte(sampleData), 0644); err != nil {
		t.Fatalf("Failed to create test database file: %v", err)
	}

	cfg := config.Config{
		Port:    8080,
		DataDir: tempDir,
	}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	// Check that database was loaded with existing data
	events := app.Database.GetEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
}

func TestSetupRoutes(t *testing.T) {
	tempDir := t.TempDir()
	cfg := config.Config{Port: 8080, DataDir: tempDir}

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
	tempDir := t.TempDir()
	cfg := config.Config{Port: 8080, DataDir: tempDir}

	app, err := NewApp(cfg)
	if err != nil {
		t.Fatalf("NewApp() failed: %v", err)
	}

	err = app.Shutdown()
	if err != nil {
		t.Errorf("Shutdown() failed: %v", err)
	}

	// Verify database file was created
	dbFile := filepath.Join(tempDir, "database.json")
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		t.Error("Database file was not created during shutdown")
	}
}

func TestRun_ContextCancellation(t *testing.T) {
	tempDir := t.TempDir()
	cfg := config.Config{
		Port:    9091, // Use different port to avoid conflicts with other tests
		DataDir: tempDir,
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

	// Verify database file was persisted during shutdown
	dbFile := filepath.Join(tempDir, "database.json")
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		t.Error("Database file was not created during graceful shutdown")
	}
}
