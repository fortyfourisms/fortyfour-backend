package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisInterface defines the contract for Redis operations
type RedisInterface interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
	Scan(pattern string) ([]string, error)
	Close() error
}

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisClient(cfg RedisConfig) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a key-value pair with expiration
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get retrieves a value by key
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Delete removes a key
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// Exists checks if a key exists
func (r *RedisClient) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	return result > 0, err
}

// Scan keys
func (r *RedisClient) Scan(pattern string) ([]string, error) {
	var keys []string
	iter := r.client.Scan(r.ctx, 0, pattern, 0).Iterator()
	for iter.Next(r.ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}
