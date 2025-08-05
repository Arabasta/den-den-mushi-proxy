package jti

import (
	"den-den-mushi-Go/pkg/config"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
	"os"
)

type Service struct {
	repo     Repository
	log      *zap.Logger
	jwtCfg   *config.JwtAudience
	hostname string
}

func New(repo Repository, log *zap.Logger, jwtCfg *config.JwtAudience) *Service {
	log.Info("Initializing JWT JTI Service...")
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Failed to get hostname", zap.Error(err))
		os.Exit(1)
	}
	return &Service{
		repo:     repo,
		log:      log,
		jwtCfg:   jwtCfg,
		hostname: hostname,
	}
}

func (s *Service) ConsumeIfNotExists(token *token.Claims) bool {
	s.log.Debug("consume token if not exists", zap.String("id", token.ID))
	jti := ToRecord(token.ID, token.ExpiresAt.Time, token.Subject, s.hostname)

	consumed, err := s.repo.consumeIfNotExists(jti)
	if err != nil {
		s.log.Error("error consuming token", zap.Any("jti", jti), zap.Error(err))
		return false // fail hard
	}
	return consumed
}
