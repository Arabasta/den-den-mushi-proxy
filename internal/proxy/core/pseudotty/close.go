package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"go.uber.org/zap"
	"io"
	"sync"
)

func (s *Session) EndSession() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	defer func() { s.closed = true }()

	s.Log.Info("Ending pty session")
	s.closeTheWorld()
}

func (s *Session) closeWs(c *client.Connection) {
	c.Close()
}

// todo: more error handling eg check if pty is arleady closed
func (s *Session) closePty() {
	err := s.Pty.Close()
	if err != nil {
		if err == io.EOF {
			s.Log.Info("PTY session ended normally")
		} else {
			s.Log.Warn("Failed to close pty", zap.Error(err))
		}
	}
	s.Log.Info("Closed pty")
}

func (s *Session) closeLogWriter() {
	if s.logWriter != nil {
		if err := s.logWriter.Close(); err != nil {
			s.Log.Warn("Failed to close log writer", zap.Error(err))
		} else {
			s.Log.Info("Closed log writer")
		}
	}
}

func (s *Session) closeTheWorld() func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			s.closeWs(s.Primary)
			for o := range s.observers {
				s.closeWs(o)
			}
			s.closePty()
			s.closeLogWriter()
		})
	}
}
