package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisCache(config *CacheConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, json, expiration).Err()
}

func (c *RedisCache) Get(key string, dest interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("key not found")
	}
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

func (c *RedisCache) DeletePattern(pattern string) error {
	iter := c.client.Scan(c.ctx, 0, pattern, 0).Iterator()
	for iter.Next(c.ctx) {
		if err := c.client.Del(c.ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func (c *RedisCache) Exists(key string) (bool, error) {
	result, err := c.client.Exists(c.ctx, key).Result()
	return result > 0, err
}

func (c *RedisCache) Increment(key string) (int64, error) {
	return c.client.Incr(c.ctx, key).Result()
}

func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(c.ctx, key, expiration).Err()
}

// Cache keys helpers
func IssuesCacheKey() string {
	return "issues:all"
}

func IssueCacheKey(id int64) string {
	return fmt.Sprintf("issue:%d", id)
}

func StatsCacheKey() string {
	return "stats:global"
}

func AnalyticsCacheKey() string {
	return "analytics:global"
}

func UserCacheKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

// Rate limiting helpers
func RateLimitKey(userID int64, endpoint string) string {
	return fmt.Sprintf("ratelimit:user:%d:%s", userID, endpoint)
}

func (c *RedisCache) CheckRateLimit(userID int64, endpoint string, limit int, window time.Duration) (bool, error) {
	key := RateLimitKey(userID, endpoint)

	count, err := c.Increment(key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		if err := c.Expire(key, window); err != nil {
			return false, err
		}
	}

	return count <= int64(limit), nil
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
