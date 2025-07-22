package proxy_host

import "den-den-mushi-Go/pkg/types"

type Model struct {
	HostName    string      `gorm:"column:HostName;primaryKey;size:191"`
	ProxyType   types.Proxy `gorm:"column:ProxyType;size:191"`
	IpAddress   string      `gorm:"column:IpAddress;size:191"`
	OSType      string      `gorm:"column:OS_TYPE;size:191"`
	Status      string      `gorm:"column:Status;size:191"`
	Environment string      `gorm:"column:Environment;size:191"`
	Country     string      `gorm:"column:Country;size:191"`

	LoadBalancerEndpoint string `gorm:"column:LoadBalancerEndpoint;size:191"` // Foreign key
}

func (Model) TableName() string {
	return "ddm_proxy_hosts"
}
