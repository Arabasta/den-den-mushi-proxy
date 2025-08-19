package host

import (
	"den-den-mushi-Go/internal/control/filters"
	dto "den-den-mushi-Go/pkg/dto/host"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
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
	err := r.db.Where("IPADDRESS = ?", ip).First(&m).Error
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

	err := r.db.Where("IPADDRESS IN ? AND PLATFORM = ?", ips, "Linux").Find(&models).Error
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
		appcode := strings.ToUpper(*f.Appcode)
		query = query.Where("APPLICATION_CODE LIKE ?", "%"+appcode+"%")
	}

	if f.Environment != nil && len(*f.Environment) > 0 {
		env := strings.ToUpper(*f.Environment)
		query = query.Where("ENVIRONMENT = ?", env)
	}

	if f.Country != nil && len(*f.Country) > 0 {
		country := strings.ToUpper(*f.Country)
		query = query.Where("COUNTRY LIKE ?", "%"+country+"%")
	}

	if f.Lob != nil && len(*f.Lob) > 0 {
		lob := strings.ToUpper(*f.Lob)
		query = query.Where("LOB LIKE ?", "%"+lob+"%")
	}

	if f.OsType != nil && len(*f.OsType) > 0 {
		query = query.Where("PLATFORM LIKE ?", "%"+*f.OsType+"%")
	}

	if f.Hostname != nil && len(*f.Hostname) > 0 {
		query = query.Where("HOSTNAME LIKE ?", "%"+*f.Hostname+"%")
	}

	//if f.Status != nil && len(*f.Status) > 0 {
	validStatus := []string{"Active", "Pre-Production", "Staging", "Tech Live", "TH_WIP", "Decommissioning"}
	query = query.Where("STATUS IN ?", validStatus)
	//}

	if f.SystemType != nil && len(*f.SystemType) > 0 {
		query = query.Where("SYSTEM_TYPE LIKE ?", "%"+*f.SystemType+"%")
	}

	query = query.Order("HOSTNAME DESC")

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

	return dto.FromModels(models), nil
}

func (r *GormRepository) CountAllByFilter(f filters.HealthcheckPtySession) (int64, error) {
	var count int64
	query := r.db.Model(&dto.Model{})

	if f.Ip != nil && len(*f.Ip) > 0 {
		query = query.Where("IPADDRESS = ?", *f.Ip)
	}

	if f.Appcode != nil && len(*f.Appcode) > 0 {
		appcode := strings.ToUpper(*f.Appcode)
		query = query.Where("APPLICATION_CODE LIKE ?", "%"+appcode+"%")
	}

	if f.Environment != nil && len(*f.Environment) > 0 {
		env := strings.ToUpper(*f.Environment)
		query = query.Where("ENVIRONMENT = ?", env)
	}

	if f.Country != nil && len(*f.Country) > 0 {
		country := strings.ToUpper(*f.Country)
		query = query.Where("COUNTRY LIKE ?", "%"+country+"%")
	}

	if f.Lob != nil && len(*f.Lob) > 0 {
		lob := strings.ToUpper(*f.Lob)
		query = query.Where("LOB LIKE ?", "%"+lob+"%")
	}

	if f.OsType != nil && len(*f.OsType) > 0 {
		query = query.Where("PLATFORM LIKE ?", "%"+*f.OsType+"%")
	}

	if f.Hostname != nil && len(*f.Hostname) > 0 {
		query = query.Where("HOSTNAME LIKE ?", "%"+*f.Hostname+"%")
	}

	if f.Status != nil && len(*f.Status) > 0 {
		validStatus := []string{"Active", "Pre-Production", "Staging", "Tech Live", "TH_WIP"}
		query = query.Where("STATUS IN ?", validStatus)
	}

	if f.SystemType != nil && len(*f.SystemType) > 0 {
		query = query.Where("SYSTEM_TYPE LIKE ?", "%"+*f.SystemType+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		r.log.Error("DB error while counting hosts by filter", zap.Any("filter", f), zap.Error(err))
		return 0, err
	}

	r.log.Debug("DB count", zap.Int64("count", count))
	return count, nil
}
