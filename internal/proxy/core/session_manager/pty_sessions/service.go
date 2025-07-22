package pty_sessions

import (
	"den-den-mushi-Go/pkg/dto/pty_sessions"
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
	log.Info("Initializing Pty Session Service...")
	return &Service{
		repo: r,
		log:  log,
	}
}

func (s *Service) FindById(id string) (*pty_sessions.Record, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	return s.repo.FindById(id)
}

func (s *Service) Save(session *pty_sessions.Record) error {
	if session == nil {
		return errors.New("session is nil")
	}
	return s.repo.Save(session)
}

func (s *Service) UpdateStateAndEndTime(id string, state types.PtySessionState) error {
	if id == "" {
		return errors.New("id is empty")
	}

	session, err := s.FindById(id)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("session not found")
	}

	if state == types.Closed {
		now := time.Now()
		session.EndTime = &now
	}
	session.State = state
	return s.repo.Save(session)
}
