package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type JWTConfig struct {
	SecretKey      string
	ExpirationHour int
}

type AppConfig struct {
	Name     string
	Version  string
	LogLevel string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "showroom_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Env:  getEnv("APP_ENV", "development"),
		},
		JWT: JWTConfig{
			SecretKey:      getEnv("JWT_SECRET_KEY", "default-secret-key"),
			ExpirationHour: getEnvAsInt("JWT_EXPIRATION_HOUR", 24),
		},
		App: AppConfig{
			Name:     getEnv("APP_NAME", "Showroom Management System"),
			Version:  getEnv("APP_VERSION", "1.0.0"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
	}
}

func (j *JWTConfig) GetExpiration() time.Duration {
	return time.Duration(j.ExpirationHour) * time.Hour
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}