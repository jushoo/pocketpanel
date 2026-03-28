package config

import (
	"os"
	"strings"
)

type Config struct {
	Port         string
	DatabasePath string
	Environment  string
	CORSOrigins  []string
	ServersPath  string
}

func Load() (*Config, error) {
	corsOrigins := getEnv("CORS_ORIGINS", "http://localhost:3000")

	return &Config{
		Port:         getEnv("PORT", ":3001"),
		DatabasePath: getEnv("DATABASE_PATH", "pocketpanel.db"),
		Environment:  getEnv("ENV", "development"),
		CORSOrigins:  strings.Split(corsOrigins, ","),
		ServersPath:  getEnv("SERVERS_PATH", "apps/api/tmp"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
