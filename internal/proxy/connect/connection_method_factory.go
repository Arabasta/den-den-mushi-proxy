package connect

import (
	"context"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/orchestrator/puppet"
	"den-den-mushi-Go/internal/proxy/pty_util"
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
	"os"
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

type ConnectionMethodFactory struct {
	deps Deps
}

func NewConnectionMethodFactory(d Deps) *ConnectionMethodFactory {
	return &ConnectionMethodFactory{deps: d}
}

type ConnectionMethod interface {
	Connect(ctx context.Context, claims *token.Claims) (*os.File, error)
}

// todo: init once and done cause they are stateless now

// Create returns the correct ConnectionMethod for the requested type
func (f *ConnectionMethodFactory) Create(t types.ConnectionMethod) ConnectionMethod {
	switch t {
	case types.LocalShell:
		return &LocalShellConnection{
			log:            f.deps.log,
			cfg:            f.deps.cfg,
			commandBuilder: f.deps.commandBuilder,
		}

	case types.SshTestKey:
		return &SshTestKeyConnection{
			log:            f.deps.log,
			cfg:            f.deps.cfg,
			commandBuilder: f.deps.commandBuilder,
		}

	case types.SshOrchestratorKey:
		return &SshOrchestratorKeyConnection{
			puppet:         f.deps.puppet,
			cfg:            f.deps.cfg,
			log:            f.deps.log,
			commandBuilder: f.deps.commandBuilder,
		}
	default:
		return nil
	}
}
