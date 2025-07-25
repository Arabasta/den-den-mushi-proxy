package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"go.uber.org/zap"
	"time"
)

// todo: refactor absolute garbage
func (s *Session) handleConnPacket(pkt protocol.Packet) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	//s.logPacket(pkt)

	var logMsg string
	var err error

	if pkt.Header == protocol.Input {
		pkt.Data = constants.StripPaste(pkt.Data)

		if constants.IsGloballyBlockedControlChar(pkt.Data) {
			s.logMsg(s.handleGloballyBlockedControlChar(pkt), pkt)
			return
		}

		// let it go through
		if constants.IsControlCharAffectsLine(pkt.Data) {
			s.updateLineEditor(pkt.Data)
		}

		// update last input for output logging to remove echo
		s.lastInput = pkt.Data

		// insert normal chars into the line editor
		if len(pkt.Data) > 0 {
			for _, r := range string(pkt.Data) {
				if isAsciiOrWhatever(r) {
					s.line.Insert(r)
				}
			}
		}

		logMsg, err = s.purpose.HandleInput(s, pkt)
		s.logAndResetLineEditorIfInputEnter(pkt)
	} else {
		logMsg, err = s.purpose.HandleOther(s, pkt)
	}

	if err != nil {
		s.log.Error("Failed to process message", zap.Error(err))
		core_helpers.SendToConn(s.activePrimary, protocol.Packet{
			Header: protocol.Warn,
			Data:   []byte("Failed to process message"),
		})
	}

	s.logMsg(logMsg, pkt)
}

func (s *Session) handleGloballyBlockedControlChar(pkt protocol.Packet) string {
	// change header and queue it
	pkt.Header = protocol.BlockedControl
	core_helpers.SendToConn(s.activePrimary, pkt)
	return "\n" + time.Now().Format(time.TimeOnly) + "[Blocked Control Character]: " + constants.GloballyBlocked[string(pkt.Data)]
}

func (s *Session) logMsg(msg string, pkt protocol.Packet) {
	if msg == "" {
		return
	}

	s.log.Debug("Message from handler", zap.String("header", string(pkt.Header)),
		zap.String("message", msg))
	s.logL(session_logging.FormatLogLine(string(pkt.Header), msg))
}

func (s *Session) updateLineEditor(b []byte) {
	switch string(b) {
	case string(constants.ArrowRight):
		s.line.MoveRight()
	case string(constants.ArrowLeft):
		s.line.MoveLeft()
	case string(constants.Backspace):
		s.line.Backspace()
	case string(constants.CtrlC):
		s.line.Reset()
	default:
		s.log.Error("updateLineEditor called with unhandled control char")
	}
}

func isAsciiOrWhatever(r rune) bool {
	return (r >= 32 && r <= 126) || r == '\t'
}
