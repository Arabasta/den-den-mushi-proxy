package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"go.uber.org/zap"
	"io"
	"sync"
	"time"
)

func (s *Session) EndSession() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	defer func() { s.closed = true }()

	s.log.Info("Ending pty session")
	s.closeTheWorld()
	s.endTime = time.Now().Format(time.RFC3339)
	s.logL(getLogFooter(s))
	s.onClose(s.id)
}

func (s *Session) closeTheWorld() func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			s.log.Debug("Closing the world")
			s.deregisterAllWsConnections()
			s.closePty()
			s.closeSessionChannels()
			s.closeLogWriter()
		})
	}
}

func (s *Session) deregisterAllWsConnections() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connDeregisterCh <- s.primary

	for o := range s.observers {
		s.connDeregisterCh <- o
	}

	s.log.Info("Closed all websocket connections")
}

func (s *Session) closeSessionChannels() {
	close(s.connRegisterCh)
	close(s.connDeregisterCh)
	close(s.outboundCh)
	s.log.Info("Closed session channels")
}

func (s *Session) closeWs(c *client.Connection) {
	c.OnceCloseWriteCh.Do(func() {
		close(c.WsWriteCh)

		if c.Sock != nil {
			s.log.Debug("Closing websocket connection", zap.String("userSessionId", c.Claims.Connection.UserSession.Id))
			err := c.Sock.Close()
			if err != nil {
				s.log.Error("Failed to close websocket connection", zap.Error(err))
				return
			}
		}
	})
}

// todo: more error handling eg check if pty is arleady closed
func (s *Session) closePty() {
	err := s.Pty.Close()
	if err != nil {
		if err != io.EOF {
			s.log.Info("PTY session ended normally")
		} else {
			s.log.Warn("Failed to close pty", zap.Error(err))
		}
	}

	s.log.Info("Closed pty")

}

func (s *Session) closeLogWriter() {
	if s.logWriter != nil {
		if err := s.logWriter.Close(); err != nil {
			s.log.Warn("Failed to close log writer", zap.Error(err))
		} else {
			s.log.Info("Closed log writer")
		}
	}
}
