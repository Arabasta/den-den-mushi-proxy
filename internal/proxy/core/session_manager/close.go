package session_manager

import (
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

func (m *Service) ptyCloseCallback(sessionId string) {
	m.log.Info("Handling pty session close callback", zap.String("sessionId", sessionId))

	m.DeletePtySession(sessionId)
	err := m.ptySessionsSvc.UpdateStateAndEndTime(sessionId, types.Closed)
	if err != nil {
		m.log.Error("Failed to update pty session on close", zap.Error(err), zap.String("sessionId", sessionId))
		return
	}

}

func (m *Service) connCloseCallback(connId string) {
	m.log.Info("Handling connection close callback", zap.String("connectionId", connId))

	err := m.connSvc.UpdateStatusAndLeaveTime(connId, types.ConnectionStatusClosed)
	if err != nil {
		m.log.Error("Failed to update connection on close", zap.Error(err))
	}
}
