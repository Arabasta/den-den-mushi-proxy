package proxy_lb

import (
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type Service struct {
	repository Repository
	log        *zap.Logger
}

func NewService(repository Repository, log *zap.Logger) *Service {
	return &Service{
		repository: repository,
		log:        log,
	}
}

func (s *Service) GetLBEndpointByProxyType(t types.Proxy) (string, error) {
	pLb, err := s.repository.FindByProxyType(t)
	if err != nil {
		return "", err
	}
	if pLb == nil {
		return "", nil
	}

	return pLb.LoadBalancerEndpoint, nil
}
