package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()
	if cfg.Port != 5000 {
		t.Errorf("expected default port 5000, got %d", cfg.Port)
	}

}

func TestLoad_EnvValues(t *testing.T) {
	os.Setenv("PORT", "8080")
	cfg := Load()
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
}

func TestLoad_InvalidIntAndLowValues(t *testing.T) {
	os.Setenv("PORT", "notanint")
	cfg := Load()
	if cfg.Port != 5000 {
		t.Errorf("expected fallback port 5000, got %d", cfg.Port)
	}
}
