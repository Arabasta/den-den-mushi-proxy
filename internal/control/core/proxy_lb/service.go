package proxy_lb

import (
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

func (s *Service) GetLBEndpointByProxyType(t types.Proxy) (string, error) {
	s.log.Debug("Finding load balancer endpoint by proxy type", zap.String("proxyType", string(t)))
	pLb, err := s.repo.FindByProxyType(t)
	if err != nil {
		return "", err
	}
	if pLb == nil {
		return "", nil
	}

	return pLb.LoadBalancerEndpoint, nil
}
