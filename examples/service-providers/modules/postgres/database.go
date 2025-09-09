package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Database implements the DatabaseInterface for PostgreSQL
type Database struct {
	config    DatabaseConfig
	db        *sql.DB
	mu        sync.RWMutex
	connected bool
}

// NewDatabase creates a new PostgreSQL database instance
func NewDatabase(config DatabaseConfig) *Database {
	return &Database{
		config: config,
	}
}

// Connect establishes a connection to the PostgreSQL database
func (d *Database) Connect(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.connected && d.db != nil {
		return nil
	}

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		d.config.Host, d.config.Port, d.config.Database,
		d.config.Username, d.config.Password, d.config.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(d.config.MaxOpenConns)
	db.SetMaxIdleConns(d.config.MaxIdleConns)
	db.SetConnMaxLifetime(d.config.ConnMaxLifetime)

	// Test the connection with context timeout
	pingCtx, cancel := context.WithTimeout(ctx, d.config.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	d.db = db
	d.connected = true
	return nil
}

// Disconnect closes the database connection
func (d *Database) Disconnect() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.db != nil {
		err := d.db.Close()
		d.db = nil
		d.connected = false
		return err
	}
	return nil
}

// Ping tests the database connection
func (d *Database) Ping(ctx context.Context) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.db == nil {
		return fmt.Errorf("database not connected")
	}

	return d.db.PingContext(ctx)
}

// IsConnected returns true if the database is connected
func (d *Database) IsConnected() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.connected && d.db != nil
}

// Query executes a query and returns results
func (d *Database) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.db == nil {
		return nil, fmt.Errorf("database not connected")
	}

	// Set query timeout
	queryCtx, cancel := context.WithTimeout(ctx, d.config.QueryTimeout)
	defer cancel()

	rows, err := d.db.QueryContext(queryCtx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		// Create slice of interface{} for scanning
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert to map
		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

// Execute runs a query without returning results
func (d *Database) Execute(ctx context.Context, query string, args ...interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.db == nil {
		return fmt.Errorf("database not connected")
	}

	// Set query timeout
	queryCtx, cancel := context.WithTimeout(ctx, d.config.QueryTimeout)
	defer cancel()

	_, err := d.db.ExecContext(queryCtx, query, args...)
	if err != nil {
		return fmt.Errorf("execute failed: %w", err)
	}

	return nil
}

// Transaction executes a function within a database transaction
func (d *Database) Transaction(ctx context.Context, fn func(tx TransactionInterface) error) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.db == nil {
		return fmt.Errorf("database not connected")
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	transaction := &Transaction{tx: tx}
	
	if err := fn(transaction); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Stats returns database connection statistics
func (d *Database) Stats() DatabaseStats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.db == nil {
		return DatabaseStats{}
	}

	stats := d.db.Stats()
	return DatabaseStats{
		OpenConnections:     stats.OpenConnections,
		InUseConnections:    stats.InUse,
		IdleConnections:     stats.Idle,
		WaitCount:          stats.WaitCount,
		WaitDuration:       stats.WaitDuration,
		MaxIdleClosed:      stats.MaxIdleClosed,
		MaxLifetimeClosed:  stats.MaxLifetimeClosed,
		MaxOpenConnections: d.config.MaxOpenConns,
		MaxIdleConnections: d.config.MaxIdleConns,
		ConnMaxLifetime:    d.config.ConnMaxLifetime,
	}
}

// SetMaxOpenConns sets the maximum number of open connections
func (d *Database) SetMaxOpenConns(n int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.config.MaxOpenConns = n
	if d.db != nil {
		d.db.SetMaxOpenConns(n)
	}
}

// SetMaxIdleConns sets the maximum number of idle connections
func (d *Database) SetMaxIdleConns(n int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.config.MaxIdleConns = n
	if d.db != nil {
		d.db.SetMaxIdleConns(n)
	}
}

// SetConnMaxLifetime sets the maximum lifetime of connections
func (d *Database) SetConnMaxLifetime(d time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.config.ConnMaxLifetime = d
	if d.db != nil {
		d.db.SetConnMaxLifetime(d)
	}
}

// Transaction implements TransactionInterface
type Transaction struct {
	tx *sql.Tx
}

// Query executes a query within the transaction
func (t *Transaction) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("transaction query failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, rows.Err()
}

// Execute runs a query within the transaction
func (t *Transaction) Execute(query string, args ...interface{}) error {
	_, err := t.tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("transaction execute failed: %w", err)
	}
	return nil
}

// Commit commits the transaction
func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}
