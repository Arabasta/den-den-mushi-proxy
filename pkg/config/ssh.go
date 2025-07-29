package config

import (
	"github.com/spf13/viper"
	"time"
)

type Ssh struct {
	IsLocalSshKeyEnabled bool
	LocalSshKeyPath      string

	SshCommand                     string
	PubKeyHostnameSuffix           string
	EphemeralKeyPath               string
	IsRemoveInjectKeyEnabled       bool
	IsRSAKeyPair                   bool
	IsCleanupEnabled               bool
	IsLogPrivateKey                bool
	ConnectDelayAfterInjectSeconds time.Duration
}

func BindSshKey(v *viper.Viper) {
	_ = v.BindEnv("Ssh.LocalSshKeyPath", "SSH_LOCAL_KEY_PATH")
	_ = v.BindEnv("Ssh.EphemeralKeyPath", "SSH_EPHEMERAL_KEY_PATH")
	_ = v.BindEnv("Ssh.PubKeyHostnameSuffix", "SSH_PUB_KEY_HOSTNAME_SUFFIX")
}
