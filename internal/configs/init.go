package configs

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/joho/godotenv"
)

var AppSettings Configs

const defaultPostgresDatabase = "analogyoutube"

func Load() error {
	log := logger.GetLogger()
	log.Info().Msg("loading application settings")

	_ = godotenv.Load(".env")

	jwtSecret, err := envRequired("JWT_SECRET")
	if err != nil {
		return err
	}

	postgresUser, err := envRequired("POSTGRES_USER")
	if err != nil {
		return err
	}

	postgresPassword, err := envRequired("POSTGRES_PASSWORD")
	if err != nil {
		return err
	}

	postgresDatabase := envString("POSTGRES_DATABASE", defaultPostgresDatabase)
	if postgresDatabase != defaultPostgresDatabase {
		return fmt.Errorf("POSTGRES_DATABASE must be %s", defaultPostgresDatabase)
	}

	AppSettings = Configs{
		AppParams: AppParams{
			ServerURL:  envString("SERVER_URL", "localhost"),
			ServerName: envString("SERVER_NAME", "GlobalServer"),
			PortRun:    envString("SERVER_PORT", "6666"),
			GinMode:    envString("GIN_MODE", "debug"),
		},
		PostgresParams: PostgresParams{
			Host:     envString("POSTGRES_HOST", "localhost"),
			Port:     envString("POSTGRES_PORT", "5432"),
			User:     postgresUser,
			Password: postgresPassword,
			Database: postgresDatabase,
		},
		AuthParams: AuthParams{
			AccessTokenTtlMinutes: envInt("ACCESS_TOKEN_TTL_MINUTES", 15),
			RefreshTokenTtlDays:   envInt("REFRESH_TOKEN_TTL_DAYS", 30),
			JwtSecret:             jwtSecret,
		},
		RedisParams: RedisParams{
			Host:     envString("REDIS_HOST", "localhost"),
			Port:     envString("REDIS_PORT", "6379"),
			Password: envString("REDIS_PASSWORD", ""),
			DB:       envInt("REDIS_DB", 0),
		},
	}

	return nil
}

func envString(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func envRequired(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("%s is empty: set it in .env or environment variables", key)
	}

	return value, nil
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsedValue
}
