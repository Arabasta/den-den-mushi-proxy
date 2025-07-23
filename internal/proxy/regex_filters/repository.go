package regex_filters

import (
	dto "den-den-mushi-Go/pkg/dto/regex_filters"
	"den-den-mushi-Go/pkg/types"
)

type Repository interface {
	FindAllEnabledByFilterType(filterType types.Filter) (*[]dto.Record, error)
}
