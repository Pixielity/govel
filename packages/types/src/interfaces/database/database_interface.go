package interfaces

import (
	"context"
	"database/sql"
)

// DatabaseInterface defines the contract for database operations in the GoVel framework.
// This interface provides comprehensive database functionality including connections,
// transactions, query building, and result handling.
//
// The interface supports multiple database drivers and provides both raw SQL execution
// and query builder patterns commonly used in Laravel-style applications.
type DatabaseInterface interface {
	// Connection Management
	
	// Connection returns the underlying database connection
	Connection() *sql.DB
	
	// SetConnection sets the database connection
	SetConnection(conn *sql.DB)
	
	// Close closes the database connection
	Close() error
	
	// Ping verifies the database connection is alive
	Ping() error
	
	// PingContext verifies the database connection is alive with context
	PingContext(ctx context.Context) error
	
	// Database Information
	
	// DatabaseName returns the name of the current database
	DatabaseName() string
	
	// DriverName returns the name of the database driver
	DriverName() string
	
	// Raw SQL Execution
	
	// Exec executes a raw SQL statement without returning rows
	Exec(query string, args ...interface{}) (sql.Result, error)
	
	// ExecContext executes a raw SQL statement with context
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	
	// Query executes a raw SQL query that returns rows
	Query(query string, args ...interface{}) (*sql.Rows, error)
	
	// QueryContext executes a raw SQL query with context
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	
	// QueryRow executes a raw SQL query that returns at most one row
	QueryRow(query string, args ...interface{}) *sql.Row
	
	// QueryRowContext executes a raw SQL query that returns at most one row with context
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	
	// Prepared Statements
	
	// Prepare creates a prepared statement
	Prepare(query string) (*sql.Stmt, error)
	
	// PrepareContext creates a prepared statement with context
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	
	// Transaction Management
	
	// Begin starts a new transaction
	Begin() (*sql.Tx, error)
	
	// BeginTx starts a new transaction with options
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	
	// Transaction executes a function within a transaction
	Transaction(fn func(*sql.Tx) error) error
	
	// TransactionContext executes a function within a transaction with context
	TransactionContext(ctx context.Context, opts *sql.TxOptions, fn func(*sql.Tx) error) error
	
	// Query Builder Interface (Laravel-style)
	
	// Table returns a new query builder for the specified table
	Table(name string) QueryBuilderInterface
	
	// Select starts a SELECT query
	Select(columns ...string) QueryBuilderInterface
	
	// From specifies the table for the query
	From(table string) QueryBuilderInterface
	
	// Where adds a WHERE condition
	Where(column string, operator string, value interface{}) QueryBuilderInterface
	
	// WhereIn adds a WHERE IN condition
	WhereIn(column string, values []interface{}) QueryBuilderInterface
	
	// Join adds an INNER JOIN
	Join(table string, first string, operator string, second string) QueryBuilderInterface
	
	// LeftJoin adds a LEFT JOIN
	LeftJoin(table string, first string, operator string, second string) QueryBuilderInterface
	
	// OrderBy adds an ORDER BY clause
	OrderBy(column string, direction string) QueryBuilderInterface
	
	// GroupBy adds a GROUP BY clause
	GroupBy(columns ...string) QueryBuilderInterface
	
	// Having adds a HAVING clause
	Having(column string, operator string, value interface{}) QueryBuilderInterface
	
	// Limit adds a LIMIT clause
	Limit(count int) QueryBuilderInterface
	
	// Offset adds an OFFSET clause
	Offset(count int) QueryBuilderInterface
	
	// Data Manipulation
	
	// Insert inserts a new record
	Insert(table string, data map[string]interface{}) (sql.Result, error)
	
	// InsertContext inserts a new record with context
	InsertContext(ctx context.Context, table string, data map[string]interface{}) (sql.Result, error)
	
	// Update updates existing records
	Update(table string, data map[string]interface{}, where map[string]interface{}) (sql.Result, error)
	
	// UpdateContext updates existing records with context
	UpdateContext(ctx context.Context, table string, data map[string]interface{}, where map[string]interface{}) (sql.Result, error)
	
	// Delete deletes records
	Delete(table string, where map[string]interface{}) (sql.Result, error)
	
	// DeleteContext deletes records with context
	DeleteContext(ctx context.Context, table string, where map[string]interface{}) (sql.Result, error)
	
	// Schema Operations
	
	// TableExists checks if a table exists
	TableExists(name string) (bool, error)
	
	// TableExistsContext checks if a table exists with context
	TableExistsContext(ctx context.Context, name string) (bool, error)
	
	// GetTableNames returns all table names in the database
	GetTableNames() ([]string, error)
	
	// GetTableNamesContext returns all table names with context
	GetTableNamesContext(ctx context.Context) ([]string, error)
	
	// GetColumns returns column information for a table
	GetColumns(table string) ([]ColumnInfo, error)
	
	// GetColumnsContext returns column information for a table with context
	GetColumnsContext(ctx context.Context, table string) ([]ColumnInfo, error)
	
	// Utility Methods
	
	// Quote quotes an identifier (table name, column name, etc.)
	Quote(identifier string) string
	
	// QuoteValue quotes a value for safe inclusion in SQL
	QuoteValue(value interface{}) string
	
	// Escape escapes special characters in a string
	Escape(value string) string
	
	// GetLastInsertID returns the last insert ID
	GetLastInsertID() (int64, error)
	
	// GetRowsAffected returns the number of rows affected by the last operation
	GetRowsAffected() (int64, error)
	
	// Statistics and Information
	
	// Stats returns database statistics
	Stats() sql.DBStats
	
	// SetMaxOpenConns sets the maximum number of open connections
	SetMaxOpenConns(n int)
	
	// SetMaxIdleConns sets the maximum number of idle connections
	SetMaxIdleConns(n int)
	
	// SetConnMaxLifetime sets the maximum connection lifetime
	SetConnMaxLifetime(d interface{}) // time.Duration
	
	// Configuration and Options
	
	// GetConfig returns the database configuration
	GetConfig() map[string]interface{}
	
	// SetConfig sets the database configuration
	SetConfig(config map[string]interface{})
	
	// GetConnectionInfo returns connection information
	GetConnectionInfo() ConnectionInfo
}

// QueryBuilderInterface defines the contract for query building operations
type QueryBuilderInterface interface {
	// Query Building
	Select(columns ...string) QueryBuilderInterface
	From(table string) QueryBuilderInterface
	Where(column string, operator string, value interface{}) QueryBuilderInterface
	WhereIn(column string, values []interface{}) QueryBuilderInterface
	WhereNull(column string) QueryBuilderInterface
	WhereNotNull(column string) QueryBuilderInterface
	Join(table string, first string, operator string, second string) QueryBuilderInterface
	LeftJoin(table string, first string, operator string, second string) QueryBuilderInterface
	RightJoin(table string, first string, operator string, second string) QueryBuilderInterface
	OrderBy(column string, direction string) QueryBuilderInterface
	GroupBy(columns ...string) QueryBuilderInterface
	Having(column string, operator string, value interface{}) QueryBuilderInterface
	Limit(count int) QueryBuilderInterface
	Offset(count int) QueryBuilderInterface
	
	// Execution
	Get() (*sql.Rows, error)
	GetContext(ctx context.Context) (*sql.Rows, error)
	First() *sql.Row
	FirstContext(ctx context.Context) *sql.Row
	Count() (int64, error)
	CountContext(ctx context.Context) (int64, error)
	Exists() (bool, error)
	ExistsContext(ctx context.Context) (bool, error)
	
	// Data Manipulation
	Insert(data map[string]interface{}) (sql.Result, error)
	InsertContext(ctx context.Context, data map[string]interface{}) (sql.Result, error)
	Update(data map[string]interface{}) (sql.Result, error)
	UpdateContext(ctx context.Context, data map[string]interface{}) (sql.Result, error)
	Delete() (sql.Result, error)
	DeleteContext(ctx context.Context) (sql.Result, error)
	
	// SQL Generation
	ToSQL() (string, []interface{}, error)
	GetBindings() []interface{}
}

// ColumnInfo represents information about a database column
type ColumnInfo struct {
	Name         string
	Type         string
	Nullable     bool
	Default      interface{}
	MaxLength    int
	Precision    int
	Scale        int
	IsPrimaryKey bool
	IsUnique     bool
	IsAutoIncrement bool
}

// ConnectionInfo represents database connection information
type ConnectionInfo struct {
	Driver          string
	Host            string
	Port            int
	Database        string
	Username        string
	Charset         string
	Collation       string
	Timezone        string
	MaxConnections  int
	IdleConnections int
	ConnMaxLifetime interface{} // time.Duration
}