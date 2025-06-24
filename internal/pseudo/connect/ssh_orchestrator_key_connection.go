package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/orchestrator/puppet"
	"den-den-mushi-Go/internal/pseudo"
	"den-den-mushi-Go/internal/pseudo/command"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type SshOrchestratorKeyConnection struct {
	puppet         *puppet.PuppetClient
	cfg            *config.Config
	log            *zap.Logger
	commandBuilder *command.Builder
}

func (c *SshOrchestratorKeyConnection) Connect(ctx context.Context, claims *token.Claims) (*os.File, error) {
	keyPath, pubKey, cleanup, err := pseudo.GenerateEphemeralKey(c.cfg, c.log)
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		cleanup()
	}()

	if err := c.puppet.PuppetKeyInject(pubKey, claims.Connection); err != nil {
		return nil, err
	}

	cmd := c.commandBuilder.BuildSshCmd(keyPath, claims.Connection)
	pty, err := pseudo.Spawn(cmd)
	if err != nil {
		c.log.Error("Failed to spawn pseudo terminal", zap.Error(err))
		return nil, err
	}

	// todo: add key removal

	return pty, err
}
