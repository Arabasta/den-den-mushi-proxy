package connect

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/orchestrator/puppet"
	"den-den-mushi-Go/pkg/connection"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type Deps struct {
	Puppet *puppet.PuppetClient
	Cfg    *config.Config
	Log    *zap.Logger
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

// Build returns the correct ConnectionMethod for the requested type
func (f *ConnectionMethodFactory) Create(t connection.ConnectionType) ConnectionMethod {
	switch t {
	case connection.LocalShell:
		return &LocalShellConnection{
			Log: f.deps.Log,
			Cfg: f.deps.Cfg,
		}

	case connection.SshTestKey:
		return &SshTestKeyConnection{
			Log: f.deps.Log,
			Cfg: f.deps.Cfg,
		}

	case connection.SshOrchestratorKey:
		return &SshOrchestratorKeyConnection{
			Puppet: f.deps.Puppet,
			Cfg:    f.deps.Cfg,
			Log:    f.deps.Log,
		}

	default:
		return nil
	}
}
