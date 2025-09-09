package postgres

import (
	"context"
	"fmt"
	"sync"
)

// ConnectionManager manages multiple database connections
type ConnectionManager struct {
	connections map[string]DatabaseInterface
	configs     map[string]DatabaseConfig
	mu          sync.RWMutex
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]DatabaseInterface),
		configs:     make(map[string]DatabaseConfig),
	}
}

// CreateConnection creates and registers a new database connection
func (cm *ConnectionManager) CreateConnection(config DatabaseConfig) (DatabaseInterface, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Generate connection name if not provided
	name := fmt.Sprintf("%s_%s_%d", config.Host, config.Database, config.Port)

	// Check if connection already exists
	if existing, exists := cm.connections[name]; exists {
		return existing, nil
	}

	// Create new database instance
	db := NewDatabase(config)

	// Store the connection and config
	cm.connections[name] = db
	cm.configs[name] = config

	return db, nil
}

// GetConnection retrieves a connection by name
func (cm *ConnectionManager) GetConnection(name string) (DatabaseInterface, error) {
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
		if err := conn.Disconnect(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close connection '%s': %w", name, err))
		}
	}

	// Clear the connections map
	cm.connections = make(map[string]DatabaseInterface)
	cm.configs = make(map[string]DatabaseConfig)

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
func (cm *ConnectionManager) GetConnectionStats(name string) (DatabaseStats, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if conn, exists := cm.connections[name]; exists {
		return conn.Stats(), nil
	}

	return DatabaseStats{}, fmt.Errorf("connection '%s' not found", name)
}

// GetAllConnectionStats returns statistics for all connections
func (cm *ConnectionManager) GetAllConnectionStats() map[string]DatabaseStats {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := make(map[string]DatabaseStats)
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

// ConnectAll connects all registered databases
func (cm *ConnectionManager) ConnectAll(ctx context.Context) map[string]error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	errors := make(map[string]error)
	for name, conn := range cm.connections {
		if err := conn.Connect(ctx); err != nil {
			errors[name] = err
		}
	}

	return errors
}
