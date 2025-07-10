package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/handler"
	"den-den-mushi-Go/internal/proxy/protocol"
	"fmt"
)

type ChangeRequestPurpose struct{}

func (p *ChangeRequestPurpose) HandleInput(s *Session, pkt protocol.Packet) (string, error) {
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for packet header: %s", string(pkt.Header))
	}

	// todo: log line on enter, need to add LineEditor to this

	return h.Handle(pkt, s.Pty, s.primary.Sock, s.primary.Claims)
}

func (p *ChangeRequestPurpose) HandleOther(s *Session, pkt protocol.Packet) (string, error) {
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for packet header: %s", string(pkt.Header))
	}

	return h.Handle(pkt, s.Pty, s.primary.Sock, s.primary.Claims)
}
