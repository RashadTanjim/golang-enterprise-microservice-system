package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// DatabaseConfig holds database configuration for a service.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// DSN returns the database connection string.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}

// Config holds migration service configuration.
type Config struct {
	LogLevel string
	UserDB   DatabaseConfig
	OrderDB  DatabaseConfig
}

// Load reads configuration from environment variables.
func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		LogLevel: getEnv("MIGRATION_SERVICE_LOG_LEVEL", "info"),
		UserDB: DatabaseConfig{
			Host:     getEnv("USER_SERVICE_DB_HOST", "localhost"),
			Port:     getEnv("USER_SERVICE_DB_PORT", "5432"),
			User:     getEnv("USER_SERVICE_DB_USER", "postgres"),
			Password: getEnv("USER_SERVICE_DB_PASSWORD", "postgres"),
			DBName:   getEnv("USER_SERVICE_DB_NAME", "appdb"),
		},
		OrderDB: DatabaseConfig{
			Host:     getEnv("ORDER_SERVICE_DB_HOST", "localhost"),
			Port:     getEnv("ORDER_SERVICE_DB_PORT", "5432"),
			User:     getEnv("ORDER_SERVICE_DB_USER", "postgres"),
			Password: getEnv("ORDER_SERVICE_DB_PASSWORD", "postgres"),
			DBName:   getEnv("ORDER_SERVICE_DB_NAME", "appdb"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
