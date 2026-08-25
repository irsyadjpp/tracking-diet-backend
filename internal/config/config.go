package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env         string
	DatabaseURL string
	Port        string
	
	// JWT Configuration
	JWTSecret     string
	JWTExpiration int // in hours
	
	// Google Genkit Configuration
	GoogleGenkitAPIKey string
	
	// Redis Configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	
	// Database Connection Pool Settings
	DBMaxIdleConns    int
	DBMaxOpenConns    int
	DBConnMaxLifetime time.Duration
	
	// Server Configuration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Env:         getEnv("ENV", "development"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost:5432/tracking_diet"),
		Port:        getEnv("PORT", "8080"),
		
		// JWT Configuration
		JWTSecret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 24), // 24 hours default
		
		// Google Genkit Configuration
		GoogleGenkitAPIKey: getEnv("GOOGLE_GENKIT_API_KEY", ""),
		
		// Redis Configuration
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
		
		// Database Connection Pool Settings
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		DBConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "1h"),
		
		// Server Configuration
		ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", "15s"),
		WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", "15s"),
		IdleTimeout:  getEnvAsDuration("IDLE_TIMEOUT", "60s"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Hour // fallback
}
