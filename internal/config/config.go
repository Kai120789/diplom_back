package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RunAddress   string
	DbURL        string
	JWTSecret    string
	Salt         string
	LogLevel     string
	Cost         int
	RefreshLive  string
	AccessLive   string
	TokensSecure bool
	AppMode      string
}

var AppConfig *Config

func GetConfig() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	cfg.DbURL = getEnvOrDefault("DB_URL", "")
	cfg.RunAddress = getEnvOrDefault("APP_ADDRESS", ":8080")
	cfg.JWTSecret = getEnvOrDefault("JWT_SECRET", "")
	cfg.Salt = getEnvOrDefault("SALT", "")
	cfg.LogLevel = getEnvOrDefault("LOG_LEVEL", "DEBUG")
	cfg.RefreshLive = getEnvOrDefault("REFRESH_LIVE", "")
	cfg.AccessLive = getEnvOrDefault("ACCESS_LIVE", "")
	cfg.AppMode = getEnvOrDefault("APP_MODE", "")

	switch cfg.AppMode {
	case "development":
		cfg.TokensSecure = false
	case "production":
		cfg.TokensSecure = true
	default:
		cfg.TokensSecure = false
	}

	cost, err := getIntEnvOrDefault("COST", 0)
	if err != nil {
		return nil, errors.New("cost is uncorrect")
	}

	cfg.Cost = cost

	AppConfig = &cfg

	return &cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) (int, error) {
	if value := os.Getenv(key); value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return intValue, nil
	}
	return defaultValue, nil
}
