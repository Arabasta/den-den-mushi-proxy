package puppet

import (
	"den-den-mushi-Go/pkg/dto"
)

const TaskInjectPublicKey puppetTask = "inject_public_key"

func (p *Client) KeyInject(publicKey string, conn dto.Connection) error {
	params := taskBody{
		Environment: p.cfg.TaskEnvironment,
		Task:        TaskInjectPublicKey,
		Params: sshKeyTaskParams{
			PublicKey:   publicKey,
			ServerIP:    conn.Server.IP,
			OSUser:      conn.Server.OSUser,
			ConnPurpose: string(conn.Purpose),
			ConnType:    string(conn.Type),
		},
		Scope: taskScope{
			Nodes: []string{conn.Server.IP},
		},
	}

	_, err := p.callPuppetTask(TaskInjectPublicKey, params)
	return err
}
