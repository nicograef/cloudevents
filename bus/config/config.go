package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration values loaded from environment variables.
type Config struct {
	Port        int      // Port for the HTTP server
	Subscribers []string // Webhook URLs to deliver messages
}

// Load reads configuration from environment variables and returns a Config.
// It returns an error if required variables are missing or invalid.
func Load() (Config, error) {
	port := parseEnvInt("PORT", 3000)
	subscriberURLs := parseEnvString("SUBSCRIBER_URLS", "")

	if strings.TrimSpace(subscriberURLs) == "" {
		return Config{}, fmt.Errorf("missing required env SUBSCRIBER_URLS (comma-separated webhook URLs)")
	}

	return Config{
		Port:        port,
		Subscribers: splitAndTrim(subscriberURLs, ","),
	}, nil
}

// parseEnvString reads an environment variable by name and returns its value, or the provided default if unset.
func parseEnvString(name, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}

	return v
}

// parseEnvInt reads an environment variable by name and converts it to int.
// If conversion fails, logs an error and returns the provided default value.
func parseEnvInt(name string, defaultValue int) int {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid %s value: %v\n", name, err)
		return defaultValue
	}

	if n < 1 {
		fmt.Fprintf(os.Stderr, "Invalid %s value: must be at least 1\n", name)
		return defaultValue
	}

	return n
}

// splitAndTrim splits a string by the given separator and trims whitespace from each element.
func splitAndTrim(s, sep string) []string {
	parts := []string{}

	for part := range strings.SplitSeq(s, sep) {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			parts = append(parts, trimmed)
		}
	}

	return parts
}
