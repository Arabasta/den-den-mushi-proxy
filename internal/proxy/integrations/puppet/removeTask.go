package puppet

import (
	"den-den-mushi-Go/pkg/dto"
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
			Nodes: []string{p.cfg.TaskNode},
		},
	}

	_, err := p.callPuppetTask(TaskRemovePublicKey, params)
	return err
}
