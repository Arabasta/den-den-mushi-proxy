package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/protocol"
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
