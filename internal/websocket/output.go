package websocket

import (
	"den-den-mushi-Go/internal/websocket/handler"
	"den-den-mushi-Go/internal/websocket/protocol"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"os"
)

func (s *Service) handleOutput(ws *websocket.Conn, pty io.ReadWriteCloser, logFile *os.File) {
	buf := make([]byte, 4096)
	outputHandler := &handler.OutputHandler{}

	for {
		n, err := pty.Read(buf)
		if err != nil {
			s.log.Error("Error reading from pty", zap.Error(err))
			s.closeAll(ws, pty)
			return
		}

		_, err = outputHandler.Handle(
			protocol.Packet{
				Header: protocol.Output,
				Data:   buf[:n],
			}, nil, ws, nil)
		if err != nil {
			s.log.Error("Error handling output packet", zap.Error(err))
			s.closeAll(ws, pty)
			return
		}

		// temp for demo, log all output from pty
		s.tempLogFunction(logFile, buf[:n])
	}
}
