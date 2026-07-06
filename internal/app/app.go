package app

import (
	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http"
	"github.com/Nonameipal/AnalogYouTube/internal/infrastructure/database"
	"github.com/Nonameipal/AnalogYouTube/internal/infrastructure/ffmpeg"
	"github.com/Nonameipal/AnalogYouTube/internal/infrastructure/storage"
	"github.com/Nonameipal/AnalogYouTube/internal/logger"
	"github.com/Nonameipal/AnalogYouTube/internal/repository/postgres"
	"github.com/Nonameipal/AnalogYouTube/internal/usecase"
)

func Run() error {
	log := logger.GetLogger()

	log.Info().Msg("starting AnalogYouTube service")

	if err := configs.Load(); err != nil {
		return err
	}

	db, err := database.InitConnection()
	if err != nil {
		return err
	}
	defer func() {
		if err := database.CloseConnection(db); err != nil {
			log.Error().Err(err).Msg("error during database connection close")
		}
	}()

	repository := postgres.NewRepository(db)
	mediaProcessor := ffmpeg.NewFFmpegSettings("ffmpeg")
	usecaseLayer := usecase.NewUsecase(repository, mediaProcessor)
	fileStorage := storage.NewVideoStorage("uploads", "uploads")
	httpHandler := httpdelivery.NewHandler(usecaseLayer, fileStorage)

	return httpHandler.InitRoutes()
}
