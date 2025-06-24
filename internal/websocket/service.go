package websocket

import (
	"context"
	"den-den-mushi-Go/internal/pseudo/connect"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"os"
)

type Service struct {
	ConnectionMethodFactory *connect.ConnectionMethodFactory
	Log                     *zap.Logger
}

func NewWebsocketService(c *connect.ConnectionMethodFactory, log *zap.Logger) *Service {
	return &Service{
		ConnectionMethodFactory: c,
		Log:                     log,
	}
}

func (s *Service) run(ctx context.Context, ws *websocket.Conn, claims *token.Claims) {
	conn := s.ConnectionMethodFactory.Create(claims.Connection.Type)
	if conn == nil {
		s.Log.Error("Unsupported connection type", zap.String("type", string(claims.Connection.Type)))
		err := ws.Close()
		if err != nil {
			s.Log.Error("Failed to close websocket connection", zap.Error(err))
		}
		return
	}

	s.Log.Info("Connected to websocket", zap.String("type", string(claims.Connection.Type)))
	pty, err := conn.Connect(ctx, claims)
	if err != nil {
		s.Log.Error("Failed to connect to pseudo terminal", zap.Error(err),
			zap.String("type", string(claims.Connection.Type)))
		err := ws.Close()
		if err != nil {
			s.Log.Error("Failed to close websocket connection", zap.Error(err))
		}
		return
	}

	s.Log.Info("Connected to pseudo terminal", zap.String("type", string(claims.Connection.Type)))

	defer func(pty *os.File) {
		err := pty.Close()
		if err != nil {
			s.Log.Error("Failed to close pseudo terminal", zap.Error(err))
		}
	}(pty)

	bridge(ws, pty)
}
