package puppet

import (
	"den-den-mushi-Go/pkg/dto"
)

func (p *Client) KeyRemove(publicKey string, conn dto.Connection) error {
	params := taskBody{
		Environment: "production",
		Task:        "task",
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

	_, err := p.callPuppetTask(TaskRemovePublicKey, params)
	return err
}
