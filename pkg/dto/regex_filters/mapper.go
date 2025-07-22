package regex_filters

import (
	"den-den-mushi-Go/pkg/util/regex"
	"gorm.io/gorm"
	"time"
)

func FromModel(m *Model) *Record {
	if m == nil {
		return nil
	}
	re, err := regex.CompilePattern(m.RegexPattern)
	if err != nil {
		return nil
	}

	return &Record{
		Id:           m.Id,
		FilterType:   m.FilterType,
		RegexPattern: *re,
		OuGroup:      m.OuGroup,
		IsEnabled:    m.IsEnabled,
		CreatedBy:    m.CreatedBy,
		DeletedBy:    m.DeletedBy,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt: func() *time.Time {
			if m.DeletedAt.Valid {
				return &m.DeletedAt.Time
			}
			return nil
		}(),
	}
}

func FromModels(models []Model) *[]Record {
	if len(models) == 0 {
		return nil
	}
	records := make([]Record, len(models))
	for i := range models {
		records[i] = *FromModel(&models[i])
	}
	return &records
}

func ToModel(r *Record) *Model {
	if r == nil {
		return nil
	}

	var deletedAt gorm.DeletedAt
	if r.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{
			Time:  *r.DeletedAt,
			Valid: true,
		}
	}

	return &Model{
		Id:           r.Id,
		FilterType:   r.FilterType,
		RegexPattern: r.RegexPattern.String(),
		OuGroup:      r.OuGroup,
		IsEnabled:    r.IsEnabled,
		CreatedBy:    r.CreatedBy,
		DeletedBy:    r.DeletedBy,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
		UpdatedBy:    r.UpdatedBy,
		DeletedAt:    deletedAt,
	}
}

func ToModels(records []*Record) *[]Model {
	if len(records) == 0 {
		return nil
	}
	models := make([]Model, len(records))
	for i, r := range records {
		models[i] = *ToModel(r)
	}
	return &models
}
