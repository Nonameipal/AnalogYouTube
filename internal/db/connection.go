package db

import (
	"fmt"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	appLogger "github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitConnection() (*sqlx.DB, error) {
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

	dbConn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	appLogger.GetLogger().Info().Msg("postgres connection established")
	return dbConn, nil
}

func CloseConnection(db *sqlx.DB) error {
	if db == nil {
		return nil
	}

	return db.Close()
}
