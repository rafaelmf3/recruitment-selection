// Package config loads application configuration from environment variables.
// Call Load() once at startup; the returned Config is safe to pass around.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration values.
type Config struct {
	Port               string
	GinMode            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	JWTExpirationHours int
	UploadDir          string
	MaxUploadSizeMB    int64
	AllowedOrigins     []string
}

// Load reads .env (if present) then environment variables.
// Returns an error if any required value is missing.
func Load() (*Config, error) {
	// .env is optional - ignore error in production containers
	_ = godotenv.Load()

	jwtExpHours, err := strconv.Atoi(getenv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("config: invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	maxUploadMB, err := strconv.ParseInt(getenv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("config: invalid MAX_UPLOAD_SIZE_MB: %w", err)
	}

	rawOrigins := getenv("ALLOWED_ORIGINS", "http://localhost:5173")
	origins := strings.Split(rawOrigins, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}

	cfg := &Config{
		Port:               getenv("PORT", "8080"),
		GinMode:            getenv("GIN_MODE", "debug"),
		DBHost:             getenv("DB_HOST", "localhost"),
		DBPort:             getenv("DB_PORT", "5432"),
		DBUser:             getenv("DB_USER", "postgres"),
		DBPassword:         getenv("DB_PASSWORD", "postgres"),
		DBName:             getenv("DB_NAME", "recruitment_selection"),
		DBSSLMode:          getenv("DB_SSLMODE", "disable"),
		JWTSecret:          getenv("JWT_SECRET", ""),
		JWTExpirationHours: jwtExpHours,
		UploadDir:          getenv("UPLOAD_DIR", "./uploads"),
		MaxUploadSizeMB:    maxUploadMB,
		AllowedOrigins:     origins,
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("config: JWT_SECRET must be set")
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
