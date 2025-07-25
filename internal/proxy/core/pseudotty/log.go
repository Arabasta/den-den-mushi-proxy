package pseudotty

import (
	"bytes"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
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

// todo: refactor
func (s *Session) initSessionLoggerForAIThing() error {
	path := fmt.Sprintf("./log/pty_sessions_input_only/%s.log", s.Id) // todo: add a config option for log path
	err := os.MkdirAll("./log/pty_sessions_input_only", 0755)
	if err != nil {
		s.log.Error("Failed to create log directory", zap.Error(err))
		return err
	}

	s.sessionLoggerForAIThing, err = session_logging.NewFileSessionLogger(path)

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
	//if pkt.Header == protocol.Resize {
	// don't log resize events
	//return
	//}
	switch pkt.Header {
	case protocol.Output:
		s.handleOutputLogging(pkt)
	case protocol.Resize:
		return
	default:
		s.logL(session_logging.FormatLogLine(pkt.Header.String(), string(pkt.Data)))
	}
}

func (s *Session) handleOutputLogging(pkt protocol.Packet) {
	if bytes.Equal(pkt.Data, s.lastInput) {
		s.log.Debug("Skipping output logging for echo input", zap.ByteString("data", pkt.Data))
		return
	}
	clean := stripNonASCII(pkt.Data)
	if len(clean) == 0 {
		return // skip logging if nothing printable
	}

	s.logL(session_logging.FormatLogLine(pkt.Header.String(), string(clean)))
}

func (s *Session) logAndResetLineEditorIfInputEnter(pkt protocol.Packet) {
	if bytes.Equal(pkt.Data, constants.Enter) {
		// log to pty session log
		s.logPacket(protocol.Packet{
			Header: protocol.Input,
			Data:   []byte(s.line.String()),
		})
		// log to yirong's ai log todo: refactor
		line := session_logging.FormatLogLine(pkt.Header.String(), s.line.String())
		if s.sessionLoggerForAIThing == nil {
			s.log.Error("Log writer is not initialized, cannot log")
			return
		}
		if err := s.sessionLoggerForAIThing.WriteLine(line); err != nil {
			s.log.Warn("Failed writing to session log", zap.Error(err), zap.String("line", line))
			return
		}
		s.line.Reset()
	}
}

func stripNonASCII(data []byte) []byte {
	out := make([]rune, 0, len(data))
	for _, r := range string(data) {
		if r >= 32 && r <= 126 || r == '\n' || r == '\r' || r == '\t' {
			out = append(out, r)
		}
	}
	return []byte(string(out))
}
