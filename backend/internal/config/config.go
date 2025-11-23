package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	CORS     CORSConfig
	Upload   UploadConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type UploadConfig struct {
	MaxSize   int64
	UploadDir string
}

type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

func Load() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	_ = godotenv.Load()

	maxSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "100"))
	rateDuration, _ := time.ParseDuration(getEnv("RATE_LIMIT_DURATION", "1m"))

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "issueapp"),
			Password: getEnv("DB_PASSWORD", "issueapp_password"),
			DBName:   getEnv("DB_NAME", "community_issues"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseOrigins(getEnv("CORS_ALLOWED_ORIGINS", "*")),
		},
		Upload: UploadConfig{
			MaxSize:   maxSize,
			UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
		},
		RateLimit: RateLimitConfig{
			Requests: rateLimit,
			Duration: rateDuration,
		},
	}

	return config, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseOrigins(origins string) []string {
	if origins == "*" {
		return []string{"*"}
	}
	// Simple comma-separated parsing
	result := []string{}
	current := ""
	for _, char := range origins {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
