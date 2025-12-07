package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort      string
	MongoURI        string
	MongoDatabase   string
	N8NWebhookURL   string
	EnableScheduler bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	config := &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase:   getEnv("MONGO_DATABASE", "restysched"),
		N8NWebhookURL:   getEnv("N8N_WEBHOOK_URL", ""),
		EnableScheduler: getEnv("ENABLE_SCHEDULER", "true") == "true",
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.N8NWebhookURL == "" {
		return fmt.Errorf("N8N_WEBHOOK_URL is required")
	}
	if c.MongoURI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}
	if c.MongoDatabase == "" {
		return fmt.Errorf("MONGO_DATABASE is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
