package lb

import (
	"den-den-mushi-Go/internal/control/lb/algo"
	"den-den-mushi-Go/pkg/dto/proxy_host"
)

func (lb *LoadBalancer) GetServer() (*proxy_host.Record2, error) {
	servers, _ := lb.snap.Load().([]*proxy_host.Record2)
	if len(servers) == 0 {
		lb.log.Error("no backend servers available")
		return nil, algo.ErrNoServersAvailable
	}

	lb.mu.RLock()
	a := lb.algo
	lb.mu.RUnlock()
	if a == nil {
		lb.log.Error("load balancer has no algorithm")
		return nil, algo.ErrNoServersAvailable
	}

	return a.GetServer(servers)
}
