package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/dto"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

func (s *Session) initLogWriter() error {
	path := fmt.Sprintf("./log/pty_sessions/%s.log", s.id)
	if err := os.MkdirAll("./log/pty_sessions", 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	s.logWriter = file
	s.Log.Info("Log writer initialized", zap.String("path", path))
	return nil
}

func (s *Session) logf(format string, args ...any) {
	if s.logWriter == nil {
		return
	}
	line := fmt.Sprintf(format, args...)

	_, err := s.logWriter.Write([]byte(line))
	if err != nil {
		s.Log.Warn("Failed writing to session log", zap.Error(err), zap.String("line", line))
	}
}

func (s *Session) logLine(h protocol.Header, data string) {
	s.logf("\n%s [%s] %s", time.Now().Format(time.TimeOnly), h, data)
}

// LogHeader to be called only once when the session starts
func (s *Session) LogHeader() {
	s.mu.Lock()
	defer s.mu.Unlock()

	claims := s.primary.Claims

	header :=
		"# Session Start Time: " + time.Now().UTC().Format(time.RFC3339) + "\n" +
			"# Created By: " + claims.Subject + "\n\n"

	header += "# Connection Details:\n" +
		"#\t- Server IP: " + claims.Connection.Server.IP + "\n" +
		"#\t- OS User: " + claims.Connection.Server.OSUser + "\n" +
		"#\t- Port: " + claims.Connection.Server.Port + "\n\n"

	header += "# Purpose: " + string(claims.Connection.Purpose) + "\n\n"

	if claims.Connection.Purpose == dto.Change {
		header +=
			"# Change Request Details:\n" +
				"#\t- Change Request ID: " + claims.Connection.ChangeRequest.Id + "\n" +
				"#\t- Implementor Group: " + claims.Connection.ChangeRequest.ImplementorGroup + "\n" +
				"#\t- End Time: " + claims.Connection.ChangeRequest.EndTime
	}

	header += "\n\n\n"

	s.logf(header)
}

func (s *Session) logFooter() {
	footer := "\n# Session End Time: " + s.endTime
	s.logf(footer)
}
