package pseudotty

import (
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func (s *Session) connLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic", zap.Any("panic", r), zap.Stack("stack"))
		}
	}()

	for {
		select {
		case c := <-s.connRegisterCh:
			s.addConn(c)
		case c := <-s.connDeregisterCh:
			s.removeConn(c)
		case <-s.ctx.Done():
			s.log.Info("connLoop: context done")
			return
		}
	}
}
