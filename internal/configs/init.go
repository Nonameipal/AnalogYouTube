package configs

import (
	"encoding/json"
	"fmt"
	"os"

	appLogger "github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/joho/godotenv"
)

var AppSettings Configs

func ReadSettings() error {
	logger := appLogger.GetLogger()
	logger.Info().Msg("starting reading settings file")

	_ = godotenv.Load(".env")

	configFile, err := os.Open("internal/configs/configs.json")
	if err != nil {
		return fmt.Errorf("couldn't open config file: %w", err)
	}

	defer func(configFile *os.File) {
		if closeErr := configFile.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("couldn't close config file")
		}
	}(configFile)

	logger.Info().Msg("starting decoding settings file")
	if err = json.NewDecoder(configFile).Decode(&AppSettings); err != nil {
		return fmt.Errorf("couldn't decode settings json file: %w", err)
	}

	// Environment variables have priority over configs.json.
	if value := os.Getenv("JWT_SECRET"); value != "" {
		AppSettings.AuthParams.JwtSecret = value
	}
	if value := os.Getenv("POSTGRES_HOST"); value != "" {
		AppSettings.PostgresParams.Host = value
	}
	if value := os.Getenv("POSTGRES_PORT"); value != "" {
		AppSettings.PostgresParams.Port = value
	}
	if value := os.Getenv("POSTGRES_USER"); value != "" {
		AppSettings.PostgresParams.User = value
	}
	if value := os.Getenv("POSTGRES_PASSWORD"); value != "" {
		AppSettings.PostgresParams.Password = value
	}
	if value := os.Getenv("POSTGRES_DATABASE"); value != "" {
		AppSettings.PostgresParams.Database = value
	}

	if AppSettings.AuthParams.JwtSecret == "" {
		return fmt.Errorf("jwt_secret is empty: set auth_params.jwt_secret in configs.json or JWT_SECRET in .env")
	}

	return nil
}
