package config

import (
	"den-den-mushi-Go/pkg/config"
)

type Config struct {
	App *config.App

	Development struct {
		UseInMemoryRepository                       bool
		SkipPolicyChecks                            bool
		IsUsingInvDb                                bool
		IsAutoMigrateEnabled                        bool
		IsBlacklistFilter                           bool
		IsLocalSshKeyIfNotIsPuppetKey               bool
		TargetSshPort                               string
		IsSMX                                       bool
		ProxyLoadbalancerEndpointForDiffProxyGroups string
		ProxyHostIpForRejoinRouting                 string
		ProxyHostNameJustForLookup                  string
	}

	Ssl *config.Ssl

	Cors *config.Cors

	Logger *config.Logger

	JwtIssuer *config.JwtIssuer

	DdmDB *config.SqlDb

	InvDB *config.SqlDb

	Redis *config.Redis

	OuGroup struct {
		IsValidationEnabled bool `json:"IsValidationEnabled"`
		Prefix              struct {
			L1   string `json:"L1"`
			L2L3 string `json:"L2_L3"`
		} `json:"Prefix_V1"`
	} `json:"OuGroup"`

	TmpAuth *config.Tmpauth
}
