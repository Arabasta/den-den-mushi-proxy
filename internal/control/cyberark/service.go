package cyberark

import (
	dto "den-den-mushi-Go/pkg/dto/cyberark"
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

func (s *Service) FindByObject(o string) (*dto.Record, error) {
	r, err := s.repo.FindByObject(o)
	if err != nil {
		s.log.Error("Failed to find CyberArk record by object", zap.String("object", o), zap.Error(err))
		return nil, err
	}
	if r == nil {
		s.log.Debug("No CyberArk record found for object", zap.String("object", o))
		return nil, nil
	}
	return r, nil
}
