package client

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"sync"
)

type Connection struct {
	Sock             *websocket.Conn
	Claims           *token.Claims
	WsWriteCh        chan protocol.Packet
	OnceCloseWriteCh sync.Once
	Close            func()
}

func (c *Connection) WriteClient(log *zap.Logger) {
	for pkt := range c.WsWriteCh {
		b := protocol.PacketToByte(pkt)

		// concurrent write ok if Write not called elsewhere
		err := c.Sock.WriteMessage(websocket.BinaryMessage, b)
		if err != nil {
			if err == io.EOF {
				log.Info("PTY session ended normally")
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Info("WebSocket closed normally") // client sends close, currently not implemented on frontend
			} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Info("WebSocket closed. Probably tab closed") // most closures are this
			} else if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Warn("WebSocket closed unexpectedly", zap.Error(err))
			} else {
				log.Error("Error handling output packet", zap.Error(err))
			}

			c.Close()
			return
		}
	}
}
