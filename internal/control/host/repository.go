package host

import dto "den-den-mushi-Go/pkg/dto/host"

// todo: add context to all service and repos
type Repository interface {
	FindByIp(ip string) (*dto.Record, error)
	FindAllByIps(ips []string) ([]*dto.Record, error)
}
