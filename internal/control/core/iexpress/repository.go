package iexpress

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/iexpress"
)

type Repository interface {
	FindByTicketNumber(num string) (*dto.Record, error)
	FindApprovedByFilter(filter filters.ListIexpress) ([]*dto.Record, error)
	CountApprovedByFilter(filter filters.ListIexpress) (int, error)
}
