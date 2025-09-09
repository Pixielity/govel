package redis

import (
	"context"
	"time"
)

/**
 * Redis Cache Interfaces
 * 
 * This file defines the contracts for Redis-based caching operations.
 * It provides a clean abstraction over Redis functionality with support for
 * various data types, expiration, and advanced operations.
 */

// CacheInterface defines the contract for cache operations
type CacheInterface interface {
	// Basic operations
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetBytes(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	
	// Advanced operations
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Increment(ctx context.Context, key string, delta int64) (int64, error)
	Decrement(ctx context.Context, key string, delta int64) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	
	// Multi-key operations
	MGet(ctx context.Context, keys ...string) ([]string, error)
	MSet(ctx context.Context, pairs map[string]interface{}) error
	DeletePattern(ctx context.Context, pattern string) (int64, error)
	
	// List operations
	ListPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	ListPop(ctx context.Context, key string) (string, error)
	ListRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ListLength(ctx context.Context, key string) (int64, error)
	
	// Hash operations
	HashSet(ctx context.Context, key, field string, value interface{}) error
	HashGet(ctx context.Context, key, field string) (string, error)
	HashGetAll(ctx context.Context, key string) (map[string]string, error)
	HashDelete(ctx context.Context, key string, fields ...string) (int64, error)
	HashExists(ctx context.Context, key, field string) (bool, error)
	
	// Set operations
	SetAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SetRemove(ctx context.Context, key string, members ...interface{}) (int64, error)
	SetMembers(ctx context.Context, key string) ([]string, error)
	SetIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	
	// Connection and health
	Ping(ctx context.Context) error
	FlushDB(ctx context.Context) error
	Close() error
	Stats() ConnectionStats
}

// ConnectionManagerInterface manages multiple Redis connections
type ConnectionManagerInterface interface {
	GetConnection(name string) (CacheInterface, error)
	CreateConnection(name string, config ConnectionConfig) (CacheInterface, error)
	CloseAllConnections() error
	HealthCheck(ctx context.Context) map[string]bool
	ListConnections() []string
}

// ConnectionConfig holds Redis connection configuration
type ConnectionConfig struct {
	// Basic connection
	Host     string
	Port     int
	Password string
	Database int
	
	// Connection pool settings
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	
	// Timeouts
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolTimeout  time.Duration
	
	// TLS configuration
	TLSConfig *TLSConfig
	
	// Sentinel configuration for high availability
	SentinelConfig *SentinelConfig
	
	// Cluster configuration
	ClusterConfig *ClusterConfig
}

// TLSConfig holds TLS configuration for secure connections
type TLSConfig struct {
	Enabled            bool
	InsecureSkipVerify bool
	CertFile           string
	KeyFile            string
	CAFile             string
}

// SentinelConfig holds Redis Sentinel configuration
type SentinelConfig struct {
	Enabled       bool
	MasterName    string
	SentinelAddrs []string
	SentinelPassword string
}

// ClusterConfig holds Redis Cluster configuration
type ClusterConfig struct {
	Enabled   bool
	Addrs     []string
	MaxRedirects int
	ReadOnly  bool
}

// ConnectionStats provides connection statistics and metrics
type ConnectionStats struct {
	// Connection pool stats
	TotalConns   int
	IdleConns    int
	StaleConns   int
	
	// Operation stats
	Hits         int64
	Misses       int64
	Timeouts     int64
	
	// Network stats
	TotalCmds    int64
	TotalErrors  int64
	
	// Memory and performance
	UsedMemory   int64
	ConnectedClients int64
	
	// Replication info
	Role         string
	ConnectedSlaves int64
}

// CacheItem represents a cached item with metadata
type CacheItem struct {
	Key        string
	Value      interface{}
	Expiration time.Duration
	CreatedAt  time.Time
	AccessedAt time.Time
	TTL        time.Duration
}

// CachePattern defines cache key patterns for bulk operations
type CachePattern struct {
	Pattern string
	Prefix  string
	Suffix  string
}

// Pipeline represents a Redis pipeline for batched operations
type Pipeline interface {
	Set(key string, value interface{}, expiration time.Duration) *PipelineResult
	Get(key string) *PipelineResult
	Delete(key string) *PipelineResult
	Increment(key string, delta int64) *PipelineResult
	Execute(ctx context.Context) error
	Discard()
	Length() int
}

// PipelineResult represents a result from a pipeline operation
type PipelineResult struct {
	Key   string
	Value interface{}
	Error error
}

// TransactionInterface provides Redis transaction support
type TransactionInterface interface {
	Multi() error
	Watch(keys ...string) error
	Unwatch() error
	Exec(ctx context.Context) ([]interface{}, error)
	Discard() error
}
