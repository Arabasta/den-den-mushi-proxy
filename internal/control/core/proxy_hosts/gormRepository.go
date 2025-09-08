package proxy_hosts

import (
	"context"
	"den-den-mushi-Go/pkg/dto/proxy_host"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	log.Info("Initializing Proxy Hosts GormRepository")
	return &GormRepository{db: db, log: log}
}

func (r *GormRepository) FindAll(ctx context.Context) ([]*proxy_host.Record2, error) {
	var models []*proxy_host.Model2
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		r.log.Error("Failed to fetch proxy hosts", zap.Error(err))
		return nil, err
	}

	return proxy_host.FromModels2(models), nil
}
