package puppet

import (
	"den-den-mushi-Go/pkg/config"
	"go.uber.org/zap"
)

type Client struct {
	cfg *config.Puppet
	log *zap.Logger
}

func NewClient(cfg *config.Puppet, log *zap.Logger) *Client {
	log.Info("Initializing Puppet client...")
	log.Debug("Puppet Client configuration", zap.Any("config", cfg))

	return &Client{
		cfg: cfg,
		log: log,
	}
}
