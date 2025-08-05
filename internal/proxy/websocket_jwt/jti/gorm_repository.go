package jti

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	log.Info("Initializing Gorm JTI Repository...")
	return &GormRepository{
		db:  db,
		log: log,
	}
}

func (r *GormRepository) consumeIfNotExists(jti *Record) (bool, error) {
	// INSERT with unique constraint on id
	result := r.db.Create(ToModel(jti))
	if result.Error != nil {
		// if duplicate key error = consumed
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			return false, nil
		}
		return false, result.Error
	}
	// can INSERT
	return true, nil
}
