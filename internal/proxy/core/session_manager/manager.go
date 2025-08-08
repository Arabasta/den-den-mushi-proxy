package session_manager

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/pseudotty"
	"den-den-mushi-Go/internal/proxy/core/session_manager/connections"
	"den-den-mushi-Go/internal/proxy/core/session_manager/pty_sessions"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/integrations/puppet"
	"go.uber.org/zap"
	"os"
	"sync"
)

type Service struct {
	mu             sync.RWMutex
	ptySessions    map[string]*pseudotty.Session
	ptySessionsSvc *pty_sessions.Service
	connSvc        *connections.Service
	log            *zap.Logger
	cfg            *config.Config
	puppetClient   *puppet.Client
	filterSvc      *filter.Service
	hostname       string
}

func New(ptySessionsSvc *pty_sessions.Service, connSvc *connections.Service, log *zap.Logger, cfg *config.Config,
	puppetClient *puppet.Client, filterSvc *filter.Service) *Service {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Failed to get hostname", zap.Error(err))
		os.Exit(1)
	}
	log.Info("Initializing Session Manager Service", zap.String("hostname", hostname))

	return &Service{
		ptySessions:    make(map[string]*pseudotty.Session),
		ptySessionsSvc: ptySessionsSvc,
		connSvc:        connSvc,
		log:            log,
		cfg:            cfg,
		filterSvc:      filterSvc,
		puppetClient:   puppetClient,
		hostname:       hostname,
	}
}
