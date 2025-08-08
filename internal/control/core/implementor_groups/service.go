package implementor_groups

import (
	dto "den-den-mushi-Go/pkg/dto/implementor_groups"
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

// FindAllByUserId retrieves all implementor groups the user is a member of
func (s *Service) FindAllByUserId(userId string) ([]*dto.Record, error) {
	r, err := s.repo.FindAllByUserId(userId)
	if err != nil {
		s.log.Error("Failed to find groups by user ID", zap.String("userId", userId), zap.Error(err))
		return nil, err
	}
	return r, nil
}
