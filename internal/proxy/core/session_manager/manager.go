package session_manager

import (
	"den-den-mushi-Go/internal/proxy/core/pseudotty"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	mu          sync.RWMutex
	ptySessions map[string]*pseudotty.Session // todo use service and repo layer
	log         *zap.Logger
}

func New(log *zap.Logger) *Service {
	return &Service{
		ptySessions: make(map[string]*pseudotty.Session),
		log:         log,
	}
}

func (m *Service) CreatePtySession(pty *os.File, log *zap.Logger) *pseudotty.Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.NewString() + strconv.FormatInt(time.Now().Unix(), 10)
	s := pseudotty.New(id, pty, log)
	s.SetOnClose(func(sessionID string) {
		m.DeletePtySession(sessionID)
	})

	m.AddPtySession(id, s)
	return s
}

func (m *Service) AddPtySession(id string, sess *pseudotty.Session) {
	m.log.Info("Adding pty session to map", zap.String("id", id))
	if _, exists := m.ptySessions[id]; exists {
		sess.Log.Error("Pty session already exists", zap.String("id", id))
	}
	m.ptySessions[id] = sess
}

func (m *Service) GetPtySession(id string) (*pseudotty.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.log.Debug("Retrieving pty session", zap.String("id", id))
	s, ok := m.ptySessions[id]
	return s, ok
}

func (m *Service) DeletePtySession(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Info("Deleting pty session", zap.String("id", id))
	if _, ok := m.ptySessions[id]; ok {
		delete(m.ptySessions, id)
	}
}

// AttachWebsocket if pty session already exists
func (m *Service) AttachWebsocket(ws *websocket.Conn, claims *token.Claims) error {
	session, exists := m.GetPtySession(claims.Connection.PtySession.Id)
	if !exists {
		return errors.New("pty session not found")
	}

	return session.RegisterConn(ws, claims)
}

// tmp  for demo
func (m *Service) GetPtySessions() []pseudotty.SessionInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.log.Info("Retrieving pty session info")
	ptySessions := make([]pseudotty.SessionInfo, 0, len(m.ptySessions))

	for s := range m.ptySessions {
		session := m.ptySessions[s]
		ptySessions = append(ptySessions, session.GetDetails())
	}

	m.log.Debug("Retrieved pty session info", zap.Any("ptySessions", ptySessions))
	return ptySessions
}
