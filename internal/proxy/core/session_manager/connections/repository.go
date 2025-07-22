package connections

import (
	"den-den-mushi-Go/pkg/dto/connections"
)

type Repository interface {
	FindById(id string) (*connections.Record, error)
	Save(*connections.Record) error
}
