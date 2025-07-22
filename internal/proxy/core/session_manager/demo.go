package session_manager

import (
	"den-den-mushi-Go/internal/proxy/core/pseudotty"
	"go.uber.org/zap"
)

// todo: tmp  for demo, remove this
func (m *Service) GetPtySessionsTmpForAdminServer() []pseudotty.SessionInfo {
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
