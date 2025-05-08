package config

import (
	"log"
	"os"
)

type Config struct {
	Port                string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	LogFilePath         string
	UploadDir           string
	SessionCookieName   string
	SessionDurationDays string
	AvatarAPIBaseURL    string
}

func LoadConfig() *Config {
	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "password"),
		DBName:              getEnv("DB_NAME", "leetboard"),
		LogFilePath:         getEnv("LOG_FILE_PATH", "logging/logging.log"),
		UploadDir:           getEnv("UPLOAD_DIR", "/data"),
		SessionCookieName:   getEnv("SESSION_COOKIE_NAME", "session_id"),
		SessionDurationDays: getEnv("SESSION_DURATION_DAYS", "7"),
		AvatarAPIBaseURL:    getEnv("AVATAR_API_BASE_URL", "https://rickandmortyapi.com/api/character"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Printf("Warning: %s not set, using default: %s", key, fallback)
		return fallback
	}
	return val
}
