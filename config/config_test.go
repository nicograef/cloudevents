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
	if cfg.ConsumerUrl != "http://localhost:4000" {
		t.Errorf("expected default consumerUrl, got %s", cfg.ConsumerUrl)
	}
}

func TestLoad_EnvValues(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("QUEUE_SIZE", "42")
	os.Setenv("CONSUMER_URL", "http://test/webhook")
	cfg := Load()
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.Capacity != 42 {
		t.Errorf("expected queue size 42, got %d", cfg.Capacity)
	}
	if cfg.ConsumerUrl != "http://test/webhook" {
		t.Errorf("expected consumerUrl http://test/webhook, got %s", cfg.ConsumerUrl)
	}
}

func TestLoad_InvalidInt(t *testing.T) {
	os.Setenv("PORT", "notanint")
	os.Setenv("QUEUE_SIZE", "badint")
	cfg := Load()
	if cfg.Port != 3000 {
		t.Errorf("expected fallback port 3000, got %d", cfg.Port)
	}
	if cfg.Capacity != 1000 {
		t.Errorf("expected fallback queue size 1000, got %d", cfg.Capacity)
	}
}
