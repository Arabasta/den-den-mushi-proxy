package session_manager

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// todo refactor, or maybe not

type LoadMonitoringScheduler struct {
	db       *gorm.DB
	hostname string
	cfg      *config.Config
	log      *zap.Logger
	r        LoadReporter
}

func newLoadMonitoringScheduler(r LoadReporter, db *gorm.DB, cfg *config.Config, log *zap.Logger) *LoadMonitoringScheduler {
	if !cfg.LoadMonitoring.IsEnabled {
		log.Info("Load Monitoring is disabled in configuration")
		return nil
	}

	if cfg.LoadMonitoring.BatchUpdateIntervalSeconds <= 0 {
		cfg.LoadMonitoring.BatchUpdateIntervalSeconds = 10
	}

	log.Info("Initializing LoadMonitoringScheduler", zap.Duration("interval", cfg.LoadMonitoring.BatchUpdateIntervalSeconds*time.Second))

	return &LoadMonitoringScheduler{
		db:       db,
		hostname: cfg.Host.Name,
		r:        r,
		cfg:      cfg,
		log:      log,
	}
}

func (s *LoadMonitoringScheduler) load(ctx context.Context) {
	// send initial load immediately
	if err := s.updateLoad(ctx); err != nil {
		s.log.Error("Initial heartbeat failed: %v", zap.Error(err))
	}

	ticker := time.NewTicker(s.cfg.LoadMonitoring.BatchUpdateIntervalSeconds * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.log.Info("stopping load monitoring")
			return
		case <-ticker.C:
			if err := s.updateLoad(ctx); err != nil {
				s.log.Error("load update failed", zap.Error(err))
			}
		}
	}
}

func (s *LoadMonitoringScheduler) updateLoad(ctx context.Context) error {
	// batch update load instead of updating on session create / delete (optimisation, totally not cause lazy)
	load := s.r.howMany()

	return s.db.WithContext(ctx).
		Model(&proxy_host.Model2{}).
		Where("HostName = ?", s.hostname).
		Updates(map[string]any{
			"ActiveSessions": load,
		}).
		Error
}

type LoadReporter interface {
	howMany() int
}

func (s *Service) howMany() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// ez way lol
	return len(s.ptySessions)
}

func (s *Service) StartLoadMonitoring(ctx context.Context, db *gorm.DB) {
	if !s.cfg.LoadMonitoring.IsEnabled {
		s.log.Info("Load monitoring disabled by config")
		return
	}
	s.loadMonitor = newLoadMonitoringScheduler(s, db, s.cfg, s.log)
	go s.loadMonitor.load(ctx)
}
