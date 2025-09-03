package heartbeat

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Scheduler struct {
	db       *gorm.DB
	hostname string
	cfg      *config.Config
	log      *zap.Logger
}

func NewScheduler(db *gorm.DB, cfg *config.Config, log *zap.Logger) *Scheduler {
	if !cfg.Heartbeat.IsEnabled {
		log.Info("Heartbeat is disabled in configuration")
		return nil
	}

	if cfg.Heartbeat.IntervalSeconds <= 0 {
		cfg.Heartbeat.IntervalSeconds = 10
	}

	log.Info("Initializing HeartbeatScheduler", zap.Duration("interval", cfg.Heartbeat.IntervalSeconds*time.Second))

	return &Scheduler{
		db:       db,
		hostname: cfg.Host.Name,
		cfg:      cfg,
		log:      log,
	}
}

func (s *Scheduler) Pulse(ctx context.Context) {
	// send initial heartbeat immediately
	if err := s.updateHeartbeat(); err != nil {
		s.log.Error("Initial heartbeat failed: %v", zap.Error(err))
	}

	ticker := time.NewTicker(s.cfg.Heartbeat.IntervalSeconds * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("stopping pulse")
			return
		case <-ticker.C:
			if err := s.updateHeartbeat(); err != nil {
				s.log.Error("pulse failed", zap.Error(err))
			}
		}
	}
}

func (s *Scheduler) updateHeartbeat() error {
	now := time.Now()

	return s.db.Model(&proxy_host.Model2{}).
		Where("HostName = ?", s.hostname).
		Updates(map[string]interface{}{
			"LastHeartbeatAt": now,
		}).
		Error
}
