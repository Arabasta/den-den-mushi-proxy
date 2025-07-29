package websocket

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/connect"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/session_manager"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Service struct {
	connectionMethodFactory *connect.ConnectionMethodFactory
	sessionManager          *session_manager.Service
	log                     *zap.Logger
	cfg                     *config.Config
}

func NewService(c *connect.ConnectionMethodFactory, sm *session_manager.Service, log *zap.Logger,
	cfg *config.Config) *Service {
	log.Info("Initializing WebSocket Service...")
	return &Service{
		connectionMethodFactory: c,
		sessionManager:          sm,
		log:                     log,
		cfg:                     cfg,
	}
}

// initial connection flow for websocket connections
// todo: handle ws close with pty close and return error to client
// todo ctx propagation to all conns
func (s *Service) run(ctx context.Context, ws *websocket.Conn, claims *token.Claims) {
	conn := client.New(ws, claims, s.cfg)

	if claims.Connection.PtySession.IsNew {
		connMethod := s.connectionMethodFactory.Create(claims.Connection.Type)
		if connMethod == nil {
			s.log.Error("Unsupported connection type", zap.String("type", string(claims.Connection.Type)))
			s.closeWs(ws)
			// todo: handle close properly
			return
		}

		s.log.Debug("Establishing pty connection", zap.String("type", string(claims.Connection.Type)))
		pty, err := connMethod.Connect(ctx, claims)
		if err != nil {
			s.log.Error("Failed to connect to pseudo terminal", zap.Error(err),
				zap.String("type", string(claims.Connection.Type)))
			s.closeWs(ws)
			// todo: close both if fail
			return
		}

		s.log.Debug("Connected to pty ", zap.String("type", string(claims.Connection.Type)))

		ptySessionId, err := s.sessionManager.CreatePtySession(pty, claims, s.log)
		if err != nil {
			s.log.Error("Failed to create pty session", zap.Error(err))
			s.closeWs(ws)
			return
		}

		s.log.Debug("Registering websocket connection to pty session")
		err = s.sessionManager.AttachConn(conn, ptySessionId)
		if err != nil {
			s.log.Error("Failed to attach websocket connection to pty session", zap.Error(err))
			s.closeWs(ws)
			// todo: close pty if fail
			return
		}
		return
	} else {
		// join existing session
		err := s.sessionManager.AttachConn(conn, claims.Connection.PtySession.Id)
		if err != nil {
			s.log.Error("Failed to attach websocket connection to existing pty session", zap.Error(err))
			s.closeWs(ws)
		}
		return
	}
}

// todo: use the close from pty session and delete this
func (s *Service) closeWs(ws *websocket.Conn) {
	if err := ws.Close(); err != nil {
		s.log.Warn("Failed to close websocket", zap.Error(err))
	} else {
		s.log.Debug("Closed websocket")
	}
}
