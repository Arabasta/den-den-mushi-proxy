package algo

import (
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"math/rand"
	"time"
)

type PowerOf2Random struct {
	name Type
	rand *rand.Rand
}

func NewPowerOf2Random() *PowerOf2Random {
	return &PowerOf2Random{
		name: AlgorithmPowerOf2Random,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (a *PowerOf2Random) String() string {
	return string(a.name)
}

func (a *PowerOf2Random) GetServer(servers []*proxy_host.Record2) (*proxy_host.Record2, error) {
	// randomly select 2 servers
	idx1, idx2 := a.getTwoRandomIndices(len(servers))
	if idx1 == -1 || idx2 == -1 {
		return nil, ErrNoServersAvailable
	}

	if idx1 == idx2 {
		// only 1 server available
		return servers[idx1], nil
	}

	server1, server2 := servers[idx1], servers[idx2]

	// return least conn
	if server1.ActiveSessions <= server2.ActiveSessions {
		return server1, nil
	}
	return server2, nil
}

// getTwoRandomIndices returns two different random indices in the range [0, n).
// If n is 0, it returns -1, -1, if n is 1, it returns 0, 0.
func (a *PowerOf2Random) getTwoRandomIndices(n int) (int, int) {
	if n == 0 {
		return -1, -1
	}

	if n == 1 {
		return 0, 0
	}

	idx1 := a.rand.Intn(n)
	idx2 := a.rand.Intn(n)

	// make sure idx2 is different from idx1
	for idx2 == idx1 {
		idx2 = a.rand.Intn(n)
	}

	return idx1, idx2
}
