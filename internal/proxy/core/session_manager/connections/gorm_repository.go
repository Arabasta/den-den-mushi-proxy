package connections

import (
	"den-den-mushi-Go/pkg/dto/connections"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewGormRepository(db *gorm.DB, log *zap.Logger) *GormRepository {
	log.Info("Initializing Pty Connections GormRepository...")
	return &GormRepository{
		db:  db,
		log: log,
	}
}

func (r *GormRepository) FindById(id string) (*connections.Record, error) {
	var model connections.Model
	err := r.db.Where("Id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No connection found for ID", zap.String("id", id))
			return nil, nil
		}
		r.log.Error("DB error while fetching connection", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return connections.FromModel(&model), nil
}

func (r *GormRepository) Save(rec *connections.Record) error {
	model := connections.ToModel(rec)
	return r.db.Save(&model).Error
}
