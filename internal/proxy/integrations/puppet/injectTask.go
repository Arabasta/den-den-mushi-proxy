package puppet

import (
	"den-den-mushi-Go/pkg/dto"
)

type taskBody struct {
	Environment string           `json:"environment"`
	Task        string           `json:"task"`
	Params      sshKeyTaskParams `json:"params"`
	Scope       taskScope        `json:"scope"`
}

type sshKeyTaskParams struct {
	PublicKey   string `json:"public_key"`
	ServerIP    string `json:"server_ip"`
	OSUser      string `json:"os_user"`
	ConnPurpose string `json:"conn_purpose"`
	ConnType    string `json:"conn_type"`
}

type taskScope struct {
	Nodes []string `json:"nodes"`
}

func (p *Client) KeyInject(publicKey string, conn dto.Connection) error {
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

	_, err := p.callPuppetTask(TaskInjectPublicKey, params)
	return err
}
