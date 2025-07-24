package pty_sessions

import (
	dto "den-den-mushi-Go/pkg/dto/pty_sessions"
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
		Where("start_conn_cr IN ?", changeIDs).
		Where("state = ?", state).
		Find(&sessions).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No pty sessions found for provided change request IDs and state",
				zap.Strings("changeIDs", changeIDs), zap.String("state", string(state)))
			return nil, nil
		}
		return nil, err
	}
	return dto.FromModels(sessions), nil
}

func (r *GormRepository) FindAllByChangeRequestIDAndServerIPs(changeRequestID string, ips []string) ([]*dto.Record, error) {
	var models []dto.Model
	err := r.db.
		Preload("Connections").
		Where("start_conn_cr_id = ?", changeRequestID).
		Where("start_conn_server_ip IN ?", ips).
		Find(&models).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No pty sessions found for provided change request ID and server IPs",
				zap.String("changeRequestID", changeRequestID), zap.Strings("ips", ips))
			return nil, nil
		}
		return nil, err
	}
	return dto.FromModels(models), nil
}

func (r *GormRepository) FindAllByStartConnServerIpsAndState(hostips []string, state *types.PtySessionState) ([]*dto.Record, error) {
	var models []dto.Model

	query := r.db.Model(&dto.Model{}).Preload("Connections")

	if state != nil {
		query = query.Where("state = ?", *state)
	}

	if len(hostips) > 0 {
		query = query.Where("start_conn_server_ip IN ?", hostips)
	}

	if err := query.Find(&models).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Debug("No pty sessions found for provided server IPs and state",
				zap.Strings("hostips", hostips))
			return nil, nil
		}
		return nil, err
	}
	return dto.FromModels(models), nil
}
