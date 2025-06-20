package puppet

import (
	"den-den-mushi-Go/internal/config"
	"go.uber.org/zap"
)

type PuppetClient struct {
	cfg    *config.Config
	logger *zap.Logger
}

func NewPuppetClient(cfg *config.Config, logger *zap.Logger) *PuppetClient {
	return &PuppetClient{cfg: cfg, logger: logger}
}
