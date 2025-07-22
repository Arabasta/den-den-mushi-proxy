package proxy_lb

import (
	"den-den-mushi-Go/pkg/dto/proxy_lb"
	"den-den-mushi-Go/pkg/types"
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

func (r *GormRepository) FindByProxyType(t types.Proxy) (*proxy_lb.Record, error) {
	var m proxy_lb.Model
	err := r.db.
		Where("type = ?", t).
		First(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.log.Error("Failed to find record by proxy type", zap.Error(err))
		return nil, err
	}

	return proxy_lb.FromModel(&m), nil
}
