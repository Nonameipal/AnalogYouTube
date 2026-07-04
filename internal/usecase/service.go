package usecase

import (
	"github.com/Nonameipal/AnalogYouTube/internal/infrastructure/ffmpeg"
	"github.com/Nonameipal/AnalogYouTube/internal/usecase/ports"
)

type Service struct {
	repository ports.RepositoryI
	ffmpegSettings *ffmpeg.FFmpegSettings
}

func NewService(repository ports.RepositoryI, ffmpegSettings *ffmpeg.FFmpegSettings) *Service {
	return &Service{
		repository: repository,
		ffmpegSettings: ffmpegSettings,
	}
}
