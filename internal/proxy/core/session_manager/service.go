package session_manager

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty"
	"den-den-mushi-Go/pkg/dto/connections"
	"den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (m *Service) CreatePtySession(pty *os.File, cmd *exec.Cmd, claims *token.Claims, log *zap.Logger) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.log.Debug("Creating pty session")

	id := uuid.NewString() + strconv.FormatInt(time.Now().Unix(), 10)
	now := time.Now()

	err := m.ptySessionsSvc.Save(&pty_sessions.Record{
		ID:                       id,
		CreatedBy:                claims.Subject,
		State:                    types.Active,
		StartTime:                &now,
		ProxyHostName:            m.hostname,
		StartConnServerIP:        claims.Connection.Server.IP,
		StartConnServerPort:      claims.Connection.Server.Port,
		StartConnServerOSUser:    claims.Connection.Server.OSUser,
		StartConnType:            claims.Connection.Type,
		StartConnPurpose:         claims.Connection.Purpose,
		StartConnUserSessionID:   claims.Connection.UserSession.Id,
		StartConnChangeRequestID: claims.Connection.ChangeRequest.Id,
		StartConnFilterType:      claims.Connection.FilterType,
	})
	if err != nil {
		return "", err
	}

	s, err := pseudotty.New(id, pty, cmd, now, m.ptyCloseCallback, log, m.cfg, m.puppetClient, m.filterSvc)
	if err != nil {
		m.log.Error("Failed to create pty session", zap.Error(err), zap.String("ptySessionId", id))
		s.ForceEndSession("Failed to create pty session: " + err.Error())
		return "", err
	}

	m.log.Info("Created pty session, adding pty session to map", zap.String("id", id))
	if _, exists := m.ptySessions[id]; exists {
		m.log.Error("Pty session already exists", zap.String("ptySessionId", id))
		return "", errors.New("pty session already exists with id: " + id)
	}

	m.ptySessions[id] = s

	err = s.Setup(claims)
	if err != nil {
		m.log.Error("Failed to setup session", zap.Error(err))
		s.ForceEndSession("Failed to setup session: " + err.Error())
		return "", err
	}

	return id, nil
}

func (m *Service) GetPtySession(id string) (*pseudotty.Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.log.Debug("Retrieving pty session", zap.String("ptySessionId", id))
	s, ok := m.ptySessions[id]
	return s, ok
}

func (m *Service) DeletePtySession(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.log.Info("Deleting pty session", zap.String("ptySessionId", id))
	delete(m.ptySessions, id)
}

func (m *Service) AttachConn(c *client.Connection, ptySessionId string) error {
	s, exists := m.GetPtySession(ptySessionId)
	if !exists || s.State == types.Closed {
		return errors.New("pty session not found or closed")
	}

	err := s.AssignRole(c)
	if err != nil {
		m.log.Warn("Failed to assign role to connection", zap.Error(err), zap.String("ptySessionId", ptySessionId))
		return err
	}

	if err := s.RegisterConn(c, m.connCloseCallback); err != nil {
		m.log.Warn("Failed to register connection to pty session", zap.Error(err), zap.String("ptySessionId", ptySessionId))
		return err
	}

	now := time.Now()
	if err := m.connSvc.Save(&connections.Record{
		ID:           c.Claims.Connection.UserSession.Id,
		UserID:       c.Claims.Subject,
		PtySessionID: s.Id,
		StartRole:    c.Claims.Connection.UserSession.StartRole,
		Status:       types.ConnectionStatusActive,
		JoinTime:     &now,
	}); err != nil {
		m.log.Error("Failed to save connection record", zap.Error(err), zap.String("ptySessionId", ptySessionId))
		s.DeregisterConn(c)
		return err
	}
	return nil
}
