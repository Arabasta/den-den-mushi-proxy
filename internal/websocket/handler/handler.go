package handler

import (
	"den-den-mushi-Go/internal/websocket/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

type Handler interface {
	Handle(packet protocol.Packet, pty io.Writer, ws *websocket.Conn, claims *token.Claims) error
}

var Get = map[protocol.Header]Handler{
	protocol.Input:  &InputHandler{},
	protocol.Output: &OutputHandler{},
	protocol.Resize: &ResizeHandler{},
	protocol.Sudo:   &SudoHandler{},
}
