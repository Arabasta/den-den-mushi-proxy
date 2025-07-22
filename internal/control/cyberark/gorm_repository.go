package cyberark

import (
	dto "den-den-mushi-Go/pkg/dto/cyberark"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	return &GormRepository{
		db:  db,
		log: log,
	}
}

func (r *GormRepository) FindByObject(o string) (*dto.Record, error) {
	var m dto.Model
	err := r.db.Where("Object = ?", o).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No CyberArk record found for object", zap.String("object", o))
			return nil, nil
		}
		r.log.Error("DB error while fetching CyberArk record", zap.String("object", o), zap.Error(err))
		return nil, err
	}

	return dto.FromModel(&m), nil
}
