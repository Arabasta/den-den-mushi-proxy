package testdata

import (
	"den-den-mushi-Go/pkg/dto/change_request"
	"den-den-mushi-Go/pkg/dto/cyberark"
	"den-den-mushi-Go/pkg/dto/implementor_groups"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/dto/proxy_lb"
	"gorm.io/gorm"
)

func CallAll(db *gorm.DB) {
	//createCR(db)
	//createCyberark(db)
	//createProxyLb(db)
	//createProxyHost(db)
	//createHosts(db)
}

func createCR(db *gorm.DB) {
	// create test crs
	db.Create(&[]change_request.Model{
		{
			// approved CR, valid time, valid implementor groups, valid cyberark objects
			ID:                1,
			CRNumber:          "CR202512314",
			Country:           "SG,HK",
			Lob:               "TechOps",
			Summary:           "System update and patching",
			Description:       "Patching EC2 instances in SG and HK region.",
			ChangeStartTime:   "2025-07-20 23:00:00",
			ChangeEndTime:     "2026-07-22 01:00:00",
			ImplementorGroups: "admin,devops",
			State:             "Approved",
			CyberArkObjects:   "54.255.144.215-ec2user-x123,127.0.0.1-root-w123",
		},
		{
			// unapproved CR
			ID:                2,
			CRNumber:          "CR202512315",
			Country:           "SG,HK",
			Lob:               "TechOps",
			Summary:           "System update and patching",
			Description:       "Patching EC2 instances in SG and HK region.",
			ChangeStartTime:   "2025-07-20 23:00:00",
			ChangeEndTime:     "2026-07-22 01:00:00",
			ImplementorGroups: "admin,devops",
			State:             "Closed",
			CyberArkObjects:   "54.255.144.215-ec2user-x123,127.0.0.1-root-w123",
		},
		{
			// invalid time
			ID:                3,
			CRNumber:          "CR202512316",
			Country:           "SG,HK",
			Lob:               "TechOps",
			Summary:           "System update and patching",
			Description:       "Patching EC2 instances in SG and HK region.",
			ChangeStartTime:   "2024-07-21 23:00:00",
			ChangeEndTime:     "2024-07-22 01:00:00",
			ImplementorGroups: "admin,devops",
			State:             "Approved",
			CyberArkObjects:   "54.255.144.215-ec2user-x123,127.0.0.1-root-w123",
		},
		{
			// invalid implementor groups
			ID:                4,
			CRNumber:          "CR202512317",
			Country:           "SG,HK",
			Lob:               "TechOps",
			Summary:           "System update and patching",
			Description:       "Patching EC2 instances in SG and HK region.",
			ChangeStartTime:   "2025-07-20 23:00:00",
			ChangeEndTime:     "2026-07-22 01:00:00",
			ImplementorGroups: "invalid,ddddevops",
			State:             "Approved",
			CyberArkObjects:   "54.255.144.215-ec2user-x123,127.0.0.1-root-w123",
		},
		{
			// invalid cyberark objects
			ID:                5,
			CRNumber:          "CR202512318",
			Country:           "SG,HK",
			Lob:               "TechOps",
			Summary:           "System update and patching",
			Description:       "Patching EC2 instances in SG and HK region.",
			ChangeStartTime:   "2025-07-20 23:00:00",
			ChangeEndTime:     "2026-07-22 01:00:00",
			ImplementorGroups: "admin,devops",
			State:             "Approved",
			CyberArkObjects:   "123-ec2user-x123,127.0.0.1-root-w123",
		},
		{
			// empty cyberark objects
			ID:                6,
			CRNumber:          "CR202512319",
			Country:           "SG,HK",
			Lob:               "TechOps",
			Summary:           "System update and patching",
			Description:       "Patching EC2 instances in SG and HK region.",
			ChangeStartTime:   "2025-07-20 23:00:00",
			ChangeEndTime:     "2026-07-22 01:00:00",
			ImplementorGroups: "admin,devops",
			State:             "Approved",
			CyberArkObjects:   "",
		},
		{
			// approved CR, valid time, valid implementor groups, valid cyberark objects
			// with cyberark ec2-user
			ID:                7,
			CRNumber:          "CR202512320",
			Country:           "SG,HK",
			Lob:               "TechOps",
			Summary:           "System update and patching",
			Description:       "Patching EC2 instances in SG and HK region.",
			ChangeStartTime:   "2025-07-20 23:00:00",
			ChangeEndTime:     "2026-07-22 01:00:00",
			ImplementorGroups: "admin,devops",
			State:             "Approved",
			CyberArkObjects:   "54.255.144.215-ec2-user-x123,127.0.0.1-root-w123",
		},
		{
			// approved CR, valid time, valid implementor groups, valid cyberark objects
			// with cyberark ec2-user
			ID:                8,
			CRNumber:          "CR202512321",
			Country:           "SG,CN",
			Lob:               "KEI",
			Summary:           "System update and patching",
			Description:       "Patching EC2 instances in SG and HK region.",
			ChangeStartTime:   "2025-07-23 23:00:00",
			ChangeEndTime:     "2026-07-24 01:00:00",
			ImplementorGroups: "kei",
			State:             "Approved",
			CyberArkObjects:   "54.255.144.215-ec2-user-x123,127.0.0.1-rootabc-w123,127.0.0.1-root-w123",
		},
	})

}

func createCyberark(db *gorm.DB) {
	// create test cyberark objects
	db.Create(&[]cyberark.Model{
		//{
		//	ID:       1,
		//	Object:   "54.255.144.215-ec2-user-x123",
		//	Hostname: "aws-ec2-",
		//	Ip:       "54.255.144.215",
		//},
		//{
		//	ID:       2,
		//	Object:   "127.0.1-root-w123",
		//	Hostname: "local-root",
		//	Ip:       "127.0.1",
		//},
		{
			ID:       3,
			Object:   "54.255.144.215-ec2-user-x456",
			Hostname: "aws-ec2-123",
			Ip:       "54.255.144.215",
		},
	})
}

func createProxyHost(db *gorm.DB) {
	db.Create(&[]proxy_host.Model{
		{
			IpAddress:            "127.0.1",
			HostName:             "ddm-proxy",
			ProxyType:            "OS",
			Status:               "Active",
			OSType:               "Linux",
			Environment:          "Production",
			Country:              "SG",
			LoadBalancerEndpoint: "localhost:45007",
		},
	})
}

func createProxyLb(db *gorm.DB) {
	// create test proxy load balancer
	db.Create(&[]proxy_lb.Model{
		{
			LoadBalancerEndpoint: "localhost:45007",
			Type:                 "OS",
			Region:               "SG",
			Environment:          "Production",
		},
	})
}

func createImplementorGroups(db *gorm.DB) {
	db.Create(&[]implementor_groups.Model{
		{
			ID:               1,
			MemberName:       "kei",
			GroupName:        "admin",
			MembershipStatus: "Active",
		},
		{
			ID:               2,
			MemberName:       "kei",
			GroupName:        "devops",
			MembershipStatus: "Active",
		},
		{
			ID:               3,
			MemberName:       "kei2",
			GroupName:        "log",
			MembershipStatus: "Active",
		},
	})
}

//func createHosts(db *gorm.DB) {
//	// create test hosts
//	db.Create(&[]host.Model{
//		{
//			IpAddress:   "127.0.0.1",
//			HostName:    "ddm-proxy",
//			Status:      "Active",
//			OSType:      "Linux",
//			Environment: "Production",
//			Country:     "SG",
//			Appcode:     "ddm",
//		},
//		{
//			IpAddress:   "54.255.144.215",
//			HostName:    "aws-ec2-123",
//			Status:      "Active",
//			OSType:      "Linux",
//			Environment: "Production",
//			Country:     "SG",
//			Appcode:     "UCHIHA",
//		},
//	})
//}
