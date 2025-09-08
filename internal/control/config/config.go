package config

import (
	"den-den-mushi-Go/internal/control/lb/algo"
	"den-den-mushi-Go/pkg/config"
	"time"
)

type Config struct {
	App *config.App

	Development struct {
		UseInMemoryRepository                       bool
		SkipPolicyChecks                            bool
		IsUsingInvDb                                bool
		IsAutoMigrateEnabled                        bool
		IsLocalSshKeyIfNotIsPuppetKey               bool
		TargetSshPort                               string
		IsSMX                                       bool
		ProxyLoadbalancerEndpointForDiffProxyGroups string
		ProxyHostIpForRejoinRouting                 string
		ProxyHostNameJustForLookup                  string
	}

	Ssl *config.Ssl

	Cors *config.Cors

	Security *config.Security

	Pprof *config.Pprof

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

	Swagger *config.Swagger

	LoadBalancer struct {
		Algorithm              algo.Type
		RefreshIntervalSeconds int `json:"RefreshIntervalSeconds"`
		Filters                struct {
			CapacityFilter struct {
				IsEnabled bool `json:"IsEnabled"`
			} `json:"CapacityFilter"`
			HealthFilter struct {
				IsEnabled                 bool          `json:"IsEnabled"`
				UnhealthyThresholdSeconds time.Duration `json:"UnhealthyThresholdSeconds"`
			} `json:"HealthFilter"`
			DeploymentColorFilter struct {
				IsEnabled     bool   `json:"IsEnabled"`
				AcceptedColor string `json:"AcceptedColor"`
			} `json:"DeploymentColorFilter"`
			DrainFilter struct {
				IsEnabled bool `json:"IsEnabled"`
			} `json:"DrainFilter"`
		} `json:"Filters"`
	}
}
