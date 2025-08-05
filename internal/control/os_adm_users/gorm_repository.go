package os_adm_users

import (
	dto "den-den-mushi-Go/pkg/dto/os_adm_users"
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

func (r *GormRepository) FindAllByUserId(userId string) ([]*dto.Record, error) {
	var records []dto.Model
	err := r.db.
		Where("UserId = ?", userId).
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	return dto.FromModels(records), nil
}
