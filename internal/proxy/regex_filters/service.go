package regex_filters

import (
	"den-den-mushi-Go/internal/proxy/config"
	dto "den-den-mushi-Go/pkg/dto/regex_filters"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
	Cfg  *config.Config
}

func NewService(r Repository, log *zap.Logger, cfg *config.Config) *Service {
	log.Info("Initializing Regex filters Service...")
	return &Service{
		repo: r,
		log:  log,
		Cfg:  cfg,
	}
}

func (s *Service) FindAllEnabledByFilterType(filterType types.Filter) (*[]dto.Record, error) {
	return s.repo.FindAllEnabledByFilterType(filterType)
}
