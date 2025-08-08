package certname

import (
	"den-den-mushi-Go/internal/control/config"
	dto "den-den-mushi-Go/pkg/dto/puppet_trusted"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
	cfg  *config.Config
}

func NewService(repo Repository, log *zap.Logger, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

func (s *Service) FindCertnameByIp(ip string) (*dto.Record, error) {
	s.log.Debug("Looking up certname by IP", zap.String("ip", ip))
	if !s.cfg.Development.IsUsingInvDb {
		s.log.Debug("Skipping certname lookup in Kei's macbook environment")
		return &dto.Record{Certname: "123123"}, nil
	}
	return s.repo.FindCertnameByIp(ip)
}
