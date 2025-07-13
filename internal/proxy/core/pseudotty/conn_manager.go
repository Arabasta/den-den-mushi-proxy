package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
)

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
		s.mu.Unlock()
		// todo: return err
		return
	}
	s.mu.Unlock()

	c.Log = s.log.With(zap.String("userSessionId", c.Claims.Connection.UserSession.Id))

	s.log.Info("Websocket connection registered")

	if c.Claims.Connection.PtySession.IsNew {
		s.log.Info("Is new pty, adding log header")
		s.logL(getLogHeader(s))
	} else {
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

	c.DoClose()

	pkt := protocol.Packet{Header: protocol.PtySessionEvent, Data: []byte(c.Claims.Subject + " has left")}
	s.logPacket(pkt)
	s.fanout(pkt, c)
}
