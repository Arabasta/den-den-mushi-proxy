package client

import (
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Sock   *websocket.Conn
	Claims *token.Claims
	Close  func()
}
