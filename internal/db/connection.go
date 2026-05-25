package db

import (
	"fmt"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
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
	fmt.Println("DEBUG Connection String:", connStr)
	dbConn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func CloseConnection(db *sqlx.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}
