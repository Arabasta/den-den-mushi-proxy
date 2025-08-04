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

func (s *Service) CleanupActiveSessionsAndConnections(proxyHost string) error {
	s.log.Info("Cleaning up old active sessions with connections", zap.String("proxyHost", proxyHost))

	sessions, err := s.repo.FindActiveByProxyHostWithConnections(proxyHost)
	if err != nil {
		return err
	}
	if len(sessions) == 0 {
		s.log.Info("No active sessions found to clean up", zap.String("proxyHost", proxyHost))
		return nil
	}

	sessionIDs := make([]string, len(sessions))
	var connectionIDs []string
	for i, sess := range sessions {
		sessionIDs[i] = sess.ID
		for _, conn := range sess.Connections {
			connectionIDs = append(connectionIDs, conn.ID)
		}
	}

	s.log.Debug("Found old active sessions and connections to close",
		zap.Int("sessionsCount", len(sessionIDs)),
		zap.Int("connectionsCount", len(connectionIDs)))

	if err := s.repo.CloseSessionsAndConnections(sessionIDs, connectionIDs); err != nil {
		s.log.Error("Failed to close old active sessions and connections", zap.Error(err))
		return err
	}

	s.log.Info("Closed old active sessions and connections",
		zap.Int("sessionsCount", len(sessionIDs)),
		zap.Int("connectionsCount", len(connectionIDs)))
	return nil
}
