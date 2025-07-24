package host

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/host"
)

// todo: add context to all service and repos
type Repository interface {
	FindByIp(ip string) (*dto.Record, error)
	FindAllByIps(ips []string) ([]*dto.Record, error)
	FindAllByFilter(f filters.HealthcheckPtySession) ([]*dto.Record, error)
}
