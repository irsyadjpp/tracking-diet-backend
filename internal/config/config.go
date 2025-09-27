package config

import (
    "os"
)

type Config struct {
    DBUrl string
    Port  string
}

func LoadConfig() *Config {
    return &Config{
        DBUrl: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/tracking_diet"),
        Port:  getEnv("PORT", "8080"),
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
