package session_manager

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	mu          sync.RWMutex
	ptySessions map[string]*pseudotty.Session // todo use service and repo layer, anyway need inmem map for retrieval
	log         *zap.Logger
	cfg         *config.Config
}

func New(log *zap.Logger, cfg *config.Config) *Service {
	return &Service{
		ptySessions: make(map[string]*pseudotty.Session),
		log:         log,
		cfg:         cfg,
	}
}

func (m *Service) CreatePtySession(pty *os.File, log *zap.Logger) (*pseudotty.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.log.Info("Creating pty session")

	id := uuid.NewString() + strconv.FormatInt(time.Now().Unix(), 10)

	s, err := pseudotty.New(id, pty, log, m.cfg)
	if err != nil {
		m.log.Error("Failed to create pty session", zap.Error(err), zap.String("id", id))
		// todo: close it
		return nil, err
	}

	s.SetOnClose(func(sessionID string) {
		m.DeletePtySession(sessionID)
	})

	m.log.Info("Created pty session", zap.String("id", id))

	if err = m.AddPtySession(id, s); err != nil {
		m.log.Error("Failed to add pty session to map", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	return s, nil
}

func (m *Service) AddPtySession(id string, s *pseudotty.Session) error {
	m.log.Info("Adding pty session to map", zap.String("id", id))
	if _, exists := m.ptySessions[id]; exists {
		m.log.Error("Pty session already exists", zap.String("id", id))
		return errors.New("pty session already exists with id: " + id)
	}

	m.ptySessions[id] = s
	return nil
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
		pty, ok := m.GetPtySession(id)
		if !ok {
			m.log.Warn("Pty session not found for deletion", zap.String("id", id))
			// todo: handle this case
			return
		}

		pty.EndSession()

		// delete from map
		delete(m.ptySessions, id)
	}
}

// AttachConnToExisting if pty session already exists
func (m *Service) AttachConnToExisting(conn *client.Connection) error {
	session, exists := m.GetPtySession(conn.Claims.Connection.PtySession.Id)
	if !exists {
		return errors.New("pty session not found")
	}

	return session.RegisterConn(conn)
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
