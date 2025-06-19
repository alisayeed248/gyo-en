package config

import (
    "os"
    "strconv"
)

type Config struct {
    RedisAddr     string
    Port          string
    CheckInterval int
    Environment   string
}

func Load() *Config {
    return &Config{
        RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
        Port:          getEnv("PORT", "8080"),
        CheckInterval: getEnvInt("CHECK_INTERVAL", 30),
        Environment:   getEnv("ENVIRONMENT", "development"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if i, err := strconv.Atoi(value); err == nil {
            return i
        }
    }
    return defaultValue
}