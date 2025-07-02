package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/pty_helpers"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type LocalShellConnection struct {
	cfg            *config.Config
	log            *zap.Logger
	commandBuilder *pty_helpers.Builder
}

func (c *LocalShellConnection) Connect(_ context.Context, _ *token.Claims) (*os.File, error) {
	cmd := c.commandBuilder.BuildBashCmd()
	return pty_helpers.Spawn(cmd)
}
