package pseudotty

import (
	"den-den-mushi-Go/internal/core/client"
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (s *Session) RegisterInitialConn(ws *websocket.Conn, claims *token.Claims) error {
	s.Log.Info("Registering initial websocket connection to pty session")
	err := s.RegisterConn(ws, claims)

	if claims.Connection.Purpose == dto.Change {
		s.purpose = &ChangeRequestPurpose{}
	} else if claims.Connection.Purpose == dto.Healthcheck {
		s.purpose = &HealthcheckPurpose{}
	} else {
		s.Log.Error("Unknown purpose for new connection", zap.String("purpose", string(claims.Connection.Purpose)))
		return errors.New("unknown purpose")
	}
	if s.purpose == nil {
		s.Log.Error("Purpose is nil")
		return errors.New("purpose is nil")
	}

	return err
}

// RegisterConn to client connection on websocket connect
func (s *Session) RegisterConn(ws *websocket.Conn, claims *token.Claims) error {
	s.Log.Info("Attaching websocket connection to pty session", zap.String("userSessionId",
		claims.Connection.UserSession.Id))

	// check if primary already exists
	if claims.Connection.UserSession.StartRole == dto.Implementor {
		s.mu.Lock()
		if s.Primary != nil {
			s.mu.Unlock()
			return errors.New("max of one primaryConn per pty session allowed")
		}
		s.mu.Unlock()
	}

	conn := &client.Connection{
		Sock:   ws,
		Claims: claims,
	}
	conn.Close = func() {
		s.connDeregisterCh <- conn
	}

	s.connRegisterCh <- conn
	return nil
}

// addConn when a new websocket connection is registered, called by the event loop
func (s *Session) addConn(c *client.Connection) {
	s.Log.Info("Registering websocket connection to pty session",
		zap.String("userSessionId", c.Claims.Connection.UserSession.Id),
		zap.String("role", string(c.Claims.Connection.UserSession.StartRole)))

	if c.Claims.Connection.UserSession.StartRole == dto.Implementor {
		s.Primary = c
	} else if c.Claims.Connection.UserSession.StartRole == dto.Observer {
		s.observers[c] = struct{}{}
	} else {
		s.Log.Error("Unknown role for websocket connection")
		return
	}

	s.Log.Info("Websocket connection registered")

	if c.Claims.Connection.UserSession.StartRole == dto.Implementor {
		s.Log.Info("Is primaryConn role, starting readClient")
		go s.readClient(c)
	}
	if c.Claims.Connection.PtySession.IsNew {
		s.Log.Info("Is new pty, adding log header")
		s.LogHeader()
	}
}

// removeConn when a new websocket connection is deregistered, called by the event loop
func (s *Session) removeConn(c *client.Connection) {
	s.Log.Info("Deregistering websocket connection from pty session",
		zap.String("userSessionId", c.Claims.Connection.UserSession.Id))

	if s.Primary == c {
		s.Primary = nil
	} else {
		delete(s.observers, c)
	}

	// todo: more error handling eg check if ws is arleady closed
	c.Sock.Close()
}
