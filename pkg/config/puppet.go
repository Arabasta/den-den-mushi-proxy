package config

import (
	"github.com/spf13/viper"
	"time"
)

type Puppet struct {
	Endpoint                string
	TaskEnvironment         string
	InjectPublicKeyTaskName string
	RemovePublicKeyTaskName string
	Token                   string
	RetryAttempts           int
	TaskRetrySeconds        time.Duration
	TaskNode                string
}

func BindPuppet(v *viper.Viper) {
	_ = v.BindEnv("puppet.token", "PUPPET_TOKEN")
	_ = v.BindEnv("puppet.taskenvironment", "PUPPET_TASK_ENVIRONMENT")
	_ = v.BindEnv("puppet.endpoint", "PUPPET_ENDPOINT")
}
