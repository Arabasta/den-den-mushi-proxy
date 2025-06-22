package handler

import (
	"den-den-mushi-Go/internal/websocket/protocol"
	"github.com/gorilla/websocket"
	"io"
)

type Handler interface {
	Handle(packet protocol.Packet, pty io.Writer, ws *websocket.Conn) error
}

var Get = map[protocol.Header]Handler{
	protocol.Input:  &InputHandler{},
	protocol.Resize: &ResizeHandler{},
}
