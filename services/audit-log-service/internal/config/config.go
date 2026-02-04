package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for audit log service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	Auth     AuthConfig
	Redis    RedisConfig
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

// RedisConfig holds Redis cache configuration
type RedisConfig struct {
	Enabled    bool
	Host       string
	Port       string
	Password   string
	DB         int
	DefaultTTL time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (optional in production)
	_ = godotenv.Load()

	rateLimit, err := strconv.Atoi(getEnv("AUDIT_LOG_SERVICE_RATE_LIMIT", "100"))
	if err != nil {
		rateLimit = 100
	}

	tokenTTLMinutes, err := strconv.Atoi(getEnv("AUTH_TOKEN_TTL_MINUTES", "60"))
	if err != nil {
		tokenTTLMinutes = 60
	}

	cacheTTLSeconds, err := strconv.Atoi(getEnv("REDIS_TTL_SECONDS", "300"))
	if err != nil {
		cacheTTLSeconds = 300
	}

	cacheDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		cacheDB = 0
	}

	config := &Config{
		Server: ServerConfig{
			Port:      getEnv("AUDIT_LOG_SERVICE_PORT", "8083"),
			RateLimit: rateLimit,
		},
		Database: DatabaseConfig{
			Host:     getEnv("AUDIT_LOG_SERVICE_DB_HOST", "localhost"),
			Port:     getEnv("AUDIT_LOG_SERVICE_DB_PORT", "5432"),
			User:     getEnv("AUDIT_LOG_SERVICE_DB_USER", "postgres"),
			Password: getEnv("AUDIT_LOG_SERVICE_DB_PASSWORD", "postgres"),
			DBName:   getEnv("AUDIT_LOG_SERVICE_DB_NAME", "appdb"),
		},
		Log: LogConfig{
			Level: getEnv("AUDIT_LOG_SERVICE_LOG_LEVEL", "info"),
		},
		Auth: AuthConfig{
			Secret:   getEnv("AUTH_JWT_SECRET", "change-me"),
			Issuer:   getEnv("AUTH_JWT_ISSUER", "enterprise-microservice-system"),
			Audience: getEnv("AUTH_JWT_AUDIENCE", "enterprise-microservice-system"),
			TokenTTL: time.Duration(tokenTTLMinutes) * time.Minute,
		},
		Redis: RedisConfig{
			Enabled:    getEnvBool("REDIS_ENABLED", true),
			Host:       getEnv("REDIS_HOST", "localhost"),
			Port:       getEnv("REDIS_PORT", "6379"),
			Password:   getEnv("REDIS_PASSWORD", ""),
			DB:         cacheDB,
			DefaultTTL: time.Duration(cacheTTLSeconds) * time.Second,
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

func getEnvBool(key string, defaultValue bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return defaultValue
	}

	switch value {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return defaultValue
	}
}
