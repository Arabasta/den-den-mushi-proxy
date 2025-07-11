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
	path := fmt.Sprintf("./log/pty_sessions/%s.log", s.id) // todo: add a config option for log path
	if err := os.MkdirAll("./log/pty_sessions", 0755); err != nil {
		s.log.Error("Failed to create log directory", zap.Error(err))
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		s.log.Error("Failed to create log file", zap.Error(err), zap.String("path", path))
		return err
	}

	s.logWriter = file
	s.log.Info("Log writer initialized", zap.String("path", path))
	return nil
}

func (s *Session) logL(line string) {
	if s.logWriter == nil {
		s.log.Error("Log writer is not initialized, cannot log")
		return
	}

	_, err := s.logWriter.Write([]byte(line + "\n"))
	if err != nil {
		s.log.Warn("Failed writing to session log", zap.Error(err), zap.String("line", line))
		return
	}
	s.log.Debug("Logging to session log", zap.String("line", line))
}

func (s *Session) logPacket(pkt protocol.Packet) {
	if pkt.Header == protocol.Resize {
		// don't log resize events
		return
	}
	s.logL(formatLogLine(pkt.Header.String(), string(pkt.Data)))
}

func formatLogLine(header, data string) string {
	return fmt.Sprintf("%s [%s] %s", time.Now().Format(time.TimeOnly), header, data)
}

func getLogHeader(s *Session) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	claims := s.primary.Claims

	header :=
		"# Session Start Time: " + time.Now().UTC().Format(time.RFC3339) + "\n" +
			"# Created By: " + claims.Subject + "\n\n"

	// todo: add ou group

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
	} else if claims.Connection.Purpose == dto.Healthcheck {
		header +=
			"# Health Check Details:\n" +
				"#\t- Filter: " + string(claims.Connection.FilterType)
		// todo: add more stuff
	} else {
		header += "# No additional details for this session purpose\n"
	}

	header += "\n\n\n"

	return header
}

func getLogFooter(s *Session) string {
	footer := "\n# Session End Time: " + s.endTime
	// todo: add list of all users
	return footer
}
