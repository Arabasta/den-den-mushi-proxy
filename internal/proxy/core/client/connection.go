package client

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync"
)

type Connection struct {
	Id     string
	Sock   *websocket.Conn
	Claims *token.Claims

	WsWriteCh        chan protocol.Packet
	OnceCloseWriteCh sync.Once

	Close   func()
	CloseDb func()
	Ctx     context.Context
	Cancel  context.CancelFunc

	Log *zap.Logger
	cfg *config.Config
}

func New(sock *websocket.Conn, claims *token.Claims, cfg *config.Config) *Connection {
	return &Connection{
		Id:        claims.Connection.UserSession.Id,
		Sock:      sock,
		Claims:    claims,
		WsWriteCh: make(chan protocol.Packet, 100), // todo: make configurable
		cfg:       cfg,
	}
}

func (c *Connection) DoClose() {
	c.OnceCloseWriteCh.Do(func() {
		if c.Cancel != nil {
			c.Cancel()
		}
		c.Log.Debug("DoClose called, closing connection", zap.String("userSessionId", c.Claims.Connection.UserSession.Id))

		if c.Sock != nil {
			c.Log.Debug("Closing websocket connection")
			err := c.Sock.Close()
			if err != nil {
				c.Log.Error("Failed to close websocket connection", zap.Error(err))
			}
		}
	})
}
