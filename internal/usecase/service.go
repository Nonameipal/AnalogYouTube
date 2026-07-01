package usecase

import (
	"github.com/Nonameipal/AnalogYouTube/internal/usecase/ports"
)

type Service struct {
	repository ports.RepositoryI
}

func NewService(repository ports.RepositoryI) *Service {
	return &Service{repository: repository}
}
