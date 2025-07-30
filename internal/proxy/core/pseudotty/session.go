package pseudotty

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/integrations/puppet"
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
	Id  string
	pty *os.File

	startClaims *token.Claims // claims from the creator of the session, this must not be modified
	crEndTime   *time.Time

	startTime time.Time
	endTime   time.Time

	purpose Purpose

	// for general logging
	log *zap.Logger

	// for logging all session events
	sessionLogger           session_logging.SessionLogger
	sessionLoggerForAIThing session_logging.SessionLogger
	lastInput               []byte

	filter filter.CommandFilter // only for health check
	line   *filter.LineEditor   // only for health check, tracks pty's current line

	activePrimary       *client.Connection
	activeObservers     map[*client.Connection]struct{}
	lifetimeConnections map[*client.Connection]struct{}

	connRegisterCh   chan *client.Connection
	connDeregisterCh chan *client.Connection

	ptyOutput                    *ds.CircularArray[protocol.Packet]
	isPtyOutputLocked            bool // lock to prevent output during Sudo just incase of password leak
	tmpMuForPtyThingTillRefactor sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	mu sync.RWMutex

	State   types.PtySessionState
	onClose func(string)
	once    sync.Once

	cfg *config.Config

	puppetClient *puppet.Client
}

func New(id string, pty *os.File, now time.Time, onClose func(string), log *zap.Logger, cfg *config.Config,
	puppetClient *puppet.Client) (*Session, error) {
	s := &Session{
		Id:        id,
		pty:       pty,
		startTime: now,
		log:       log.With(zap.String("ptySession", id)),
		cfg:       cfg,

		line: new(filter.LineEditor),

		activeObservers:     make(map[*client.Connection]struct{}),
		lifetimeConnections: make(map[*client.Connection]struct{}),

		connRegisterCh:   make(chan *client.Connection),
		connDeregisterCh: make(chan *client.Connection),

		State:   types.Created,
		onClose: onClose,

		ptyOutput:         ds.NewCircularArray[protocol.Packet](500), // todo: make configurable capa and maybe track line or something
		isPtyOutputLocked: false,
		puppetClient:      puppetClient,
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	if err := s.initSessionLogger(); err != nil {
		s.log.Error("Failed to create session log", zap.Error(err))
		return s, err
	}

	if err := s.initSessionLoggerForAIThing(); err != nil {
		s.log.Error("Failed to create session logger for ai thing", zap.Error(err))
		return s, err
	}

	return s, nil
}

func (s *Session) Setup(claims *token.Claims) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.log.Debug("Setting up session", zap.String("id", s.Id))
	s.startClaims = claims

	if s.startClaims.Connection.Purpose == types.Change {
		s.crEndTime = &s.startClaims.Connection.ChangeRequest.EndTime
	}

	err := setPurpose(s, s.startClaims.Connection.Purpose)
	if err != nil {
		s.log.Error("Failed to set up session", zap.String("id", s.Id), zap.Error(err))
		return err
	}

	s.log.Debug("Session purpose", zap.String("purpose", string(s.startClaims.Connection.Purpose)))
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

	if s.crEndTime != nil {
		go s.monitorCrEndTime()
	}

	go s.connLoop()
	go s.readPtyLoop()

	return nil
}
