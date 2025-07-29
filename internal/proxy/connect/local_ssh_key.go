package connect

import (
	"context"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"den-den-mushi-Go/pkg/config"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"go.uber.org/zap"
	"os"
)

type LocalSshKeyConnection struct {
	cfg            *config.Ssh
	log            *zap.Logger
	commandBuilder *pty_util.Builder
}

func (c *LocalSshKeyConnection) Connect(_ context.Context, claims *token.Claims) (*os.File, error) {
	if !c.cfg.IsLocalSshKeyEnabled {
		return nil, errors.New("LocalSshKey is not supported in this build")
	}

	cmd := c.commandBuilder.BuildSshCmd(c.cfg.LocalSshKeyPath, claims.Connection.Server)
	return pty_util.Spawn(cmd)
}
