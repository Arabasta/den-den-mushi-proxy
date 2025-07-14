package pseudotty

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/ds"
	"den-den-mushi-Go/pkg/dto"
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
	EndTime     string

	purpose Purpose

	// for general logging
	log *zap.Logger

	// for logging all session events
	sessionLogger session_logging.SessionLogger

	filter filter.CommandFilter // only for health check
	line   *filter.LineEditor   // only for health check, tracks pty's current line

	primary   *client.Connection
	observers map[*client.Connection]struct{}

	connRegisterCh   chan *client.Connection
	connDeregisterCh chan *client.Connection

	ptyOutput *ds.CircularArray[protocol.Packet]

	ctx    context.Context
	cancel context.CancelFunc

	mu      sync.RWMutex
	closed  bool // todo: change to state
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

		observers: make(map[*client.Connection]struct{}),

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

type SessionInfo struct {
	SessionID   string          `json:"session_id"`
	StartClaims *token.Claims   `json:"start_claims"`
	Primary     *token.Claims   `json:"current_implementor"`
	Observers   []*token.Claims `json:"observers"`
}

func (s *Session) GetDetails() SessionInfo {
	s.mu.Lock()
	defer s.mu.Unlock()

	var primary = (*token.Claims)(nil)

	if s.primary != nil {
		primary = s.primary.Claims
	}

	observers := make([]*token.Claims, 0, len(s.observers))
	for o := range s.observers {
		observers = append(observers, o.Claims)
	}

	return SessionInfo{ // todo: remove redunant fields, use another DTO
		SessionID:   s.id,
		StartClaims: s.startClaims,
		Primary:     primary,
		Observers:   observers,
	}
}

type Participants struct {
	Primary   *token.Claims  `json:"primary"`
	Observers []token.Claims `json:"observers"`
}

type ParticipantInfo struct {
	UserSessionID string          `json:"user_session_id"`
	UserID        string          `json:"user_id"`
	StartRole     types.StartRole `json:"start_role,required"`
}

func participantInfoFromClaims(claims *token.Claims) ParticipantInfo {
	return ParticipantInfo{
		UserSessionID: claims.Connection.UserSession.Id,
		UserID:        claims.Subject,
		StartRole:     claims.Connection.UserSession.StartRole,
		// todo: add join time
	}
}

type SessionInfo2 struct {
	SessionID              string            `json:"session_id"`
	ProxyDetails           ProxyDetails      `json:"proxy_details"`
	StartConnectionDetails dto.Connection    `json:"start_connection_details"`
	StartTime              string            `json:"start_time"`
	EndTime                string            `json:"end_time,omitempty"`
	State                  string            `json:"state,omitempty"`         // todo: use enum
	LastActivity           string            `json:"last_activity,omitempty"` // ISO 8601 format
	Participants           []ParticipantInfo `json:"participants"`
}

type ProxyDetails struct {
	Hostname    string `json:"hostname"`
	IP          string `json:"ip"`
	Type        string `json:"type"`
	Region      string `json:"region"`
	Environment string `json:"environment"`
}
