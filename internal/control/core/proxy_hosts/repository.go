package proxy_hosts

import (
	"context"
	"den-den-mushi-Go/pkg/dto/proxy_host"
)

type Repository interface {
	FindAll(ctx context.Context) ([]*proxy_host.Record2, error)
	FindByHostname(ctx context.Context, hostname string) (*proxy_host.Record2, error)
}
