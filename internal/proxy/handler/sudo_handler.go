package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

type SudoHandler struct{}

func (h *SudoHandler) Handle(pkt protocol.Packet, pty io.Writer, _ *websocket.Conn, claims *token.Claims) (string, error) {

	if true { // todo: this should check user's claims for sudo
		user := pkt.Data
		command := "su " + string(user) + " -i\n"

		_, err := pty.Write([]byte(command))
		return "", err
	}
	return "", nil
}
