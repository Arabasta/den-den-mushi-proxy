package handler

import (
	"den-den-mushi-Go/internal/websocket/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

type InputHandler struct{}

func (h *InputHandler) Handle(pkt protocol.Packet, pty io.Writer, _ *websocket.Conn, _ *token.Claims) error {
	_, err := pty.Write(pkt.Data)
	return err
}
