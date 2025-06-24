package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

func (s *Service) handlePing(ws *websocket.Conn, missedPongs *int32) {
	ticker := time.NewTicker(s.cfg.Websocket.PingPong.PingIntervalSeconds * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := ws.WriteControl(
				websocket.PingMessage,
				[]byte("ping"),
				time.Now().Add(s.cfg.Websocket.PingPong.PingTimeoutSeconds*time.Second),
			)
			if err != nil {
				newMissed := atomic.AddInt32(missedPongs, 1)
				s.log.Debug("WebSocket ping failed",
					zap.Int32("missedPongs", newMissed),
				)

				if newMissed >= int32(s.cfg.Websocket.PingPong.MaxPingMissed) {
					s.log.Warn("Max missed pings exceeded, closing WebSocket")
					s.closeWs(ws)
					return
				}
			}
			s.log.Debug("Sent ping")
		}
	}
}

func (s *Service) handlePong(ws *websocket.Conn, missedPongs *int32) {
	ws.SetPongHandler(func(appData string) error {
		s.log.Debug("Received pong from client", zap.String("data", appData))
		atomic.StoreInt32(missedPongs, 0) // reset missed counter
		return ws.SetReadDeadline(time.Now().Add(s.cfg.Websocket.PingPong.PongWaitSeconds * time.Second))
	})
}
