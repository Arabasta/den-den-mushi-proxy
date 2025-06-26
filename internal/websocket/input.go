package websocket

import (
	"den-den-mushi-Go/internal/websocket/handler"
	"den-den-mushi-Go/internal/websocket/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"os"
)

func (s *Service) handleInput(ws *websocket.Conn, pty io.ReadWriteCloser, claims *token.Claims, logFile *os.File) {
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			// todo: refactor once pty decoupled from websocket
			if err == io.EOF {
				s.log.Info("PTY session ended normally")

			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				s.log.Info("WebSocket closed normally") // client sends close, currently not implemented in frontend
			} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				s.log.Info("WebSocket closed. Probably tab closed") // most closures are this
			} else if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				s.log.Warn("WebSocket closed unexpectedly", zap.Error(err))
			} else {
				s.log.Error("Error reading message from websocket", zap.Error(err))
			}
			s.closeAll(ws, pty)
			return
		}

		if msgType != websocket.BinaryMessage || len(msg) == 0 {
			continue
		}

		packet := protocol.Parse(msg)
		if packet.Header == protocol.ParseError {
			s.log.Error("Received invalid message from websocket", zap.Any("packet", packet))
			s.closeAll(ws, pty)
			return
		}

		// todo: validate header

		h, exists := handler.Get[packet.Header]
		if !exists {
			s.log.Error("No handler found for packet header", zap.Any("header", packet.Header))
			s.closeAll(ws, pty)
			return
		}

		logMsg, err := h.Handle(packet, pty, ws, claims)
		if packet.Header == protocol.Input {
			if err != nil {
				s.log.Error("Error handling input packet", zap.Error(err), zap.String("message", logMsg))
				s.closeAll(ws, pty)
				return
			}
			if logMsg != "" && logFile != nil {
				logFile.WriteString(logMsg)
			}
		}
	}
}
