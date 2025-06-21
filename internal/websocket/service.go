package websocket

import (
	"context"
	"den-den-mushi-Go/internal/pseudo/connect"
	"den-den-mushi-Go/internal/websocket/io"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"os"
)

type Service struct {
	ConnectionMethodFactory *connect.ConnectionMethodFactory
}

func NewWebsocketService(c *connect.ConnectionMethodFactory) *Service {
	return &Service{
		ConnectionMethodFactory: c,
	}
}

func (s *Service) Run(ctx context.Context, ws *websocket.Conn, claims *token.Claims) {
	conn := s.ConnectionMethodFactory.Create(claims.Connection.Type)
	if conn == nil {
		err := ws.Close()
		if err != nil {
			// todo: handle error
		}
		return
	}

	pty, err := conn.Connect(ctx, claims)
	if err != nil {
		err := ws.Close()
		if err != nil {
			// todo: handle error
		}
		return
	}

	defer func(pty *os.File) {
		err := pty.Close()
		if err != nil {

		}
	}(pty)

	io.Bridge(ws, pty)
}
