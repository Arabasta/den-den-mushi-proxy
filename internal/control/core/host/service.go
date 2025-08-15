package host

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/host"
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

func (s *Service) FindAllLinuxOsByIps(ips []string) ([]*dto.Record, error) {
	return s.repo.FindAllLinuxOsByIps(ips)
}

func (s *Service) FindAllByFilter(f filters.HealthcheckPtySession) ([]*dto.Record, error) {
	s.log.Debug("Finding hosts by filter", zap.Any("filter", f))
	return s.repo.FindAllByFilter(f)
}

func (s *Service) CountAllByFilter(f filters.HealthcheckPtySession) (int64, error) {
	return s.repo.CountAllByFilter(f)
}
