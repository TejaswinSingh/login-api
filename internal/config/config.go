package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	db "github.com/TejaswinSingh/login-api/internal/db/postgres"
)

type Config struct {
	Env       GoEnv
	HttpPort  int
	SecretKey []byte
	DbConfig  db.DbConfig
}

func NewEnvConfig() Config {

	secretKey := []byte(requireEnv("SECRET_KEY"))
	if len(secretKey) < 32 {
		panic(fmt.Errorf("SECRET_KEY must be at least 32 bytes, got %d", len(secretKey)))
	}

	return Config{
		Env:       ValidateEnv(GoEnv(os.Getenv("ENV"))),
		HttpPort:  envInt("HTTP_PORT", default_http_port),
		SecretKey: secretKey,
		DbConfig: db.DbConfig{
			DbUser:     requireEnv("DB_USER"),
			DbPassword: requireEnv("DB_PASSWORD"),
			DbHost:     requireEnv("DB_HOST"),
			DbPort:     requireEnvInt("DB_PORT"),
			DbName:     requireEnv("DB_NAME"),
			DbSslMode:  requireEnv("DB_SSL_MODE"),
		},
	}
}

func envInt(key string, fallback int) int {
	val, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key)))
	if err != nil {
		return fallback
	}
	return val
}

func requireEnv(key string) string {
	val, ok := os.LookupEnv(key)
	val = strings.TrimSpace(val)
	if !ok || val == "" {
		panic(fmt.Errorf("invalid environment config: %s is required", key))
	}
	return val
}

func requireEnvInt(key string) int {
	val := requireEnv(key)
	valInt, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("invalid environment config: %s should be of type integer", key))
	}
	return valInt
}
