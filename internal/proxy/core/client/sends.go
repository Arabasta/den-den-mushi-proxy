package client

import "den-den-mushi-Go/internal/proxy/protocol"

// SendToConn sends packet to a specific connection, used for targeted messages
func SendToConn(c *Connection, pkt protocol.Packet) {
	if c == nil {
		return
	}
	Send(c.WsWriteCh, pkt)
}

// Send just sends a packet to a channel
func Send(ch chan protocol.Packet, pkt protocol.Packet) {
	select {
	case ch <- pkt:
	default:
		// queue full, todo discuss drop or log or error or what
	}
}
