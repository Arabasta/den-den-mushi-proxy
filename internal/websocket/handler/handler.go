package handler

import (
	"den-den-mushi-Go/internal/websocket/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

type Handler interface {
	Handle(pkt protocol.Packet, pty io.Writer, ws *websocket.Conn, claims *token.Claims) (string, error)
}

// todo: preinit once they are stateless
var Get = map[protocol.Header]Handler{
	protocol.Input:  &InputHandler{},
	protocol.Output: &OutputHandler{},
	protocol.Resize: &ResizeHandler{},
	protocol.Sudo:   &SudoHandler{},
}
