package client

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Sock      *websocket.Conn
	Claims    *token.Claims
	WsWriteCh chan protocol.Packet
	Close     func()
}

func (c *Connection) WriteClient() {
	for pkt := range c.WsWriteCh {
		b := protocol.PacketToByte(pkt)

		// concurrent write ok if Write not called elsewhere
		err := c.Sock.WriteMessage(websocket.BinaryMessage, b)

		if err != nil {
			// todo handle err
			//if err == io.EOF {
			//	s.Log.Info("PTY session ended normally")
			//} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			//	s.Log.Info("WebSocket closed normally") // client sends close, currently not implemented in frontend
			//} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			//	s.Log.Info("WebSocket closed. Probably tab closed") // most closures are this
			//} else if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
			//	s.Log.Warn("WebSocket closed unexpectedly", zap.Error(err))
			//} else {
			//	s.Log.Error("Error handling output packet", zap.Error(err))
			//}
			//s.closeAll(ws, pty)
			return
		}
	}
}
