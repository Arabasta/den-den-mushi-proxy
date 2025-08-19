package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
	"io"
	"syscall"
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
	s.terminateSshProcess()
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

func (s *Session) terminateSshProcess() {
	if s.cmd == nil || s.cmd.Process == nil {
		s.log.Warn("SSH process is nil, cannot terminate")
		return
	}

	pid := s.cmd.Process.Pid
	s.log.Info("Terminating SSH process", zap.Int("pid", pid))

	err := s.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		s.log.Warn("Failed to send SIGTERM to SSH process", zap.Error(err))
	}

	done := make(chan error, 1)
	go func() {
		done <- s.cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			s.log.Info("SSH process terminated with error", zap.Error(err))
		} else {
			s.log.Info("SSH process terminated gracefully")
		}
		return
	case <-time.After(10 * time.Second):
		s.log.Warn("SSH process didn't terminate gracefully, forcing kill")
	}

	// force kill if graceful termination failed
	err = s.cmd.Process.Kill()
	if err != nil {
		s.log.Error("Failed to force kill SSH process", zap.Error(err))
	} else {
		s.log.Info("SSH process force killed")
	}

	go func() {
		err := s.cmd.Wait()
		if err != nil {
			s.log.Error("SSH process wait returned error after kill", zap.Error(err))
			return
		}
	}()
}
