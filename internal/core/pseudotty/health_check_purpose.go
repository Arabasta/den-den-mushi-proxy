package pseudotty

import (
	"bytes"
	"den-den-mushi-Go/internal/handler"
	"den-den-mushi-Go/internal/protocol"
	"fmt"
)

type HealthcheckPurpose struct{}

/**
control characters
- check if banned else write to pty
- if stuff like ctrl+c, clear input buffer

left right arrows
- need to handle cursor movement

enter
- use filter
- if allow write to pty and clear input buffer

backspace
- remove last character from input buffer and write to pty

*/

func (p *HealthcheckPurpose) HandleInput(s *Session, pkt protocol.Packet) (string, error) {

	data := pkt.Data

	switch {
	case bytes.Equal(data, Enter):
		return p.handleEnter(s, pkt)
	case bytes.Equal(data, Backspace):
		return p.handleBackspace(s, pkt)
	case bytes.Equal(data, ArrowLeft):
		return p.HandleLeft(s, pkt)
	case bytes.Equal(data, ArrowRight):
		return p.HandleRight(s, pkt)

	default:
		// todo: insert all char not just first rune
		s.LineEditor.Insert(rune(data[0]))
		_, err := s.Pty.Write(data)
		return "", err
	}
}

func (p *HealthcheckPurpose) handleEnter(s *Session, pkt protocol.Packet) (string, error) {
	defer s.LineEditor.Reset()

	reason, allowed := s.Filter.Feed(s.LineEditor.String())
	if !allowed {
		return fmt.Sprintf("\n[Command Blocked] %s", reason), nil
	}
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for header: %s", pkt.Header)
	}
	return h.Handle(pkt, s.Pty, s.Primary.Sock, s.Primary.Claims)

}

func (p *HealthcheckPurpose) handleBackspace(s *Session, pkt protocol.Packet) (string, error) {
	s.LineEditor.Backspace()
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for header: %s", pkt.Header)
	}
	return h.Handle(pkt, s.Pty, s.Primary.Sock, s.Primary.Claims)
}

func (p *HealthcheckPurpose) HandleLeft(s *Session, pkt protocol.Packet) (string, error) {
	s.LineEditor.MoveLeft()
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for header: %s", pkt.Header)
	}
	return h.Handle(pkt, s.Pty, s.Primary.Sock, s.Primary.Claims)
}

func (p *HealthcheckPurpose) HandleRight(s *Session, pkt protocol.Packet) (string, error) {
	s.LineEditor.MoveRight()
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for header: %s", pkt.Header)
	}
	return h.Handle(pkt, s.Pty, s.Primary.Sock, s.Primary.Claims)
}

// ctrl r, up down etc that wants to be blocked
func (p *HealthcheckPurpose) handleBlockedControlChar(s *Session, pkt protocol.Packet) (string, error) {
	return "", nil

}

// ctrl c u etc that are allowed and terminating
func (p *HealthcheckPurpose) handleTerminatingControlChar(s *Session, pkt protocol.Packet) (string, error) {
	defer s.LineEditor.Reset()
	return "", nil

}

func (p *HealthcheckPurpose) HandleOther(s *Session, pkt protocol.Packet) (string, error) {

	return "", nil
}

var (
	ArrowUp          = []byte{27, 91, 65}
	ArrowDown        = []byte{27, 91, 66}
	ArrowRight       = []byte{27, 91, 67}
	ArrowLeft        = []byte{27, 91, 68}
	CtrlR            = []byte{18}
	Enter            = []byte{13}
	Backspace        = []byte{127}
	PasteStart       = []byte{27, 91, 50, 48, 48, 126}
	PasteEnd         = []byte{27, 91, 50, 48, 49, 126}
	SemiColon        = []byte{59}
	Ampersand        = []byte{38}
	Pipe             = []byte{124}
	LeftParenthesis  = []byte{40}
	RightParenthesis = []byte{41}
)
