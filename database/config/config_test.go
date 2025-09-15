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
	if cfg.DataDir != "." {
		t.Errorf("expected default data directory '.', got %s", cfg.DataDir)
	}
}

func TestLoad_EnvValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "8080"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("DATA_DIR", "/tmp/testdata"); err != nil {
		t.Fatalf("Failed to set DATA_DIR: %v", err)
	}

	cfg := Load()

	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
	if cfg.DataDir != "/tmp/testdata" {
		t.Errorf("expected data directory '/tmp/testdata', got %s", cfg.DataDir)
	}
}

func TestLoad_InvalidIntAndLowValues(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "notanint"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}
	if err := os.Setenv("DATA_DIR", ""); err != nil { // Empty string should use default
		t.Fatalf("Failed to set DATA_DIR: %v", err)
	}

	cfg := Load()

	if cfg.Port != 5000 {
		t.Errorf("expected fallback port 5000, got %d", cfg.Port)
	}
	if cfg.DataDir != "." {
		t.Errorf("expected fallback data directory '.', got %s", cfg.DataDir)
	}
}

func TestLoad_LowPortValue(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "0"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}

	cfg := Load()

	// Should fallback to default due to validation (must be at least 1)
	if cfg.Port != 5000 {
		t.Errorf("expected fallback port 5000 for invalid low value, got %d", cfg.Port)
	}
}

func TestLoad_NegativePortValue(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "-1"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}

	cfg := Load()

	// Should fallback to default due to validation (must be at least 1)
	if cfg.Port != 5000 {
		t.Errorf("expected fallback port 5000 for negative value, got %d", cfg.Port)
	}
}

func TestLoad_ValidPortEdgeCase(t *testing.T) {
	os.Clearenv()

	if err := os.Setenv("PORT", "1"); err != nil {
		t.Fatalf("Failed to set PORT: %v", err)
	}

	cfg := Load()

	if cfg.Port != 1 {
		t.Errorf("expected port 1, got %d", cfg.Port)
	}
}
