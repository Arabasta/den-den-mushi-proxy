package proxy_lb

import "den-den-mushi-Go/pkg/types"

type Repository interface {
	FindByProxyType(t types.Proxy) (*Entity, error)
}
