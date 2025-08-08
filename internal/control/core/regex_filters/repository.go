package regex_filters

import (
	dto "den-den-mushi-Go/pkg/dto/regex_filters"
	"den-den-mushi-Go/pkg/types"
)

type Repository interface {
	Save(r *dto.Record) (*dto.Record, error)
	FindAllByFilterTypeAndOuGroup(filterType types.Filter, ouGroup string) (*[]dto.Record, error)
	FindById(id uint) (*dto.Record, error)
	Delete(r *dto.Record) error
}
