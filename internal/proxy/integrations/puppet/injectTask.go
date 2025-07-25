package puppet

import (
	"den-den-mushi-Go/pkg/dto"
)

type injectKeyParams struct {
	PublicKey   string `json:"public_key"`
	ServerIP    string `json:"server_ip"`
	OSUser      string `json:"os_user"`
	ConnPurpose string `json:"conn_purpose"`
	ConnType    string `json:"conn_type"`
}

func (p *Client) KeyInject(publicKey string, conn dto.Connection) error {
	params := injectKeyParams{
		PublicKey:   publicKey,
		ServerIP:    conn.Server.IP,
		OSUser:      conn.Server.OSUser,
		ConnPurpose: string(conn.Purpose),
		ConnType:    string(conn.Type),
	}

	_, err := p.callPuppetTask(TaskInjectPublicKey, params)
	return err
}
