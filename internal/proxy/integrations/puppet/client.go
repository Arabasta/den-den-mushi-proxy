package puppet

import (
	config2 "den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/pkg/config"
	"go.uber.org/zap"
)

type Client struct {
	cfg    *config.Puppet
	log    *zap.Logger
	Cfgtmp *config2.Config
}

func NewClient(cfg *config.Puppet, cfgtmp *config2.Config, log *zap.Logger) *Client {
	log.Info("Initializing Puppet client...")
	log.Debug("Puppet Client configuration", zap.Any("config", cfg))

	return &Client{
		cfg:    cfg,
		Cfgtmp: cfgtmp,
		log:    log,
	}
}
