package proxy_lb

import "den-den-mushi-Go/pkg/types"

type Entity struct {
	Type                 types.Proxy `json:"type"`
	Region               string      `json:"region"`
	Environment          string      `json:"environment"`
	LoadBalancerEndpoint string      `json:"load_balancer_endpoint"`
	Hostnames            []string    `json:"hostnames"`
	IPs                  []string    `json:"ips"`
}
