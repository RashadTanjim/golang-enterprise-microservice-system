package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for repository service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	Auth     AuthConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port      string
	RateLimit int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Secret   string
	Issuer   string
	Audience string
	TokenTTL time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (optional in production)
	_ = godotenv.Load()

	rateLimit, err := strconv.Atoi(getEnv("REPOSITORY_SERVICE_RATE_LIMIT", "100"))
	if err != nil {
		rateLimit = 100
	}

	tokenTTLMinutes, err := strconv.Atoi(getEnv("AUTH_TOKEN_TTL_MINUTES", "60"))
	if err != nil {
		tokenTTLMinutes = 60
	}

	config := &Config{
		Server: ServerConfig{
			Port:      getEnv("REPOSITORY_SERVICE_PORT", "8083"),
			RateLimit: rateLimit,
		},
		Database: DatabaseConfig{
			Host:     getEnv("REPOSITORY_SERVICE_DB_HOST", "localhost"),
			Port:     getEnv("REPOSITORY_SERVICE_DB_PORT", "5432"),
			User:     getEnv("REPOSITORY_SERVICE_DB_USER", "postgres"),
			Password: getEnv("REPOSITORY_SERVICE_DB_PASSWORD", "postgres"),
			DBName:   getEnv("REPOSITORY_SERVICE_DB_NAME", "repositorydb"),
		},
		Log: LogConfig{
			Level: getEnv("REPOSITORY_SERVICE_LOG_LEVEL", "info"),
		},
		Auth: AuthConfig{
			Secret:   getEnv("AUTH_JWT_SECRET", "change-me"),
			Issuer:   getEnv("AUTH_JWT_ISSUER", "enterprise-microservice-system"),
			Audience: getEnv("AUTH_JWT_AUDIENCE", "enterprise-microservice-system"),
			TokenTTL: time.Duration(tokenTTLMinutes) * time.Minute,
		},
	}

	return config, nil
}

// DSN returns the database connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
