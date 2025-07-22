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
		UseInMemoryRepository bool
		SshTestKeyPath        string
	}

	Logger *config.Logger

	Puppet *config.Puppet

	DdmDB *config.SqlDb

	InvDB *config.SqlDb

	Token struct {
		Issuer      string
		Audience    string
		Secret      string
		ExpectedTyp string
		Ttl         int
	}

	Websocket struct {
		PingPong struct {
			PingIntervalSeconds time.Duration
			PingTimeoutSeconds  time.Duration
			PongWaitSeconds     time.Duration
			MaxPingMissed       int
		}
	}
}
