package pty_sessions

import (
	dto "den-den-mushi-Go/pkg/dto/pty_sessions"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	log.Info("Initializing pty_sessions GormRepository...")
	return &GormRepository{
		db:  db,
		log: log,
	}
}

func (r *GormRepository) FindAllByChangeRequestID(changeRequestID string) ([]*dto.Record, error) {
	var models []dto.Model
	err := r.db.
		Where("start_conn_cr_id = ?", changeRequestID).
		Order("start_time ASC").
		Find(&models).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No pty sessions found for provided change request ID and server IPs",
				zap.String("changeRequestID", changeRequestID))
			return nil, nil
		}
		return nil, err
	}
	return dto.FromModels(models), nil
}
