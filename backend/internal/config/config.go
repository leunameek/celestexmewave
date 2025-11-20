package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBURL      string // Full connection URL for Supabase
	DBDriver   string // pgx for pgx/v5 driver

	// Server
	ServerPort int
	ServerEnv  string
	ServerHost string

	// JWT
	JWTSecret              string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration

	// Email
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string

	// Frontend
	FrontendURL string

	// File Upload
	UploadDir     string
	MaxUploadSize int64
}

var cfg *Config

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		// Database - Prioritize Supabase (DATABASE_URL) over individual parameters
		DBHost:     getEnv("DB_HOST", ""),
		DBPort:     getEnvInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "postgres"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "require"),
		DBURL:      getEnv("DATABASE_URL", ""),      // Supabase connection URL (REQUIRED)
		DBDriver:   getEnv("DB_DRIVER", "postgres"), // postgres or pgx

		// Server
		ServerPort: getEnvInt("SERVER_PORT", 8080),
		ServerEnv:  getEnv("SERVER_ENV", "development"),
		ServerHost: getEnv("SERVER_HOST", "localhost"),

		// JWT
		JWTSecret:              getEnv("JWT_SECRET", "your_secret_key_change_in_production"),
		JWTExpiration:          parseDuration(getEnv("JWT_EXPIRATION", "24h")),
		RefreshTokenExpiration: parseDuration(getEnv("REFRESH_TOKEN_EXPIRATION", "7d")),

		// Email
		SMTPHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort: getEnvInt("SMTP_PORT", 587),
		SMTPUser: getEnv("SMTP_USER", ""),
		SMTPPass: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom: getEnv("SMTP_FROM", "noreply@celestexmewave.com"),

		// Frontend
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		// File Upload
		UploadDir:     getEnv("UPLOAD_DIR", "../assets/images"),
		MaxUploadSize: getEnvInt64("MAX_UPLOAD_SIZE", 5242880), // 5MB
	}

	cfg = config
	return config, nil
}

// Get returns the global config instance
func Get() *Config {
	if cfg == nil {
		Load()
	}
	return cfg
}

// GetDSN returns the PostgreSQL connection string
// Prioritizes Supabase (DATABASE_URL) over traditional connection parameters
func (c *Config) GetDSN() string {
	// Supabase is the primary option - DATABASE_URL must be set
	if c.DBURL != "" {
		return c.DBURL
	}

	// Fallback to individual parameters (for local development only)
	if c.DBHost != "" && c.DBUser != "" && c.DBPassword != "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.DBHost,
			c.DBPort,
			c.DBUser,
			c.DBPassword,
			c.DBName,
			c.DBSSLMode,
		)
	}

	// No valid configuration found
	return ""
}

// Helper functions
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvInt64(key string, defaultVal int64) int64 {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseInt(valStr, 10, 64); err == nil {
		return val
	}
	return defaultVal
}

func parseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour // default to 24 hours
	}
	return duration
}
