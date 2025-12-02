package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config centralizes application configuration loaded from env vars.
type Config struct {
	AppPort      string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	DefaultAdmin string
	DefaultPass  string
	AllowOrigins string
}

// Load reads configuration values while providing safe fallbacks.
func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "mysql"),
		DBPort:       getEnv("DB_PORT", "3306"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", "root"),
		DBName:       getEnv("DB_NAME", "enterprise_kpi"),
		JWTSecret:    getEnv("JWT_SECRET", "change-me"),
		DefaultAdmin: getEnv("DEFAULT_ADMIN_EMAIL", "admin@enterprise.local"),
		DefaultPass:  getEnv("DEFAULT_ADMIN_PASSWORD", "ChangeMe123!"),
		AllowOrigins: getEnv("ALLOW_ORIGINS", "*"),
	}
}

// DSN returns a mysql formatted data source name.
func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
