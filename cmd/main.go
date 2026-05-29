package main

import (
	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	"github.com/Nonameipal/AnalogYouTube/internal/controller"
	"github.com/Nonameipal/AnalogYouTube/internal/db"
	appLogger "github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/Nonameipal/AnalogYouTube/internal/repository"
	"github.com/Nonameipal/AnalogYouTube/internal/service"
)

func main() {
	logger := appLogger.GetLogger()

	logger.Info().Msg("Starting AnalogYouTube service")

	if err := configs.ReadSettings(); err != nil {
		logger.Error().Err(err).Msg("error reading settings")
		return
	}

	dbConn, err := db.InitConnection()
	if err != nil {
		logger.Error().Err(err).Msg("error during database connection initialization")
		return
	}
	defer func() {
		if err := db.CloseConnection(dbConn); err != nil {
			logger.Error().Err(err).Msg("error during database connection close")
		}
	}()

	repo := repository.NewRepository(dbConn)
	svc := service.NewService(repo)
	ctrl := controller.NewController(svc)

	if err = ctrl.InitRoutes(); err != nil {
		logger.Error().Err(err).Msg("error during http-service initialization")
		return
	}
}
