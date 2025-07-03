package puppet

import (
	"den-den-mushi-Go/internal/proxy/config"
	"go.uber.org/zap"
)

type Client struct {
	cfg *config.Config
	log *zap.Logger
}

func NewClient(cfg *config.Config, log *zap.Logger) *Client {
	return &Client{cfg: cfg, log: log}
}
