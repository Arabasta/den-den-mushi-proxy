package proxy_lb

import "den-den-mushi-Go/pkg/types"

type InMemRepository struct {
	proxies map[types.Proxy]*Entity
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		//  hardcode for now, todo: store in db or config
		proxies: map[types.Proxy]*Entity{
			types.OS: {
				Region:               "Singapore",
				Environment:          "uat",
				LoadBalancerEndpoint: "127.0.1:45007",
				Hostnames:            []string{"proxy1"},
				IPs:                  []string{"127.0.1"},
			},
		},
	}
}

func (r *InMemRepository) FindByProxyType(t types.Proxy) (*Entity, error) {
	return r.proxies[t], nil
}
