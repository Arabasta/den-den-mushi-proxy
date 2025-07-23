package testdata

import (
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/dto/proxy_lb"
	"gorm.io/gorm"
)

func CreateSMXTestData(db *gorm.DB) {
	createProxyLb2(db)
	createProxyHost2(db)
}

func createProxyLb2(db *gorm.DB) {
	// create proxy load balancer
	db.Create(&[]proxy_lb.Model{
		{
			// needs changing
			LoadBalancerEndpoint: "localhost:45007",

			Type:        "OS",
			Region:      "SG",
			Environment: "UAT",
		},
	})
}

func createProxyHost2(db *gorm.DB) {
	db.Create(&[]proxy_host.Model{
		{
			// needs changing
			IpAddress:            "127.0.1",
			HostName:             "ddm-proxy",
			LoadBalancerEndpoint: "localhost:45007", // foreign key to proxy_lb

			ProxyType:   "OS",
			Status:      "Active",
			OSType:      "Linux",
			Environment: "UAT",
			Country:     "SG",
		},
	})
}
