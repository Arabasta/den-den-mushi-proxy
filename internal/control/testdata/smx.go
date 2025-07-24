package testdata

import (
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/dto/proxy_lb"
	"gorm.io/gorm"
)

func CreateProxyHostAndLb(db *gorm.DB, cfg *config.Config) {
	createProxyLb2(db, cfg)
	createProxyHost2(db, cfg)
}

func createProxyLb2(db *gorm.DB, cfg *config.Config) {
	// create proxy load balancer
	db.Create(&[]proxy_lb.Model{
		{
			// needs changing
			LoadBalancerEndpoint: cfg.Development.ProxyLoadbalancerEndpointForDiffProxyGroups,

			Type:        "OS",
			Region:      "SG",
			Environment: "UAT",
		},
	})
}

func createProxyHost2(db *gorm.DB, cfg *config.Config) {
	db.Create(&[]proxy_host.Model{
		{
			// needs changing
			IpAddress:            cfg.Development.ProxyHostIpForRejoinRouting,
			HostName:             cfg.Development.ProxyHostNameJustForLookup,
			LoadBalancerEndpoint: "localhost:45007", // foreign key to proxy_lb

			ProxyType:   "OS",
			Status:      "Active",
			OSType:      "Linux",
			Environment: "UAT",
			Country:     "SG",
		},
	})
}
