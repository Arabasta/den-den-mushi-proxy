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

func (c *SshOrchestratorKeyConnection) Connect(_ context.Context, claims *token.Claims) (*os.File, error) {
	keyPath, pubKey, cleanup, err := pty_util.GenerateEphemeralKey(c.cfg, c.log)
	if err != nil {
		return nil, err
	}

	if err := c.puppet.KeyInject(pubKey, claims.Connection); err != nil {
		c.log.Error("Error injecting key", zap.Error(err))
		cleanup()
		_ = c.puppet.KeyRemove(pubKey, claims.Connection)
		return nil, err
	}

	c.log.Debug("Puppet Key inject successful. Spawning pseudo terminal for ephemeral SSH key")

	cmd := c.commandBuilder.BuildSshCmd(keyPath, claims.Connection.Server)
	pty, err := pty_util.Spawn(cmd)
	if err != nil {
		c.log.Error("Failed to spawn pseudo terminal", zap.Error(err))
		cleanup()
		_ = c.puppet.KeyRemove(pubKey, claims.Connection)
		return nil, err
	}

	go func() {
		time.Sleep(60 * time.Second)

		cleanup()
		if err := c.puppet.KeyRemove(pubKey, claims.Connection); err != nil {
			c.log.Error("Failed to remove remote key", zap.Error(err))
		} else {
			c.log.Info("Ephemeral SSH key removed from server")
		}
	}()

	return pty, nil
}
