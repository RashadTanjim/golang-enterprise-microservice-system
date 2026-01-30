package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for order service
type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	Log            LogConfig
	UserService    UserServiceConfig
	CircuitBreaker CircuitBreakerConfig
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

// UserServiceConfig holds user service configuration
type UserServiceConfig struct {
	URL string
}

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	MaxRequests uint32
	Interval    time.Duration
	Timeout     time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file (optional in production)
	_ = godotenv.Load()

	rateLimit, err := strconv.Atoi(getEnv("ORDER_SERVICE_RATE_LIMIT", "100"))
	if err != nil {
		rateLimit = 100
	}

	maxRequests, err := strconv.ParseUint(getEnv("CIRCUIT_BREAKER_MAX_REQUESTS", "3"), 10, 32)
	if err != nil {
		maxRequests = 3
	}

	interval, err := strconv.Atoi(getEnv("CIRCUIT_BREAKER_INTERVAL", "60"))
	if err != nil {
		interval = 60
	}

	timeout, err := strconv.Atoi(getEnv("CIRCUIT_BREAKER_TIMEOUT", "30"))
	if err != nil {
		timeout = 30
	}

	config := &Config{
		Server: ServerConfig{
			Port:      getEnv("ORDER_SERVICE_PORT", "8082"),
			RateLimit: rateLimit,
		},
		Database: DatabaseConfig{
			Host:     getEnv("ORDER_SERVICE_DB_HOST", "localhost"),
			Port:     getEnv("ORDER_SERVICE_DB_PORT", "5432"),
			User:     getEnv("ORDER_SERVICE_DB_USER", "postgres"),
			Password: getEnv("ORDER_SERVICE_DB_PASSWORD", "postgres"),
			DBName:   getEnv("ORDER_SERVICE_DB_NAME", "orderdb"),
		},
		Log: LogConfig{
			Level: getEnv("ORDER_SERVICE_LOG_LEVEL", "info"),
		},
		UserService: UserServiceConfig{
			URL: getEnv("ORDER_SERVICE_USER_SERVICE_URL", "http://localhost:8081"),
		},
		CircuitBreaker: CircuitBreakerConfig{
			MaxRequests: uint32(maxRequests),
			Interval:    time.Duration(interval) * time.Second,
			Timeout:     time.Duration(timeout) * time.Second,
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
