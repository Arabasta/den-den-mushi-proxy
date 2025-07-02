package puppet

import (
	"den-den-mushi-Go/pkg/dto"
	"errors"
	"fmt"
)

type PuppetTask string

const (
	PuppetTaskInjectPublicKey PuppetTask = "inject_public_key"
	PuppetTaskRemovePublicKey PuppetTask = "remove_public_key"
)

func (pc *Client) callPuppetTask(task PuppetTask, payload interface{}) (interface{}, error) {
	switch task {
	case PuppetTaskInjectPublicKey:
		p := payload.(injectKeyParams)
		fmt.Println(p)
		// todo: implement
		return "ok", nil

	case PuppetTaskRemovePublicKey:
		// todo: implement
		return nil, errors.New("not implemented yet")

	default:
		return nil, fmt.Errorf("unsupported puppet task: %s", task)
	}
}

type injectKeyParams struct {
	PublicKey   string
	ServerIP    string
	OSUser      string
	ConnPurpose dto.ConnectionPurpose
	ConnType    dto.ConnectionType
}

func (pc *Client) PuppetKeyInject(publicKey string, conn dto.Connection) error {
	params := injectKeyParams{
		PublicKey:   publicKey,
		ServerIP:    conn.Server.IP,
		OSUser:      conn.Server.OSUser,
		ConnPurpose: conn.Purpose,
		ConnType:    conn.Type,
	}

	_, err := pc.callPuppetTask(PuppetTaskInjectPublicKey, params)
	return err
}
