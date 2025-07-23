package connection

import (
	"den-den-mushi-Go/pkg/dto/connections"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(r Repository, log *zap.Logger) *Service {
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) FindById(id string) (*connections.Record, error) {
	c, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	return c, nil
}

func (s *Service) FindActiveImplementorByPtySessionId(ptySessionId string) (*connections.Record, error) {
	return s.repo.FindActiveImplementorByPtySessionId(ptySessionId)
}
