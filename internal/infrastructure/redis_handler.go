// Package infrastructure provides concrete implementations of interfaces defined in domain.
package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tyriis/rest-go/internal/domain"
)

// RedisHandler implements lock storage using Redis.
type RedisHandler struct {
	client *redis.Client
	ctx    context.Context
	logger domain.Logger
	config Config
}

// NewRedisHandler creates a new RedisHandler with the given configuration and logger.
func NewRedisHandler(config Config, logger domain.Logger) *RedisHandler {
	redisJSON, err := json.Marshal(config.Redis)
	if err == nil {
		const msg = "NewRedisHandler - json.Marshal > config.Redis > %s"
		logger.Debug(fmt.Sprintf(msg, string(redisJSON)))
	}
	client := redis.NewClient(&redis.Options{
		Addr: config.Redis.Host + ":" + config.Redis.Port,
	})
	return &RedisHandler{
		client: client,
		ctx:    context.Background(),
		logger: logger,
		config: config,
	}
}

// Set stores a lock with the given key, value, and TTL.
func (h *RedisHandler) Set(key string, value string, ttl time.Duration) error {
	return h.client.Set(h.ctx, h.config.Redis.Prefix+key, value, ttl).Err()
}

// Get retrieves a lock by key. If key is "*", returns all locks.
func (h *RedisHandler) Get(key string) ([]string, error) {
	if key == "*" {
		// Get all keys with prefix
		keys := h.client.Keys(h.ctx, h.config.Redis.Prefix+"*").Val()
		if len(keys) == 0 {
			return []string{}, nil
		}

		// Get all values
		values, err := h.GetMultiple(keys)
		if err != nil {
			return nil, fmt.Errorf("failed to get multiple keys: %w", err)
		}
		return values, nil
	}
	val, err := h.client.Get(h.ctx, h.config.Redis.Prefix+key).Result()
	if err != nil {
		return nil, err
	}
	return []string{val}, nil
}

// Del removes a lock by key.
func (h *RedisHandler) Del(key string) error {
	return h.client.Del(h.ctx, h.config.Redis.Prefix+key).Err()
}

// GetMultiple retrieves multiple locks by their keys.
func (h *RedisHandler) GetMultiple(keys []string) ([]string, error) {
	h.logger.Debug(fmt.Sprintf("RedisHandler.GetMultiple - keys: %v", keys))
	// Fetch multiple keys in one call
	results, err := h.client.MGet(h.ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("RedisHandler.GetMultiple - MGet failed: %w", err)
	}

	// Convert interface{} slice to string slice
	values := make([]string, 0, len(results))
	for _, result := range results {
		if result != nil {
			values = append(values, result.(string))
		}
	}

	h.logger.Debug(fmt.Sprintf("RedisHandler.GetMultiple - values: %v", values))
	return values, nil
}
