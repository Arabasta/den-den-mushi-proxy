package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/pseudoty"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type LocalShellConnection struct {
	Cfg *config.Config
	Log *zap.Logger
}

func (c *LocalShellConnection) Connect(_ context.Context, _ *token.Claims) (*os.File, error) {
	cmd := pseudo.BuildBashCmd()
	return pseudo.Spawn(cmd)
}
