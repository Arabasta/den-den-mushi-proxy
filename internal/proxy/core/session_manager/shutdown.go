package session_manager

import (
	"context"
	"den-den-mushi-Go/internal/proxy/core/pseudotty"
	"go.uber.org/zap"
)

func (m *Service) CleanupActiveSessionsAndConnections() error {
	m.log.Info("Cleaning up old active sessions with connections", zap.String("proxyHost", m.hostname))
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.ptySessionsSvc.CleanupActiveSessionsAndConnections(m.hostname)
}

func (m *Service) Shutdown(ctx context.Context) {
	m.log.Info("Shutting down all pty sessions and associated connections...")

	m.mu.Lock()
	sessions := make([]*pseudotty.Session, 0, len(m.ptySessions))
	for _, s := range m.ptySessions {
		sessions = append(sessions, s)
	}
	m.mu.Unlock()

	for _, s := range sessions {
		select {
		case <-ctx.Done():
			m.log.Warn("Shutdown timed out")
			return
		default:
			m.connCloseCallback(s.ActivePrimary.Id)
			for conn := range s.ActiveObservers {
				m.connCloseCallback(conn.Id)
			}
			s.ForceEndSession("Proxy shutdown initiated") // will trigger onClose once
		}
	}
	m.log.Info("All pty sessions and connections have been shut down successfully.")
}
