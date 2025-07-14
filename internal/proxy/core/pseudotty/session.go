package pseudotty

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/ds"
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

type Session struct {
	id          string
	pty         *os.File
	startClaims *token.Claims // claims from the creator of the session, this must not be modified
	startTime   string
	endTime     string

	purpose Purpose

	// for general logging
	log *zap.Logger

	// for logging all session events
	sessionLogger session_logging.SessionLogger

	filter filter.CommandFilter // only for health check
	line   *filter.LineEditor   // only for health check, tracks pty's current line

	activePrimary       *client.Connection
	activeObservers     map[*client.Connection]struct{}
	livetimeConnections map[*client.Connection]struct{}

	connRegisterCh   chan *client.Connection
	connDeregisterCh chan *client.Connection

	ptyOutput *ds.CircularArray[protocol.Packet]

	ctx    context.Context
	cancel context.CancelFunc

	mu      sync.RWMutex
	Closed  bool // todo: change to state
	onClose func(string)
	once    sync.Once

	cfg *config.Config
}

func New(id string, pty *os.File, log *zap.Logger, cfg *config.Config) (*Session, error) {
	s := &Session{
		id:        id,
		pty:       pty,
		startTime: time.Now().Format(time.RFC3339),
		log:       log.With(zap.String("ptySession", id)),
		cfg:       cfg,

		line: new(filter.LineEditor),

		activeObservers:     make(map[*client.Connection]struct{}),
		livetimeConnections: make(map[*client.Connection]struct{}),

		connRegisterCh:   make(chan *client.Connection),
		connDeregisterCh: make(chan *client.Connection),

		ptyOutput: ds.NewCircularArray[protocol.Packet](500), // todo: make configurable capa and maybe track line or something
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	if err := s.initSessionLogger(); err != nil {
		s.log.Error("Failed to create session log", zap.Error(err))
		return s, err
	}

	return s, nil
}

func (s *Session) Setup(claims *token.Claims) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.log.Debug("Setting up session", zap.String("id", s.id))
	s.startClaims = claims

	err := setPurpose(s, s.startClaims.Connection.Purpose)
	if err != nil {
		return err
	}

	if s.startClaims.Connection.Purpose == types.Healthcheck {
		s.log.Info("Setting healthcheck filter")
		s.filter = filter.GetFilter(s.startClaims.Connection.FilterType)
		if s.filter == nil {
			err = errors.New("invalid filter type")
			s.log.Error("Failed to register initial connection", zap.Error(err))
			return err
		}
	}

	s.log.Info("Initializing conn loop and pty reader")
	go s.connLoop()
	go s.readPtyLoop()

	return nil
}

func (s *Session) SetOnClose(f func(string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onClose = f
}
