package redis

import (
	"context"
	"fmt"
	"sync"
)

/**
 * Redis Connection Manager
 * 
 * This file provides connection management for multiple Redis instances,
 * supporting different Redis configurations (standalone, sentinel, cluster).
 */

// ConnectionManager manages multiple Redis cache connections
type ConnectionManager struct {
	connections map[string]CacheInterface
	configs     map[string]ConnectionConfig
	mu          sync.RWMutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]CacheInterface),
		configs:     make(map[string]ConnectionConfig),
	}
}

// CreateConnection creates a new Redis connection with the given configuration
func (cm *ConnectionManager) CreateConnection(name string, config ConnectionConfig) (CacheInterface, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Check if connection already exists
	if existing, exists := cm.connections[name]; exists {
		return existing, nil
	}

	// Set default values if not provided
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 6379
	}
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 2
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	// Create new cache instance
	cache, err := NewCache(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis cache for '%s': %w", name, err)
	}

	// Store the connection and config
	cm.connections[name] = cache
	cm.configs[name] = config

	return cache, nil
}

// GetConnection retrieves a connection by name
func (cm *ConnectionManager) GetConnection(name string) (CacheInterface, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if conn, exists := cm.connections[name]; exists {
		return conn, nil
	}

	return nil, fmt.Errorf("connection '%s' not found", name)
}

// CloseAllConnections closes all managed connections
func (cm *ConnectionManager) CloseAllConnections() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var errs []error
	for name, conn := range cm.connections {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection '%s': %w", name, err))
		}
	}

	// Clear the connections map
	cm.connections = make(map[string]CacheInterface)
	cm.configs = make(map[string]ConnectionConfig)

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}

	return nil
}

// HealthCheck performs health checks on all connections
func (cm *ConnectionManager) HealthCheck(ctx context.Context) map[string]bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	results := make(map[string]bool)
	for name, conn := range cm.connections {
		err := conn.Ping(ctx)
		results[name] = err == nil
	}

	return results
}

// ListConnections returns a list of connection names
func (cm *ConnectionManager) ListConnections() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	names := make([]string, 0, len(cm.connections))
	for name := range cm.connections {
		names = append(names, name)
	}

	return names
}

// GetConnectionStats returns statistics for a specific connection
func (cm *ConnectionManager) GetConnectionStats(name string) (ConnectionStats, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if conn, exists := cm.connections[name]; exists {
		return conn.Stats(), nil
	}

	return ConnectionStats{}, fmt.Errorf("connection '%s' not found", name)
}

// GetAllConnectionStats returns statistics for all connections
func (cm *ConnectionManager) GetAllConnectionStats() map[string]ConnectionStats {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := make(map[string]ConnectionStats)
	for name, conn := range cm.connections {
		stats[name] = conn.Stats()
	}

	return stats
}

// RemoveConnection removes a connection from management (but doesn't close it)
func (cm *ConnectionManager) RemoveConnection(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.connections[name]; !exists {
		return fmt.Errorf("connection '%s' not found", name)
	}

	delete(cm.connections, name)
	delete(cm.configs, name)
	return nil
}

// TestAllConnections tests connectivity to all managed connections
func (cm *ConnectionManager) TestAllConnections(ctx context.Context) map[string]error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	errors := make(map[string]error)
	for name, conn := range cm.connections {
		if err := conn.Ping(ctx); err != nil {
			errors[name] = err
		}
	}

	return errors
}

// FlushAllDatabases flushes all databases for all connections
func (cm *ConnectionManager) FlushAllDatabases(ctx context.Context) map[string]error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	errors := make(map[string]error)
	for name, conn := range cm.connections {
		if err := conn.FlushDB(ctx); err != nil {
			errors[name] = err
		}
	}

	return errors
}

// GetConnectionConfig returns the configuration for a specific connection
func (cm *ConnectionManager) GetConnectionConfig(name string) (ConnectionConfig, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if config, exists := cm.configs[name]; exists {
		return config, nil
	}

	return ConnectionConfig{}, fmt.Errorf("connection '%s' not found", name)
}

// UpdateConnectionConfig updates the configuration for an existing connection
// Note: This will recreate the connection with new settings
func (cm *ConnectionManager) UpdateConnectionConfig(name string, config ConnectionConfig) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Check if connection exists
	if _, exists := cm.connections[name]; !exists {
		return fmt.Errorf("connection '%s' not found", name)
	}

	// Close existing connection
	if err := cm.connections[name].Close(); err != nil {
		return fmt.Errorf("failed to close existing connection: %w", err)
	}

	// Create new connection with updated config
	cache, err := NewCache(config)
	if err != nil {
		return fmt.Errorf("failed to create new connection: %w", err)
	}

	// Update stored connection and config
	cm.connections[name] = cache
	cm.configs[name] = config

	return nil
}

// ConnectionInfo provides information about a connection
type ConnectionInfo struct {
	Name   string
	Config ConnectionConfig
	Stats  ConnectionStats
	Healthy bool
}

// GetConnectionInfo returns detailed information about a connection
func (cm *ConnectionManager) GetConnectionInfo(ctx context.Context, name string) (*ConnectionInfo, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conn, exists := cm.connections[name]
	if !exists {
		return nil, fmt.Errorf("connection '%s' not found", name)
	}

	config, exists := cm.configs[name]
	if !exists {
		return nil, fmt.Errorf("config for connection '%s' not found", name)
	}

	// Test health
	healthy := conn.Ping(ctx) == nil

	return &ConnectionInfo{
		Name:    name,
		Config:  config,
		Stats:   conn.Stats(),
		Healthy: healthy,
	}, nil
}

// GetAllConnectionInfo returns detailed information about all connections
func (cm *ConnectionManager) GetAllConnectionInfo(ctx context.Context) []*ConnectionInfo {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var infos []*ConnectionInfo
	for name := range cm.connections {
		if info, err := cm.GetConnectionInfo(ctx, name); err == nil {
			infos = append(infos, info)
		}
	}

	return infos
}
