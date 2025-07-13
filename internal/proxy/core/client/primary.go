package client

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strings"
)

// PrimaryReadLoop should only be accessible by the primary connection
func (c *Connection) PrimaryReadLoop(onPacket func(protocol.Packet)) {
	for {
		if c.Ctx.Err() != nil {
			c.Log.Info("PrimaryReadLoop: context done")
			return
		}

		msgType, msg, err := c.Sock.ReadMessage()
		if err != nil {
			c.logWsError(err)
			// close conn on any error
			if c.Close != nil {
				c.Close()
			}
			return
		}

		pkt := protocol.Parse(msgType, msg)
		if pkt.Header == protocol.ParseError {
			c.Log.Error("Received invalid message from websocket", zap.Any("msg", msg))
			continue
		}

		c.Log.Debug("Received packet from client", zap.Any("packet", pkt))

		onPacket(pkt)
	}
}

func (c *Connection) logWsError(err error) {
	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure) || isUseOfClosedConn(err):
		c.Log.Info("WebSocket closed normally")
	case websocket.IsCloseError(err, websocket.CloseGoingAway):
		c.Log.Info("WebSocket closed. Probably tab closed")
	case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
		c.Log.Warn("WebSocket abnormal closure", zap.Error(err))
	default:
		c.Log.Error("Error reading from websocket", zap.Error(err))
	}
}

func isUseOfClosedConn(err error) bool {
	return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}
