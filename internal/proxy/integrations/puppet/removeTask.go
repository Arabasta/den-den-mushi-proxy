package puppet

import (
	"den-den-mushi-Go/pkg/dto"
)

const TaskRemovePublicKey puppetTask = "remove_public_key"

func (p *Client) KeyRemove(publicKey string, conn dto.Connection) error {
	params := taskBody{
		Environment: p.cfg.TaskEnvironment,
		Task:        TaskRemovePublicKey,
		Params: sshKeyTaskParams{
			PublicKey:    publicKey,
			OSUser:       conn.Server.OSUser,
			ConnectionID: conn.UserSession.Id,
		},
		Scope: taskScope{
			Nodes: []string{conn.Server.IP},
		},
	}

	_, err := p.callPuppetTask(TaskRemovePublicKey, params)
	return err
}
