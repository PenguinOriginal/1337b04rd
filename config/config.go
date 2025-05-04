package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	LogFilePath string
	UploadDir   string
}

func LoadConfig() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "password"),
		DBName:      getEnv("DB_NAME", "leetboard"),
		LogFilePath: getEnv("LOG_FILE_PATH", "logging/logging.log"),
		UploadDir:   getEnv("UPLOAD_DIR", "/data"),
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
