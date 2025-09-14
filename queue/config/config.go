package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration values loaded from environment variables.
type Config struct {
	Port             int    // Port for the HTTP server
	Capacity         int    // Maximum number of messages in the queue
	ConsumerURL      string // Webhook URL to deliver messages
	DeliveryAttempts int    // Number of attempts for delivering a message
}

// Load reads configuration from environment variables and returns a Config struct.
// Defaults: PORT=3000 CAPACITY=1000, CONSUMER_URL="http://localhost:4000" DELIVERY_ATTEMPTS=3
func Load() Config {
	port := parseEnvInt("PORT", 3000)
	capacity := parseEnvInt("CAPACITY", 1000)
	consumerURL := parseEnvString("CONSUMER_URL", "http://localhost:4000")
	deliveryAttempts := parseEnvInt("DELIVERY_ATTEMPTS", 3)

	return Config{
		Port:             port,
		Capacity:         capacity,
		ConsumerURL:      consumerURL,
		DeliveryAttempts: deliveryAttempts,
	}
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
