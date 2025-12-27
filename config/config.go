package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	Pagination PaginationConfig `json:"pagination"`
	Validation ValidationConfig `json:"validation"`
	Logging    LoggingConfig    `json:"logging"`
	JWT        JWTConfig        `json:"jwt"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	Host         string `json:"host"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
}

type DatabaseConfig struct {
	Host            string `json:"host"`
	Port            string `json:"port"`
	User            string `json:"user"`
	Password        string `json:"-"` // No serializar en JSON
	Name            string `json:"name"`
	SSLMode         string `json:"ssl_mode"`
	MaxConnections  int    `json:"max_connections"`
	MaxIdleConns    int    `json:"max_idle_connections"`
	ConnMaxLifetime string `json:"conn_max_lifetime"`
}

type PaginationConfig struct {
	DefaultLimit int `json:"default_limit"`
	MaxLimit     int `json:"max_limit"`
}

type ValidationConfig struct {
	MaxNameLength        int `json:"max_name_length"`
	MaxDescriptionLength int `json:"max_description_length"`
}

type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

type JWTConfig struct {
	SecretKey      string `json:"-"`
	ExpirationTime string `json:"expiration_time"`
}

func Load(configPath string) (*Config, error) {
	_ = godotenv.Load()

	cleanPath := filepath.Clean(configPath)
	file, err := os.Open(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Error closing file, but nothing we can do at this point
			_ = err
		}
	}()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	config.Database.Password = os.Getenv("DB_PASSWORD")
	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	config.Database.Host = getEnvOrDefault("DB_HOST", config.Database.Host)
	config.Database.Port = getEnvOrDefault("DB_PORT", config.Database.Port)
	config.Database.User = getEnvOrDefault("DB_USER", config.Database.User)
	config.Database.Name = getEnvOrDefault("DB_NAME", config.Database.Name)
	config.Database.SSLMode = getEnvOrDefault("DB_SSLMODE", config.Database.SSLMode)

	config.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")
	if config.JWT.SecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is required")
	}

	return &config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) GetConnMaxLifetime() (time.Duration, error) {
	return time.ParseDuration(c.ConnMaxLifetime)
}
