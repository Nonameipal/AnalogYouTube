package service

import (
	"github.com/Nonameipal/AnalogYouTube/internal/contracts"
	"github.com/Nonameipal/AnalogYouTube/internal/repository"
)

type Service struct {
	repository contracts.RepositoryI
}

func NewService(repository *repository.Repository) *Service {
	return &Service{repository: repository}
}
