package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

// Handler should not have any websocket writes
type Handler interface {
	Handle(pkt protocol.Packet, pty io.Writer, ws *websocket.Conn, claims *token.Claims) (string, error)
}

var (
	inputHandler  Handler = &InputHandler{}
	resizeHandler Handler = &ResizeHandler{}
	sudoHandler   Handler = &SudoHandler{}
)

var Get = map[protocol.Header]Handler{
	protocol.Input:  inputHandler,
	protocol.Resize: resizeHandler,
	protocol.Sudo:   sudoHandler,
}
