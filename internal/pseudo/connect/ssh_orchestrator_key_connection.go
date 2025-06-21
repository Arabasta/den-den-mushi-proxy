package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/orchestrator/puppet"
	"den-den-mushi-Go/internal/pseudo"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type SshOrchestratorKeyConnection struct {
	Puppet *puppet.PuppetClient
	Cfg    *config.Config
	Log    *zap.Logger
}

func (c *SshOrchestratorKeyConnection) Connect(ctx context.Context, claims *token.Claims) (*os.File, error) {
	keyPath, pubKey, cleanup, err := pseudo.GenerateEphemeralKey()
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		cleanup()
	}()

	if err := c.Puppet.PuppetKeyInject(pubKey, claims.Connection); err != nil {
		return nil, err
	}

	cmd := pseudo.BuildSshCmd(keyPath, claims.Connection)
	pty, err := pseudo.Spawn(cmd)
	if err != nil {
		// handle error
	} else {
		// todo: add key removal
	}

	return pty, err
}
