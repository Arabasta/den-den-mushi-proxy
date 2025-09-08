package algo

import (
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"errors"
)

type Rithm interface {
	String() string
	GetServer(servers []*proxy_host.Record2) (*proxy_host.Record2, error)
}

var (
	ErrNoServersAvailable = errors.New("no proxy servers available")
)
