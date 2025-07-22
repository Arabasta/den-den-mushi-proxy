package proxy_host

import "den-den-mushi-Go/pkg/types"

type Record struct {
	IpAddress            string
	ProxyType            types.Proxy
	HostName             string
	OSType               string
	Status               string
	Environment          string
	Country              string
	LoadBalancerEndpoint string
}
