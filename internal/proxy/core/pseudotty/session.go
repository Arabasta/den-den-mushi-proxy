package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/protocol"
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

	outboundCh     chan protocol.Packet
	ptyLastPackets []protocol.Packet

	mu      sync.Mutex
	closed  bool // todo: change to state and atomic
	onClose func(string)
}

func New(id string, pty *os.File, log *zap.Logger) (*Session, error) {
	s := &Session{
		id:        id,
		Pty:       pty,
		startTime: time.Now().Format(time.RFC3339),
		log:       log.With(zap.String("ptySession", id)),

		line: new(filter.LineEditor),

		observers: make(map[*client.Connection]struct{}),

		outboundCh: make(chan protocol.Packet, 100), // todo: make configurable

		connRegisterCh:   make(chan *client.Connection),
		connDeregisterCh: make(chan *client.Connection),

		ptyLastPackets: types.NewCircularArray[protocol.Packet](100), // todo: make configurable capa and maybe track line or something
	}

	if err := s.initLogWriter(); err != nil {
		s.log.Error("Failed to create session log", zap.Error(err))
		return s, err
	}

	s.log.Info("Initializing event loop and pty reader")

	go s.eventLoop()
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
