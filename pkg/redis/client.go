package redis

import (
	"context"
	"den-den-mushi-Go/pkg/config"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

func Client(cfg config.Redis, log *zap.Logger) (*redis.ClusterClient, error) {
	log.Info("Connecting to Redis...")
	log.Debug("Connection parameters",
		zap.Strings("Addrs", cfg.Addrs),
		zap.Int("PoolSize", cfg.PoolSize))

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.Addrs,
		Password: cfg.Password,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ping
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis cluster ping failed: %w", err)
	}

	log.Info("Connected to Redis cluster")
	return client, nil
}
