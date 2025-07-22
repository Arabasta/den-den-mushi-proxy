package change_request

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/change_request"
)

type Repository interface {
	FindByTicketNumber(num string) (*dto.Record, error)
	FindChangeRequests(filter filters.ListCR) ([]*dto.Record, error)
}
