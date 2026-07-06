package main

import (
	"github.com/Nonameipal/AnalogYouTube/internal/app"
	"github.com/Nonameipal/AnalogYouTube/internal/logger"
)

func main() {
	log := logger.GetLogger()

	if err := app.Run(); err != nil {
		log.Error().Err(err).Msg("application stopped")
	}
}
