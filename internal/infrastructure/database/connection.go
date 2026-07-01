package database

import (
	"fmt"

	"context"
	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	appLogger "github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitConnection() (*pgxpool.Pool, error) {
	connectionConfigs := configs.AppSettings.PostgresParams
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		connectionConfigs.User,
		connectionConfigs.Password,
		connectionConfigs.Host,
		connectionConfigs.Port,
		connectionConfigs.Database)

	appLogger.GetLogger().Debug().
		Str("host", connectionConfigs.Host).
		Str("port", connectionConfigs.Port).
		Str("database", connectionConfigs.Database).
		Str("user", connectionConfigs.User).
		Msg("connecting to postgres")

	pgxpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	appLogger.GetLogger().Info().Msg("postgres connection established")
	return pgxpool, nil
}

func CloseConnection(db *pgxpool.Pool) error {
	if db == nil {
		return nil
	}

	db.Close()
	return nil
}
