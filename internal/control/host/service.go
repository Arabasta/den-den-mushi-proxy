package host

import (
	dto "den-den-mushi-Go/pkg/dto/host"
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

func (s *Service) FindByIp(ip string) (*dto.Record, error) {
	h, err := s.repo.FindByIp(ip)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Service) FindTypeByIp(ip string) (types.Proxy, error) {
	//_, _ := s.repo.FindByIp(ip)
	// todo
	return "OS", nil
}

func (s *Service) FindAllByIps(ips []string) ([]*dto.Record, error) {
	return s.repo.FindAllByIps(ips)
}

// todo
func (s *Service) FindFilterTypeByHostType(h types.Proxy) (types.Filter, error) {
	return types.Blacklist, nil
}
