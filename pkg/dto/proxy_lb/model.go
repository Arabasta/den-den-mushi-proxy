package proxy_lb

import (
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/types"
)

type Model struct {
	LoadBalancerEndpoint string              `gorm:"primaryKey;column:load_balancer_endpoint;size:191"`
	Type                 types.Proxy         `gorm:"column:type;size:191"`
	Region               string              `gorm:"column:region;size:191"`
	Environment          string              `gorm:"column:environment;size:191"`
	ProxyHosts           []*proxy_host.Model `gorm:"foreignKey:LoadBalancerEndpoint;references:LoadBalancerEndpoint"`
}

func (Model) TableName() string {
	return "ddm_proxy_load_balancers"
}
