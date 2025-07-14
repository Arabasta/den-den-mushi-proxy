package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/protocol"
	"go.uber.org/zap"
)

func (s *Session) handleConnPacket(pkt protocol.Packet) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.logPacket(pkt)

	var logMsg string
	var err error

	if pkt.Header == protocol.Input {
		logMsg, err = s.purpose.HandleInput(s, pkt)
	} else {
		logMsg, err = s.purpose.HandleOther(s, pkt)
	}

	if err != nil {
		s.log.Error("Failed to process message", zap.Error(err))
		core_helpers.SendToConn(s.primary, protocol.Packet{
			Header: protocol.Warn,
			Data:   []byte("Failed to process message"),
		})
	}

	if logMsg != "" {
		s.log.Info("Message from handler", zap.String("header", string(pkt.Header)),
			zap.String("message", logMsg))
		s.logL(session_logging.FormatLogLine(string(pkt.Header), logMsg))
	}
}
