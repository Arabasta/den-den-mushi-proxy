package websocket

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/connect"
	"den-den-mushi-Go/internal/proxy/core/session_manager"
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Service struct {
	connectionMethodFactory *connect.ConnectionMethodFactory
	sessionManager          *session_manager.SessionManager
	log                     *zap.Logger
	cfg                     *config.Config
}

func NewWebsocketService(c *connect.ConnectionMethodFactory, sm *session_manager.SessionManager, log *zap.Logger,
	cfg *config.Config) *Service {
	return &Service{
		connectionMethodFactory: c,
		sessionManager:          sm,
		log:                     log,
		cfg:                     cfg,
	}
}

// initial connection flow for websocket connections
// todo: handle ws close with pty close
func (s *Service) run(ctx context.Context, ws *websocket.Conn, claims *token.Claims) {

	if claims.Connection.PtySession.IsNew {
		conn := s.connectionMethodFactory.Create(claims.Connection.Type)
		if conn == nil {
			s.log.Error("Unsupported connection type", zap.String("type", string(claims.Connection.Type)))
			s.closeWs(ws)
			// todo: handle close properly
			return
		}

		s.log.Info("Connected to websocket, establishing pty connection", zap.String("type", string(claims.Connection.Type)))

		pty, err := conn.Connect(ctx, claims)
		if err != nil {
			s.log.Error("Failed to connect to pseudo terminal", zap.Error(err),
				zap.String("type", string(claims.Connection.Type)))
			s.closeWs(ws)
			// todo: close both if fail
			return
		}
		s.log.Info("Connected to pty, creating pty session", zap.String("type", string(claims.Connection.Type)))

		session := s.sessionManager.CreatePtySession(pty, s.log)

		s.log.Info("Registering websocket connection to pty session")
		err = session.RegisterInitialConn(ws, claims)
		if err != nil {
			s.closeWs(ws)
			// todo: close both if fail
			return
		}
		return
	} else {
		// join existing session
		err := s.sessionManager.AttachWebsocket(ws, claims)
		if err != nil {
			s.closeWs(ws)
		}
	}
}

// todo: use the close from pty session and delete this
func (s *Service) closeWs(ws *websocket.Conn) {
	if err := ws.Close(); err != nil {
		s.log.Warn("Failed to close websocket", zap.Error(err))
	} else {
		s.log.Info("Closed websocket")
	}
}
