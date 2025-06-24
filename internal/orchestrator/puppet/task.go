package puppet

import (
	"den-den-mushi-Go/pkg/dto/connection"
	"errors"
	"fmt"
)

type PuppetTask string

const (
	PuppetTaskInjectPublicKey PuppetTask = "inject_public_key"
	PuppetTaskRemovePublicKey PuppetTask = "remove_public_key"
)

func (pc *PuppetClient) callPuppetTask(task PuppetTask, payload interface{}) (interface{}, error) {
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
	ConnPurpose connection.ConnectionPurpose
	ConnType    connection.ConnectionType
}

func (pc *PuppetClient) PuppetKeyInject(publicKey string, conn connection.Connection) error {
	params := injectKeyParams{
		PublicKey:   publicKey,
		ServerIP:    conn.ServerIP,
		OSUser:      conn.OSUser,
		ConnPurpose: conn.Purpose,
		ConnType:    conn.Type,
	}

	_, err := pc.callPuppetTask(PuppetTaskInjectPublicKey, params)
	return err
}
