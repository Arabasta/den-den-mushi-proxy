package config

import "github.com/spf13/viper"

type SshKey struct {
	IsLocalSshKeyEnabled bool
	LocalSshKeyPath      string
}

func BindSshKey(v *viper.Viper) {
	_ = v.BindEnv("SshKey.LocalSshKeyPath", "SSHKEY_LOCAL_KEY_PATH")
}
