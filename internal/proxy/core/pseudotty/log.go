package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/protocol"
	"fmt"
	"go.uber.org/zap"
	"os"
)

func (s *Session) initSessionLogger() error {
	path := fmt.Sprintf("./log/pty_sessions/%s.log", s.Id) // todo: add a config option for log path
	err := os.MkdirAll("./log/pty_sessions", 0755)
	if err != nil {
		s.log.Error("Failed to create log directory", zap.Error(err))
		return err
	}

	s.sessionLogger, err = session_logging.NewFileSessionLogger(path)

	if err != nil {
		s.log.Error("Failed to initialize session logger", zap.Error(err), zap.String("path", path))
		return err
	}
	s.log.Info("Log writer initialized", zap.String("path", path))
	return nil
}

func (s *Session) logL(line string) {
	if s.sessionLogger == nil {
		s.log.Error("Log writer is not initialized, cannot log")
		return
	}

	if err := s.sessionLogger.WriteLine(line); err != nil {
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
	s.logL(session_logging.FormatLogLine(pkt.Header.String(), string(pkt.Data)))
}
