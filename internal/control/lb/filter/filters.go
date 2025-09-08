package filter

import (
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"time"

	"go.uber.org/zap"
)

func (f *Factory) IsHealthy(unhealthyThreshold time.Duration) Func {
	f.log.Info("Using IsHealthyFilter", zap.Duration("unhealthy_threshold_seconds", unhealthyThreshold))
	return func(h *proxy_host.Record2) bool {
		timeSinceLastSeen := time.Now().Sub(h.LastHeartbeatAt)
		f.log.Debug("IsHealthyFilter", zap.Duration("seconds_since_last_seen", timeSinceLastSeen),
			zap.Duration("unhealthy_threshold_seconds", unhealthyThreshold))
		return timeSinceLastSeen <= unhealthyThreshold
	}
}

func (f *Factory) HasCapacity() Func {
	f.log.Info("Using HasCapacityFilter")
	return func(h *proxy_host.Record2) bool {
		f.log.Debug("HasCapacityFilter", zap.Int("active_sessions", h.ActiveSessions),
			zap.Int("max_sessions", h.MaxSessions))
		return h.ActiveSessions < h.MaxSessions
	}
}

func (f *Factory) NotDraining() Func {
	f.log.Info("Using NotDrainingFilter")
	return func(h *proxy_host.Record2) bool {
		now := time.Now()

		if h.DrainingStartAt == nil {
			return true // not draining
		}

		if h.DrainingEndAt == nil {
			// no end, forever draining if drain start after now
			return h.DrainingStartAt.After(now)
		}

		inDrainWindow := (h.DrainingStartAt.Before(now) || h.DrainingStartAt.Equal(now)) &&
			(now.Before(*h.DrainingEndAt) || now.Equal(*h.DrainingEndAt))
		return !inDrainWindow
	}
}

func (f *Factory) MatchingDeploymentColor(acceptedColor proxy_host.DeploymentColor) Func {
	f.log.Info("Using MatchingDeploymentColorFilter", zap.String("acceptedColor", string(acceptedColor)))
	return func(h *proxy_host.Record2) bool {
		if acceptedColor == proxy_host.AllDeployment {
			return true
		}
		return h.DeploymentColor == acceptedColor
	}
}
