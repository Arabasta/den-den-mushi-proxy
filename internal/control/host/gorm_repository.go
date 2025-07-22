package host

import (
	dto "den-den-mushi-Go/pkg/dto/host"
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

func (r *GormRepository) FindByIp(ip string) (*dto.Record, error) {
	var m dto.Model
	err := r.db.Where("IpAddress = ?", ip).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		r.log.Debug("No host found for IP", zap.String("ip", ip))
		return nil, nil
	}
	if err != nil {
		r.log.Error("DB error while fetching host", zap.String("ip", ip), zap.Error(err))
		return nil, err
	}

	return dto.FromModel(&m), nil
}
