package pty_sessions

import (
	dto "den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/types"
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

func (s *Service) FindById(id string) (*dto.Record, error) {
	return s.repo.FindById(id)
}

func (s *Service) FindByStartConnChangeRequestIdsAndState(changeIDs []string,
	state types.PtySessionState) ([]*dto.Record, error) {
	return s.repo.FindByStartConnChangeRequestIdsAndState(changeIDs, state)
}

func (s *Service) FindAllByChangeRequestIDAndServerIPs(changeRequestID string, serverIPs []string) ([]*dto.Record, error) {
	return s.repo.FindAllByChangeRequestIDAndServerIPs(changeRequestID, serverIPs)
}

func (s *Service) FindAllByStartConnServerIpsAndState(hostips []string, state *types.PtySessionState) ([]*dto.Record, error) {
	return s.repo.FindAllByStartConnServerIpsAndState(hostips, state)
}
