package iexpress

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/iexpress"
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

func (s *Service) FindByTicketNumber(num string) (*dto.Record, error) {
	s.log.Debug("Finding iexpress request by ticket number", zap.String("ticketNumber", num))
	return s.repo.FindByTicketNumber(num)
}

func (s *Service) FindApprovedByFilter(f filters.ListIexpress) ([]*dto.Record, error) {
	s.log.Debug("Finding iexpress requests by filter", zap.Any("filter", f))
	return s.repo.FindApprovedByFilter(f)
}

func (s *Service) CountApprovedByFilter(f filters.ListIexpress) (int, error) {
	return s.repo.CountApprovedByFilter(f)
}
