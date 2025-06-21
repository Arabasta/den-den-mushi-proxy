package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/pseudoty"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type SshTestKeyConnection struct {
	Cfg *config.Config
	Log *zap.Logger
}

func (c *SshTestKeyConnection) Connect(_ context.Context, claims *token.Claims) (*os.File, error) {
	keyPath := c.Cfg.Development.SshTestKeyPath
	cmd := pseudo.BuildSshCmd(keyPath, claims.Connection)
	return pseudo.Spawn(cmd)
}
