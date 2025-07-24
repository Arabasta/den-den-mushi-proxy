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

//var healthcheckAllowedControlHandlers map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error)

var healthcheckBlockedControlHandlers map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error)

func init() {
	// explicitly define allowed for now
	//healthcheckAllowedControlHandlers = map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error){
	//	string(constants.Enter): (*HealthcheckPurpose).handleEnter,
	//
	//	string(constants.Backspace):  (*HealthcheckPurpose).handleBackspace,
	//	string(constants.ArrowLeft):  (*HealthcheckPurpose).HandleLeft,
	//	string(constants.ArrowRight): (*HealthcheckPurpose).HandleRight,
	//	string(constants.CtrlC):      (*HealthcheckPurpose).handleTerminatingControlChar,
	//}

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
		string(constants.OutputRedirection): (*HealthcheckPurpose).handleBlockedControlChar,
		string(constants.InputRedirection):  (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.SingleQuote):       (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.DoubleQuote):       (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.Backtick):          (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.Comma):             (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.Colon):             (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.ExclamationMark):   (*HealthcheckPurpose).handleBlockedControlChar,

		// for now, block control characters that changes the line on pty without
		// explicitly sending the command through the input handler
		// filtering these will require reading the command from the output handler
		// todo: remove cause blocked in global blocked
		//string(constants.ArrowUp):   (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.ArrowDown): (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.CtrlR):     (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.CtrlU):     (*HealthcheckPurpose).handleBlockedControlChar,
		//string(constants.CtrlZ):     (*HealthcheckPurpose).handleBlockedControlChar,

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

//func isAllowedNormalRune(r rune) bool {
//	return (r >= 'a' && r <= 'z') ||
//		(r >= 'A' && r <= 'Z') ||
//		(r >= '0' && r <= '9') ||
//		r == '-' || r == '_' || r == '.' || r == ' ' || r == '/'
//}

// HandleInput for packets with header Input 0x00
func (p *HealthcheckPurpose) HandleInput(s *Session, pkt protocol.Packet) (string, error) {
	if pkt.Header != protocol.Input {
		return "", fmt.Errorf("expected Input header, got %s", string(pkt.Header))
	}

	//data := constants.StripPaste(pkt.Data) // handled in conn_packet.go
	data := pkt.Data

	// check for control char
	//if handlerFunc, ok := healthcheckAllowedControlHandlers[string(data)]; ok {
	//	return handlerFunc(p, s, pkt)
	//}

	if handlerFunc, ok := healthcheckBlockedControlHandlers[string(data)]; ok {
		return handlerFunc(p, s, pkt)
	}

	//if len(data) > 0 {
	//	for _, r := range string(data) {
	//		if isAllowedNormalRune(r) {
	//			s.line.Insert(r)
	//		}
	//	}
	if bytes.Equal(pkt.Data, constants.Enter) {
		return p.handleEnter(s, pkt)
	}

	return handler.Get[protocol.Input].Handle(pkt, s.pty)
	//	}

	//	return "Unhandled input", nil
}

func (p *HealthcheckPurpose) handleEnter(s *Session, pkt protocol.Packet) (string, error) {
	//defer s.line.Reset()
	s.log.Debug("Handling Enter key press. Checking command: ", zap.String("line", s.line.String()))

	msg, allowed := s.filter.IsValid(s.line.String(), s.activePrimary.Claims.OuGroup)
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

		return fmt.Sprintf("\n%s [Command Blocked] %s", time.Now().Format(time.TimeOnly), msg), nil
	}

	return handler.Get[protocol.Input].Handle(pkt, s.pty)
}

//func (p *HealthcheckPurpose) handleBackspace(s *Session, pkt protocol.Packet) (string, error) {
//	s.line.Backspace()
//	return handler.Get[protocol.Input].Handle(pkt, s.pty)
//}
//
//func (p *HealthcheckPurpose) HandleLeft(s *Session, pkt protocol.Packet) (string, error) {
//	s.line.MoveLeft()
//	return handler.Get[protocol.Input].Handle(pkt, s.pty)
//}
//
//func (p *HealthcheckPurpose) HandleRight(s *Session, pkt protocol.Packet) (string, error) {
//	s.line.MoveRight()
//	return handler.Get[protocol.Input].Handle(pkt, s.pty)
//}

//// handleTerminatingControlChar for CtrlC etc that are allowed and terminating
//func (p *HealthcheckPurpose) handleTerminatingControlChar(s *Session, pkt protocol.Packet) (string, error) {
//	defer s.line.Reset()
//	return handler.Get[protocol.Input].Handle(pkt, s.pty)
//}

// ctrl r, up down etc that wants to be blocked
func (p *HealthcheckPurpose) handleBlockedControlChar(s *Session, pkt protocol.Packet) (string, error) {
	// change header and queue it
	pkt.Header = protocol.BlockedControl
	core_helpers.SendToConn(s.activePrimary, pkt)
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
