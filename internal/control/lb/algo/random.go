package algo

import (
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"math/rand"
	"time"
)

type Random struct {
	name Type
	rand *rand.Rand
}

func NewRandom() *Random {
	return &Random{
		name: AlgorithmRandom,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (a *Random) String() string {
	return string(a.name)
}

func (a *Random) GetServer(servers []*proxy_host.Record2) (*proxy_host.Record2, error) {
	if len(servers) == 0 {
		return nil, ErrNoServersAvailable
	}

	index := a.rand.Intn(len(servers))
	return servers[index], nil
}
