package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/pty_helpers"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type SshTestKeyConnection struct {
	cfg            *config.Config
	log            *zap.Logger
	commandBuilder *pty_helpers.Builder
}

func (c *SshTestKeyConnection) Connect(_ context.Context, claims *token.Claims) (*os.File, error) {
	keyPath := c.cfg.Development.SshTestKeyPath
	cmd := c.commandBuilder.BuildSshCmd(keyPath, claims.Connection.Server)
	return pty_helpers.Spawn(cmd)
}
