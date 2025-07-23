package connection

import (
	"den-den-mushi-Go/pkg/dto/connections"
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

func (r *GormRepository) FindById(id string) (*connections.Record, error) {
	var m connections.Model
	err := r.db.Where("Id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No connection found for ID", zap.String("id", id))
			return nil, nil
		}
		r.log.Error("DB error while fetching connection", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return connections.FromModel(&m), nil
}

func (r *GormRepository) FindActiveImplementorByPtySessionId(ptySessionId string) (*connections.Record, error) {
	var m *connections.Model
	r.log.Debug("Finding active implementor connections",
		zap.String("pty_session_id", ptySessionId),
		zap.String("status", string(types.ConnectionStatusActive)),
		zap.String("startRole", string(types.Implementor)))

	err := r.db.Where("pty_session_id = ? AND Status = ? AND Start_role = ?",
		ptySessionId, types.ConnectionStatusActive, types.Implementor).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No active implementor connection found for pty session", zap.String("ptySessionId", ptySessionId))
			return nil, nil
		}
		r.log.Error("DB error while fetching active connections", zap.String("ptySessionId", ptySessionId), zap.Error(err))
		return nil, err
	}

	return connections.FromModel(m), nil
}
