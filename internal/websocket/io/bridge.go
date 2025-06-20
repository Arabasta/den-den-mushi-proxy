package io

import (
	"github.com/gorilla/websocket"
	"io"
	"sync"
)

func Bridge(ws *websocket.Conn, pty io.ReadWriteCloser, opts ...Option) {
	bridgeConfig := &bridgeConfig{}
	for _, o := range opts {
		o(bridgeConfig)
	}

	closeOnce := sync.Once{}
	closeAll := func() {
		closeOnce.Do(func() {
			_ = ws.Close()
			_ = pty.Close()
		})
	}

	// pty > websocket
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := pty.Read(buf)
			if err != nil {
				closeAll()
				return
			}
			out := applyFilters(buf[:n], bridgeConfig)
			if len(out) == 0 {
				continue
			}
			if err := ws.WriteMessage(websocket.TextMessage, out); err != nil {
				closeAll()
				return
			}
		}
	}()

	// websocket > pty
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			closeAll()
			return
		}
		in := applyFilters(msg, bridgeConfig)
		if len(in) == 0 {
			continue
		}
		if _, err := pty.Write(in); err != nil {
			closeAll()
			return
		}
	}
}
