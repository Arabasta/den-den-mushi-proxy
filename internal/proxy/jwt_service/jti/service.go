package jti

import (
	"den-den-mushi-Go/pkg/config"
	"den-den-mushi-Go/pkg/token"
	"go.uber.org/zap"
)

type Service struct {
	repo    Repository
	log     *zap.Logger
	jwtCfg  *config.JwtAudience
	hostCfg *config.Host
}

func New(repo Repository, log *zap.Logger, jwtCfg *config.JwtAudience, hostCfg *config.Host) *Service {
	log.Info("Initializing JWT JTI Service...")
	return &Service{
		repo:    repo,
		log:     log,
		jwtCfg:  jwtCfg,
		hostCfg: hostCfg,
	}
}

func (s *Service) ConsumeIfNotExists(token *token.Claims) bool {
	s.log.Debug("consume token if not exists", zap.String("id", token.ID))
	jti := ToRecord(token.ID, token.ExpiresAt.Time, token.Subject, s.hostCfg.Ip)

	consumed, err := s.repo.consumeIfNotExists(jti)
	if err != nil {
		s.log.Error("error consuming token", zap.Any("jti", jti), zap.Error(err))
		return false // fail hard
	}
	return consumed
}
