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
		t.Errorf("expected default queue size 1000, got %d", cfg.Capacity)
	}
	if cfg.ConsumerURL != "http://localhost:4000" {
		t.Errorf("expected default ConsumerURL, got %s", cfg.ConsumerURL)
	}
	if cfg.DeliveryAttempts != 3 {
		t.Errorf("expected default DeliveryAttempts 3, got %d", cfg.DeliveryAttempts)
	}
}

func TestLoad_EnvValues(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("CAPACITY", "42")
	os.Setenv("CONSUMER_URL", "http://test/webhook")
	os.Setenv("DELIVERY_ATTEMPTS", "5")
	cfg := Load()
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.Capacity != 42 {
		t.Errorf("expected queue size 42, got %d", cfg.Capacity)
	}
	if cfg.ConsumerURL != "http://test/webhook" {
		t.Errorf("expected ConsumerURL http://test/webhook, got %s", cfg.ConsumerURL)
	}
	if cfg.DeliveryAttempts != 5 {
		t.Errorf("expected DeliveryAttempts 5, got %d", cfg.DeliveryAttempts)
	}
}

func TestLoad_InvalidIntAndLowValues(t *testing.T) {
	os.Setenv("PORT", "notanint")
	os.Setenv("CAPACITY", "badint")
	os.Setenv("DELIVERY_ATTEMPTS", "0")
	cfg := Load()
	if cfg.Port != 3000 {
		t.Errorf("expected fallback port 3000, got %d", cfg.Port)
	}
	if cfg.Capacity != 1000 {
		t.Errorf("expected fallback queue size 1000, got %d", cfg.Capacity)
	}
	if cfg.DeliveryAttempts != 3 {
		t.Errorf("expected fallback DeliveryAttempts 3, got %d", cfg.DeliveryAttempts)
	}
}
