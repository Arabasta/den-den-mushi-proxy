package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/protocol"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var WarningThresholds = []time.Duration{
	30 * time.Minute,
	20 * time.Minute,
	10 * time.Minute,
	5 * time.Minute,
	1 * time.Minute,
}

func (s *Session) monitorCrEndTime() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic", zap.Any("panic", r), zap.Stack("stack"))
		}
	}()

	now := time.Now()

	if s.crEndTime.After(now) {
		timeTillEnd := s.crEndTime.Sub(now)

		s.log.Info("Ticket session timeout scheduled", zap.Duration("timeTillEnd", timeTillEnd))

		for _, threshold := range WarningThresholds {
			if timeTillEnd > threshold {
				// schedule warning relative to CR end time
				delay := timeTillEnd - threshold
				go s.scheduleWarning(delay, threshold)
			}
		}

		// timer will send signal on its channel when the time is up
		timer := time.NewTimer(timeTillEnd)

		select {
		case <-timer.C: // on timer ch fire
			s.log.Info("Ticket end time reached, closing PTY session")

			pkt := protocol.Packet{
				Header: protocol.PtyCRTimeout,
				Data:   []byte(s.Id),
			}

			s.ptyOutput.Add(pkt)
			s.fanout(pkt, nil)
			time.Sleep(time.Millisecond * 3000) // allow pkt to be sent

			s.EndSession()
		case <-s.ctx.Done():
			s.log.Info("Session cancelled before Ticket timeout")
			timer.Stop()
		}
	} else {
		s.log.Warn("Ticket end time already passed, closing PTY session immediately")
		s.EndSession()
	}
}

func (s *Session) scheduleWarning(delay time.Duration, threshold time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("panic", zap.Any("panic", r), zap.Stack("stack"))
		}
	}()

	select {
	case <-time.After(delay):
		s.log.Debug("Sending Ticket timeout warning, session ending in ", zap.Duration("minutes", threshold))

		// send warning packet
		pkt := protocol.Packet{
			Header: protocol.PtyCRTimeoutWarning,
			Data:   []byte(strconv.Itoa(int(threshold.Minutes()))),
		}
		s.ptyOutput.Add(pkt)
		s.fanout(pkt, nil)

	case <-s.ctx.Done():
		return // session closed early
	}
}
