package lb

import (
	"context"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/core/proxy_hosts"
	"den-den-mushi-Go/internal/control/lb/algo"
	"den-den-mushi-Go/internal/control/lb/filter"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type LoadBalancer struct {
	algo algo.Rithm

	filters       []filter.Func
	filterFactory *filter.Factory

	provider *proxy_hosts.Service
	snap     atomic.Value

	mu sync.RWMutex

	log *zap.Logger
	cfg *config.Config
}

func NewLoadBalancer(provider *proxy_hosts.Service, cfg *config.Config, log *zap.Logger) *LoadBalancer {
	log.Info("Initializing LoadBalancer")

	lb := &LoadBalancer{
		provider:      provider,
		log:           log,
		cfg:           cfg,
		filterFactory: filter.NewFactory(cfg, log),
	}

	lb.setAlgorithm(cfg.LoadBalancer.Algorithm)
	lb.setFilters()

	lb.snap.Store([]*proxy_host.Record2(nil))
	lb.startRefreshSnapshot(context.TODO(), time.Duration(cfg.LoadBalancer.RefreshIntervalSeconds)*time.Second)

	return lb
}

func (lb *LoadBalancer) setAlgorithm(a algo.Type) {
	lb.mu.Lock()
	lb.algo = algo.Rithms[a]
	lb.mu.Unlock()
	lb.log.Info("Set LoadBalancer algorithm", zap.String("type", lb.algo.String()))
}

func (lb *LoadBalancer) setFilters() {
	filters := lb.filterFactory.BuildFilters()

	lb.mu.Lock()
	lb.filters = filters
	lb.mu.Unlock()

	lb.log.Info("Refreshed filters",
		zap.Int("filter_count", len(filters)),
		zap.Bool("health_enabled", lb.cfg.LoadBalancer.Filters.HealthFilter.IsEnabled),
		zap.Bool("capacity_enabled", lb.cfg.LoadBalancer.Filters.CapacityFilter.IsEnabled),
		zap.Bool("color_enabled", lb.cfg.LoadBalancer.Filters.DeploymentColorFilter.IsEnabled),
		zap.Bool("drain_enabled", lb.cfg.LoadBalancer.Filters.DrainFilter.IsEnabled),
	)
}
