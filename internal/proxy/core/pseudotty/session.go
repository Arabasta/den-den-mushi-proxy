package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/ds"
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
	"io"
	"os"
	"sync"
	"time"
)

type Session struct {
	id          string
	Pty         *os.File
	startClaims *token.Claims // claims from the creator of the session, this must not be modified
	startTime   string
	endTime     string

	purpose Purpose

	log       *zap.Logger
	logWriter io.WriteCloser

	filter filter.CommandFilter // only for health check
	line   *filter.LineEditor   // only for health check, tracks pty's current line

	primary   *client.Connection
	observers map[*client.Connection]struct{}

	connRegisterCh   chan *client.Connection
	connDeregisterCh chan *client.Connection

	ptyOutput *ds.CircularArray[protocol.Packet]

	mu      sync.RWMutex
	closed  bool // todo: change to state and atomic
	onClose func(string)

	cfg *config.Config
}

func New(id string, pty *os.File, log *zap.Logger, cfg *config.Config) (*Session, error) {
	s := &Session{
		id:        id,
		Pty:       pty,
		startTime: time.Now().Format(time.RFC3339),
		log:       log.With(zap.String("ptySession", id)),
		cfg:       cfg,

		line: new(filter.LineEditor),

		observers: make(map[*client.Connection]struct{}),

		connRegisterCh:   make(chan *client.Connection),
		connDeregisterCh: make(chan *client.Connection),

		ptyOutput: ds.NewCircularArray[protocol.Packet](500), // todo: make configurable capa and maybe track line or something
	}

	if err := s.initLogWriter(); err != nil {
		s.log.Error("Failed to create session log", zap.Error(err))
		return s, err
	}

	s.log.Info("Initializing event loop and pty reader")

	go s.connLoop()
	go s.readPty()
	return s, nil
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
