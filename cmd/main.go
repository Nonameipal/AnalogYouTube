package main

import (
	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	httpdelivery "github.com/Nonameipal/AnalogYouTube/internal/delivery/http"
	"github.com/Nonameipal/AnalogYouTube/internal/infrastructure/database"
	appLogger "github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/Nonameipal/AnalogYouTube/internal/repository/postgres"
	"github.com/Nonameipal/AnalogYouTube/internal/usecase"
)

func main() {
	logger := appLogger.GetLogger()

	logger.Info().Msg("Starting AnalogYouTube service")

	if err := configs.ReadSettings(); err != nil {
		logger.Error().Err(err).Msg("error reading settings")
		return
	}

	pgxpool, err := database.InitConnection()
	if err != nil {
		logger.Error().Err(err).Msg("error during database connection initialization")
		return
	}
	defer func() {
		if err := database.CloseConnection(pgxpool); err != nil {
			logger.Error().Err(err).Msg("error during database connection close")
		}
	}()

	repo := postgres.NewRepository(pgxpool)
	svc := usecase.NewService(repo)
	ctrl := httpdelivery.NewController(svc)

	if err = ctrl.InitRoutes(); err != nil {
		logger.Error().Err(err).Msg("error during http-service initialization")
		return
	}
}
