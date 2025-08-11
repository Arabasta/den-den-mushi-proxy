package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/handler"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"den-den-mushi-Go/pkg/types"
	"den-den-mushi-Go/pkg/util/cyberark"
	"errors"
	"go.uber.org/zap"
	"slices"
	"time"
)

// todo: refactor absolute garbage
func (s *Session) handleConnPacket(pkt protocol.Packet) {
	var released bool
	s.mu.RLock()
	defer func() {
		if !released {
			s.mu.RUnlock()
		}
	}()

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
			// should not log, should not add to last input, should not insert into line editor
			logMsg, err = s.purpose.HandleInput(s, pkt)
			return
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
		if err != nil && errors.Is(err, CommandBlockedError) {
			// skip
		} else {
			s.logAndResetLineEditorIfInputEnter(pkt)
		}
	} else if pkt.Header == protocol.Sudo && s.startClaims.Connection.Purpose == types.Change {
		s.logPacket(pkt)
		/// todo: refactor this garbage
		//// packet should contain username to sudo to
		//targetUser := string(pkt.Data)
		//s.log.Debug("Handling Sudo packet", zap.String("target OS user", targetUser))

		// packet should contain cyberark object to draw password for
		cyberarkObject := string(pkt.Data)
		s.log.Debug("Handling Sudo packet", zap.String("cyberark object", cyberarkObject))

		// extract ip and os user from cyberarkObject
		ip := cyberark.ExtractIPFromObject(cyberarkObject)
		targetUser := cyberark.ExtractOsUserFromObject(cyberarkObject)

		if s.puppetClient.Cfgtmp.PuppetTasks.CyberarkPasswordDraw.IsValidationEnabled {
			// check ip against initial claims
			if ip != s.startClaims.Connection.Server.IP {
				s.logL(session_logging.FormatLogLine(pkt.Header.String(), "Unauthorized sudo attempt IP mismatch"))
				s.log.Error("Unauthorized sudo attempt IP mismatch", zap.String("ip", ip),
					zap.String("expectedIP", s.startClaims.Connection.Server.IP))
				return
			}

			// check targetUser against initial claims
			if !slices.Contains(s.startClaims.Connection.AllowedSuOsUsers, targetUser) {
				s.logL(session_logging.FormatLogLine(pkt.Header.String(), "Unauthorized sudo attempt Target User Mismatch"))
				s.log.Error("Unauthorized sudo attempt Target User Mismatch", zap.String("user", targetUser),
					zap.Strings("allowedUsers", s.startClaims.Connection.AllowedSuOsUsers))
				return
			}
		}

		// draw password from cyberark
		password, err := s.puppetClient.DrawCyberarkKey(cyberarkObject, s.startClaims.Connection.ServerFQDNTmpTillRefactor)
		if err != nil {
			s.logL(session_logging.FormatLogLine(pkt.Header.String(), "Failed to draw cyberark key"))
			s.log.Error("Failed to draw cyberark key", zap.Error(err), zap.String("cyberark object", cyberarkObject))
			return
		}

		s.log.Debug("Successfully drew cyberark key, gonna su now", zap.String("cyberark object", cyberarkObject),
			zap.String("target OS user", targetUser))
		passwordPacket := protocol.Packet{
			Header: protocol.SudoInputPassword,
			Data:   []byte(password),
		}

		userPacket := protocol.Packet{
			Header: protocol.SudoInputUser,
			Data:   []byte(targetUser),
		}

		s.tmpMuForPtyThingTillRefactor.Lock()
		s.isPtyOutputLocked = true

		defer func() {
			s.logL(session_logging.FormatLogLine(pkt.Header.String(), "Sudo Operation for "+targetUser+" complete"))
			s.isPtyOutputLocked = false
			logMsg, err = handler.Get[protocol.Input].Handle(protocol.Packet{
				Header: protocol.Input,
				Data:   append(constants.CtrlC, constants.Enter...),
			}, s.pty)
			time.Sleep(1000 * time.Millisecond)

			// clear the screen just in case password leaks
			s.log.Debug("Clearing screen after sudo operation", zap.String("target OS user", targetUser))
			logMsg, err = handler.Get[protocol.Input].Handle(protocol.Packet{
				Header: protocol.Input,
				Data:   append([]byte("clear"), constants.Enter...),
			}, s.pty)
			s.tmpMuForPtyThingTillRefactor.Unlock()
		}()

		// call SudoUsernameHandler
		s.log.Debug("Handling Sudo username packet", zap.String("target OS user", targetUser))
		logMsg, err = handler.Get[userPacket.Header].Handle(userPacket, s.pty)
		if err != nil {
			s.logL(session_logging.FormatLogLine(pkt.Header.String(), "Failed to handle SudoUsername packet"))
			s.log.Error("Failed to handle SudoUsername packet", zap.Error(err))
			return
		}

		// call SudoPasswordHandler
		s.log.Debug("Handling Sudo password packet", zap.String("target OS user", targetUser))
		logMsg, err = handler.Get[passwordPacket.Header].Handle(passwordPacket, s.pty)
		if err != nil {
			s.logL(session_logging.FormatLogLine(pkt.Header.String(), "Failed to handle SudoPassword packet"))
			s.log.Error("Failed to handle SudoPassword packet", zap.Error(err))
			return
		}
	} else if pkt.Header == protocol.ClientEndSession {
		s.log.Debug("Handling ClientEndSession packet")
		released = true
		s.mu.RUnlock()
		s.logL(session_logging.FormatLogLine(pkt.Header.String(), "Received end session from client"))
		s.EndSession()
		return
	} else {
		logMsg, err = s.purpose.HandleOther(s, pkt)
	}

	if err != nil {
		s.log.Error("Failed to process message", zap.Error(err))
		core_helpers.SendToConn(s.ActivePrimary, protocol.Packet{
			Header: protocol.Warn,
			Data:   []byte("Failed to process message"),
		})
	}

	s.logMsg(logMsg, pkt)
}

func (s *Session) handleGloballyBlockedControlChar(pkt protocol.Packet) string {
	// change header and queue it
	pkt.Header = protocol.BlockedControl
	core_helpers.SendToConn(s.ActivePrimary, pkt)
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
	case string(constants.Home):
		s.line.MoveStart()
	case string(constants.End):
		s.line.MoveEnd()
	default:
		s.log.Error("updateLineEditor called with unhandled control char")
	}
}

func isAsciiOrWhatever(r rune) bool {
	return (r >= 32 && r <= 126) || r == '\t'
}
