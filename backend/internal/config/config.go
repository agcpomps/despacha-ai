package config

import (
	"os"
	"strings"
)

type Config struct {
	AppEnv             string
	Port               string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	JWTSecret          string
	CORSAllowedOrigins []string
	MigrationsPath     string
	PublicBaseURL      string
}

func Load() *Config {
	return &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "despacha_ai"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET", "change_me"),
		CORSAllowedOrigins: getCSVEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:4173"),
		MigrationsPath:     getEnv("MIGRATIONS_PATH", "migrations"),
		// origem pública para construir URLs de imagens (ex: https://despachaai.com).
		// Vazio em dev → usa o host do pedido (localhost:8080).
		PublicBaseURL: getEnv("PUBLIC_BASE_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getCSVEnv(key, fallback string) []string {
	value := getEnv(key, fallback)
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
