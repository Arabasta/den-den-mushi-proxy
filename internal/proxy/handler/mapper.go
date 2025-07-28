package handler

import (
	"den-den-mushi-Go/internal/proxy/protocol"
)

var (
	inputHandler      Handler = &Input{}
	resizeHandler     Handler = &Resize{}
	sudoInputUsername Handler = &SudoUsername{}
	sudoInputPassword Handler = &SudoPassword{}
)

var Get = map[protocol.Header]Handler{
	protocol.Input:             inputHandler,
	protocol.Resize:            resizeHandler,
	protocol.SudoInputUser:     sudoInputUsername,
	protocol.SudoInputPassword: sudoInputPassword,
}
