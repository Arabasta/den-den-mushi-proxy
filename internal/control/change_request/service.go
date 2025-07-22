package change_request

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/change_request"
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
	s.log.Debug("Finding change request by ticket number", zap.String("ticketNumber", num))
	return s.repo.FindByTicketNumber(num)
}

func (s *Service) FindChangeRequests(filter filters.ListCR) ([]*dto.Record, error) {
	s.log.Debug("Finding change requests with filter", zap.Any("filter", filter))
	return s.repo.FindChangeRequests(filter)
}
