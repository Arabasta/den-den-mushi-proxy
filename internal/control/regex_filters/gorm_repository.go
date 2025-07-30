package regex_filters

import (
	dto "den-den-mushi-Go/pkg/dto/regex_filters"
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

func (r *GormRepository) Save(rec *dto.Record) (*dto.Record, error) {
	m := dto.ToModel(rec)

	err := r.db.Save(m).Error
	if err != nil {
		r.log.Error("Failed to save regex filter", zap.Error(err))
		return nil, err
	}
	return dto.FromModel(m), nil
}

func (r *GormRepository) FindAllByFilterTypeAndOuGroup(filterType types.Filter, ouGroup string) (*[]dto.Record, error) {
	var models []dto.Model

	err := r.db.Where("type = ? AND ougroup = ?", filterType, ouGroup).
		Order("CreatedAt DESC").
		Find(&models).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No regex filters found for type and OU group", zap.String("type", string(filterType)), zap.String("ou_group", ouGroup))
			return nil, nil
		}
		r.log.Error("DB error while fetching regex filters by type and OU group", zap.String("type", string(filterType)), zap.String("ou_group", ouGroup), zap.Error(err))
		return nil, err
	}

	if len(models) == 0 {
		r.log.Debug("No regex filters found for type and OU group", zap.String("type", string(filterType)), zap.String("ou_group", ouGroup))
		return nil, nil
	}

	return dto.FromModels(models), nil
}

func (r *GormRepository) FindById(id uint) (*dto.Record, error) {
	var model dto.Model
	err := r.db.First(&model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No regex filter found for ID", zap.Uint("id", id))
			return nil, gorm.ErrRecordNotFound
		}
		r.log.Error("DB error while fetching regex filter by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return dto.FromModel(&model), nil
}

func (r *GormRepository) Delete(rec *dto.Record) error {
	return r.db.Delete(dto.ToModel(rec)).Error
}
