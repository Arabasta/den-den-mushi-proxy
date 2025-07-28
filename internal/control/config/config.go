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
		TargetSshPort                               string
		IsSMX                                       bool
		ProxyLoadbalancerEndpointForDiffProxyGroups string
		ProxyHostIpForRejoinRouting                 string
		ProxyHostNameJustForLookup                  string
		HealthcheckOsUsers                          []string
		IsLocalSshKeyIfNotIsPuppetKey               bool
		IsTmpAuthCookieEnabled                      bool
	}

	CookieTmp struct {
		Name      string
		Redirect  string
		UserIdKey string
		Secret    string
	}

	Ssl *config.Ssl

	Cors *config.Cors

	Logger *config.Logger

	JwtIssuer *config.JwtIssuer

	DdmDB *config.SqlDb

	InvDB *config.SqlDb

	Redis *config.Redis
}
