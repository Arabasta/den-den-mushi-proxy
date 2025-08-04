package pty_sessions

import (
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
