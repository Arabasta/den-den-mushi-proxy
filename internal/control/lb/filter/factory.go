package filter

import (
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"time"

	"go.uber.org/zap"
)

type Factory struct {
	cfg *config.Config
	log *zap.Logger
}

func NewFactory(cfg *config.Config, log *zap.Logger) *Factory {
	return &Factory{
		cfg: cfg,
		log: log,
	}
}

func (f *Factory) BuildFilters() []Func {
	f.log.Info("Building filters for load balancer based on configuration")
	var filters []Func

	filterCfg := f.cfg.LoadBalancer.Filters

	if filterCfg.HealthFilter.IsEnabled {
		filters = append(filters, f.IsHealthy(filterCfg.HealthFilter.UnhealthyThresholdSeconds*time.Second))
	}

	if filterCfg.CapacityFilter.IsEnabled {
		filters = append(filters, f.HasCapacity())
	}

	if filterCfg.DeploymentColorFilter.IsEnabled {
		filters = append(filters, f.MatchingDeploymentColor(proxy_host.DeploymentColor(filterCfg.DeploymentColorFilter.AcceptedColor)))
	}

	if filterCfg.DrainFilter.IsEnabled {
		filters = append(filters, f.NotDraining())
	}

	return filters
}
