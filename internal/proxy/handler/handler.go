package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"io"
)

// Handler should not have any websocket writes
type Handler interface {
	Handle(pkt protocol.Packet, pty io.Writer) (string, error)
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
