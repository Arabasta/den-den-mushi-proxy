package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/client"
	"den-den-mushi-Go/internal/proxy/filter"
	"den-den-mushi-Go/internal/proxy/protocol"
	"go.uber.org/zap"
	"io"
	"os"
	"sync"
)

type Session struct {
	id      string
	Pty     *os.File
	purpose Purpose

	Log       *zap.Logger
	logWriter io.WriteCloser

	Filter     *filter.CommandFilter // only for health check
	LineEditor *filter.LineEditor    // only for health check, tracks pty's current line

	Primary   *client.Connection
	observers map[*client.Connection]struct{}

	connRegisterCh   chan *client.Connection
	connDeregisterCh chan *client.Connection

	outboundCh chan protocol.Packet

	mu     sync.Mutex
	closed bool
}

func New(id string, pty *os.File, log *zap.Logger) *Session {
	s := &Session{
		id:  id,
		Pty: pty,
		Log: log.With(zap.String("ptySession", id)),

		observers: make(map[*client.Connection]struct{}),

		outboundCh:       make(chan protocol.Packet, 100),
		connRegisterCh:   make(chan *client.Connection),
		connDeregisterCh: make(chan *client.Connection),
	}

	if err := s.initLogWriter(); err != nil {
		s.Log.Error("Failed to create session log", zap.Error(err))
		// close
	}

	s.Log.Info("Created new pty session, initializing event loop and pty reader")

	go s.eventLoop()
	go s.readPty()
	return s
}

// eventLoop should be the only goroutine mutating Primary and observers, other code needs to lock if reading
func (s *Session) eventLoop() {
	for {
		select {
		case pkt, ok := <-s.outboundCh:
			if !ok {
				return
			}
			s.fanout(pkt)
		case c := <-s.connRegisterCh:
			s.addConn(c)
		case c := <-s.connDeregisterCh:
			s.removeConn(c)
		}
	}
}

// readPty and add data to outbound channel
func (s *Session) readPty() {
	buf := make([]byte, 4096)
	for {
		n, err := s.Pty.Read(buf)
		if err != nil {
			if err == io.EOF {
				s.Log.Info("PTY session ended normally")
			} else {
				s.Log.Error("Error reading from pty", zap.Error(err))
			}
			close(s.outboundCh)
			return
		}

		data := append([]byte{}, buf[:n]...)
		s.outboundCh <- protocol.Packet{
			Header: protocol.Output,
			Data:   data,
		}

		s.Log.Info("Pty Output", zap.ByteString("data", data))
	}
}
