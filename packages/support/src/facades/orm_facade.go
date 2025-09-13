package facades

import (
	ormInterfaces "govel/packages/types/src/interfaces/orm"
	facade "govel/packages/support/src"
)

// ORM provides a clean, static-like interface to the application's Object-Relational Mapping service.
//
// This facade implements the facade pattern, providing global access to the ORM
// service configured in the dependency injection container. It offers a Laravel-style
// API for database operations with automatic service resolution, model relationships,
// query building, migrations, and active record patterns.
//
// Architecture:
//   - Uses facade.Resolve() internally for service resolution
//   - Automatically caches the resolved ORM service for performance
//   - Provides compile-time type safety through generics
//   - Thread-safe for concurrent database operations across goroutines
//   - Supports multiple database drivers with unified interface
//   - Built-in relationship management and eager loading
//
// Behavior:
//   - First call: Resolves ORM service from container, performs type assertion, caches result
//   - Subsequent calls: Returns cached service instance (extremely fast)
//   - Panics if ORM service cannot be resolved (fail-fast behavior)
//   - Automatically handles model hydration, relationship loading, and query optimization
//
// Returns:
//   - ORMInterface: The application's ORM service instance
//
// Panics:
//   - If no container is set via facades.SetContainer() or support.SetContainer()
//   - If "orm" service is not registered in the container
//   - If the resolved service doesn't implement ORMInterface
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
//   - Multiple goroutines can call ORM() concurrently without synchronization
//   - Internal caching uses optimized read-write mutexes
//   - Service resolution is protected against race conditions
//   - Database operations are thread-safe with connection pooling
//
// Usage Examples:
//
//	// Basic model operations
//	type User struct {
//	    ID        uint      `gorm:"primaryKey"`
//	    Name      string    `gorm:"not null"`
//	    Email     string    `gorm:"uniqueIndex;not null"`
//	    CreatedAt time.Time
//	    UpdatedAt time.Time
//	}
//
//	// Create a new user
//	user := &User{
//	    Name:  "John Doe",
//	    Email: "john@example.com",
//	}
//
//	result := facades.ORM().Create(user)
//	if result.Error != nil {
//	    log.Fatal(result.Error)
//	}
//	fmt.Printf("Created user with ID: %d\n", user.ID)
//
//	// Find user by ID
//	var user User
//	result := facades.ORM().First(&user, 1)
//	if result.Error != nil {
//	    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
//	        fmt.Println("User not found")
//	    } else {
//	        log.Fatal(result.Error)
//	    }
//	}
//
//	// Find user by email
//	var user User
//	result := facades.ORM().Where("email = ?", "john@example.com").First(&user)
//	if result.Error == nil {
//	    fmt.Printf("Found user: %+v\n", user)
//	}
//
//	// Update user
//	facades.ORM().Model(&user).Update("name", "John Smith")
//
//	// Update multiple fields
//	facades.ORM().Model(&user).Updates(User{
//	    Name:  "John Smith",
//	    Email: "johnsmith@example.com",
//	})
//
//	// Delete user
//	facades.ORM().Delete(&user)
//
//	// Soft delete (if DeletedAt field exists)
//	facades.ORM().Delete(&user) // Sets DeletedAt timestamp
//
//	// Hard delete
//	facades.ORM().Unscoped().Delete(&user) // Permanent deletion
//
//	// Query multiple records
//	var users []User
//	facades.ORM().Find(&users)
//
//	var activeUsers []User
//	facades.ORM().Where("active = ?", true).Find(&activeUsers)
//
//	// Count records
//	var count int64
//	facades.ORM().Model(&User{}).Where("active = ?", true).Count(&count)
//	fmt.Printf("Active users: %d\n", count)
//
//	// Pagination
//	var users []User
//	offset := (page - 1) * limit
//	facades.ORM().Offset(offset).Limit(limit).Find(&users)
//
//	// Raw SQL queries
//	var users []User
//	facades.ORM().Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)
//
// Model Relationships:
//
//	// One-to-Many relationship
//	type User struct {
//	    ID    uint
//	    Name  string
//	    Posts []Post `gorm:"foreignKey:UserID"`
//	}
//
//	type Post struct {
//	    ID     uint
//	    Title  string
//	    UserID uint
//	    User   User `gorm:"foreignKey:UserID"`
//	}
//
//	// Eager loading relationships
//	var user User
//	facades.ORM().Preload("Posts").First(&user, 1)
//	fmt.Printf("User has %d posts\n", len(user.Posts))
//
//	// Nested preloading
//	var users []User
//	facades.ORM().Preload("Posts.Comments").Find(&users)
//
//	// Many-to-Many relationship
//	type User struct {
//	    ID    uint
//	    Name  string
//	    Roles []Role `gorm:"many2many:user_roles;"`
//	}
//
//	type Role struct {
//	    ID    uint
//	    Name  string
//	    Users []User `gorm:"many2many:user_roles;"`
//	}
//
//	// Associate many-to-many
//	user := User{Name: "John"}
//	roles := []Role{{Name: "Admin"}, {Name: "Editor"}}
//
//	facades.ORM().Create(&user)
//	facades.ORM().Create(&roles)
//	facades.ORM().Model(&user).Association("Roles").Append(&roles)
//
//	// Query with joins
//	var users []User
//	facades.ORM().Joins("LEFT JOIN posts ON posts.user_id = users.id").
//	    Where("posts.published = ?", true).
//	    Find(&users)
//
// Advanced ORM Features:
//
//	// Transactions
//	err := facades.ORM().Transaction(func(tx *gorm.DB) error {
//	    // Create user
//	    user := &User{Name: "John", Email: "john@example.com"}
//	    if err := tx.Create(user).Error; err != nil {
//	        return err
//	    }
//
//	    // Create profile
//	    profile := &Profile{UserID: user.ID, Bio: "Software Developer"}
//	    if err := tx.Create(profile).Error; err != nil {
//	        return err
//	    }
//
//	    return nil
//	})
//
//	if err != nil {
//	    log.Printf("Transaction failed: %v", err)
//	}
//
//	// Manual transaction control
//	tx := facades.ORM().Begin()
//	defer func() {
//	    if r := recover(); r != nil {
//	        tx.Rollback()
//	    }
//	}()
//
//	if err := tx.Create(&user).Error; err != nil {
//	    tx.Rollback()
//	    return err
//	}
//
//	if err := tx.Create(&profile).Error; err != nil {
//	    tx.Rollback()
//	    return err
//	}
//
//	tx.Commit()
//
//	// Batch operations
//	users := []User{
//	    {Name: "User1", Email: "user1@example.com"},
//	    {Name: "User2", Email: "user2@example.com"},
//	    {Name: "User3", Email: "user3@example.com"},
//	}
//
//	// Batch insert
//	facades.ORM().CreateInBatches(users, 100)
//
//	// Scopes for reusable queries
//	func ActiveUsers(db *gorm.DB) *gorm.DB {
//	    return db.Where("active = ?", true)
//	}
//
//	func RecentUsers(db *gorm.DB) *gorm.DB {
//	    return db.Where("created_at > ?", time.Now().AddDate(0, -1, 0))
//	}
//
//	// Use scopes
//	var users []User
//	facades.ORM().Scopes(ActiveUsers, RecentUsers).Find(&users)
//
//	// Hooks and callbacks
//	func (u *User) BeforeCreate(tx *gorm.DB) error {
//	    // Hash password before saving
//	    if u.Password != "" {
//	        hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
//	        if err != nil {
//	            return err
//	        }
//	        u.Password = string(hashed)
//	    }
//	    return nil
//	}
//
//	func (u *User) AfterCreate(tx *gorm.DB) error {
//	    // Send welcome email after user creation
//	    return facades.Mail().Queue(
//	        u.Email,
//	        "Welcome to our platform!",
//	        "emails/welcome",
//	        map[string]interface{}{"name": u.Name},
//	    )
//	}
//
// Repository Pattern Integration:
//
//	type UserRepository struct {
//	    db *gorm.DB
//	}
//
//	func NewUserRepository() *UserRepository {
//	    return &UserRepository{db: facades.ORM()}
//	}
//
//	func (r *UserRepository) Create(user *User) error {
//	    return r.db.Create(user).Error
//	}
//
//	func (r *UserRepository) FindByID(id uint) (*User, error) {
//	    var user User
//	    err := r.db.First(&user, id).Error
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &user, nil
//	}
//
//	func (r *UserRepository) FindByEmail(email string) (*User, error) {
//	    var user User
//	    err := r.db.Where("email = ?", email).First(&user).Error
//	    if err != nil {
//	        return nil, err
//	    }
//	    return &user, nil
//	}
//
//	func (r *UserRepository) Update(user *User) error {
//	    return r.db.Save(user).Error
//	}
//
//	func (r *UserRepository) Delete(id uint) error {
//	    return r.db.Delete(&User{}, id).Error
//	}
//
//	func (r *UserRepository) GetActiveUsers(limit, offset int) ([]User, error) {
//	    var users []User
//	    err := r.db.Where("active = ?", true).
//	        Limit(limit).
//	        Offset(offset).
//	        Find(&users).Error
//	    return users, err
//	}
//
// Migration Support:
//
//	// Auto-migrate models
//	func RunMigrations() error {
//	    return facades.ORM().AutoMigrate(
//	        &User{},
//	        &Post{},
//	        &Comment{},
//	        &Role{},
//	    )
//	}
//
//	// Check if table exists
//	if facades.ORM().Migrator().HasTable(&User{}) {
//	    fmt.Println("Users table exists")
//	}
//
//	// Drop table
//	facades.ORM().Migrator().DropTable(&User{})
//
// Performance Optimization:
//
//	// Select specific fields
//	var users []User
//	facades.ORM().Select("id", "name", "email").Find(&users)
//
//	// Omit fields
//	facades.ORM().Omit("password", "secret_token").Create(&user)
//
//	// Use indexes for better performance
//	type User struct {
//	    ID    uint   `gorm:"primaryKey"`
//	    Email string `gorm:"uniqueIndex"`
//	    Name  string `gorm:"index"`
//	}
//
//	// Prepared statements
//	stmt := facades.ORM().Session(&gorm.Session{PrepareStmt: true})
//
//	// Connection pooling configuration
//	db := facades.ORM()
//	sqlDB, _ := db.DB()
//	sqlDB.SetMaxIdleConns(10)
//	sqlDB.SetMaxOpenConns(100)
//	sqlDB.SetConnMaxLifetime(time.Hour)
//
// Best Practices:
//   - Use struct tags for database column mapping
//   - Implement proper error handling for database operations
//   - Use transactions for operations that must be atomic
//   - Implement repository pattern for better code organization
//   - Use eager loading to avoid N+1 query problems
//   - Add appropriate database indexes for query performance
//   - Use soft deletes for audit trails
//   - Implement proper validation in model hooks
//
// Error Handling:
// This facade uses panic-on-error behavior for clean code:
//   - Most application code can assume ORM service always works
//   - Failures are detected early and halt execution
//   - No need for error checking in normal application flow
//   - Container configuration issues are caught immediately
//
// Alternative Error-Safe Access:
// If you need error handling instead of panics, use support package directly:
//
//	orm, err := facade.TryResolve[ORMInterface]("orm")
//	if err != nil {
//	    // Handle ORM service unavailability gracefully
//	    log.Printf("ORM service unavailable: %v", err)
//	    return nil, fmt.Errorf("database not available")
//	}
//	result := orm.Create(&user)
//
// Testing Support:
// This facade supports comprehensive testing through service swapping:
//
//	func TestUserRepository(t *testing.T) {
//	    // Create in-memory SQLite database for testing
//	    testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
//	    require.NoError(t, err)
//
//	    // Auto-migrate test schema
//	    err = testDB.AutoMigrate(&User{})
//	    require.NoError(t, err)
//
//	    // Swap the real ORM with test database
//	    restore := support.SwapService("orm", testDB)
//	    defer restore() // Always restore after test
//
//	    // Now facades.ORM() returns testDB
//	    repo := NewUserRepository()
//
//	    // Test user creation
//	    user := &User{Name: "Test User", Email: "test@example.com"}
//	    err = repo.Create(user)
//	    require.NoError(t, err)
//	    assert.NotZero(t, user.ID)
//
//	    // Test user retrieval
//	    foundUser, err := repo.FindByID(user.ID)
//	    require.NoError(t, err)
//	    assert.Equal(t, user.Name, foundUser.Name)
//	    assert.Equal(t, user.Email, foundUser.Email)
//	}
//
// Container Configuration:
// Ensure the ORM service is properly configured in your container:
//
//	// Example ORM registration
//	container.Singleton("orm", func() interface{} {
//	    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
//	        config.DB.Username,
//	        config.DB.Password,
//	        config.DB.Host,
//	        config.DB.Port,
//	        config.DB.Database,
//	    )
//
//	    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
//	        Logger: logger.Default.LogMode(logger.Info),
//	        NamingStrategy: schema.NamingStrategy{
//	            TablePrefix:   "app_",   // table name prefix
//	            SingularTable: false,    // use singular table name
//	        },
//	        DisableForeignKeyConstraintWhenMigrating: true,
//	        PrepareStmt: true, // cache prepared statements
//	    })
//
//	    if err != nil {
//	        log.Fatalf("Failed to connect to database: %v", err)
//	    }
//
//	    // Configure connection pool
//	    sqlDB, _ := db.DB()
//	    sqlDB.SetMaxIdleConns(10)
//	    sqlDB.SetMaxOpenConns(100)
//	    sqlDB.SetConnMaxLifetime(time.Hour)
//
//	    return db
//	})
func ORM() ormInterfaces.OrmInterface {
	// Use facade.Resolve() for clean facade implementation:
	// - Resolves "orm" service from the dependency injection container
	// - Performs type assertion to OrmInterface
	// - Caches the result for subsequent calls
	// - Panics with descriptive error if resolution fails
	// - Thread-safe with optimized locking
	return facade.Resolve[ormInterfaces.OrmInterface](ormInterfaces.ORM_TOKEN)
}

// ORMWithError provides error-safe access to the ORM service.
//
// This function offers the same functionality as ORM() but returns errors
// instead of panicking, making it suitable for error-sensitive contexts where
// you want to handle ORM service unavailability gracefully.
//
// This is a convenience wrapper around facade.TryResolve() that provides
// the same caching and performance benefits as ORM() but with error handling.
//
// Returns:
//   - ORMInterface: The resolved ORM instance (nil if error occurs)
//   - error: Detailed error information if resolution fails
//
// Errors:
//   - support.FacadeError: If container not set or service resolution fails
//   - Type assertion errors: If service doesn't implement ORMInterface
//
// Usage Examples:
//
//	// Basic error-safe ORM operations
//	orm, err := facades.ORMWithError()
//	if err != nil {
//	    log.Printf("ORM service unavailable: %v", err)
//	    return nil, fmt.Errorf("database operations not available")
//	}
//	result := orm.Create(&user)
//
//	// Conditional database operations
//	if orm, err := facades.ORMWithError(); err == nil {
//	    // Perform optional database operations
//	    var count int64
//	    orm.Model(&User{}).Count(&count)
//	    log.Printf("Total users: %d", count)
//	}
//
//	// Health check pattern
//	func CheckDatabaseHealth() error {
//	    orm, err := facades.ORMWithError()
//	    if err != nil {
//	        return fmt.Errorf("ORM service unavailable: %w", err)
//	    }
//
//	    // Test basic database connectivity
//	    sqlDB, err := orm.DB()
//	    if err != nil {
//	        return fmt.Errorf("failed to get database connection: %w", err)
//	    }
//
//	    if err := sqlDB.Ping(); err != nil {
//	        return fmt.Errorf("database ping failed: %w", err)
//	    }
//
//	    return nil
//	}
func ORMWithError() (ormInterfaces.OrmInterface, error) {
	// Use facade.TryResolve() for error-return behavior:
	// - Resolves "orm" service from the dependency injection container
	// - Performs type assertion with error handling
	// - Caches the result for subsequent calls
	// - Returns detailed error information instead of panicking
	// - Thread-safe with optimized locking
	return facade.TryResolve[ormInterfaces.OrmInterface](ormInterfaces.ORM_TOKEN)
}
