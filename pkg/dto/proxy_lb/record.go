package proxy_lb

import "den-den-mushi-Go/pkg/types"

type Record struct {
	LoadBalancerEndpoint string
	Type                 types.Proxy
	Region               string
	Environment          string
	ProxyHosts           []string
}
