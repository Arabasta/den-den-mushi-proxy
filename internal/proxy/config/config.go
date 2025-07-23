package config

import (
	"den-den-mushi-Go/pkg/config"
	"time"
)

type Config struct {
	App *config.App

	Ssl *config.Ssl

	Cors *config.Cors

	Host *config.Host

	Development struct {
		UseSqlJtiRepo        bool
		SshTestKeyPath       string
		UseRedis             bool
		IsUsingInvDb         bool
		IsAutoMigrateEnabled bool
	}

	Logger *config.Logger

	Puppet *config.Puppet

	DdmDB *config.SqlDb

	InvDB *config.SqlDb

	Redis *config.Redis

	JwtAudience *config.JwtAudience

	Websocket struct {
		PingPong struct {
			PingIntervalSeconds time.Duration
			PingTimeoutSeconds  time.Duration
			PongWaitSeconds     time.Duration
			MaxPingMissed       int
		}
	}
}
