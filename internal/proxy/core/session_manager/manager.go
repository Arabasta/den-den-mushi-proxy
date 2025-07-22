package session_manager

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/pseudotty"
	"den-den-mushi-Go/internal/proxy/core/session_manager/connections"
	"den-den-mushi-Go/internal/proxy/core/session_manager/pty_sessions"
	"go.uber.org/zap"
	"sync"
)

type Service struct {
	mu             sync.RWMutex
	ptySessions    map[string]*pseudotty.Session
	ptySessionsSvc *pty_sessions.Service
	connSvc        *connections.Service
	log            *zap.Logger
	cfg            *config.Config
}

func New(ptySessionsSvc *pty_sessions.Service, connSvc *connections.Service, log *zap.Logger, cfg *config.Config) *Service {
	return &Service{
		ptySessions:    make(map[string]*pseudotty.Session),
		ptySessionsSvc: ptySessionsSvc,
		connSvc:        connSvc,
		log:            log,
		cfg:            cfg,
	}
}
