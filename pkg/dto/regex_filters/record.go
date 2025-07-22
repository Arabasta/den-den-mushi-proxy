package regex_filters

import (
	"den-den-mushi-Go/pkg/types"
	"regexp"
	"time"
)

type Record struct {
	Id           uint
	FilterType   types.Filter
	RegexPattern regexp.Regexp
	OuGroup      string
	IsEnabled    bool
	CreatedBy    string
	DeletedBy    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    *time.Time // soft delete, nullable
}
