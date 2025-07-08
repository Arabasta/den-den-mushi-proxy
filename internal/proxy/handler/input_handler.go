package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

type InputHandler struct{}

func (h *InputHandler) Handle(pkt protocol.Packet, pty io.Writer, _ *websocket.Conn, _ *token.Claims) (string, error) {
	_, err := pty.Write(pkt.Data)
	if err != nil {
		return "", err
	}
	return "", err
}
