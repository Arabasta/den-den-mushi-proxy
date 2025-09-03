package pseudotty

import (
	"den-den-mushi-Go/internal/proxy/core/activity"
	"den-den-mushi-Go/internal/proxy/core/core_helpers"
	"den-den-mushi-Go/internal/proxy/protocol"
	"den-den-mushi-Go/pkg/types"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func (s *Session) initActivityTracker() {
	if s.startClaims.Connection.Purpose != types.Healthcheck {
		s.log.Info("Not a healthcheck session, not creating activity tracker")
	}

	if s.activityTracker != nil {
		s.log.Warn("Activity tracker already exists")
		return
	}

	if !s.cfg.Healthcheck.IsInactiveTimeoutEnabled {
		s.log.Info("Inactive timeout is disabled, not creating activity tracker")
		return
	}

	if s.cfg.Healthcheck.InactiveTimeoutSeconds <= 0 {
		s.log.Info("Inactive timeout is not set, not creating activity tracker")
		return
	}
	s.log.Info("Initializing activity tracker", zap.Int("timeoutSeconds", int(s.cfg.Healthcheck.InactiveTimeoutSeconds)))

	onTimeout := func() {
		s.log.Info("Session inactive for too long, ending...")
		if s.ActivePrimary != nil {
			core_helpers.SendToConn(s.ActivePrimary, protocol.Packet{
				Header: protocol.InactiveTimeout,
				Data:   []byte(s.Id),
			})
		}
		time.Sleep(2000 * time.Millisecond)
		s.ForceEndSession("Read-only session inactive timeout")
	}

	onWarning := func(remainingSecs int) {
		s.log.Debug("Session inactive warning", zap.Int("remainingSeconds", remainingSecs))
		if s.ActivePrimary != nil {
			core_helpers.SendToConn(s.ActivePrimary, protocol.Packet{
				Header: protocol.InactiveWarning,
				Data:   []byte(strconv.Itoa(remainingSecs)),
			})
		}
	}

	s.activityTracker = activity.NewTracker(s.ctx, s.cfg.Healthcheck.InactiveTimeoutSeconds,
		onTimeout, onWarning)
	s.log.Info("Created activity tracker for healthcheck session")
}

func (s *Session) touchActivity() {
	if s.activityTracker != nil {
		s.activityTracker.Touch()
	}
}
