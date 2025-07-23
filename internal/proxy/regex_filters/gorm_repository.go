package regex_filters

import (
	dto "den-den-mushi-Go/pkg/dto/regex_filters"
	"den-den-mushi-Go/pkg/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	log.Info("Initializing GormRepository for regex filters")
	return &GormRepository{
		db:  db,
		log: log,
	}
}

func (r *GormRepository) FindAllEnabledByFilterType(filterType types.Filter) (*[]dto.Record, error) {
	var models []dto.Model
	err := r.db.Where("Type = ? AND IsEnabled = ?", filterType, true).
		Find(&models).Error
	if err != nil {
		r.log.Error("Failed to find regex filters by type", zap.Error(err), zap.String("type", string(filterType)))
		return nil, err
	}

	if len(models) == 0 {
		r.log.Debug("No regex filters found for type", zap.String("type", string(filterType)))
		return nil, nil
	}

	return dto.FromModels(models), nil
}
