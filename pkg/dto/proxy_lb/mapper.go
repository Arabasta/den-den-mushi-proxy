package proxy_lb

import "den-den-mushi-Go/pkg/dto/proxy_host"

func ToModel(r Record) *Model {
	return &Model{
		LoadBalancerEndpoint: r.LoadBalancerEndpoint,
		Type:                 r.Type,
		Region:               r.Region,
		Environment:          r.Environment,
		ProxyHosts:           make([]*proxy_host.Model, len(r.ProxyHosts)),
	}
}

func FromModel(m *Model) *Record {
	proxyHosts := make([]string, len(m.ProxyHosts))
	for i, host := range m.ProxyHosts {
		proxyHosts[i] = host.HostName
	}

	return &Record{
		LoadBalancerEndpoint: m.LoadBalancerEndpoint,
		Type:                 m.Type,
		Region:               m.Region,
		Environment:          m.Environment,
		ProxyHosts:           proxyHosts,
	}
}
