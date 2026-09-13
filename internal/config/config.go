package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	ServerPort string
	JWTSecret  []byte
	JWTTTL     time.Duration
}

func Load(envPath string) (*Config, error) {
	_ = godotenv.Load(envPath) // читает .env

	cfg := &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 chars")
	}
	cfg.JWTSecret = []byte(secret)

	ttlHours, _ := strconv.Atoi(getEnv("JWT_TTL_HOURS", "24"))
	cfg.JWTTTL = time.Duration(ttlHours) * time.Hour

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		c.DBUser, c.DBPassword, c.DBName, c.DBHost, c.DBPort,
	)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
