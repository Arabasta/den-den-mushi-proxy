package connections

import (
	"den-den-mushi-Go/pkg/dto/connections"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	repo Repository
	log  *zap.Logger
}

func NewService(r Repository, log *zap.Logger) *Service {
	log.Info("Initializing Pty Connections Service...")
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) FindById(id string) (*connections.Record, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	return s.repo.FindById(id)
}

func (s *Service) Save(r *connections.Record) error {
	if err := s.repo.Save(r); err != nil {
		s.log.Error("Failed to save connection record", zap.String("id", r.ID), zap.Error(err))
		return err
	}
	s.log.Info("Connection record saved successfully", zap.String("id", r.ID))
	return nil
}

func (s *Service) UpdateStatusAndLeaveTime(id string, status types.ConnectionStatus) error {
	if id == "" || status == "" {
		return errors.New("id or status is empty")
	}

	c, err := s.FindById(id)
	if err != nil {
		return err
	}
	if c == nil {
		return errors.New("connection not found")
	}

	if status == types.ConnectionStatusClosed {
		now := time.Now()
		c.LeaveTime = &now
	}

	c.Status = status
	return s.repo.Save(c)
}
