package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
	"io"
	"time"
)

func (s *Session) ForceEndSession(reason string) {
	s.log.Warn("Forcing end of pty session")
	s.logL(session_logging.FormatLogLine("System", "Force shutdown of session... Reason: "+reason))
	s.EndSession()
}

func (s *Session) EndSession() {
	s.once.Do(func() {
		s.cancel() // exit conn loop
		if s.State == types.Closed {
			return
		}
		s.State = types.Closed

		s.log.Info("Ending pty session")
		s.closeTheWorld()

		if s.onClose != nil {
			s.onClose(s.Id)
		}
		s.log.Debug("Session closed")
	})
}

func (s *Session) closeTheWorld() {
	s.log.Debug("Closing the world")
	s.deregisterAllWsConnections()
	s.closePty()
	s.endTime = time.Now()
	s.logL(session_logging.FormatFooter(s.endTime))
	s.closeLogWriter()
}

func (s *Session) deregisterAllWsConnections() {
	s.mu.RLock()
	// don't use deregisterCh
	primary := s.ActivePrimary
	observers := make([]*client.Connection, 0, len(s.ActiveObservers))
	for o := range s.ActiveObservers {
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

// todo: more error handling eg check if pty is arleady Closed
func (s *Session) closePty() {
	if s.pty == nil {
		s.log.Warn("PTY is nil, cannot close")
		return
	}
	err := s.pty.Close()
	if err != nil {
		if err == io.EOF {
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
