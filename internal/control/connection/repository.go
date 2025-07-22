package connection

import "den-den-mushi-Go/pkg/dto/connections"

type Repository interface {
	FindById(id string) (*connections.Record, error)
}
