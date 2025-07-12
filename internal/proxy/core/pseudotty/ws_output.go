package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/protocol"
)

// fanout to primary and all observers' channels, called in event loop
func (s *Session) fanout(pkt protocol.Packet) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.primary != nil {
		client.Send(s.primary.WsWriteCh, pkt)
	}

	for o := range s.observers {
		client.Send(o.WsWriteCh, pkt)
	}
}

func (s *Session) fanoutExcept(pkt protocol.Packet, except *client.Connection) {
	if s.primary != except {
		client.SendToConn(s.primary, pkt)
	}
	for o := range s.observers {
		if o != except {
			client.SendToConn(o, pkt)
		}
	}
}
