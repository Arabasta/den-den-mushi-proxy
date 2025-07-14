package host

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

func (s *Service) FindByIp(ip string) (*Entity, error) {
	h, err := s.Repository.FindByIp(ip)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Service) FindTypeByIp(ip string) (types.Proxy, error) {
	h, err := s.Repository.FindByIp(ip)
	if err != nil {
		return "", err
	}
	if h == nil {
		return "", nil
	}

	// todo
	return "OS", nil
}

// todo
func (s *Service) FindFilterTypeByHostType(h types.Proxy) (types.Filter, error) {
	return types.Blacklist, nil
}
