package proxy_hosts

import (
	"context"
	"den-den-mushi-Go/pkg/dto/proxy_host"

	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(r Repository, log *zap.Logger) *Service {
	log.Info("Initializing Proxy Hosts Service")
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) FindAll(ctx context.Context) ([]*proxy_host.Record2, error) {
	return s.repo.FindAll(ctx)
}
