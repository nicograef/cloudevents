package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoad_MissingSubscribers(t *testing.T) {
	os.Clearenv()

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "SUBSCRIBER_URLS") {
		t.Fatalf("expected error about SUBSCRIBER_URLS, got %v", err)
	}
}

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("SUBSCRIBER_URLS", "http://test/webhook"); err != nil {
		t.Fatalf("Failed to set SUBSCRIBER_URL: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Port != 3000 {
		t.Errorf("expected default port 3000, got %d", cfg.Port)
	}
}

func TestLoad_EnvValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "8080"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("SUBSCRIBER_URLS", "http://test/webhook"); err != nil {
		t.Fatalf("Failed to set SUBSCRIBER_URL: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if len(cfg.Subscribers) != 1 || cfg.Subscribers[0] != "http://test/webhook" {
		t.Errorf("expected subscriber URL 'http://test/webhook', got %v", cfg.Subscribers)
	}

}

func TestLoad_InvalidInt(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("SUBSCRIBER_URLS", "http://test/webhook"); err != nil {
		t.Fatalf("Failed to set SUBSCRIBER_URL: %v", err)
	}
	if err := os.Setenv("PORT", "notanint"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Port != 3000 {
		t.Errorf("expected fallback port 3000, got %d", cfg.Port)
	}
}

func TestLoad_NegativeValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("SUBSCRIBER_URLS", "http://test/webhook"); err != nil {
		t.Fatalf("Failed to set SUBSCRIBER_URL: %v", err)
	}
	if err := os.Setenv("PORT", "-1"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Should fallback to defaults due to validation (must be at least 1)
	if cfg.Port != 3000 {
		t.Errorf("expected fallback port 3000 for negative value, got %d", cfg.Port)
	}
}

func TestLoad_ValidEdgeCases(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("SUBSCRIBER_URLS", "http://test/webhook"); err != nil {
		t.Fatalf("Failed to set SUBSCRIBER_URL: %v", err)
	}
	if err := os.Setenv("PORT", "1"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Port != 1 {
		t.Errorf("expected port 1, got %d", cfg.Port)
	}
}

func TestLoad_HighValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "65535"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("SUBSCRIBER_URLS", "https://example.com/very/long/webhook/url/path"); err != nil {
		t.Fatalf("Failed to set SUBSCRIBER_URLS: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Port != 65535 {
		t.Errorf("expected port 65535, got %d", cfg.Port)
	}
	if cfg.Subscribers[0] != "https://example.com/very/long/webhook/url/path" {
		t.Errorf("expected long subscriber URL, got %s", cfg.Subscribers[0])
	}
}
