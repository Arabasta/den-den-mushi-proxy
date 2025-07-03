package jti

import (
	"go.uber.org/zap"
	"sync"
)

type InMemRepository struct {
	jtiStore sync.Map
	log      *zap.Logger
}

func NewInMemRepository(log *zap.Logger) *InMemRepository {
	return &InMemRepository{
		jtiStore: sync.Map{},
		log:      log,
	}
}

func (r *InMemRepository) tryGetJti(jti string) (bool, error) {
	_, found := r.jtiStore.Load(jti)
	if found {
		r.log.Debug("JTI already exists", zap.String("jti", jti))
		return true, nil
	}

	return false, nil
}

func (r *InMemRepository) addJti(jti string) error {
	r.jtiStore.Store(jti, struct{}{})
	r.log.Debug("JTI added", zap.String("jti", jti))
	return nil
}
