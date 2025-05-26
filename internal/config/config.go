package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
    GRPCPort int
	Env      string
    LogLevel string
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() (*Config, error) {
	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		portStr = "50051" // default port if not set
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid GRPC_PORT: %v", err)
	}

	return &Config{
		GRPCPort: port,
	}, nil
}

// getEnv fetches an environment variable or returns the fallback
func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}