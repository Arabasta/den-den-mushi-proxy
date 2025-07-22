package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
)

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
