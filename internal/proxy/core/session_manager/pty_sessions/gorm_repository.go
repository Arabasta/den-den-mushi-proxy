package pty_sessions

import (
	"den-den-mushi-Go/pkg/dto/connections"
	"den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/types"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
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

// todo: this will not cleanup active connections with a closed session
// todo: this won't scale well
func (r *GormRepository) FindActiveByProxyHostWithConnections(proxyHost string) ([]*pty_sessions.Record, error) {
	var sessions []pty_sessions.Model
	err := r.db.
		Preload("Connections").
		Where("proxy_host_name = ? AND state = ?", proxyHost, types.Active).
		Find(&sessions).Error

	if err != nil {
		r.log.Error("DB error while fetching active PTY sessions with connections", zap.String("proxyHost", proxyHost), zap.Error(err))
		return nil, err
	}
	return pty_sessions.FromModels(sessions), nil
}

func (r *GormRepository) CloseSessionsAndConnections(sessionIDs []string, connectionIDs []string) error {
	now := time.Now()

	// close sessions
	if err := r.db.Model(&pty_sessions.Model{}).
		Where("id IN ?", sessionIDs).
		Updates(map[string]interface{}{
			"state":    types.Closed,
			"end_time": now,
		}).Error; err != nil {
		return err
	}

	// close connections
	if len(connectionIDs) > 0 {
		if err := r.db.Model(&connections.Model{}).
			Where("id IN ?", connectionIDs).
			Updates(map[string]interface{}{
				"status":     types.ConnectionStatusClosed,
				"leave_time": now,
			}).Error; err != nil {
			return err
		}
	}

	return nil
}
