package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()

	cfg := Load()

	if cfg.Port != 3000 {
		t.Errorf("expected default port 3000, got %d", cfg.Port)
	}
	if cfg.Capacity != 1000 {
		t.Errorf("expected default capacity 1000, got %d", cfg.Capacity)
	}
	if cfg.ConsumerURL != "http://localhost:4000" {
		t.Errorf("expected default consumer URL 'http://localhost:4000', got %s", cfg.ConsumerURL)
	}
	if cfg.DeliveryAttempts != 3 {
		t.Errorf("expected default delivery attempts 3, got %d", cfg.DeliveryAttempts)
	}
}

func TestLoad_EnvValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "8080"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("CAPACITY", "42"); err != nil {
		t.Fatalf("Failed to set CAPACITY: %v", err)
	}
	if err := os.Setenv("CONSUMER_URL", "http://test/webhook"); err != nil {
		t.Fatalf("Failed to set CONSUMER_URL: %v", err)
	}
	if err := os.Setenv("DELIVERY_ATTEMPTS", "5"); err != nil {
		t.Fatalf("Failed to set DELIVERY_ATTEMPTS: %v", err)
	}

	cfg := Load()

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.Capacity != 42 {
		t.Errorf("expected capacity 42, got %d", cfg.Capacity)
	}
	if cfg.ConsumerURL != "http://test/webhook" {
		t.Errorf("expected consumer URL 'http://test/webhook', got %s", cfg.ConsumerURL)
	}
	if cfg.DeliveryAttempts != 5 {
		t.Errorf("expected delivery attempts 5, got %d", cfg.DeliveryAttempts)
	}
}

func TestLoad_InvalidIntAndLowValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "notanint"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("CAPACITY", "badint"); err != nil {
		t.Fatalf("Failed to set CAPACITY: %v", err)
	}
	if err := os.Setenv("DELIVERY_ATTEMPTS", "0"); err != nil {
		t.Fatalf("Failed to set DELIVERY_ATTEMPTS: %v", err)
	}
	if err := os.Setenv("CONSUMER_URL", ""); err != nil { // Empty string should use default
		t.Fatalf("Failed to set CONSUMER_URL: %v", err)
	}

	cfg := Load()

	if cfg.Port != 3000 {
		t.Errorf("expected fallback port 3000, got %d", cfg.Port)
	}
	if cfg.Capacity != 1000 {
		t.Errorf("expected fallback capacity 1000, got %d", cfg.Capacity)
	}
	if cfg.DeliveryAttempts != 3 {
		t.Errorf("expected fallback delivery attempts 3, got %d", cfg.DeliveryAttempts)
	}
	if cfg.ConsumerURL != "http://localhost:4000" {
		t.Errorf("expected fallback consumer URL 'http://localhost:4000', got %s", cfg.ConsumerURL)
	}
}

func TestLoad_NegativeValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "-1"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("CAPACITY", "-10"); err != nil {
		t.Fatalf("Failed to set CAPACITY: %v", err)
	}
	if err := os.Setenv("DELIVERY_ATTEMPTS", "-2"); err != nil {
		t.Fatalf("Failed to set DELIVERY_ATTEMPTS: %v", err)
	}

	cfg := Load()

	// Should fallback to defaults due to validation (must be at least 1)
	if cfg.Port != 3000 {
		t.Errorf("expected fallback port 3000 for negative value, got %d", cfg.Port)
	}
	if cfg.Capacity != 1000 {
		t.Errorf("expected fallback capacity 1000 for negative value, got %d", cfg.Capacity)
	}
	if cfg.DeliveryAttempts != 3 {
		t.Errorf("expected fallback delivery attempts 3 for negative value, got %d", cfg.DeliveryAttempts)
	}
}

func TestLoad_ValidEdgeCases(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "1"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("CAPACITY", "1"); err != nil {
		t.Fatalf("Failed to set CAPACITY: %v", err)
	}
	if err := os.Setenv("DELIVERY_ATTEMPTS", "1"); err != nil {
		t.Fatalf("Failed to set DELIVERY_ATTEMPTS: %v", err)
	}

	cfg := Load()

	if cfg.Port != 1 {
		t.Errorf("expected port 1, got %d", cfg.Port)
	}
	if cfg.Capacity != 1 {
		t.Errorf("expected capacity 1, got %d", cfg.Capacity)
	}
	if cfg.DeliveryAttempts != 1 {
		t.Errorf("expected delivery attempts 1, got %d", cfg.DeliveryAttempts)
	}
}

func TestLoad_HighValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "65535"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("CAPACITY", "10000"); err != nil {
		t.Fatalf("Failed to set CAPACITY: %v", err)
	}
	if err := os.Setenv("DELIVERY_ATTEMPTS", "100"); err != nil {
		t.Fatalf("Failed to set DELIVERY_ATTEMPTS: %v", err)
	}
	if err := os.Setenv("CONSUMER_URL", "https://example.com/very/long/webhook/url/path"); err != nil {
		t.Fatalf("Failed to set CONSUMER_URL: %v", err)
	}

	cfg := Load()

	if cfg.Port != 65535 {
		t.Errorf("expected port 65535, got %d", cfg.Port)
	}
	if cfg.Capacity != 10000 {
		t.Errorf("expected capacity 10000, got %d", cfg.Capacity)
	}
	if cfg.DeliveryAttempts != 100 {
		t.Errorf("expected delivery attempts 100, got %d", cfg.DeliveryAttempts)
	}
	if cfg.ConsumerURL != "https://example.com/very/long/webhook/url/path" {
		t.Errorf("expected long consumer URL, got %s", cfg.ConsumerURL)
	}
}
