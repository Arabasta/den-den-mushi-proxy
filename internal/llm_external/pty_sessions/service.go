package pty_sessions

import (
	oapi "den-den-mushi-Go/openapi/llm_external"
	"go.uber.org/zap"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(r Repository, log *zap.Logger) *Service {
	log.Info("Initializing pty_sessions Service...")
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) FindAllByChangeRequestID(changeRequestID string) ([]oapi.GetPtySessionResponse, error) {
	records, err := s.repo.FindAllByChangeRequestID(changeRequestID)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []oapi.GetPtySessionResponse{}, nil
	}

	responses := make([]oapi.GetPtySessionResponse, len(records))
	for i, r := range records {
		responses[i] = oapi.GetPtySessionResponse{
			PtySessionId:             &r.ID,
			TicketId:                 &r.StartConnChangeRequestID,
			SessionCreatedBy:         &r.CreatedBy,
			SessionConnectedServerIp: &r.StartConnServerIP,
		}
	}

	return responses, nil
}
