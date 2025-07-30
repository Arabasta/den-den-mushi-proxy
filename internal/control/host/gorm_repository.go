package host

import (
	"den-den-mushi-Go/internal/control/filters"
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

func (r *GormRepository) FindAllLinuxOsByIps(ips []string) ([]*dto.Record, error) {
	var models []dto.Model

	err := r.db.Where("IpAddress IN ? AND PLATFORM = ?", ips, "Linux").Find(&models).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No hosts found for provided IPs", zap.Strings("ips", ips))
			return nil, nil
		}
		r.log.Error("DB error while fetching hosts", zap.Strings("ips", ips), zap.Error(err))
		return nil, err
	}

	return dto.FromModels(models), nil
}

func (r *GormRepository) FindAllByFilter(f filters.HealthcheckPtySession) ([]*dto.Record, error) {
	var models []dto.Model
	query := r.db.Model(&dto.Model{})

	if f.Ip != nil && len(*f.Ip) > 0 {
		query = query.Where("IPADDRESS = ?", *f.Ip)
	}

	if f.Appcode != nil && len(*f.Appcode) > 0 {
		query = query.Where("APPLICATION_CODE = ?", *f.Appcode)
	}

	if f.Environment != nil && len(*f.Environment) > 0 {
		query = query.Where("ENVIRONMENT = ?", *f.Environment)
	}

	if f.Country != nil && len(*f.Country) > 0 {
		query = query.Where("COUNTRY = ?", *f.Country)
	}

	if f.Lob != nil && len(*f.Lob) > 0 {
		query = query.Where("LOB = ?", *f.Lob)
	}

	if f.OsType != nil && len(*f.OsType) > 0 {
		query = query.Where("PLATFORM = ?", *f.OsType)
	}

	if f.Hostname != nil && len(*f.Hostname) > 0 {
		query = query.Where("HOSTNAME = ?", *f.Hostname)
	}

	if f.Status != nil && len(*f.Status) > 0 {
		query = query.Where("STATUS = ?", *f.Status)
	}

	if f.SystemType != nil && len(*f.SystemType) > 0 {
		query = query.Where("SYSTEM_TYPE = ?", *f.SystemType)
	}

	page := f.Page
	if page < 1 {
		page = 1
	}

	pageSize := f.PageSize
	if pageSize <= 0 || pageSize > 1000 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	err := query.Find(&models).Error
	if err != nil {
		r.log.Error("DB error while fetching hosts by filter", zap.Any("filter", f), zap.Error(err))
		return nil, err
	}

	r.log.Debug("Fetched hosts", zap.Int("Count", len(models)))
	return dto.FromModels(models), nil
}
