package config

import "time"

type Config struct {
	App struct {
		Name        string
		Environment string
		Version     string
		Port        int
	}

	Development struct {
		SshTestKeyPath string
	}

	Logging struct {
		Level    string
		Format   string
		Output   string
		FilePath string
	}

	Puppet struct {
		Endpoints          []string
		TaskTimeoutSeconds int
		RetryAttempts      int
	}

	Token struct {
		Secret string
		Expiry int
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
