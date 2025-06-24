package puppet

import (
	"den-den-mushi-Go/internal/config"
	"go.uber.org/zap"
)

type PuppetClient struct {
	cfg *config.Config
	log *zap.Logger
}

func NewPuppetClient(cfg *config.Config, log *zap.Logger) *PuppetClient {
	return &PuppetClient{cfg: cfg, log: log}
}
