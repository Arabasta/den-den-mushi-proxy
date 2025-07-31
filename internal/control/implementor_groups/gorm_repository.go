package implementor_groups

import (
	"den-den-mushi-Go/internal/control/config"
	dto "den-den-mushi-Go/pkg/dto/implementor_groups"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.Logger
	cfg *config.Config
}

func NewGormRepository(db *gorm.DB, log *zap.Logger, cfg *config.Config) *GormRepository {
	return &GormRepository{
		db:  db,
		log: log,
		cfg: cfg,
	}
}

func (r *GormRepository) FindAllByUserId(userId string) ([]*dto.Record, error) {
	var models []dto.Model

	email := userId + "@%"

	if err := r.db.Where("MemberEmail LIKE ? AND GroupMembershhipStatus = ?",
		email, "Active").Find(&models).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No implementor groups found for user ID", zap.String("user_id", userId))
			return nil, nil
		}
		r.log.Error("Failed to find implementor groups by user ID", zap.String("user_id", userId), zap.Error(err))
		return nil, err
	}
	return dto.FromModels(models), nil
}
