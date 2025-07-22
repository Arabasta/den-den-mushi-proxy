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
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) Save(r *dto.Record) (*dto.Record, error) {
	res, err := s.repo.Save(r)
	if err != nil {
		s.log.Error("Failed to save regex filter", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *Service) FindAllByFilterTypeAndOuGroup(filterType types.Filter, ouGroup string) (*[]dto.Record, error) {
	return s.repo.FindAllByFilterTypeAndOuGroup(filterType, ouGroup)
}

func (s *Service) Update(r *dto.Record) (*dto.Record, error) {
	existing, err := s.repo.FindById(r.Id)
	if err != nil {
		s.log.Error("Failed to find regex filter for update", zap.Uint("id", r.Id), zap.Error(err))
		return nil, err
	}

	s.log.Debug("Updating regex filter", zap.String("id", r.UpdatedBy))
	existing.RegexPattern = r.RegexPattern
	existing.UpdatedBy = r.UpdatedBy
	existing.IsEnabled = r.IsEnabled

	existing, err = s.repo.Save(existing)
	if err != nil {
		s.log.Error("Failed to update regex filter", zap.Uint("id", r.Id), zap.Error(err))
		return nil, err
	}

	return existing, nil
}

func (s *Service) SoftDelete(r *dto.Record) (*dto.Record, error) {
	existing, err := s.repo.FindById(r.Id)
	if err != nil {
		s.log.Error("Failed to find regex filter for update", zap.Uint("id", r.Id), zap.Error(err))
		return nil, err
	}

	existing.DeletedBy = r.DeletedBy

	existing, err = s.repo.Save(existing)
	if err != nil {
		s.log.Error("Failed to delete regex filter", zap.Uint("id", r.Id), zap.Error(err))
		return nil, err
	}

	if err := s.repo.Delete(existing); err != nil {
		s.log.Error("Failed to delete regex filter", zap.Uint("id", r.Id), zap.Error(err))
		return nil, err
	}

	return existing, nil
}
