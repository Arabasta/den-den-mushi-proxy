package websocket

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/pseudo/connect"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"os"
)

type Service struct {
	connectionMethodFactory *connect.ConnectionMethodFactory
	log                     *zap.Logger
	cfg                     *config.Config
}

func NewWebsocketService(c *connect.ConnectionMethodFactory, log *zap.Logger, cfg *config.Config) *Service {
	return &Service{
		connectionMethodFactory: c,
		log:                     log,
		cfg:                     cfg,
	}
}

func (s *Service) run(ctx context.Context, ws *websocket.Conn, claims *token.Claims) {
	conn := s.connectionMethodFactory.Create(claims.Connection.Type)
	if conn == nil {
		s.log.Error("Unsupported connection type", zap.String("type", string(claims.Connection.Type)))
		s.closeWs(ws)
		return
	}
	s.log.Info("Connected to websocket, establishing pty connection", zap.String("type", string(claims.Connection.Type)))

	pty, err := conn.Connect(ctx, claims)
	if err != nil {
		s.log.Error("Failed to connect to pseudo terminal", zap.Error(err),
			zap.String("type", string(claims.Connection.Type)))
		s.closeWs(ws)
		return
	}
	s.log.Info("Connected to pseudo terminal, bridging connection", zap.String("type", string(claims.Connection.Type)))

	defer func(pty *os.File) {
		s.closePty(pty)
	}(pty)

	s.bridge(ws, pty, claims)
}
