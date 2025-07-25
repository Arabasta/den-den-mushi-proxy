package puppet

import (
	"den-den-mushi-Go/pkg/dto"
)

type removeKeyParams struct {
	PublicKey   string `json:"public_key"`
	ServerIP    string `json:"server_ip"`
	OSUser      string `json:"os_user"`
	ConnPurpose string `json:"conn_purpose"`
	ConnType    string `json:"conn_type"`
}

func (p *Client) KeyRemove(publicKey string, conn dto.Connection) error {
	params := removeKeyParams{
		PublicKey:   publicKey,
		ServerIP:    conn.Server.IP,
		OSUser:      conn.Server.OSUser,
		ConnPurpose: string(conn.Purpose),
		ConnType:    string(conn.Type),
	}

	_, err := p.callPuppetTask(TaskRemovePublicKey, params)
	return err
}
