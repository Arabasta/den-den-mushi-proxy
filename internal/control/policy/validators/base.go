package validators

import (
	"den-den-mushi-Go/internal/control/config"
	"go.uber.org/zap"
)

type Validator struct {
	log *zap.Logger
	cfg *config.Config
}

func NewValidator(log *zap.Logger, cfg *config.Config) *Validator {
	log.Info("Initializing Policy Validator...")
	return &Validator{
		log: log,
		cfg: cfg,
	}
}
