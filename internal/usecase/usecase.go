package usecase

import (
	"github.com/Nonameipal/AnalogYouTube/internal/infrastructure/ffmpeg"
	"github.com/Nonameipal/AnalogYouTube/internal/usecase/ports"
)

type Usecase struct {
	repository     ports.RepositoryI
	ffmpegSettings *ffmpeg.FFmpegSettings
}

func NewUsecase(repository ports.RepositoryI, ffmpegSettings *ffmpeg.FFmpegSettings) *Usecase {
	return &Usecase{
		repository:     repository,
		ffmpegSettings: ffmpegSettings,
	}
}
