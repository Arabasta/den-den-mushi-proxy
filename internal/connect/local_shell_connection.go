package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/pty_util"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type LocalShellConnection struct {
	cfg            *config.Config
	log            *zap.Logger
	commandBuilder *pty_util.Builder
}

func (c *LocalShellConnection) Connect(_ context.Context, _ *token.Claims) (*os.File, error) {
	cmd := c.commandBuilder.BuildBashCmd()
	return pty_util.Spawn(cmd)
}
