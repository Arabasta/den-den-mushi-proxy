package pty_sessions

import (
	"den-den-mushi-Go/pkg/dto/pty_sessions"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	log.Info("Initializing Pty Session GormRepository...")
	return &GormRepository{
		db:  db,
		log: log,
	}
}

func (r *GormRepository) FindById(id string) (*pty_sessions.Record, error) {
	var model pty_sessions.Model
	err := r.db.Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No PTY session found for ID", zap.String("id", id))
			return nil, nil
		}
		r.log.Error("DB error while fetching PTY session", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return pty_sessions.FromModel(&model), nil
}

func (r *GormRepository) Save(session *pty_sessions.Record) error {
	model := pty_sessions.ToModel(session)
	if err := r.db.Save(&model).Error; err != nil {
		r.log.Error("Failed to save PTY session", zap.String("id", session.ID), zap.Error(err))
		return err
	}
	return nil
}
