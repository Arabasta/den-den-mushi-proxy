package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (s *Session) RegisterInitialConn(ws *websocket.Conn, claims *token.Claims) error {
	s.startClaims = claims

	s.log.Info("Registering initial websocket connection to pty session")
	err := s.RegisterConn(ws, claims)
	if err != nil {
		s.log.Error("Failed to register initial connection", zap.Error(err))
		return err
	}

	s.log.Info("Setting purpose for pty session", zap.String("purpose", string(claims.Connection.Purpose)))
	err = setPurpose(s, claims.Connection.Purpose)
	if err != nil {
		s.log.Error("Failed to register initial connection", zap.Error(err))
		return err
	}

	if claims.Connection.Purpose == types.Healthcheck {
		s.log.Info("Setting healthcheck filter")
		s.filter = filter.GetFilter(claims.Connection.FilterType)
		if s.filter == nil {
			err = errors.New("invalid filter type")
			s.log.Error("Failed to register initial connection", zap.Error(err))
			return err
		}
	}

	return nil
}

// RegisterConn to client connection on websocket connect
func (s *Session) RegisterConn(ws *websocket.Conn, claims *token.Claims) error {
	s.log.Info("Attaching websocket connection to pty session", zap.String("userSessionId",
		claims.Connection.UserSession.Id))

	if claims.Connection.UserSession.StartRole == types.Implementor {
		// check if primary already exists
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
	s.log.Info("Registering websocket connection to pty session",
		zap.String("userSessionId", c.Claims.Connection.UserSession.Id),
		zap.String("role", string(c.Claims.Connection.UserSession.StartRole)))

	s.mu.Lock()
	if c.Claims.Connection.UserSession.StartRole == types.Implementor {
		s.primary = c
	} else if c.Claims.Connection.UserSession.StartRole == types.Observer {
		s.observers[c] = struct{}{}
	} else {
		s.log.Error("Unknown role for websocket connection")
		// todo: return err
		return
	}
	s.mu.Unlock()

	s.log.Info("Websocket connection registered")

	if c.Claims.Connection.PtySession.IsNew {
		s.log.Info("Is new pty, adding log header")
		s.logL(getLogHeader(s))
	} else {
		// is joining existing pty session
		ptyLastPackets := s.ptyOutput.GetAll()
		for i := range ptyLastPackets {
			client.SendToConn(c, ptyLastPackets[i])
		}

		pkt := protocol.Packet{Header: protocol.PtySessionEvent,
			Data: []byte(c.Claims.Subject + " joined as " + string(c.Claims.Connection.UserSession.StartRole))}
		s.logPacket(pkt)
		s.fanoutExcept(pkt, c)
	}

	// start sending messages to the client
	go c.WriteClient(s.log)

	// doesn't handle swapping roles for now
	if c.Claims.Connection.UserSession.StartRole == types.Implementor {
		s.log.Info("Is implementor role, starting readClient")
		go s.readClient(c)
	} else {
		s.log.Info("Is observer role, starting ObserverReadLoop")
		go c.ObserverReadLoop()
	}
}

// removeConn when a new websocket connection is deregistered, called by the event loop
func (s *Session) removeConn(c *client.Connection) {
	s.log.Info("Deregistering websocket connection from pty session",
		zap.String("userSessionId", c.Claims.Connection.UserSession.Id))

	s.mu.Lock()
	if s.primary == c {
		s.primary = nil
	} else {
		delete(s.observers, c)
	}
	s.mu.Unlock()

	s.closeWs(c)

	pkt := protocol.Packet{Header: protocol.PtySessionEvent, Data: []byte(c.Claims.Subject + " has left")}
	s.logPacket(pkt)
	s.fanoutExcept(pkt, c)
}
