package certname

import (
	"den-den-mushi-Go/pkg/dto/host"
	dto "den-den-mushi-Go/pkg/dto/puppet_trusted"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	return &GormRepository{db: db, log: log}
}

// temporary for demo as usual todo refactor
func (r *GormRepository) FindCertnameByIp(ip string) (*dto.Record, error) {
	var m dto.Model

	var hostModel host.Model
	hostTable := hostModel.TableName()
	puppetTable := m.TableName()

	query := fmt.Sprintf(`
		SELECT p.certname
		FROM %s h
		JOIN %s p ON h.hostname = p.hostname
		WHERE h.ip_address = ?
	`, hostTable, puppetTable)

	err := r.db.Raw(query, ip).Scan(&m).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No certname found for IP", zap.String("ip", ip))
			return nil, nil
		}
		r.log.Error("DB error while fetching certname", zap.String("ip", ip), zap.Error(err))
		return nil, err
	}
	return dto.FromModel(&m), nil
}
