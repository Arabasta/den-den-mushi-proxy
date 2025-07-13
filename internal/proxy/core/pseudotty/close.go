package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/logging"
	"go.uber.org/zap"
	"io"
	"time"
)

func (s *Session) EndSession() {
	s.once.Do(func() {
		s.cancel() // exit conn loop
		if s.closed {
			return
		}
		defer func() { s.closed = true }()

		s.log.Info("Ending pty session")
		s.closeTheWorld()
		s.EndTime = time.Now().Format(time.RFC3339)
		s.logL(logging.FormatFooter(s.EndTime))
		if s.onClose != nil {
			s.onClose(s.id)
		}
	})
}

func (s *Session) closeTheWorld() {
	s.log.Debug("Closing the world")
	s.deregisterAllWsConnections()
	s.closePty()
	s.closeLogWriter()
}

func (s *Session) deregisterAllWsConnections() {
	s.mu.RLock()
	// don't use deregisterCh
	primary := s.primary
	observers := make([]*client.Connection, 0, len(s.observers))
	for o := range s.observers {
		observers = append(observers, o)
	}
	s.mu.RUnlock()

	if primary != nil {
		s.removeConn(primary)
	}
	for _, o := range observers {
		s.removeConn(o)
	}

	s.log.Info("Closed all websocket connections")
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
	if s.sessionLogger != nil {
		if err := s.sessionLogger.Close(); err != nil {
			s.log.Warn("Failed to close log writer", zap.Error(err))
		} else {
			s.log.Info("Closed log writer")
		}
	}
}
