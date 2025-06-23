package websocket

import (
	"den-den-mushi-Go/internal/websocket/handler"
	"den-den-mushi-Go/internal/websocket/protocol"
	"github.com/gorilla/websocket"
	"io"
	"sync"
)

func Bridge(ws *websocket.Conn, pty io.ReadWriteCloser) {
	closeOnce := sync.Once{}
	closeAll := func() {
		closeOnce.Do(func() {
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
				closeAll()
				return
			}

			if err := outputHandler.Handle(
				protocol.Packet{
					Header: protocol.Output,
					Data:   buf[:n],
				},
				nil, ws); err != nil {
				closeAll()
				return
			}
		}
	}()

	// websocket > pty
	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			// todo: handle error
			closeAll()
			return
		}

		if msgType != websocket.BinaryMessage || len(msg) == 0 {
			continue
		}

		packet := protocol.Parse(msg)
		if packet.Header == protocol.ParseError {
			continue
		}

		h, exists := handler.Get[packet.Header]
		if !exists {
			continue
		}

		_ = h.Handle(packet, pty, ws)
	}
}
