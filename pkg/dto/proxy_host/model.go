package proxy_host

import (
	"den-den-mushi-Go/pkg/types"
	"time"
)

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

type Model2 struct {
	HostName  string      `gorm:"column:HostName;primaryKey;size:191"`
	ProxyType types.Proxy `gorm:"column:ProxyType;size:191"`

	Url         string           `gorm:"column:Url;size:191;not null;unique"`
	IpAddress   string           `gorm:"column:IpAddress;size:191"`
	Environment ProxyEnvironment `gorm:"column:Environment;size:191;not null;default:'development'"`
	Country     string           `gorm:"column:Country;size:191;not null;default:'unknown'"`

	LastHeartbeatAt time.Time `gorm:"column:LastHeartbeatAt;type:timestamp"`

	DrainingStartAt *time.Time      `gorm:"column:DrainingStartAt;type:timestamp;nullable;default:NULL"`
	DrainingEndAt   *time.Time      `gorm:"column:DrainingEndAt;type:timestamp;nullable;default:NULL"`
	DeploymentColor DeploymentColor `gorm:"column:DeploymentColor;size:191;not null;default:'blue'"`

	MaxSessions    int `gorm:"column:MaxSessions;default:100"`
	ActiveSessions int `gorm:"column:ActiveSessions;default:0"`
}

func (Model2) TableName() string {
	return "ddm_proxies"
}

type ProxyEnvironment string

const (
	Development ProxyEnvironment = "development"
	UAT         ProxyEnvironment = "uat"
	Production  ProxyEnvironment = "production"
)

type DeploymentColor string

const (
	// AllDeployment only used for control to indicate that all deployments are being considered
	AllDeployment DeploymentColor = "all"

	BlueDeployment   DeploymentColor = "blue"
	GreenDeployment  DeploymentColor = "green"
	RedDeployment    DeploymentColor = "red"
	YellowDeployment DeploymentColor = "yellow"
)
