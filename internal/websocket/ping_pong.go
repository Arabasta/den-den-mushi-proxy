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
				s.log.Error("WebSocket ping failed", zap.Error(err))
				s.closeWs(ws)
				return
			}
			newMissed := atomic.AddInt32(missedPongs, 1)
			s.log.Debug("Sent ping", zap.Int32("missedPongs", newMissed))

			if newMissed >= int32(s.cfg.Websocket.PingPong.MaxPingMissed) {
				s.log.Warn("Max missed pongs exceeded", zap.Int32("missedPongs", newMissed))
				s.closeWs(ws)
				return
			}
		}
	}
}

func (s *Service) handlePong(ws *websocket.Conn, missedPongs *int32) {
	ws.SetPongHandler(func(appData string) error {
		s.log.Debug("Received pong from client", zap.String("data", appData))
		atomic.StoreInt32(missedPongs, 0) // reset on pong

		err := ws.SetReadDeadline(time.Now().Add(s.cfg.Websocket.PingPong.PongWaitSeconds * time.Second))
		if err != nil {
			s.log.Error("Failed to set read deadline after pong", zap.Error(err))
			s.closeWs(ws)
			return err
		}

		return nil
	})

	if err := ws.SetReadDeadline(time.Now().Add(s.cfg.Websocket.PingPong.PongWaitSeconds * time.Second)); err != nil {
		s.log.Error("Failed to set initial read deadline", zap.Error(err))
		s.closeWs(ws)
	}
}
