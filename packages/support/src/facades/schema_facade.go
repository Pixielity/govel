package facades

import (
	schemaInterfaces "govel/types/interfaces/schema"
	facade "govel/support"
)

// Schema provides a clean, static-like interface to the application's database schema management service.
//
// This facade implements the facade pattern, providing global access to the schema
// service configured in the dependency injection container. It offers a Laravel-style
// API for database schema operations, migrations, table management, and database
// structure modifications with automatic service resolution and type safety.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved schema service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent schema operations across goroutines
//   - Supports multiple database drivers and schema formats
//   - Built-in migration tracking and rollback functionality
//
// Behavior:
//   - First call: Resolves schema service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if schema service cannot be resolved (fail-fast behavior)
//   - Automatically handles schema compilation, validation, and execution
//
// Returns:
//   - SchemaInterface: The application's schema management service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "schema" service is not registered in the container
//   - If the resolved service doesn't implement SchemaInterface
//   - If container resolution fails for any reason
//
// Performance Characteristics:
//   - First call: ~100-1000ns (depending on container and service complexity)
//   - Subsequent calls: ~10-50ns (cached lookup with atomic operations)
//   - Memory: Minimal overhead, shared cache across all facade calls
//   - Concurrency: Optimized read-write locks minimize contention
//
// Thread Safety:
// This facade is completely thread-safe:
//   - Multiple goroutines can call Schema() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Schema operations are thread-safe and transactional
//
// Usage Examples:
//
//	// Table creation and structure
//	facades.Schema().Create("users", func(table *Table) {
//	    table.ID()
//	    table.String("name", 100)
//	    table.String("email", 150).Unique()
//	    table.Timestamp("email_verified_at").Nullable()
//	    table.String("password")
//	    table.RememberToken()
//	    table.Timestamps()
//	})
//
//	facades.Schema().Create("posts", func(table *Table) {
//	    table.ID()
//	    table.String("title")
//	    table.Text("content")
//	    table.String("slug").Unique()
//	    table.ForeignID("user_id").Constrained()
//	    table.Boolean("published").Default(false)
//	    table.Timestamps()
//	    table.SoftDeletes()
//	})
//
//	// Table modifications
//	facades.Schema().Table("users", func(table *Table) {
//	    table.String("phone", 20).After("email")
//	    table.Date("birth_date").Nullable()
//	    table.Index(["name", "email"])
//	})
//
//	facades.Schema().Table("posts", func(table *Table) {
//	    table.Text("excerpt").After("title")
//	    table.Integer("view_count").Default(0)
//	    table.DropColumn("old_column")
//	    table.RenameColumn("summary", "description")
//	})
//
//	// Column types and constraints
//	facades.Schema().Create("products", func(table *Table) {
//	    table.ID()
//	    table.String("name", 200)
//	    table.Decimal("price", 8, 2)
//	    table.Integer("stock_quantity")
//	    table.Boolean("active").Default(true)
//	    table.JSON("metadata")
//	    table.Enum("status", []string{"draft", "active", "archived"})
//	    table.UUID("external_id").Nullable()
//	    table.Timestamps()
//	})
//
//	// Foreign key relationships
//	facades.Schema().Create("orders", func(table *Table) {
//	    table.ID()
//	    table.ForeignID("user_id").Constrained().OnDelete("cascade")
//	    table.ForeignID("product_id").Constrained("products")
//	    table.Integer("quantity")
//	    table.Decimal("total_amount", 10, 2)
//	    table.Timestamp("ordered_at")
//	    table.Timestamps()
//	})
//
//	// Indexes and performance optimization
//	facades.Schema().Table("posts", func(table *Table) {
//	    table.Index("title") // Single column index
//	    table.Index(["user_id", "created_at"]) // Composite index
//	    table.Unique(["user_id", "slug"]) // Unique composite index
//	    table.SpatialIndex("location") // Spatial index for geographic data
//	    table.FullText(["title", "content"]) // Full-text search index
//	})
//
//	// Index management
//	facades.Schema().Table("users", func(table *Table) {
//	    table.DropIndex("users_email_index")
//	    table.DropUnique("users_username_unique")
//	    table.RenameIndex("old_index", "new_index")
//	})
//
//	// Table operations
//	if facades.Schema().HasTable("users") {
//	    log.Println("Users table exists")
//	}
//
//	if facades.Schema().HasColumn("users", "email") {
//	    log.Println("Email column exists in users table")
//	}
//
//	facades.Schema().Drop("old_table")
//	facades.Schema().DropIfExists("temp_table")
//	facades.Schema().Rename("old_name", "new_name")
//
//	// Database information
//	tables := facades.Schema().GetTables()
//	columns := facades.Schema().GetColumns("users")
//	indexes := facades.Schema().GetIndexes("posts")
//
//	// Migration operations
//	facades.Schema().EnableForeignKeyConstraints()
//	facades.Schema().DisableForeignKeyConstraints()
//	facades.Schema().DropAllTables()
//	facades.Schema().DropAllViews()
//
//	// Views and stored procedures
//	facades.Schema().CreateView("active_users", `
//	    SELECT id, name, email FROM users WHERE active = true
//	`)
//
//	facades.Schema().DropView("old_view")
//
// Advanced Schema Patterns:
//
//	// Polymorphic relationships
//	facades.Schema().Create("comments", func(table *Table) {
//	    table.ID()
//	    table.Text("content")
//	    table.MorphsTo("commentable") // Creates commentable_type and commentable_id
//	    table.ForeignID("user_id").Constrained()
//	    table.Timestamps()
//	})
//
//	// Pivot tables for many-to-many relationships
//	facades.Schema().Create("user_roles", func(table *Table) {
//	    table.ID()
//	    table.ForeignID("user_id").Constrained().OnDelete("cascade")
//	    table.ForeignID("role_id").Constrained().OnDelete("cascade")
//	    table.Timestamp("assigned_at").UseCurrent()
//	    table.Unique(["user_id", "role_id"])
//	})
//
//	// Conditional schema modifications
//	func MigrateUserTable() {
//	    if !facades.Schema().HasColumn("users", "avatar") {
//	        facades.Schema().Table("users", func(table *Table) {
//	            table.String("avatar").Nullable().After("email")
//	        })
//	    }
//
//	    if facades.Config().GetBool("features.user_preferences") {
//	        facades.Schema().Table("users", func(table *Table) {
//	            table.JSON("preferences").Default("{}")
//	        })
//	    }
//	}
//
//	// Database-specific features
//	switch facades.DB().GetDriverName() {
//	case "postgres":
//	    facades.Schema().Create("locations", func(table *Table) {
//	        table.ID()
//	        table.String("name")
//	        table.Point("coordinates") // PostGIS point type
//	        table.Geometry("boundary") // PostGIS geometry
//	    })
//
//	case "mysql":
//	    facades.Schema().Create("logs", func(table *Table) {
//	        table.ID()
//	        table.String("level")
//	        table.LongText("message")
//	        table.Timestamp("created_at").UseCurrent()
//	    })
//	}
//
//	// Schema blueprints and reusable patterns
//	type UserTableBlueprint struct{}
//
//	func (b UserTableBlueprint) Apply(table *Table) {
//	    table.ID()
//	    table.String("name", 100)
//	    table.String("email", 150).Unique()
//	    table.Timestamp("email_verified_at").Nullable()
//	    table.String("password")
//	    table.RememberToken()
//	    table.Timestamps()
//	    table.SoftDeletes()
//	}
//
//	facades.Schema().Create("users", UserTableBlueprint{}.Apply)
//	facades.Schema().Create("admins", UserTableBlueprint{}.Apply)
//
// Migration Patterns:
//
//	// Safe column additions
//	func AddColumnSafely(tableName, columnName string, definition func(*Table)) {
//	    if !facades.Schema().HasColumn(tableName, columnName) {
//	        facades.Schema().Table(tableName, definition)
//	    }
//	}
//
//	// Batch schema operations
//	func CreateApplicationTables() {
//	    tables := map[string]func(*Table){
//	        "users":    createUsersTable,
//	        "posts":    createPostsTable,
//	        "comments": createCommentsTable,
//	    }
//
//	    for tableName, definition := range tables {
//	        if !facades.Schema().HasTable(tableName) {
//	            facades.Schema().Create(tableName, definition)
//	        }
//	    }
//	}
//
// Best Practices:
//   - Always check for table/column existence before modifications
//   - Use appropriate data types for optimal storage and performance
//   - Create indexes on frequently queried columns
//   - Use foreign key constraints to maintain data integrity
//   - Name indexes and constraints consistently
//   - Use migrations for schema version control
//   - Test schema changes in development before production
//   - Consider database-specific features when available
//
// Schema Design Principles:
//  1. Normalize data to reduce redundancy
//  2. Use appropriate data types and constraints
//  3. Plan for scalability with proper indexing
//  4. Maintain referential integrity with foreign keys
//  5. Document complex schema decisions
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume schema operations always work
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	schema, err := facade.TryResolve[SchemaInterface]("schema")
//	if err != nil {
//	    // Handle schema unavailability gracefully
//	    return fmt.Errorf("schema service unavailable: %w", err)
//	}
//	schema.Create("table", definition)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestTableCreation(t *testing.T) {
//	    // Create a test schema builder
//	    testSchema := &TestSchema{
//	        tables: make(map[string]*TableDefinition),
//	    }
//
//	    // Swap the real schema with test schema
//	    restore := support.SwapService("schema", testSchema)
//	    defer restore() // Always restore after test
//
//	    // Now facades.Schema() returns testSchema
//	    facades.Schema().Create("test_table", func(table *Table) {
//	        table.ID()
//	        table.String("name")
//	    })
//
//	    // Verify table creation
//	    assert.True(t, testSchema.HasTable("test_table"))
//	    table := testSchema.GetTable("test_table")
//	    assert.True(t, table.HasColumn("id"))
//	    assert.True(t, table.HasColumn("name"))
//	}
//
// Container Configuration:
// Ensure the schema service is properly configured in your container:
//
//	// Example schema registration
//	container.Singleton("schema", func() interface{} {
//	    config := schema.Config{
//	        // Database connection
//	        Connection: facades.DB().Connection("default"),
//
//	        // Schema configuration
//	        DefaultStringLength: 255,
//	        DefaultDecimalPrecision: 8,
//	        DefaultDecimalScale: 2,
//
//	        // Migration settings
//	        MigrationTable: "migrations",
//	        MigrationPath:  "/database/migrations",
//
//	        // Schema caching
//	        CacheSchema: facades.App().IsProduction(),
//	        CacheTTL:    30 * time.Minute,
//
//	        // Database-specific settings
//	        EnableForeignKeys: true,
//	        CheckConstraints:  true,
//	        EnableUUID:        true,
//	    }
//
//	    return schema.NewSchemaBuilder(config)
//	})
func Schema() schemaInterfaces.SchemaInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "schema" service from the dependency injection container
	// - Performs type assertion to SchemaInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[schemaInterfaces.SchemaInterface](schemaInterfaces.SCHEMA_TOKEN)
}

// SchemaWithError provides error-safe access to the schema management service.
//
// This function offers the same functionality as Schema() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle schema service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as Schema() but with error handling.
//
// Returns:
//   - SchemaInterface: The resolved schema instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement SchemaInterface
//
// Usage Examples:
//
//	// Basic error-safe schema operations
//	schema, err := facades.SchemaWithError()
//	if err != nil {
//	    log.Printf("Schema service unavailable: %v", err)
//	    return fmt.Errorf("database schema operations not available")
//	}
//	schema.Create("test_table", func(table *Table) {
//	    table.ID()
//	    table.String("name")
//	})
//
//	// Conditional schema modifications
//	if schema, err := facades.SchemaWithError(); err == nil {
//	    if !schema.HasTable("optional_table") {
//	        schema.Create("optional_table", tableDefinition)
//	    }
//	}
func SchemaWithError() (schemaInterfaces.SchemaInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "schema" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[schemaInterfaces.SchemaInterface](schemaInterfaces.SCHEMA_TOKEN)
}
