package pseudotty

import (
	"context"
	"den-den-mushi-Go/internal/proxy/core/client"
	"go.uber.org/zap"
)

// RegisterConn client connection to pty session
func (s *Session) RegisterConn(c *client.Connection) error {
	s.log.Info("Attaching connection to pty session", zap.String("userSessionId",
		c.Claims.Connection.UserSession.Id))

	c.Ctx, c.Cancel = context.WithCancel(s.ctx)

	c.Close = func() {
		s.connDeregisterCh <- c
	}

	s.connRegisterCh <- c
	return nil
}
