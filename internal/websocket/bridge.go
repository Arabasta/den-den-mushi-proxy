package websocket

import (
	"den-den-mushi-Go/internal/websocket/handler"
	"den-den-mushi-Go/internal/websocket/protocol"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"sync"
)

func (s *Service) bridge(ws *websocket.Conn, pty io.ReadWriteCloser) {
	closeOnce := sync.Once{}
	closeAll := func() {
		closeOnce.Do(func() {
			s.log.Info("Closing websocket and pty connections")
			_ = ws.Close()
			_ = pty.Close()
		})
	}

	// pty > websocket
	go func() {
		buf := make([]byte, 4096)
		outputHandler := &handler.OutputHandler{}

		for {
			n, err := pty.Read(buf)
			if err != nil {
				s.log.Error("Error reading from pty", zap.Error(err))
				closeAll()
				return
			}

			if err := outputHandler.Handle(
				protocol.Packet{
					Header: protocol.Output,
					Data:   buf[:n],
				},
				nil, ws); err != nil {
				s.log.Error("Error handling output packet", zap.Error(err))
				closeAll()
				return
			}
		}
	}()

	// websocket > pty
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			s.log.Error("Error reading message from websocket", zap.Error(err))
			closeAll()
			return
		}

		if msgType != websocket.BinaryMessage || len(msg) == 0 {
			continue
		}

		packet := protocol.Parse(msg)
		if packet.Header == protocol.ParseError {
			s.log.Error("Received invalid message from websocket", zap.Any("packet", packet))
			closeAll()
			return
		}

		h, exists := handler.Get[packet.Header]
		if !exists {
			s.log.Error("No handler found for packet header", zap.Any("header", packet.Header))
			closeAll()
			return
		}

		_ = h.Handle(packet, pty, ws)
	}
}
