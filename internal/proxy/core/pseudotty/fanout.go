package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/protocol"
)

func (s *Session) fanout(pkt protocol.Packet, except *client.Connection) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.ActivePrimary != except {
		core_helpers.SendToConn(s.ActivePrimary, pkt)
	}
	for o := range s.ActiveObservers {
		if o != except {
			core_helpers.SendToConn(o, pkt)
		}
	}
}
