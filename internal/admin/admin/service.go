package pty_sessions

import "go.uber.org/zap"

type Service struct {
	log *zap.Logger
}

func NewService(log *zap.Logger) *Service {
	log.Info("Initializing pty_sessions Service...")
	return &Service{
		log: log,
	}
}
