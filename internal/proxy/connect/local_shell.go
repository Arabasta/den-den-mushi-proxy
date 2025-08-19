package connect

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"go.uber.org/zap"
	"os"
	"os/exec"
)

type LocalShellConnection struct {
	cfg            *config.Config
	log            *zap.Logger
	commandBuilder *pty_util.Builder
}

func (c *LocalShellConnection) Connect(_ context.Context, _ *token.Claims) (*os.File, *exec.Cmd, error) {
	if !c.cfg.Ssh.IsLocalSshKeyEnabled {
		return nil, nil, errors.New("LocalShell is not supported in this build")
	}

	cmd := c.commandBuilder.BuildBashCmd()
	return pty_util.Spawn(cmd, c.log)
}
