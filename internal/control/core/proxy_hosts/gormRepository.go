package proxy_hosts

import (
	"context"
	dto "den-den-mushi-Go/pkg/dto/proxy_host"
	"errors"

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

func (r *GormRepository) FindAll(ctx context.Context) ([]*dto.Record2, error) {
	var models []*dto.Model2
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		r.log.Error("Failed to fetch proxy hosts", zap.Error(err))
		return nil, err
	}

	return dto.FromModels2(models), nil
}

func (r *GormRepository) FindByHostname(ctx context.Context, hostname string) (*dto.Record2, error) {
	var m dto.Model2
	err := r.db.WithContext(ctx).
		Where("HostName = ?", hostname).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No proxy host found for HostName", zap.String("hostname", hostname))
			// must return err
			return nil, errors.New(hostname + " not found")
		}
		r.log.Error("DB error while fetching proxy hostname", zap.String("hostname", hostname), zap.Error(err))
		return nil, err
	}

	return dto.FromModel2(&m), nil
}
