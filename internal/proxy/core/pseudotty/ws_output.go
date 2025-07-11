package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/protocol"
)

// fanout sends packet to primary and all observers' channels, called by event loop
func (s *Session) fanout(pkt protocol.Packet) {
	if s.primary != nil {
		sendToConn(s.primary, pkt)
	}
	for o := range s.observers {
		sendToConn(o, pkt)
	}
}

func (s *Session) fanoutExcept(pkt protocol.Packet, except *client.Connection) {
	if s.primary != nil && s.primary != except {
		sendToConn(s.primary, pkt)
	}
	for o := range s.observers {
		if o != except {
			sendToConn(o, pkt)
		}
	}
}

// sendToConn sends packet to a specific connection, used for targeted messages
func sendToConn(c *client.Connection, pkt protocol.Packet) {
	if c == nil {
		return
	}
	send(c.WsWriteCh, pkt)
}

// send just sends a packet to a channel
func send(ch chan protocol.Packet, pkt protocol.Packet) {
	select {
	case ch <- pkt:
	default:
		// queue full, todo discuss drop or log or error or what
	}
}
