package jti

import (
	"go.uber.org/zap"
	"time"
)

type Service struct {
	repo repository
	ttl  time.Duration
	log  *zap.Logger
}

func New(repo repository, ttl time.Duration, log *zap.Logger) *Service {
	return &Service{
		repo: repo,
		ttl:  ttl,
		log:  log,
	}
}

func (s *Service) IsConsumed(id string) bool {
	found, err := s.repo.tryGetJti(id)
	if found {
		s.log.Error("already consumed token", zap.String("id", id))
		return true
	}
	if err != nil {
		s.log.Error("error checking token", zap.String("id", id), zap.Error(err))
		// if error checking token, assume it's consumed
		return true
	}

	s.log.Debug("token not consumed", zap.String("id", id))
	return false
}

func (s *Service) Consume(id string) bool {
	s.log.Debug("consuming token", zap.String("id", id))
	err := s.repo.addJti(id)
	if err != nil {
		s.log.Error("error consuming token", zap.String("id", id), zap.Error(err))
		// fail hard if can't consume
		return false
	}

	// for now just delete after ttl, in prod maybe don't cleanup?
	// todo: discuss this
	// go func() { <-time.After(s.ttl); s.repo.Delete(id) }()
	s.log.Info("consumed token", zap.String("id", id))
	return true
}
