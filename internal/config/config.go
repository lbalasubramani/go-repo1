package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port     int
	LogLevel string
	Env      string
}

func Load() (*Config, error) {
	port := getEnvAsInt("PORT", 8080)
	logLevel := getEnv("LOG_LEVEL", "info")
	env := getEnv("ENV", "development")

	return &Config{
		Port:     port,
		LogLevel: logLevel,
		Env:      env,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
