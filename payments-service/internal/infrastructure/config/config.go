package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	Port        string
	KafkaBroker string
	KafkaTopic  string
}

func Load() (*Config, error) {
	// In production, real env vars are injected directly (Docker/K8s secrets),
	// so a missing .env.local file is expected and NOT an error — only log it.
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Println("no .env.local file found, relying on real environment variables")
	}

	cfg := &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		Port:        getEnv("PORT", "8080"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "payments"),
	}

	var missing []string

	cfg.DBUser = requireEnv("DB_USER", &missing)
	cfg.DBPassword = requireEnv("DB_PASSWORD", &missing)
	cfg.DBName = requireEnv("DB_NAME", &missing)

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func requireEnv(key string, missing *[]string) string {
	value := os.Getenv(key)
	if value == "" {
		*missing = append(*missing, key)
	}
	return value
}
