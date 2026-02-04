package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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
	Auth           AuthConfig
	Redis          RedisConfig
	AuditLog       AuditLogConfig
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

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Secret         string
	Issuer         string
	Audience       string
	TokenTTL       time.Duration
	ServiceSubject string
	ServiceRoles   []string
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

// AuditLogConfig holds audit log service configuration.
type AuditLogConfig struct {
	Enabled bool
	URL     string
	Timeout time.Duration
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

	auditTimeoutSeconds, err := strconv.Atoi(getEnv("AUDIT_LOG_SERVICE_TIMEOUT_SECONDS", "3"))
	if err != nil {
		auditTimeoutSeconds = 3
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
			DBName:   getEnv("ORDER_SERVICE_DB_NAME", "appdb"),
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
		Auth: AuthConfig{
			Secret:         getEnv("AUTH_JWT_SECRET", "change-me"),
			Issuer:         getEnv("AUTH_JWT_ISSUER", "enterprise-microservice-system"),
			Audience:       getEnv("AUTH_JWT_AUDIENCE", "enterprise-microservice-system"),
			TokenTTL:       time.Duration(tokenTTLMinutes) * time.Minute,
			ServiceSubject: getEnv("AUTH_SERVICE_SUBJECT", "order-service"),
			ServiceRoles:   getEnvList("AUTH_SERVICE_ROLES", []string{"service"}),
		},
		Redis: RedisConfig{
			Enabled:    getEnvBool("REDIS_ENABLED", true),
			Host:       getEnv("REDIS_HOST", "localhost"),
			Port:       getEnv("REDIS_PORT", "6379"),
			Password:   getEnv("REDIS_PASSWORD", ""),
			DB:         cacheDB,
			DefaultTTL: time.Duration(cacheTTLSeconds) * time.Second,
		},
		AuditLog: AuditLogConfig{
			Enabled: getEnvBool("AUDIT_LOG_SERVICE_ENABLED", true),
			URL:     getEnv("AUDIT_LOG_SERVICE_URL", "http://localhost:8083"),
			Timeout: time.Duration(auditTimeoutSeconds) * time.Second,
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

func getEnvList(key string, defaultValues []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValues
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return defaultValues
	}

	return result
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
