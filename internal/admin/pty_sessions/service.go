package pty_sessions

import (
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
