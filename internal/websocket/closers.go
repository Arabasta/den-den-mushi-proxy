package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"sync"
)

func (s *Service) closeWs(ws *websocket.Conn) {
	if err := ws.Close(); err != nil {
		s.log.Warn("Failed to close websocket", zap.Error(err))
	} else {
		s.log.Info("Closed websocket")
	}
}

func (s *Service) closePty(pty io.Closer) {
	err := pty.Close()
	if err != nil {
		if err == io.EOF {
			s.log.Info("PTY session ended normally")
		} else {
			s.log.Warn("Failed to close pty", zap.Error(err))
		}
	}
	s.log.Info("Closed pty")

}

func (s *Service) closeAll(ws *websocket.Conn, pty io.ReadWriteCloser) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			s.closeWs(ws)
			s.closePty(pty)
		})
	}
}
