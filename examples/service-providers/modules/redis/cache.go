package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

/**
 * Redis Cache Implementation
 * 
 * This file provides a comprehensive Redis cache implementation with support for
 * connection pooling, various data types, and advanced Redis operations.
 */

// Cache implements the CacheInterface for Redis
type Cache struct {
	client   redis.UniversalClient
	config   ConnectionConfig
	mu       sync.RWMutex
	stats    ConnectionStats
	closed   bool
}

// NewCache creates a new Redis cache instance
func NewCache(config ConnectionConfig) (*Cache, error) {
	// Build Redis options based on configuration
	var client redis.UniversalClient
	
	if config.ClusterConfig != nil && config.ClusterConfig.Enabled {
		// Redis Cluster mode
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        config.ClusterConfig.Addrs,
			Password:     config.Password,
			MaxRedirects: config.ClusterConfig.MaxRedirects,
			ReadOnly:     config.ClusterConfig.ReadOnly,
			PoolSize:     config.PoolSize,
			MinIdleConns: config.MinIdleConns,
			MaxRetries:   config.MaxRetries,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolTimeout:  config.PoolTimeout,
		})
	} else if config.SentinelConfig != nil && config.SentinelConfig.Enabled {
		// Redis Sentinel mode
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       config.SentinelConfig.MasterName,
			SentinelAddrs:    config.SentinelConfig.SentinelAddrs,
			SentinelPassword: config.SentinelConfig.SentinelPassword,
			Password:         config.Password,
			DB:               config.Database,
			PoolSize:         config.PoolSize,
			MinIdleConns:     config.MinIdleConns,
			MaxRetries:       config.MaxRetries,
			DialTimeout:      config.DialTimeout,
			ReadTimeout:      config.ReadTimeout,
			WriteTimeout:     config.WriteTimeout,
			PoolTimeout:      config.PoolTimeout,
		})
	} else {
		// Standard single Redis instance
		addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
		client = redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     config.Password,
			DB:           config.Database,
			PoolSize:     config.PoolSize,
			MinIdleConns: config.MinIdleConns,
			MaxRetries:   config.MaxRetries,
			DialTimeout:  config.DialTimeout,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			PoolTimeout:  config.PoolTimeout,
		})
	}

	return &Cache{
		client: client,
		config: config,
		stats:  ConnectionStats{},
	}, nil
}

// Basic operations

// Set stores a value with expiration
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("cache connection is closed")
	}

	// Convert value to string
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	case int, int8, int16, int32, int64:
		strValue = fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		strValue = fmt.Sprintf("%d", v)
	case float32, float64:
		strValue = fmt.Sprintf("%f", v)
	case bool:
		strValue = fmt.Sprintf("%t", v)
	default:
		// JSON encode complex types
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		strValue = string(jsonBytes)
	}

	err := c.client.Set(ctx, key, strValue, expiration).Err()
	if err != nil {
		c.stats.TotalErrors++
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return nil
}

// Get retrieves a value as string
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return "", fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.stats.Misses++
			return "", fmt.Errorf("key not found: %s", key)
		}
		c.stats.TotalErrors++
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	c.stats.Hits++
	c.stats.TotalCmds++
	return result, nil
}

// GetBytes retrieves a value as bytes
func (c *Cache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	str, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// Delete removes a key
func (c *Cache) Delete(ctx context.Context, key string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("cache connection is closed")
	}

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.stats.TotalErrors++
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return nil
}

// Exists checks if a key exists
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return false, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		c.stats.TotalErrors++
		return false, fmt.Errorf("failed to check key existence %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result > 0, nil
}

// Advanced operations

// SetNX sets a key only if it doesn't exist
func (c *Cache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return false, fmt.Errorf("cache connection is closed")
	}

	strValue := fmt.Sprintf("%v", value)
	result, err := c.client.SetNX(ctx, key, strValue, expiration).Result()
	if err != nil {
		c.stats.TotalErrors++
		return false, fmt.Errorf("failed to setnx key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// Increment increments a key by delta
func (c *Cache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.IncrBy(ctx, key, delta).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// Decrement decrements a key by delta
func (c *Cache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return c.Increment(ctx, key, -delta)
}

// Expire sets expiration for a key
func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("cache connection is closed")
	}

	err := c.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		c.stats.TotalErrors++
		return fmt.Errorf("failed to set expiration for key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return nil
}

// TTL returns the time-to-live for a key
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// Multi-key operations

// MGet retrieves multiple keys
func (c *Cache) MGet(ctx context.Context, keys ...string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		c.stats.TotalErrors++
		return nil, fmt.Errorf("failed to mget keys: %w", err)
	}

	// Convert interface{} slice to string slice
	values := make([]string, len(result))
	for i, v := range result {
		if v != nil {
			values[i] = fmt.Sprintf("%v", v)
			c.stats.Hits++
		} else {
			c.stats.Misses++
		}
	}

	c.stats.TotalCmds++
	return values, nil
}

// MSet sets multiple key-value pairs
func (c *Cache) MSet(ctx context.Context, pairs map[string]interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("cache connection is closed")
	}

	// Convert map to slice of interface{}
	args := make([]interface{}, 0, len(pairs)*2)
	for key, value := range pairs {
		args = append(args, key, fmt.Sprintf("%v", value))
	}

	err := c.client.MSet(ctx, args...).Err()
	if err != nil {
		c.stats.TotalErrors++
		return fmt.Errorf("failed to mset: %w", err)
	}

	c.stats.TotalCmds++
	return nil
}

// DeletePattern deletes keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("cache connection is closed")
	}

	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to get keys for pattern %s: %w", pattern, err)
	}

	if len(keys) == 0 {
		return 0, nil
	}

	result, err := c.client.Del(ctx, keys...).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to delete keys: %w", err)
	}

	c.stats.TotalCmds += 2 // KEYS + DEL commands
	return result, nil
}

// List operations

// ListPush pushes values to a list
func (c *Cache) ListPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.LPush(ctx, key, values...).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to lpush to key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// ListPop pops a value from a list
func (c *Cache) ListPop(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return "", fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.LPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("list is empty: %s", key)
		}
		c.stats.TotalErrors++
		return "", fmt.Errorf("failed to lpop from key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// ListRange returns a range of elements from a list
func (c *Cache) ListRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.LRange(ctx, key, start, stop).Result()
	if err != nil {
		c.stats.TotalErrors++
		return nil, fmt.Errorf("failed to lrange key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// ListLength returns the length of a list
func (c *Cache) ListLength(ctx context.Context, key string) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.LLen(ctx, key).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to llen key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// Hash operations

// HashSet sets a field in a hash
func (c *Cache) HashSet(ctx context.Context, key, field string, value interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("cache connection is closed")
	}

	err := c.client.HSet(ctx, key, field, value).Err()
	if err != nil {
		c.stats.TotalErrors++
		return fmt.Errorf("failed to hset key %s field %s: %w", key, field, err)
	}

	c.stats.TotalCmds++
	return nil
}

// HashGet gets a field from a hash
func (c *Cache) HashGet(ctx context.Context, key, field string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return "", fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			c.stats.Misses++
			return "", fmt.Errorf("field not found: %s in key %s", field, key)
		}
		c.stats.TotalErrors++
		return "", fmt.Errorf("failed to hget key %s field %s: %w", key, field, err)
	}

	c.stats.Hits++
	c.stats.TotalCmds++
	return result, nil
}

// HashGetAll gets all fields from a hash
func (c *Cache) HashGetAll(ctx context.Context, key string) (map[string]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		c.stats.TotalErrors++
		return nil, fmt.Errorf("failed to hgetall key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	if len(result) > 0 {
		c.stats.Hits++
	} else {
		c.stats.Misses++
	}
	return result, nil
}

// HashDelete deletes fields from a hash
func (c *Cache) HashDelete(ctx context.Context, key string, fields ...string) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.HDel(ctx, key, fields...).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to hdel key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// HashExists checks if a field exists in a hash
func (c *Cache) HashExists(ctx context.Context, key, field string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return false, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.HExists(ctx, key, field).Result()
	if err != nil {
		c.stats.TotalErrors++
		return false, fmt.Errorf("failed to hexists key %s field %s: %w", key, field, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// Set operations

// SetAdd adds members to a set
func (c *Cache) SetAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.SAdd(ctx, key, members...).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to sadd key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// SetRemove removes members from a set
func (c *Cache) SetRemove(ctx context.Context, key string, members ...interface{}) (int64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return 0, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.SRem(ctx, key, members...).Result()
	if err != nil {
		c.stats.TotalErrors++
		return 0, fmt.Errorf("failed to srem key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// SetMembers returns all members of a set
func (c *Cache) SetMembers(ctx context.Context, key string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.SMembers(ctx, key).Result()
	if err != nil {
		c.stats.TotalErrors++
		return nil, fmt.Errorf("failed to smembers key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// SetIsMember checks if a member is in a set
func (c *Cache) SetIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return false, fmt.Errorf("cache connection is closed")
	}

	result, err := c.client.SIsMember(ctx, key, member).Result()
	if err != nil {
		c.stats.TotalErrors++
		return false, fmt.Errorf("failed to sismember key %s: %w", key, err)
	}

	c.stats.TotalCmds++
	return result, nil
}

// Connection and health methods

// Ping tests the connection
func (c *Cache) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("cache connection is closed")
	}

	err := c.client.Ping(ctx).Err()
	if err != nil {
		c.stats.TotalErrors++
		return fmt.Errorf("ping failed: %w", err)
	}

	c.stats.TotalCmds++
	return nil
}

// FlushDB flushes the current database
func (c *Cache) FlushDB(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("cache connection is closed")
	}

	err := c.client.FlushDB(ctx).Err()
	if err != nil {
		c.stats.TotalErrors++
		return fmt.Errorf("flushdb failed: %w", err)
	}

	c.stats.TotalCmds++
	return nil
}

// Close closes the connection
func (c *Cache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	err := c.client.Close()
	c.closed = true
	return err
}

// Stats returns connection statistics
func (c *Cache) Stats() ConnectionStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Get Redis server info if connection is active
	if !c.closed {
		if info, err := c.client.Info(context.Background()).Result(); err == nil {
			c.updateStatsFromInfo(info)
		}
	}

	return c.stats
}

// updateStatsFromInfo parses Redis INFO command output and updates stats
func (c *Cache) updateStatsFromInfo(info string) {
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			
			key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			
			switch key {
			case "used_memory":
				if mem, err := strconv.ParseInt(value, 10, 64); err == nil {
					c.stats.UsedMemory = mem
				}
			case "connected_clients":
				if clients, err := strconv.ParseInt(value, 10, 64); err == nil {
					c.stats.ConnectedClients = clients
				}
			case "role":
				c.stats.Role = value
			case "connected_slaves":
				if slaves, err := strconv.ParseInt(value, 10, 64); err == nil {
					c.stats.ConnectedSlaves = slaves
				}
			}
		}
	}
}
