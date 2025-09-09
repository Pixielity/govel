package postgres

import (
	"context"
	"time"
)

// DatabaseInterface defines the contract for database operations
type DatabaseInterface interface {
	// Connection management
	Connect(ctx context.Context) error
	Disconnect() error
	Ping(ctx context.Context) error
	IsConnected() bool

	// Query operations
	Query(ctx context.Context, sql string, args ...interface{}) ([]map[string]interface{}, error)
	Execute(ctx context.Context, sql string, args ...interface{}) error
	Transaction(ctx context.Context, fn func(tx TransactionInterface) error) error

	// Health and statistics
	Stats() DatabaseStats
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

// TransactionInterface defines operations available within a transaction
type TransactionInterface interface {
	Query(sql string, args ...interface{}) ([]map[string]interface{}, error)
	Execute(sql string, args ...interface{}) error
	Commit() error
	Rollback() error
}

// ConnectionManagerInterface manages database connection pools and health
type ConnectionManagerInterface interface {
	CreateConnection(config DatabaseConfig) (DatabaseInterface, error)
	GetConnection(name string) (DatabaseInterface, error)
	CloseAllConnections() error
	HealthCheck(ctx context.Context) map[string]bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnectTimeout  time.Duration
	QueryTimeout    time.Duration
}

// DatabaseStats provides database connection statistics
type DatabaseStats struct {
	OpenConnections     int
	InUseConnections    int
	IdleConnections     int
	WaitCount           int64
	WaitDuration        time.Duration
	MaxIdleClosed       int64
	MaxLifetimeClosed   int64
	MaxOpenConnections  int
	MaxIdleConnections  int
	ConnMaxLifetime     time.Duration
}

// QueryResult represents a query result with metadata
type QueryResult struct {
	Rows         []map[string]interface{}
	RowsAffected int64
	LastInsertID int64
	ExecutionTime time.Duration
}
