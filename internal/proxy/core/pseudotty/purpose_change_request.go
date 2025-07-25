package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/handler"
	"den-den-mushi-Go/internal/proxy/protocol"
	"fmt"
	"go.uber.org/zap"
)

type ChangeRequestPurpose struct{}

func (p *ChangeRequestPurpose) HandleInput(s *Session, pkt protocol.Packet) (string, error) {
	if pkt.Header != protocol.Input {
		s.log.Error("Invalid input packet", zap.Any("pkt", pkt))
		return "", fmt.Errorf("expected Input header, got %s", string(pkt.Header))
	}

	return handler.Get[protocol.Input].Handle(pkt, s.pty)
}

func (p *ChangeRequestPurpose) HandleOther(s *Session, pkt protocol.Packet) (string, error) {
	h, exists := handler.Get[pkt.Header]
	if !exists {
		return "", fmt.Errorf("no handler found for packet header: %s", string(pkt.Header))
	}

	return h.Handle(pkt, s.pty)
}
