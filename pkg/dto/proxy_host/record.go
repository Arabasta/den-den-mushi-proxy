package proxy_host

import (
	"den-den-mushi-Go/pkg/types"
	"time"
)

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

type Record2 struct {
	HostName  string
	ProxyType types.Proxy

	Url         string
	IpAddress   string
	Environment ProxyEnvironment
	Country     string

	LastHeartbeatAt time.Time

	DrainingStartAt *time.Time
	DrainingEndAt   *time.Time
	DeploymentColor DeploymentColor

	MaxSessions    int
	ActiveSessions int
}
