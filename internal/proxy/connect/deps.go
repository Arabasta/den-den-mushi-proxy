package connect

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/integrations/puppet"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"go.uber.org/zap"
)

type Deps struct {
	puppet         *puppet.Client
	commandBuilder *pty_util.Builder
	cfg            *config.Config
	log            *zap.Logger
}

func NewDeps(puppet *puppet.Client, commandBuilder *pty_util.Builder, cfg *config.Config, log *zap.Logger) Deps {
	return Deps{
		puppet:         puppet,
		commandBuilder: commandBuilder,
		cfg:            cfg,
		log:            log,
	}
}
