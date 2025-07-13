package core_helpers

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/protocol"
	"go.uber.org/zap"
)

// SendToConn sends packet to a specific connection, used for targeted messages
func SendToConn(c *client.Connection, pkt protocol.Packet) {
	if c == nil {
		return
	}

	Send(c.WsWriteCh, pkt)
}

// Send just sends a packet to a channel
func Send(ch chan protocol.Packet, pkt protocol.Packet) {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("Recovered in Send", zap.Any("error", r))
		}
	}()

	select {
	case ch <- pkt:
	default:
		// queue full
		zap.L().Warn("WebSocket write channel is full, dropping packet", zap.Any("packet", pkt))
		//todo drop oldest
	}
}
