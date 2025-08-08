package certname

import dto "den-den-mushi-Go/pkg/dto/puppet_trusted"

type Repository interface {
	FindCertnameByIp(ip string) (*dto.Record, error)
}
