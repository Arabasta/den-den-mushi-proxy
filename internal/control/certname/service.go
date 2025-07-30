package certname

import (
	dto "den-den-mushi-Go/pkg/dto/puppet_trusted"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(repo Repository, log *zap.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) FindCertnameByIp(ip string) (*dto.Record, error) {
	s.log.Debug("Looking up certname by IP", zap.String("ip", ip))
	return s.repo.FindCertnameByIp(ip)
}
