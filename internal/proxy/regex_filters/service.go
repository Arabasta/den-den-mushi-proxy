package regex_filters

import (
	dto "den-den-mushi-Go/pkg/dto/regex_filters"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(r Repository, log *zap.Logger) *Service {
	log.Info("Initializing Regex filters Service...")
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) FindAllEnabledByFilterType(filterType types.Filter) (*[]dto.Record, error) {
	return s.repo.FindAllEnabledByFilterType(filterType)
}
