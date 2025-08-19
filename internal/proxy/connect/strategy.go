package connect

import (
	"den-den-mushi-Go/pkg/types"
	"errors"
)

type Strategy struct {
	m map[types.ConnectionMethod]ConnectionMethod
}

func NewRegistry(deps Deps) *Strategy {
	return &Strategy{
		m: map[types.ConnectionMethod]ConnectionMethod{
			types.LocalShell: &LocalShellConnection{
				log: deps.log, cfg: deps.cfg, commandBuilder: deps.commandBuilder,
			},
			types.LocalSshKey: &LocalSshKeyConnection{
				log: deps.log, cfg: deps.cfg.Ssh, commandBuilder: deps.commandBuilder,
			},
			types.SshOrchestratorKey: &SshOrchestratorKeyConnection{
				puppet: deps.puppet, cfg: deps.cfg, log: deps.log, commandBuilder: deps.commandBuilder,
			},
		},
	}
}

var ErrUnsupportedMethod = errors.New("unsupported connection method")

func (r *Strategy) Get(t types.ConnectionMethod) (ConnectionMethod, error) {
	s, ok := r.m[t]
	if !ok {
		return nil, ErrUnsupportedMethod
	}
	return s, nil
}
