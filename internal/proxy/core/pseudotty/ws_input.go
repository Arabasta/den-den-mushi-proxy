package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/protocol"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// readClient should only be accessible by the primary connection
func (s *Session) readClient(c *client.Connection) {
	for {
		msgType, msg, err := c.Sock.ReadMessage()
		if err != nil {
			s.handleReadError(err)
			// close conn on any error
			s.connDeregisterCh <- c
			return
		}

		pkt := protocol.Parse(msgType, msg)
		if pkt.Header == protocol.ParseError {
			s.log.Error("Received invalid message from websocket", zap.Any("msg", msg))
			s.logLine(pkt.Header, string(msg))
			continue
		}

		s.log.Debug("Received packet from client", zap.Any("packet", pkt))
		s.logLine(pkt.Header, string(pkt.Data))

		s.processClientMsg(pkt)
	}
}

func (s *Session) processClientMsg(pkt protocol.Packet) {
	var logMsg string
	var err error

	if pkt.Header == protocol.Input {
		logMsg, err = s.purpose.HandleInput(s, pkt)
	} else {
		logMsg, err = s.purpose.HandleOther(s, pkt)
	}

	if err != nil {
		s.log.Error("Failed to process message", zap.Error(err))
		sendToConn(s.primary, protocol.Packet{
			Header: protocol.Warn,
			Data:   []byte("Failed to process message"),
		})
	}

	if logMsg != "" {
		s.log.Info("Message from handler", zap.String("header", string(pkt.Header)),
			zap.String("message", logMsg))
		s.logLine(pkt.Header, logMsg)
	}
}

func (s *Session) handleReadError(err error) {
	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure):
		s.log.Info("WebSocket closed normally")
	case websocket.IsCloseError(err, websocket.CloseGoingAway):
		s.log.Info("WebSocket closed. Probably tab closed")
	case websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure):
		s.log.Warn("WebSocket closed unexpectedly", zap.Error(err))
	default:
		s.log.Error("Error reading from websocket", zap.Error(err))
	}
}
