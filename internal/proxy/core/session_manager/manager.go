package session_manager

import (
	"den-den-mushi-Go/internal/proxy/core/pseudotty"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"os"
	"sync"
)

type SessionManager struct {
	mu          sync.RWMutex
	ptySessions map[string]*pseudotty.Session // todo use service and repo layer
	log         *zap.Logger
}

func New(log *zap.Logger) *SessionManager {
	return &SessionManager{
		ptySessions: make(map[string]*pseudotty.Session),
		log:         log,
	}
}

func (m *SessionManager) CreatePtySession(pty *os.File, log *zap.Logger) *pseudotty.Session {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.NewString()
	s := pseudotty.New(id, pty, log)

	m.AddPtySession(id, s)
	return s
}

func (m *SessionManager) AddPtySession(id string, sess *pseudotty.Session) {
	m.log.Info("Adding pty session to map", zap.String("id", id))
	if _, exists := m.ptySessions[id]; exists {
		sess.Log.Error("Pty session already exists", zap.String("id", id))
	}
	m.ptySessions[id] = sess
}

func (m *SessionManager) GetPtySession(id string) (*pseudotty.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.log.Debug("Retrieving pty session", zap.String("id", id))
	s, ok := m.ptySessions[id]
	return s, ok
}

func (m *SessionManager) DeletePtySession(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Info("Deleting pty session", zap.String("id", id))
	if sess, ok := m.ptySessions[id]; ok {
		sess.EndSession() // close the pty and its ws connections
		delete(m.ptySessions, id)
	}
}

// if pty session already exists, attach the websocket to the existing session
func (m *SessionManager) AttachWebsocket(ws *websocket.Conn, claims *token.Claims) error {
	session, exists := m.GetPtySession(claims.Connection.PtySession.Id)
	if !exists {
		return errors.New("pty session not found")
	}

	return session.RegisterConn(ws, claims)
}
