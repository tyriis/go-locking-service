package infrastructure

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHandler struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisHandler(addr string) *RedisHandler {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisHandler{
		client: client,
		ctx:    context.Background(),
	}
}

func (h *RedisHandler) Set(key string, value string, ttl time.Duration) error {
	return h.client.Set(h.ctx, key, value, ttl).Err()
}

func (h *RedisHandler) Get(key string) (string, error) {
	return h.client.Get(h.ctx, key).Result()
}

func (h *RedisHandler) Del(key string) error {
	return h.client.Del(h.ctx, key).Err()
}
