package connect

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/integrations/puppet"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
	"time"
)

type SshOrchestratorKeyConnection struct {
	puppet         *puppet.Client
	cfg            *config.Config
	log            *zap.Logger
	commandBuilder *pty_util.Builder
}

// todo refactor garbage
func (c *SshOrchestratorKeyConnection) Connect(_ context.Context, claims *token.Claims) (*os.File, error) {
	keyPath, pubKey, cleanup, err := pty_util.GenerateEphemeralKey(c.cfg, c.log)
	if err != nil {
		c.log.Error("failed to generate ephemeral key", zap.Error(err))
		return nil, err
	}

	err = c.puppet.KeyInject(pubKey, claims.Connection)
	// always assume success for now cause cant get any actual success response unless
	// we query the task id
	//if err != nil {
	//	c.log.Error("Error injecting key", zap.Error(err))
	//	cleanup()
	//	if c.cfg.Ssh.IsRemoveInjectKeyEnabled {
	//		_ = c.puppet.KeyRemove(pubKey, claims.Connection)
	//	}
	//	return nil, err
	//}

	c.log.Debug("Puppet Key inject called. Waiting to spawn pseudo terminal for ephemeral SSH key")

	time.Sleep(c.cfg.Ssh.ConnectDelayAfterInjectSeconds * time.Second) // wait for key to be injected

	var pty *os.File
	for i := 0; i <= c.cfg.Pty.SpawnRetryCount; i++ {
		c.log.Debug("Spawning pseudo terminal for SSH connection", zap.Int("attempt", i+1))
		cmd := c.commandBuilder.BuildSshCmd(keyPath, claims.Connection.Server)

		pty, err = pty_util.Spawn(cmd)
		if err == nil {
			c.log.Info("Pseudo terminal spawned successfully", zap.String("keyPath", keyPath))
			break
		}

		time.Sleep(c.cfg.Pty.SpawnRetryIntervalSeconds * time.Second)
	}
	if err != nil {
		cleanup()
		if c.cfg.Ssh.IsRemoveInjectKeyEnabled {
			_ = c.puppet.KeyRemove(pubKey, claims.Connection)
		}
		return nil, err
	}

	go func() {
		time.Sleep(60 * time.Second)

		cleanup()
		if c.cfg.Ssh.IsRemoveInjectKeyEnabled {
			if err := c.puppet.KeyRemove(pubKey, claims.Connection); err != nil {
				c.log.Error("Failed to remove remote key", zap.Error(err))
			} else {
				c.log.Info("Ephemeral SSH key removed from server")
			}
		}
	}()

	return pty, nil
}
