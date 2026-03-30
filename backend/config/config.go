package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port              string
	MaxFileSize       int64 // in bytes
	AllowedMimeTypes  []string
	MaxImageDimension int
}

func Load() Config {
	return Config{
		Port:              getEnv("PORT", "4000"),
		MaxFileSize:       getEnvInt64("MAX_FILE_SIZE", 10*1024*1024),
		AllowedMimeTypes:  getEnvSlice("ALLOWED_MIME_TYPES", "image/jpeg,image/png,image/gif,image/webp,image/avif,image/svg+xml"),
		MaxImageDimension: int(getEnvInt64("MAX_IMAGE_DIMENSION", 1920)),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if val := os.Getenv(key); val != "" {
		if num, err := strconv.ParseInt(val, 10, 64); err == nil {
			return num
		}
	}
	return fallback
}

func getEnvSlice(key, fallback string) []string {
	raw := getEnv(key, fallback)
	parts := strings.Split(raw, ",")
	variables := make([]string, 0, len(parts))
	for _, part := range parts {
		if variable := strings.TrimSpace(part); variable != "" {
			variables = append(variables, variable)
		}
	}
	return variables
}
