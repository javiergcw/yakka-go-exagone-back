package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Logging  LoggingConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port        string
	Environment string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Determine which environment file to load
	env := getEnv("ENVIRONMENT", "development")

	// Load environment-specific .env file
	switch env {
	case "development":
		_ = godotenv.Load(".env.dev")
	case "production":
		_ = godotenv.Load(".env.prod")
	default:
		// Fallback to .env if no specific environment is set
		_ = godotenv.Load()
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnvAsInt("DB_PORT", 0),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
			SSLMode:  getEnv("DB_SSLMODE", ""),
		},
		Server: ServerConfig{
			Port:        getEnv("PORT", ""),
			Environment: env,
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", ""),
		},
	}

	// Validate required configuration
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// validateConfig validates that all required configuration is present
func validateConfig(config *Config) error {
	// Validate database configuration
	if config.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if config.Database.Port == 0 {
		return fmt.Errorf("DB_PORT is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if config.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if config.Database.SSLMode == "" {
		return fmt.Errorf("DB_SSLMODE is required")
	}

	// Validate server configuration
	if config.Server.Port == "" {
		return fmt.Errorf("PORT is required")
	}

	// Validate logging configuration
	if config.Logging.Level == "" {
		return fmt.Errorf("LOG_LEVEL is required")
	}

	return nil
}
