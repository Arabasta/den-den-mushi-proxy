package puppet

import (
	"den-den-mushi-Go/pkg/dto"
	"go.uber.org/zap"
)

const TaskRemovePublicKey PuppetTask = "remove_public_key"

func (p *Client) KeyRemove(publicKey string, conn dto.Connection) error {
	params := taskBody{
		Environment: p.cfg.TaskEnvironment,
		Task:        p.cfg.RemovePublicKeyTaskName,
		Params: sshKeyTaskParams{
			PublicKey:    publicKey,
			OSUser:       conn.Server.OSUser,
			ConnectionID: conn.UserSession.Id,
		},
		Scope: taskScope{
			Nodes: []string{conn.ServerFQDNTmpTillRefactor},
		},
	}

	res, err := p.callPuppetTask(TaskRemovePublicKey, params)
	p.log.Debug("Remove key task called", zap.String("Task name", p.getPuppetTaskName(res)))
	return err
}
