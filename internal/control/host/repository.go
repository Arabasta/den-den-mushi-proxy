package host

import host2 "den-den-mushi-Go/pkg/dto/host"

// todo: add context to all service and repos
type Repository interface {
	FindByIp(ip string) (*host2.Record, error)
}
