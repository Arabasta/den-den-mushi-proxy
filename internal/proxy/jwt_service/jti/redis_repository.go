package jti

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type RedisRepository struct {
	redis *redis.Client
	log   *zap.Logger
}

func NewRedisRepository(redis *redis.Client, log *zap.Logger) *RedisRepository {
	log.Info("Initializing Redis JTI Repository...")
	return &RedisRepository{
		redis: redis,
		log:   log,
	}
}

// consumeIfNotExists inserts JTI if not present and sets expiry.
// returns true if token was not consumed before, false if already consumed
func (r *RedisRepository) consumeIfNotExists(jti *Record) (bool, error) {
	ttl := time.Until(jti.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Minute // fallback TTL if expiry already passed
	}

	ok, err := r.redis.SetNX(context.Background(), jti.Id, "Subject: "+jti.Subject+" FromProxy: "+jti.FromProxy, ttl).Result()
	if err != nil {
		r.log.Error("Failed to consume JTI in Redis", zap.String("jti", jti.Id), zap.Error(err))
		return false, err
	}
	return ok, nil
}
