package websocket

import (
	"den-den-mushi-Go/pkg/token"
	"github.com/gorilla/websocket"
	"io"
)

func (s *Service) bridge(ws *websocket.Conn, pty io.ReadWriteCloser, claims *token.Claims) {
	s.closeAll(ws, pty)

	var missedPongs int32

	go s.handlePing(ws, &missedPongs)
	s.handlePong(ws, &missedPongs)

	go s.handleOutput(ws, pty)
	s.handleInput(ws, pty, claims)
}
