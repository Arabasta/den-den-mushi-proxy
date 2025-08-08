package proxy_lb

import (
	dto "den-den-mushi-Go/pkg/dto/proxy_lb"
	"den-den-mushi-Go/pkg/types"
)

type Repository interface {
	FindByProxyType(t types.Proxy) (*dto.Record, error)
}
