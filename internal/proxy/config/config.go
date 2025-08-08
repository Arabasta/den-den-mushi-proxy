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
		UseRedis             bool
		IsUsingInvDb         bool
		IsAutoMigrateEnabled bool
		IsSMX                bool
	}

	Pty struct {
		SpawnRetryCount           int
		SpawnRetryIntervalSeconds time.Duration
	}

	Logger *config.Logger

	Puppet *config.Puppet

	DdmDB *config.SqlDb

	InvDB *config.SqlDb

	Redis *config.Redis

	JwtAudience *config.JwtAudience

	Ssh *config.Ssh

	Websocket struct {
		PingPong struct {
			PingIntervalSeconds time.Duration
			PingTimeoutSeconds  time.Duration
			PongWaitSeconds     time.Duration
			MaxPingMissed       int
		}
	}

	PuppetTasks struct {
		QueryJobs struct {
			OrchestratorEndpoint   string        `json:"OrchestratorEndpoint"`
			WaitBeforeQuerySeconds time.Duration `json:"WaitBeforeQuerySeconds"`
			MaxQueryAttempts       int           `json:"MaxQueryAttempts"`
			QueryIntervalSeconds   time.Duration `json:"QueryIntervalSeconds"`
		} `json:"QueryJobs"`
		CyberarkPasswordDraw struct {
			Environment         string `json:"Environment"`
			TaskName            string `json:"TaskName"`
			CybidA              string `json:"CybidA"`
			CybidB              string `json:"CybidB"`
			SafeA               string `json:"SafeA"`
			SafeB               string `json:"SafeB"`
			IsValidationEnabled bool
		}
	}

	TmpAuth *config.Tmpauth

	Filters struct {
		DbPollIntervalSeconds time.Duration

		IsHealthcheckBlacklistEnabled bool
		IsHealthcheckWhitelistEnabled bool
		IsChangeBlacklistEnabled      bool
	}
}
