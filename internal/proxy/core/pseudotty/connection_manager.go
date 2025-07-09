package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (s *Session) RegisterInitialConn(ws *websocket.Conn, claims *token.Claims) error {
	s.startClaims = claims

	s.Log.Info("Registering initial websocket connection to pty session")
	err := s.RegisterConn(ws, claims)
	if err != nil {
		s.Log.Error("Failed to register initial connection", zap.Error(err))
		return err
	}

	s.Log.Info("Setting purpose for pty session", zap.String("purpose", string(claims.Connection.Purpose)))
	err = setPurpose(s, claims.Connection.Purpose)
	if err != nil {
		s.Log.Error("Failed to register initial connection", zap.Error(err))
		return err
	}

	if claims.Connection.Purpose == dto.Healthcheck {
		s.Log.Info("Setting healthcheck filter")
		s.filter = filter.GetFilter(claims.Connection.FilterType)
		if s.filter == nil {
			err = errors.New("invalid filter type")
			s.Log.Error("Failed to register initial connection", zap.Error(err))
			return err
		}
	}

	return nil
}

// RegisterConn to client connection on websocket connect
func (s *Session) RegisterConn(ws *websocket.Conn, claims *token.Claims) error {
	s.Log.Info("Attaching websocket connection to pty session", zap.String("userSessionId",
		claims.Connection.UserSession.Id))

	// check if primary already exists
	if claims.Connection.UserSession.StartRole == dto.Implementor {
		s.mu.Lock()
		if s.primary != nil {
			s.mu.Unlock()
			return errors.New("max of one primaryConn per pty session allowed")
		}
		s.mu.Unlock()
	}

	conn := &client.Connection{
		Sock:      ws,
		Claims:    claims,
		WsWriteCh: make(chan protocol.Packet, 100), // todo: make configurable
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
		s.primary = c
	} else if c.Claims.Connection.UserSession.StartRole == dto.Observer {
		s.observers[c] = struct{}{}
	} else {
		s.Log.Error("Unknown role for websocket connection")
		return
	}

	s.Log.Info("Websocket connection registered")

	// doesn't handle swapping roles for now
	if c.Claims.Connection.UserSession.StartRole == dto.Implementor {
		s.Log.Info("Is primary role, starting readClient")
		go s.readClient(c)
	}
	if c.Claims.Connection.PtySession.IsNew {
		s.Log.Info("Is new pty, adding log header")
		s.LogHeader()
	} else {
		// is joining existing pty session, notify everyone
		if c.Claims.Connection.UserSession.StartRole == dto.Observer {
			s.logf("[%s] %s joined as observer", c.Claims.Subject, c.Claims.Connection.UserSession.Id)
			pkt := protocol.Packet{Header: protocol.PtySessionEvent, Data: []byte(c.Claims.Subject + " joined as observer")}
			s.outboundCh <- pkt
		} else if c.Claims.Connection.UserSession.StartRole == dto.Implementor {
			s.logf("[%s] %s joined as implementor", c.Claims.Subject, c.Claims.Connection.UserSession.Id)
			pkt := protocol.Packet{Header: protocol.PtySessionEvent, Data: []byte(c.Claims.Subject + " joined as implementor")}
			s.outboundCh <- pkt
		}

		for i := range s.ptyLastPackets {
			sendToConn(c, s.ptyLastPackets[i])
		}
	}

	// start sending messages to the client
	go c.WriteClient()
}

// removeConn when a new websocket connection is deregistered, called by the event loop
func (s *Session) removeConn(c *client.Connection) {
	s.Log.Info("Deregistering websocket connection from pty session",
		zap.String("userSessionId", c.Claims.Connection.UserSession.Id))

	if s.primary == c {
		s.primary = nil
	} else {
		delete(s.observers, c)
	}

	if s.primary != nil && s.observers != nil {
		s.logf("[%s] %s has left the session", c.Claims.Subject, c.Claims.Connection.UserSession.Id)
		pkt := protocol.Packet{Header: protocol.PtySessionEvent, Data: []byte(c.Claims.Subject + " has left")}
		s.outboundCh <- pkt
	}
}
