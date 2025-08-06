package pseudotty

import (
	"context"
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/core/pseudotty/session_logging"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
)

// RegisterConn client connection to pty session
func (s *Session) RegisterConn(c *client.Connection, onClose func(string)) error {
	s.log.Info("Attaching connection to pty session", zap.String("userSessionId",
		c.Claims.Connection.UserSession.Id))

	if s.State == types.Closed {
		return errors.New("pty session is closed, cannot register connection")
	}

	c.Ctx, c.Cancel = context.WithCancel(s.ctx)
	c.Close = func() {
		s.connDeregisterCh <- c
		if onClose != nil {
			onClose(c.Claims.Connection.UserSession.Id)
		}
	}

	s.connRegisterCh <- c
	return nil
}

func (s *Session) DeregisterConn(c *client.Connection) {
	s.log.Info("Deregistering connection from pty session", zap.String("userSessionId",
		c.Claims.Connection.UserSession.Id))
	c.Close()
}

// addConn when a new websocket connection is registered, called by the event loop
func (s *Session) addConn(c *client.Connection) {
	s.log.Info("Registering websocket connection to pty session",
		zap.String("userSessionId", c.Claims.Connection.UserSession.Id),
		zap.String("role", string(c.Claims.Connection.UserSession.StartRole)))

	c.Log = s.log.With(zap.String("userSessionId", c.Claims.Connection.UserSession.Id))

	s.log.Info("Websocket connection registered")

	if c.Claims.Connection.PtySession.IsNew {
		s.log.Info("Is new pty, adding log header")
		s.logL(session_logging.FormatHeader(s.ActivePrimary.Claims))

		core_helpers.SendToConn(c, protocol.Packet{
			Header: protocol.PtyConnectionSuccess,
			Data:   []byte(s.Id),
		})
	} else {
		core_helpers.SendToConn(c, protocol.Packet{
			Header: protocol.PtyConnectionSuccess,
			Data:   []byte(s.Id),
		})

		// is joining existing pty session
		ptyLastPackets := s.ptyOutput.GetAll()
		for i := range ptyLastPackets {
			core_helpers.SendToConn(c, ptyLastPackets[i])
		}

		pkt := protocol.Packet{Header: protocol.PtySessionEvent,
			Data: []byte(c.Claims.Subject + " joined as " + string(c.Claims.Connection.UserSession.StartRole))}
		s.logPacket(pkt)
		s.fanout(pkt, c)
	}

	// start sending messages to client
	go c.WriteClient()

	// doesn't handle swapping roles for now
	if c.Claims.Connection.UserSession.StartRole == types.Implementor {
		s.log.Info("Is implementor role, starting primaryReadLoop")
		go c.PrimaryReadLoop(s.handleConnPacket)
	} else {
		s.log.Info("Is observer role, starting ObserverReadLoop")
		go c.ObserverReadLoop()
	}
}

func (s *Session) AssignRole(c *client.Connection) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lifetimeConnections[c] = struct{}{}

	if c.Claims.Connection.UserSession.StartRole == types.Implementor {
		// check if primary already exists
		if s.ActivePrimary != nil {
			s.log.Warn("Primary connection already exists, cannot register another primary",
				zap.String("userSessionId", c.Claims.Connection.UserSession.Id))
			return errors.New("max of one primaryConn per pty session allowed")
		}

		s.ActivePrimary = c
	} else if c.Claims.Connection.UserSession.StartRole == types.Observer {
		if _, exists := s.ActiveObservers[c]; exists {
			return errors.New("already registered as observer")
		}
		s.ActiveObservers[c] = struct{}{}
	} else {
		return errors.New("unknown role for websocket connection")
	}

	return nil
}

// removeConn when a new websocket connection is deregistered, called by the event loop
func (s *Session) removeConn(c *client.Connection) {
	s.log.Info("Deregistering websocket connection from pty session",
		zap.String("userSessionId", c.Claims.Connection.UserSession.Id))

	s.mu.Lock()
	if s.ActivePrimary == c {
		s.ActivePrimary = nil
	} else {
		delete(s.ActiveObservers, c)
	}
	s.mu.Unlock()

	c.DoClose()

	pkt := protocol.Packet{Header: protocol.PtySessionEvent, Data: []byte(c.Claims.Subject + " has left")}
	s.logPacket(pkt)
	s.fanout(pkt, c)
}
