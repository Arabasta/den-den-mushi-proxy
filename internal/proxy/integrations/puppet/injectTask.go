package puppet

import (
	"den-den-mushi-Go/pkg/dto"
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

	_, err := p.callPuppetTask(TaskInjectPublicKey, params)
	return err
}
