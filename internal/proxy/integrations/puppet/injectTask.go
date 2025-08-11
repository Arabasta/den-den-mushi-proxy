package puppet

import (
	"den-den-mushi-Go/pkg/dto"
	"go.uber.org/zap"
)

const TaskInjectPublicKey PuppetTask = "inject_public_key"

func (p *Client) KeyInject(publicKey string, conn dto.Connection) error {
	params := taskBody{
		Environment: p.cfg.TaskEnvironment,
		Task:        p.cfg.InjectPublicKeyTaskName,
		Params: sshKeyTaskParams{
			PublicKey:    publicKey,
			OSUser:       conn.Server.OSUser,
			ConnectionID: conn.UserSession.Id,
		},
		Scope: taskScope{
			Nodes: []string{conn.ServerFQDNTmpTillRefactor},
		},
	}

	p.log.Debug("Preparing to inject public key", zap.Any("params", params))

	res, err := p.callPuppetTask(TaskInjectPublicKey, params)
	p.log.Debug("Inject key task called", zap.String("Task name", p.getPuppetTaskName(res)))
	return err
}
