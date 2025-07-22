package cyberark

import dto "den-den-mushi-Go/pkg/dto/cyberark"

type Repository interface {
	FindByObject(o string) (*dto.Record, error)
}
