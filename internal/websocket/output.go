package websocket

import (
	"den-den-mushi-Go/internal/websocket/handler"
	"den-den-mushi-Go/internal/websocket/protocol"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
)

func (s *Service) handleOutput(ws *websocket.Conn, pty io.ReadWriteCloser) {
	buf := make([]byte, 4096)
	outputHandler := &handler.OutputHandler{}

	for {
		n, err := pty.Read(buf)
		if err != nil {
			s.log.Error("Error reading from pty", zap.Error(err))
			s.closeAll(ws, pty)
			return
		}

		if err := outputHandler.Handle(
			protocol.Packet{
				Header: protocol.Output,
				Data:   buf[:n],
			},
			nil, ws, nil); err != nil {
			s.log.Error("Error handling output packet", zap.Error(err))
			s.closeAll(ws, pty)
			return
		}
	}
}
