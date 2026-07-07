package main

import (
	_ "github.com/Nonameipal/AnalogYouTube/docs"
	"github.com/Nonameipal/AnalogYouTube/internal/app"
	"github.com/Nonameipal/AnalogYouTube/internal/logger"
)

// @title AnalogYouTube API
// @version 1.0
// @description API для AnalogYouTube: пользователи, видео, чаты, донаты и остальное.
// @description Для закрытых запросов нажмите Authorize и вставьте: Bearer <token>.
// @contact.name AnalogYouTube Backend
// @contact.email analogyoutube@example.com
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	log := logger.GetLogger()

	if err := app.Run(); err != nil {
		log.Error().Err(err).Msg("application stopped")
	}
}
