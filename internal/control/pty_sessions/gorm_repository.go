package pty_sessions

import (
	dto "den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/types"
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

func (r *GormRepository) FindById(id string) (*dto.Record, error) {
	var session dto.Model
	err := r.db.
		Preload("Connections").
		Preload("ProxyDetails").
		First(&session, "id = ?", id).Error

	if err != nil {
		return nil, err
	}
	return dto.FromModel(&session), nil
}

func (r *GormRepository) FindByStartConnChangeRequestIdsAndState(changeIDs []string, state types.PtySessionState) ([]*dto.Record, error) {
	var sessions []dto.Model
	err := r.db.
		Preload("Connections").
		Preload("ProxyDetails").
		Where("start_conn_change_request_id IN ?", changeIDs).
		Where("state = ?", state).
		Find(&sessions).Error

	if err != nil {
		return nil, err
	}
	return dto.FromModels(sessions), nil
}
