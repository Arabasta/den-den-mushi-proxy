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

func (r *GormRepository) FindAllByIps(ips []string) ([]*dto.Record, error) {
	var models []dto.Model
	err := r.db.Where("IpAddress IN ?", ips).Find(&models).Error
	if err != nil {
		r.log.Error("DB error while fetching hosts", zap.Strings("ips", ips), zap.Error(err))
		return nil, err
	}

	if len(models) == 0 {
		r.log.Debug("No hosts found for provided IPs", zap.Strings("ips", ips))
		return nil, nil
	}

	return dto.FromModels(models), nil
}
