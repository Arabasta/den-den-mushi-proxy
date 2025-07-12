package pseudotty

import (
	"bytes"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/handler"
	"den-den-mushi-Go/internal/proxy/protocol"
	"fmt"
	"go.uber.org/zap"
)

type HealthcheckPurpose struct{}

var healthcheckAllowedControlHandlers map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error)

var healthcheckBlockedControlHandlers map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error)

func init() {
	// explicitly define allowed for now
	healthcheckAllowedControlHandlers = map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error){
		string(Enter): (*HealthcheckPurpose).handleEnter,

		string(Backspace):  (*HealthcheckPurpose).handleBackspace,
		string(ArrowLeft):  (*HealthcheckPurpose).HandleLeft,
		string(ArrowRight): (*HealthcheckPurpose).HandleRight,
		string(CtrlC):      (*HealthcheckPurpose).handleTerminatingControlChar,
	}

	healthcheckBlockedControlHandlers = map[string]func(*HealthcheckPurpose, *Session, protocol.Packet) (string, error){
		// for now block unhandled
		string(RightBackslash):    (*HealthcheckPurpose).handleBlockedControlChar,
		string(LeftBackslash):     (*HealthcheckPurpose).handleBlockedControlChar,
		string(LeftParenthesis):   (*HealthcheckPurpose).handleBlockedControlChar,
		string(RightParenthesis):  (*HealthcheckPurpose).handleBlockedControlChar,
		string(LeftBracket):       (*HealthcheckPurpose).handleBlockedControlChar,
		string(RightBracket):      (*HealthcheckPurpose).handleBlockedControlChar,
		string(LeftBrace):         (*HealthcheckPurpose).handleBlockedControlChar,
		string(RightBrace):        (*HealthcheckPurpose).handleBlockedControlChar,
		string(EqualSign):         (*HealthcheckPurpose).handleBlockedControlChar,
		string(OutputRedirection): (*HealthcheckPurpose).handleBlockedControlChar,
		string(InputRedirection):  (*HealthcheckPurpose).handleBlockedControlChar,
		string(SingleQuote):       (*HealthcheckPurpose).handleBlockedControlChar,
		string(DoubleQuote):       (*HealthcheckPurpose).handleBlockedControlChar,
		string(Backtick):          (*HealthcheckPurpose).handleBlockedControlChar,
		string(Comma):             (*HealthcheckPurpose).handleBlockedControlChar,
		string(Colon):             (*HealthcheckPurpose).handleBlockedControlChar,
		string(ExclamationMark):   (*HealthcheckPurpose).handleBlockedControlChar,

		// for now, block control characters that changes the line on pty without
		// explicitly sending the command through the input handler
		// filtering these will require reading the command from the output handler
		string(ArrowUp):   (*HealthcheckPurpose).handleBlockedControlChar,
		string(ArrowDown): (*HealthcheckPurpose).handleBlockedControlChar,
		string(CtrlR):     (*HealthcheckPurpose).handleBlockedControlChar,
		string(CtrlU):     (*HealthcheckPurpose).handleBlockedControlChar,
		string(CtrlZ):     (*HealthcheckPurpose).handleBlockedControlChar,

		// block pipe, sequential, background execution
		string(Pipe):      (*HealthcheckPurpose).handleBlockedControlChar,
		string(SemiColon): (*HealthcheckPurpose).handleBlockedControlChar,
		string(Ampersand): (*HealthcheckPurpose).handleBlockedControlChar,

		// block $ for env vars
		// this is required as users can use env vars to bypass the filter
		// e.g.,
		// a=su
		// $a
		string(DollarSign): (*HealthcheckPurpose).handleBlockedControlChar,
	}
}

func isAllowedNormalRune(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_' || r == '.' || r == ' ' || r == '/'
}

// HandleInput for packets with header Input 0x00
func (p *HealthcheckPurpose) HandleInput(s *Session, pkt protocol.Packet) (string, error) {
	if pkt.Header != protocol.Input {
		return "", fmt.Errorf("expected Input header, got %s", string(pkt.Header))
	}

	data := pkt.Data

	// for now block pasting
	if bytes.Equal(data, PasteStart) || bytes.Equal(data, PasteEnd) {
		// todo: handle paste
		return handler.Get[protocol.BlockedControl].Handle(pkt, s.Pty)
	}

	// check for control char
	if handlerFunc, ok := healthcheckAllowedControlHandlers[string(data)]; ok {
		return handlerFunc(p, s, pkt)
	}

	if handlerFunc, ok := healthcheckBlockedControlHandlers[string(data)]; ok {
		return handlerFunc(p, s, pkt)
	}

	// check for normal text
	if len(data) == 1 && isAllowedNormalRune(rune(data[0])) {
		s.line.Insert(rune(data[0]))
		return handler.Get[protocol.Input].Handle(pkt, s.Pty)
	}

	return "Unhandled input", nil
}

func (p *HealthcheckPurpose) handleEnter(s *Session, pkt protocol.Packet) (string, error) {
	defer s.line.Reset()

	s.log.Info("Handling Enter key press. Checking command: ", zap.String("line", s.line.String()))

	msg, allowed := s.filter.IsValid(s.line.String())
	if !allowed {
		errPkt := protocol.Packet{Header: protocol.BlockedCommand, Data: []byte(s.line.String())}
		s.ptyOutput.Add(errPkt)
		s.fanout(errPkt)

		// send Ctrl +C to clear pty
		ctrlCPkt := protocol.Packet{Header: protocol.Input, Data: CtrlC}
		_, err := handler.Get[protocol.Input].Handle(ctrlCPkt, s.Pty)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("\n[Command Blocked] %s", msg), nil
	}

	return handler.Get[protocol.Input].Handle(pkt, s.Pty)
}

func (p *HealthcheckPurpose) handleBackspace(s *Session, pkt protocol.Packet) (string, error) {
	s.line.Backspace()
	return handler.Get[protocol.Input].Handle(pkt, s.Pty)
}

func (p *HealthcheckPurpose) HandleLeft(s *Session, pkt protocol.Packet) (string, error) {
	s.line.MoveLeft()
	return handler.Get[protocol.Input].Handle(pkt, s.Pty)
}

func (p *HealthcheckPurpose) HandleRight(s *Session, pkt protocol.Packet) (string, error) {
	s.line.MoveRight()
	return handler.Get[protocol.Input].Handle(pkt, s.Pty)
}

// handleTerminatingControlChar for CtrlC etc that are allowed and terminating
func (p *HealthcheckPurpose) handleTerminatingControlChar(s *Session, pkt protocol.Packet) (string, error) {
	defer s.line.Reset()
	return handler.Get[protocol.Input].Handle(pkt, s.Pty)
}

// ctrl r, up down etc that wants to be blocked
func (p *HealthcheckPurpose) handleBlockedControlChar(s *Session, pkt protocol.Packet) (string, error) {
	// change header and queue it
	pkt.Header = protocol.BlockedControl
	client.SendToConn(s.primary, pkt)
	return "\n[Blocked Control Character]: " + string(pkt.Data), nil
}

// HandleOther for pkt headers that are not input 0x00
func (p *HealthcheckPurpose) HandleOther(s *Session, pkt protocol.Packet) (string, error) {
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for packet header: %s", string(pkt.Header))
	}

	return h.Handle(pkt, s.Pty)
}

var (
	Enter     = []byte{13}
	Backspace = []byte{127}

	ArrowUp    = []byte{27, 91, 65}
	ArrowDown  = []byte{27, 91, 66}
	ArrowRight = []byte{27, 91, 67}
	ArrowLeft  = []byte{27, 91, 68}

	CtrlR = []byte{18}
	CtrlC = []byte{3}
	CtrlZ = []byte{26}
	CtrlU = []byte{21}

	PasteStart = []byte{27, 91, 50, 48, 48, 126}
	PasteEnd   = []byte{27, 91, 50, 48, 49, 126}

	SemiColon = []byte{59}
	Ampersand = []byte{38}
	Pipe      = []byte{124}

	// todo: handle >>, 2>, &>, >&2, <<<, <<, 2>&1, etc
	OutputRedirection = []byte{62} // >
	InputRedirection  = []byte{60} // <

	SingleQuote = []byte{39} // '
	DoubleQuote = []byte{34} // "
	Backtick    = []byte{96} // `

	Comma           = []byte{44} // ,
	Colon           = []byte{58} // :
	ExclamationMark = []byte{33} // !

	LeftParenthesis  = []byte{40}  // (
	RightParenthesis = []byte{41}  // )
	LeftBracket      = []byte{91}  // [
	RightBracket     = []byte{93}  // ]
	LeftBrace        = []byte{123} // {
	RightBrace       = []byte{125} // }

	DollarSign     = []byte{36} // $
	EqualSign      = []byte{61} // =
	RightBackslash = []byte{92} // \
	LeftBackslash  = []byte{47} // /
)
