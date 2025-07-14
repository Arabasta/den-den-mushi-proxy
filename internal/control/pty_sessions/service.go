package pty_sessions

import (
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type Service struct {
	Repository Repository
	log        *zap.Logger
}

func NewService(repository Repository, log *zap.Logger) *Service {
	return &Service{
		Repository: repository,
		log:        log,
	}
}

func (s *Service) FindById(id string) (*Entity, error) {
	return s.Repository.FindById(id)
}

func (s *Service) FindByState(st types.PtySessionState) ([]*Entity, error) {
	return s.Repository.FindAllByState(st)
}

func (s *Service) FindAllByIp(ip string) ([]*Entity, error) {
	return s.Repository.FindAllByIp(ip)
}
