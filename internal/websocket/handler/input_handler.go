package handler

import (
	"den-den-mushi-Go/internal/websocket/protocol"
	"github.com/gorilla/websocket"
	"io"
)

type InputHandler struct{}

func (h *InputHandler) Handle(pkt protocol.Packet, pty io.Writer, _ *websocket.Conn) error {
	_, err := pty.Write(pkt.Data)
	return err
}
