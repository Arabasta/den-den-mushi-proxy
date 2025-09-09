package connect

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

type SshTectiaKeyConnection struct {
	cfg            *config.Config
	log            *zap.Logger
	commandBuilder *pty_util.Builder
}

func (c *SshTectiaKeyConnection) Connect(_ context.Context, claims *token.Claims) (*os.File, *exec.Cmd, error) {
	// todo: implement Tectia key handling
	return nil, nil, errors.New("not implemented")
}
