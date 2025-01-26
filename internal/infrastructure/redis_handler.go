// Package infrastructure provides concrete implementations of interfaces defined in domain.
package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tyriis/go-locking-service/internal/domain"
)

// RedisHandler implements lock storage using Redis.
type RedisHandler struct {
	client *redis.Client
	ctx    context.Context
	logger domain.Logger
	config domain.Config
}

// NewRedisHandler creates a new RedisHandler with the given configuration and logger.
func NewRedisHandler(config domain.Config, logger domain.Logger) *RedisHandler {
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
	if h.Ping() != nil {
		return fmt.Errorf("failed to connect to Redis")
	}
	return h.client.Set(h.ctx, h.config.Redis.Prefix+key, value, ttl).Err()
}

// Get retrieves a lock by key. If key is "*", returns all locks.
func (h *RedisHandler) Get(key string) ([]string, error) {
	if h.Ping() != nil {
		return nil, fmt.Errorf("failed to connect to Redis")
	}
	if key == "*" {
		// Get all keys with prefix
		cmd := h.client.Keys(h.ctx, h.config.Redis.Prefix+"*")
		if cmd.Err() != nil {
			return nil, cmd.Err()
		}
		keys := cmd.Val()
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
		switch err {
		case redis.Nil:
			return nil, nil
		default:
			return nil, err
		}
	}
	return []string{val}, nil
}

// Del removes a lock by key.
func (h *RedisHandler) Del(key string) error {
	if h.Ping() != nil {
		return fmt.Errorf("failed to connect to Redis")
	}
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

// Ping checks if the Redis server is accessible.
func (h *RedisHandler) Ping() error {
	return h.client.Ping(h.ctx).Err()
}

func (h *RedisHandler) Close() error {
	return h.client.Close()
}

// Count returns the number of locks stored in Redis.
func (h *RedisHandler) Count() (int, error) {
	var cursor uint64
	var count int

	for {
		var keys []string
		var err error
		keys, cursor, err = h.client.Scan(context.Background(), cursor, h.config.Redis.Prefix+"*", 1000).Result()
		if err != nil {
			return 0, err
		}

		count += len(keys)

		if cursor == 0 {
			break
		}
	}

	return count, nil
}
