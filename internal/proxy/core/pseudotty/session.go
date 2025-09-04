package pseudotty

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/activity"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/integrations/puppet"
	"den-den-mushi-Go/internal/proxy/line"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/ds"
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"os"
	"os/exec"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Session struct {
	Id  string
	pty *os.File
	cmd *exec.Cmd

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
	line   *line.Editor         // only for health check, tracks pty's current line

	ActivePrimary       *client.Connection
	ActiveObservers     map[*client.Connection]struct{}
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
	filterSvc    *filter.Service

	activityTracker *activity.Tracker
}

func New(id string, pty *os.File, cmd *exec.Cmd, now time.Time, onClose func(string), log *zap.Logger, cfg *config.Config,
	puppetClient *puppet.Client, filterSvc *filter.Service) (*Session, error) {
	s := &Session{
		Id:        id,
		pty:       pty,
		cmd:       cmd,
		startTime: now,
		log:       log.With(zap.String("ptySession", id)),
		cfg:       cfg,

		line: new(line.Editor),

		ActiveObservers:     make(map[*client.Connection]struct{}),
		lifetimeConnections: make(map[*client.Connection]struct{}),

		connRegisterCh:   make(chan *client.Connection),
		connDeregisterCh: make(chan *client.Connection),

		State:   types.Created,
		onClose: onClose,

		ptyOutput:         ds.NewCircularArray[protocol.Packet](500), // todo: make configurable capa and maybe track line or something
		isPtyOutputLocked: false,
		puppetClient:      puppetClient,
		filterSvc:         filterSvc,
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	if err := s.initSessionLogger(); err != nil {
		s.log.Error("Failed to create session log", zap.Error(err))
		s.EndSession()
		return s, err
	}

	if err := s.initSessionLoggerForAIThing(); err != nil {
		s.log.Error("Failed to create session logger for ai thing", zap.Error(err))
		s.EndSession()
		return s, err
	}

	s.logSessionDetails()
	return s, nil
}

func (s *Session) Setup(claims *token.Claims) error {
	s.mu.Lock()
	defer func() {
		if r := recover(); r != nil {
			s.log.Error("panic", zap.Any("panic", r), zap.Stack("stack"))
		}
		s.mu.Unlock()
	}()

	s.log.Debug("Setting up session", zap.String("id", s.Id))
	s.startClaims = claims

	if s.startClaims.Connection.Purpose == types.Change || s.startClaims.Connection.Purpose == types.IExpress {
		if s.crEndTime != nil {
			s.log.Debug("Setting up cr")
			s.crEndTime = &s.startClaims.Connection.ChangeRequest.EndTime

			// if cr ending in 1min or already ended, don't allow session to start
			// else it will crash the conn loop
			if s.crEndTime.Before(time.Now().Add(1 * time.Minute)) {
				s.log.Error("Ticket already expired or ending soon", zap.Time("endTime", *s.crEndTime))
				s.logL("Ticket already expired or ending soon")
				s.EndSession()
				return errors.New("ticket already expired")
			}
		}
	}

	err := setPurpose(s, s.startClaims.Connection.Purpose)
	if err != nil {
		s.log.Error("invalid session purpose", zap.String("id", s.Id), zap.Error(err))
		s.logL("Invalid session purpose")
		s.EndSession()
		return err
	}

	s.log.Debug("Session purpose", zap.String("purpose", string(s.startClaims.Connection.Purpose)))
	s.filter = s.filterSvc.GetFilter(s.startClaims.Connection.FilterType, s.startClaims.Connection.Purpose)
	if s.filter == nil {
		err = errors.New("invalid filter type")
		s.log.Error("Failed to register initial connection", zap.Error(err))
		s.logL("Invalid filter type")
		s.EndSession()
		return err
	}

	s.log.Info("Initializing conn loop and pty reader")

	go s.connLoop()
	go s.readPtyLoop()

	if s.crEndTime != nil {
		s.log.Debug("Starting cr end time monitor", zap.Time("endTime", *s.crEndTime))
		go s.monitorCrEndTime()
	}

	s.initActivityTracker()

	return nil
}
