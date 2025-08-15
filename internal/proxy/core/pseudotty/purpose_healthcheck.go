package pseudotty

import (
	"bytes"
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/handler"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// todo: very very bad need to refactor
type HealthcheckPurpose struct{}

var healthcheckBlockedControlHandlers map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error)

// HandleInput for packets with header Input 0x00
func (p *HealthcheckPurpose) HandleInput(s *Session, pkt protocol.Packet) (string, error) {
	if pkt.Header != protocol.Input {
		s.log.Error("Invalid input packet", zap.Any("pkt", pkt))
		return "", fmt.Errorf("expected Input header, got %s", string(pkt.Header))
	}

	for _, b := range pkt.Data {

		if handlerFunc, ok := healthcheckBlockedControlHandlers[string([]byte{b})]; ok {
			return handlerFunc(p, s, pkt)
		}

		if bytes.Equal([]byte{b}, constants.Enter) {
			return p.handleEnter(s, pkt)
		}

		return handler.Get[protocol.Input].Handle(pkt, s.pty)
	}
	return "", nil
}

func (p *HealthcheckPurpose) handleEnter(s *Session, pkt protocol.Packet) (string, error) {
	s.log.Debug("Handling Enter key press. Checking command: ", zap.String("line", s.line.String()))

	s.log.Debug("Calling filter.IsValid for command", zap.String("line", s.line.String()), zap.String("ouGroup", s.ActivePrimary.Claims.OuGroup))
	msg, allowed := s.filter.IsValid(s.line.String(), s.ActivePrimary.Claims.OuGroup)
	s.log.Debug("Filter result", zap.String("line", s.line.String()), zap.Bool("allowed", allowed), zap.String("message", msg))
	if !allowed {
		errPkt := protocol.Packet{Header: protocol.BlockedCommand, Data: []byte(s.line.String())}
		s.ptyOutput.Add(errPkt)
		s.fanout(errPkt, nil)

		// send Ctrl +C to clear pty
		ctrlCPkt := protocol.Packet{Header: protocol.Input, Data: constants.CtrlC}
		_, err := handler.Get[protocol.Input].Handle(ctrlCPkt, s.pty)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("\n%s [Command Blocked] %s", time.Now().Format(time.TimeOnly), msg), CommandBlockedError
	}

	s.log.Debug("Checked command", zap.String("line", s.line.String()), zap.Bool("allowed", allowed))

	return handler.Get[protocol.Input].Handle(pkt, s.pty)
}

// ctrl r, up down etc that wants to be blocked
func (p *HealthcheckPurpose) handleBlockedControlChar(s *Session, pkt protocol.Packet) (string, error) {
	// change header and queue it
	pkt.Header = protocol.BlockedControl
	core_helpers.SendToConn(s.ActivePrimary, pkt)
	return "\n" + time.Now().Format(time.TimeOnly) + "[Blocked Control Character]: " + string(pkt.Data), nil
}

// HandleOther for pkt headers that are not input 0x00
func (p *HealthcheckPurpose) HandleOther(s *Session, pkt protocol.Packet) (string, error) {
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for packet header: %s", string(pkt.Header))
	}

	return h.Handle(pkt, s.pty)
}

func init() {
	healthcheckBlockedControlHandlers = map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error){
		// for now block unhandled
		string(constants.RightBackslash): (*HealthcheckPurpose).handleBlockedControlChar, // \
		//string(constants.LeftBackslash):     (*HealthcheckPurpose).handleBlockedControlChar, // /
		//string(constants.LeftParenthesis):   (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.RightParenthesis):  (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.LeftBracket):       (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.RightBracket):      (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.LeftBrace):         (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.RightBrace):        (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.EqualSign):         (*HealthcheckPurpose).handleBlockedControlChar,

		// Edwin said to not block <>
		// string(constants.OutputRedirection): (*HealthcheckPurpose).handleBlockedControlChar,
		// string(constants.InputRedirection):  (*HealthcheckPurpose).handleBlockedControlChar,

		//string(constants.SingleQuote):       (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.DoubleQuote):       (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.Backtick):          (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.Comma):             (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.Colon):             (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.ExclamationMark):   (*HealthcheckPurpose).handleBlockedControlChar,

		// block pipe, sequential, background execution
		//string(constants.Pipe):      (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.SemiColon): (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.Ampersand): (*HealthcheckPurpose).handleBlockedControlChar,

		// block $ for env vars
		// this is required as users can use env vars to bypass the filter
		// e.g.,
		// a=su
		// $a
		//string(constants.DollarSign): (*HealthcheckPurpose).handleBlockedControlChar,
	}
}
