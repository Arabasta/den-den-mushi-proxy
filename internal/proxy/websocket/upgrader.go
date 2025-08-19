package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
)

//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//	CheckOrigin:     func(r *http.Request) bool { return true }, // todo:
//	Subprotocols:    []string{"X-Proxy-Session-Token"},
//}

func (s *Service) websocketUpgrader() websocket.Upgrader {
	allowed := map[string]struct{}{}
	for _, o := range s.cfg.Websocket.AllowedOrigins {
		u, err := url.Parse(o)
		if err != nil {
			continue
		}
		key := strings.ToLower(u.Scheme + "://" + u.Host) // includes port if present
		allowed[key] = struct{}{}
	}

	s.log.Info("Setting up websocket upgrader",
		zap.Int("readBufferSize", s.cfg.Websocket.ReadBufferSize),
		zap.Int("writeBufferSize", s.cfg.Websocket.WriteBufferSize),
		zap.String("subprotocols", s.cfg.Websocket.Subprotocols),
		zap.Any("allowedOrigins", s.cfg.Websocket.AllowedOrigins),
		zap.Any("allowedOriginsMap", allowed))

	return websocket.Upgrader{
		ReadBufferSize:  s.cfg.Websocket.ReadBufferSize,
		WriteBufferSize: s.cfg.Websocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			if len(s.cfg.Websocket.AllowedOrigins) > 0 && s.cfg.Websocket.AllowedOrigins[0] == "*" {
				return true // allow all origins
			}

			origin := r.Header.Get("Origin")
			if origin == "" || origin == "null" {
				// must have origin header
				return false
			}
			u, err := url.Parse(origin)
			if err != nil {
				s.log.Error("Failed to parse origin URL", zap.Error(err), zap.String("origin", origin))
				return false
			}
			_, ok := allowed[strings.ToLower(u.Scheme+"://"+u.Host)]
			if !ok {
				s.log.Warn("Origin not allowed", zap.String("origin", origin), zap.String("request_url", r.URL.String()))
				return false
			}
			return true
		},
		Subprotocols: []string{s.cfg.Websocket.Subprotocols},
	}
}
