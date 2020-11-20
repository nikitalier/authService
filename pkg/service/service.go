package service

import (
	"github.com/nikitalier/authService/pkg/repository"
	"github.com/rs/zerolog"
)

type Service struct {
	repository *repository.Repository
	logger     *zerolog.Logger
}

func New(rep *repository.Repository, logger *zerolog.Logger) *Service {
	return &Service{
		repository: rep,
		logger:     logger,
	}
}
