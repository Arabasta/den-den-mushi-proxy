package client

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
)

func (c *Connection) WriteClient() {
	for {
		if c.Ctx.Err() != nil {
			c.Log.Info("WriteClient: context done")
			return
		}

		pkt, ok := <-c.WsWriteCh
		if !ok {
			return
		}

		// concurrent write ok if Write not called elsewhere
		if err := c.Sock.WriteMessage(websocket.BinaryMessage, protocol.PacketToByte(pkt)); err != nil {
			if err == io.EOF {
				c.Log.Info("PTY session ended normally")
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				c.Log.Info("WebSocket closed normally") // client sends close, currently not implemented on frontend
			} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				c.Log.Info("WebSocket closed. Probably tab closed") // most closures are this
			} else if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				c.Log.Warn("WebSocket closed unexpectedly", zap.Error(err))
			} else {
				c.Log.Error("Error handling output packet", zap.Error(err))
			}

			if c.Close != nil {
				c.Close()
			}
			return
		}
	}
}
