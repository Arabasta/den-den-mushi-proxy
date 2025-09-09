package lb

import (
	"context"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func (lb *LoadBalancer) startRefreshSnapshot(ctx context.Context, interval time.Duration) {
	if err := lb.refreshSnapshot(ctx); err != nil {
		lb.log.Warn("refresh snapshot failed", zap.Error(err))
	}

	t := time.NewTicker(interval)

	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := lb.refreshSnapshot(ctx); err != nil {
					lb.log.Warn("refresh snapshot failed", zap.Error(err))
				}
			}
		}
	}()
}

func (lb *LoadBalancer) refreshSnapshot(ctx context.Context) error {
	lb.log.Debug("refresh snapshot started")
	raw, err := lb.provider.FindAll(ctx)
	if err != nil {
		lb.log.Error("snapshot fetch failed", zap.Error(err))
		return err
	}

	filtered := lb.applyFilters(raw)
	if len(filtered) == 0 {
		lb.snap.Store([]*proxy_host.Record2{})
		lb.log.Warn("no hosts passed filters")
		return ErrNoHostsPassedFilters
	}

	lb.snap.Store(filtered)
	lb.log.Debug("snapshot updated", zap.Int("eligible_hosts", len(filtered)))
	return nil
}

var ErrNoHostsPassedFilters = fmt.Errorf("no hosts passed filters")
