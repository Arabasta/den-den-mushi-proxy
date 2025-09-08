package lb

import (
	"den-den-mushi-Go/internal/control/lb/filter"
	"den-den-mushi-Go/pkg/dto/proxy_host"

	"go.uber.org/zap"
)

func (lb *LoadBalancer) applyFilters(in []*proxy_host.Record2) []*proxy_host.Record2 {
	lb.mu.RLock()
	fs := append([]filter.Func(nil), lb.filters...)
	lb.mu.RUnlock()

	if len(fs) == 0 {
		// defensive clone
		out := make([]*proxy_host.Record2, len(in))
		copy(out, in)
		return out
	}

	out := make([]*proxy_host.Record2, 0, len(in))
	for _, h := range in {
		passes := true
		for _, f := range fs {
			if !f(h) {
				passes = false
				break
			}
		}
		if passes {
			out = append(out, h)
		}
	}

	lb.log.Debug("Applied filters",
		zap.Int("input_hosts", len(in)),
		zap.Int("output_hosts", len(out)),
		zap.Int("filter_count", len(fs)),
	)

	return out
}
