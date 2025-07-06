package pty_sessions

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	prefix string
}

func NewRedisRepository(client *redis.Client, prefix string) *RedisRepository {
	return &RedisRepository{
		client: client,
		prefix: prefix,
	}
}

func (r *RedisRepository) key(id string) string {
	return r.prefix + ":" + id
}

func (r *RedisRepository) Save(session *Entity) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return r.client.Set(context.Background(), r.key(session.ID), data, 0).Err()
}

func (r *RedisRepository) Get(id string) (*Entity, error) {
	val, err := r.client.Get(context.Background(), r.key(id)).Result()
	if err == redis.Nil {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var sess Entity
	if err := json.Unmarshal([]byte(val), &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

func (r *RedisRepository) Delete(id string) error {
	return r.client.Del(context.Background(), r.key(id)).Err()
}
