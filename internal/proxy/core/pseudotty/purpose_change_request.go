package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/handler"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/constants"
	"fmt"
)

// todo: very very bad need to refactor
type ChangeRequestPurpose struct{}

var crBlockedControlHandlers map[string]func(*ChangeRequestPurpose, *Session, protocol.Packet) (string, error)

func init() {
	crBlockedControlHandlers = map[string]func(*ChangeRequestPurpose, *Session, protocol.Packet) (string, error){

		// for now, block control characters that changes the line on pty without
		// explicitly sending the command through the input handler
		// filtering these will require reading the command from the output handler
		string(constants.ArrowUp):   (*ChangeRequestPurpose).handleBlockedControlChar,
		string(constants.ArrowDown): (*ChangeRequestPurpose).handleBlockedControlChar,
		string(constants.CtrlR):     (*ChangeRequestPurpose).handleBlockedControlChar,
		string(constants.CtrlU):     (*ChangeRequestPurpose).handleBlockedControlChar,
		string(constants.CtrlZ):     (*ChangeRequestPurpose).handleBlockedControlChar}
}

func (p *ChangeRequestPurpose) HandleInput(s *Session, pkt protocol.Packet) (string, error) {
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for packet header: %s", string(pkt.Header))
	}

	if handlerFunc, ok := crBlockedControlHandlers[string(pkt.Data)]; ok {
		return handlerFunc(p, s, pkt)
	}

	// todo: log line on enter, need to add LineEditor to this

	return h.Handle(pkt, s.pty)
}
func (p *ChangeRequestPurpose) handleBlockedControlChar(s *Session, pkt protocol.Packet) (string, error) {
	// change header and queue it
	pkt.Header = protocol.BlockedControl
	core_helpers.SendToConn(s.activePrimary, pkt)
	return "\n[Blocked Control Character]: " + string(pkt.Data), nil
}

func (p *ChangeRequestPurpose) HandleOther(s *Session, pkt protocol.Packet) (string, error) {
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for packet header: %s", string(pkt.Header))
	}

	return h.Handle(pkt, s.pty)
}
