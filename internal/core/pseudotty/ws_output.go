package pseudotty

import (
	"den-den-mushi-Go/internal/core/client"
	"den-den-mushi-Go/internal/handler"
	"den-den-mushi-Go/internal/protocol"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
)

// fanout sends packet to Primary and all observers, called by event loop
func (s *Session) fanout(pkt protocol.Packet) {
	if s.Primary != nil {
		s.handlePacket(pkt, s.Primary)
	}
	for o := range s.observers {
		s.handlePacket(pkt, o)
	}
}

func (s *Session) handlePacket(pkt protocol.Packet, c *client.Connection) {
	h, ok := handler.Get[pkt.Header]
	if !ok {
		s.Log.Error("No handler found for packet", zap.String("header", string(pkt.Header)))
		return
	}

	msg, err := h.Handle(pkt, nil, c.Sock, c.Claims)
	if err != nil {
		// todo: refactor once pty decoupled from websocket
		if err == io.EOF {
			s.Log.Info("PTY session ended normally")
		} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			s.Log.Info("WebSocket closed normally") // client sends close, currently not implemented in frontend
		} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			s.Log.Info("WebSocket closed. Probably tab closed") // most closures are this
		} else if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
			s.Log.Warn("WebSocket closed unexpectedly", zap.Error(err))
		} else {
			s.Log.Error("Error handling output packet", zap.Error(err))
		}
		//s.closeAll(ws, pty)
		return
	}

	if msg != "" {
		s.Log.Info("Message from handler", zap.String("message", msg))
		s.logLine(pkt.Header, msg)
	}
}
