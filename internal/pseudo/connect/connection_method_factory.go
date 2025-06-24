package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/orchestrator/puppet"
	"den-den-mushi-Go/internal/pseudo/command"
	"den-den-mushi-Go/pkg/dto/connection"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type Deps struct {
	puppet         *puppet.PuppetClient
	commandBuilder *command.Builder
	cfg            *config.Config
	log            *zap.Logger
}

func NewDeps(puppet *puppet.PuppetClient, commandBuilder *command.Builder, cfg *config.Config, log *zap.Logger) Deps {
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

// Create returns the correct ConnectionMethod for the requested type
func (f *ConnectionMethodFactory) Create(t connection.ConnectionType) ConnectionMethod {
	switch t {
	case connection.LocalShell:
		return &LocalShellConnection{
			log:            f.deps.log,
			cfg:            f.deps.cfg,
			commandBuilder: f.deps.commandBuilder,
		}

	case connection.SshTestKey:
		return &SshTestKeyConnection{
			log:            f.deps.log,
			cfg:            f.deps.cfg,
			commandBuilder: f.deps.commandBuilder,
		}

	case connection.SshOrchestratorKey:
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
